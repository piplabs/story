//go:build integration

package dkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

// RunCase 执行单条用例；若该用例绑定了场景且 Harness.Driver 非空，会先 SetupScenario 再 Run、最后 TeardownScenario（含 panic 时也会恢复环境）。
func (h *Harness) RunCase(t *testing.T, tc TestCase) {
	t.Helper()
	ctx := context.Background()

	// 打印用例描述与期望
	scenario := ScenarioNameForCase(tc.ID)
	t.Logf("[%s] %s", tc.ID, tc.Description)
	if tc.Expected != "" {
		t.Logf("  Expected: %s", tc.Expected)
	}
	if scenario != "" {
		t.Logf("  Scenario: %s", scenario)
	}

	if scenario != "" && h.Driver != nil {
		if err := h.Driver.SetupScenario(ctx, scenario); err != nil {
			t.Fatalf("SetupScenario(%s): %v", scenario, err)
		}
		defer func() {
			if err := h.Driver.TeardownScenario(ctx, scenario); err != nil {
				t.Logf("TeardownScenario(%s): %v", scenario, err)
			}
			// 场景恢复后确保所有 validator 健康
			ensureAllHealthy(t)
		}()
		// 验证场景前置条件
		verifyScenarioPrecondition(t, h, scenario)
	}
	if tc.Run != nil {
		tc.Run(t, h)
	}
}

// TestCase 表示一条 DKG 集成测试用例，与 DKG_Integration_Test_Cases.md 中的 ID 一一对应。
type TestCase struct {
	ID          string // 如 "IT-E2E-01", "IT-BB-01"
	Priority    string // P0, P1, P2, P3
	Description string
	Expected    string // 期望结果描述，如 "stage=Active, GlobalPubKey non-empty"
	// SkipIfLive 非空时，在仅连接 live devnet 的情况下跳过（需故障注入或特殊拓扑时填写）。
	SkipIfLive string
	// Run 执行断言；若 SkipIfLive 为空则必须可仅通过查询链上状态完成验证。
	Run func(t *testing.T, h *Harness)
	// NeedsRoundWait 为 true 时，该测试需要等待 DKG 阶段转换（建议配合短周期 devnet 运行）。
	// false 表示仅查询当前链上状态，不依赖轮次推进，可瞬时完成。
	NeedsRoundWait bool
	// NeedsLongPeriod 为 true 时，该测试需要在长周期下运行（验证周期内不提前转换等时序行为）。
	NeedsLongPeriod bool
}

// check 统一打印期望值与实际值，并在不匹配时调用 t.Errorf。
func check(t *testing.T, field string, expected, actual interface{}) bool {
	t.Helper()
	pass := fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
	status := "PASS"
	if !pass {
		status = "FAIL"
	}
	t.Logf("  %s: expected=%v actual=%v → %s", field, expected, actual, status)
	if !pass {
		t.Errorf("%s mismatch: expected %v, got %v", field, expected, actual)
	}
	return pass
}

// checkTrue 打印条件断言结果。
func checkTrue(t *testing.T, field string, condition bool, actual string) bool {
	t.Helper()
	status := "PASS"
	if !condition {
		status = "FAIL"
	}
	t.Logf("  %s: %s → %s", field, actual, status)
	if !condition {
		t.Errorf("%s: %s", field, actual)
	}
	return condition
}

// AllCases 返回全部用例（含原有 156 条 + 新增 CL-* 系列）。
func AllCases() []TestCase {
	var out []TestCase
	out = append(out, P0Cases()...)
	out = append(out, P1Cases()...)
	out = append(out, P2Cases()...)
	out = append(out, P3Cases()...)
	out = append(out, CLCases()...)
	return out
}

// upgradeScenarioCaseIDs lists ALL test case IDs that call scheduleUpgrade (directly or via alias).
// These trigger upgrade resharing rounds and must run LAST to avoid polluting other tests.
// Source: grep for scheduleUpgradeHelper / ScheduleDKGUpgrade in scenarios_run.go.
var upgradeScenarioCaseIDs = map[string]bool{
	"IT-BB-03": true, "IT-BB-09": true,
	"IT-REG-03": true, "IT-REG-13": true, "IT-REG-14": true,
	"IT-DL-09": true, "IT-ACT-05": true, "IT-E2E-06": true,
	"IT-UPG-01": true, "IT-UPG-02": true, "IT-UPG-03": true, "IT-UPG-04": true, "IT-UPG-05": true,
}

// OrderedCases 返回按安全执行顺序排列的全部用例：
//  1. 快照查询类（不改链上状态，不等阶段转换）
//  2. 阶段转换类（不改链上状态，但等阶段转换）
//  3. 造景类（改链上状态，但不含 upgrade）
//  4. upgrade_scheduled 类（最后跑，因为会触发 upgrade resharing 影响后续轮次）
func OrderedCases() []TestCase {
	all := AllCases()
	var snapshot, stageWait, scenario, upgrade []TestCase
	for _, c := range all {
		if upgradeScenarioCaseIDs[c.ID] {
			upgrade = append(upgrade, c)
		} else if ScenarioNameForCase(c.ID) != "" {
			scenario = append(scenario, c)
		} else if c.NeedsRoundWait {
			stageWait = append(stageWait, c)
		} else {
			snapshot = append(snapshot, c)
		}
	}
	var out []TestCase
	out = append(out, snapshot...)
	out = append(out, stageWait...)
	out = append(out, scenario...)
	out = append(out, upgrade...)
	return out
}

// CasesByPriority 返回按优先级过滤的用例。
func CasesByPriority(priority string) []TestCase {
	var out []TestCase
	for _, c := range AllCases() {
		if c.Priority == priority {
			out = append(out, c)
		}
	}
	return out
}

// ShortPeriodCases 返回需要等待 DKG 阶段转换的用例。
// 这些测试需要将 devnet 配置为短周期（如 10s/10s/10s/30s），否则等待时间过长。
func ShortPeriodCases() []TestCase {
	var out []TestCase
	for _, c := range AllCases() {
		if c.NeedsRoundWait {
			out = append(out, c)
		}
	}
	return out
}

// DefaultPeriodCases 返回仅查询链上状态的用例（快照式），不依赖阶段转换。
// 这些测试在默认长周期（86400s）下即可运行，无需调整 devnet 配置。
func DefaultPeriodCases() []TestCase {
	var out []TestCase
	for _, c := range AllCases() {
		if !c.NeedsRoundWait {
			out = append(out, c)
		}
	}
	return out
}

// assertStage 要求当前最新轮阶段为 s；否则报错。
func assertStage(t *testing.T, h *Harness, s dkgtypes.DKGStage) {
	ctx := context.Background()
	net, err := h.GetLatestDKGNetwork(ctx)
	if err != nil {
		t.Fatalf("GetLatestDKGNetwork: %v", err)
	}
	if net == nil {
		t.Fatal("GetLatestDKGNetwork: no network")
	}
	if net.Stage != s {
		t.Errorf("expected stage %s, got %s", s.String(), net.Stage.String())
	}
}

// verifyScenarioPrecondition 验证场景 setup 后前置条件是否真正生效。
func verifyScenarioPrecondition(t *testing.T, h *Harness, scenario string) {
	t.Helper()
	switch scenario {
	case "tee_down_one_node":
		// 验证目标 kernel 确实停了
		targetNode := 2 // 0-indexed, 对应 validator 3
		out, err := grepValidatorLog("Connected to kernel", targetNode)
		if err == nil && len(out) > 0 {
			// 有日志不代表 kernel 在跑，可能是历史日志；改用检查进程
			out2, _ := grepValidatorLog("kernel client Finalize request failed", targetNode)
			if len(out2) > 0 {
				t.Logf("  ✓ Precondition: validator %d kernel failure confirmed (log shows TEE errors)", targetNode+1)
			}
		}
		// SSH 检查进程是否真的停了
		output, sshErr := sshCheckService(targetNode, "story-kernel")
		if sshErr != nil {
			t.Logf("  ⚠ Cannot verify kernel status via SSH: %v", sshErr)
		} else {
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, fmt.Sprintf("validator %d kernel stopped", targetNode+1), !isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}

	case "insufficient_verified":
		// 验证 2/3 的 kernel 确实停了，只有 node 1 有 kernel
		for i := 1; i < validatorCount(); i++ {
			output, err := sshCheckService(i, "story-kernel")
			if err != nil {
				t.Logf("  ⚠ Cannot verify node %d kernel: %v", i+1, err)
				continue
			}
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, fmt.Sprintf("validator %d kernel stopped", i+1), !isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}
		// node 1 kernel 应该还在
		output, err := sshCheckService(0, "story-kernel")
		if err == nil {
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, "validator 1 kernel running", isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}

	case "insufficient_finalized":
		// 同 insufficient_verified，验证 node 2+ 的 kernel 停了
		for i := 1; i < validatorCount(); i++ {
			output, err := sshCheckService(i, "story-kernel")
			if err != nil {
				t.Logf("  ⚠ Cannot verify node %d kernel: %v", i+1, err)
				continue
			}
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, fmt.Sprintf("validator %d kernel stopped", i+1), !isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}

	case "dkg_disabled_one_node":
		// 验证目标 node DKG 被禁用（检查 [dkg] section 下 enable = false）
		targetNode := validatorCount() - 1 // last validator, 0-indexed
		output, err := sshRunCmd(targetNode, "sed -n '/^\\[dkg\\]/,/^\\[/{/^\\s*enable/p}' ~/.story/story/config/story.toml | head -1")
		if err != nil {
			t.Logf("  ⚠ Cannot verify DKG config: %v", err)
		} else {
			disabled := strings.Contains(output, "false")
			checkTrue(t, fmt.Sprintf("validator %d DKG disabled", targetNode+1), disabled,
				fmt.Sprintf("config=%q", strings.TrimSpace(output)))
		}

	case "resharing_second_round":
		// 验证当前有 active round（第一轮已完成）
		ctx := context.Background()
		net, err := h.GetLatestDKGNetwork(ctx)
		if err == nil && net != nil {
			checkTrue(t, "round >= 1", net.Round >= 1, fmt.Sprintf("round=%d", net.Round))
		}

	case "kernel_reconnect":
		// 验证目标 kernel 确实停了，但 story 还在运行
		targetNode := validatorCount() - 1
		output, err := sshCheckService(targetNode, "story-kernel")
		if err == nil {
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, fmt.Sprintf("validator %d kernel stopped", targetNode+1), !isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}
		output2, err2 := sshCheckService(targetNode, "story")
		if err2 == nil {
			isActive := strings.TrimSpace(output2) == "active"
			checkTrue(t, fmt.Sprintf("validator %d story still running", targetNode+1), isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output2)))
		}

	case "story_crash_one_validator":
		// 验证 validator 2 的 story 确实停了
		output, err := sshCheckService(1, "story")
		if err == nil {
			isActive := strings.TrimSpace(output) == "active"
			checkTrue(t, "validator 2 story stopped", !isActive,
				fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
		}

	case "all_nodes_crash":
		// 验证所有节点的 story 都停了
		for i := 0; i < validatorCount(); i++ {
			output, err := sshCheckService(i, "story")
			if err == nil {
				isActive := strings.TrimSpace(output) == "active"
				checkTrue(t, fmt.Sprintf("validator %d story stopped", i+1), !isActive,
					fmt.Sprintf("systemctl status=%q", strings.TrimSpace(output)))
			}
		}

	default:
		t.Logf("  (no precondition check for scenario %q)", scenario)
	}
}

// sshCheckService 通过 SSH 检查 systemd 服务状态，返回 "active (running)" 或其他。
func sshCheckService(nodeIndex int, service string) (string, error) {
	return sshRunCmd(nodeIndex, fmt.Sprintf("systemctl is-active %s 2>/dev/null || echo inactive", service))
}

// sshRunCmd 通过 SSH 在指定 validator 上执行命令。
func sshRunCmd(nodeIndex int, cmd string) (string, error) {
	targets := strings.TrimSpace(os.Getenv("DKG_SSH_TARGETS"))
	if targets == "" {
		return "", fmt.Errorf("DKG_SSH_TARGETS not set")
	}
	keyPath := strings.TrimSpace(os.Getenv("DKG_SSH_KEY"))
	if keyPath == "" {
		return "", fmt.Errorf("DKG_SSH_KEY not set")
	}
	targetList := strings.Split(targets, ",")
	if nodeIndex < 0 || nodeIndex >= len(targetList) {
		return "", fmt.Errorf("nodeIndex %d out of range (have %d)", nodeIndex, len(targetList))
	}
	target := strings.TrimSpace(targetList[nodeIndex])
	sshCmd := exec.Command("ssh", "-i", keyPath, "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=5", target, cmd)
	output, err := sshCmd.CombinedOutput()
	return string(output), err
}

// validatorCount 返回 validator 数量。
func validatorCount() int {
	s := os.Getenv("DKG_VALIDATOR_COUNT")
	if s == "" {
		return 3
	}
	n := 3
	fmt.Sscanf(s, "%d", &n)
	return n
}

// assertActiveWithGlobalKey 要求存在 Active 轮且 GlobalPublicKey 非空。
func assertActiveWithGlobalKey(t *testing.T, h *Harness) {
	t.Helper()
	ctx := context.Background()
	net, err := h.GetLatestActiveDKGNetwork(ctx)
	if err != nil {
		t.Fatalf("GetLatestActiveDKGNetwork: %v", err)
	}
	if net == nil {
		t.Fatal("no active DKG network")
		return
	}
	check(t, "stage", "DKG_STAGE_ACTIVE", net.Stage.String())
	checkTrue(t, "GlobalPublicKey", len(net.GlobalPublicKey) > 0,
		fmt.Sprintf("len=%d (need >0)", len(net.GlobalPublicKey)))
	t.Logf("  Result: round=%d globalPubKey=%d bytes", net.Round, len(net.GlobalPublicKey))
}
