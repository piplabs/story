//go:build integration

package dkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Driver 用于在测试前/后"造景"：控制 validator 的 DKG 开关、TEE 启停等。
// 若你拥有 validator 与 TEE 的控制权限，可实现本接口并通过环境变量注入，测试会先 SetupScenario 再跑断言、最后 TeardownScenario。
type Driver interface {
	// SetupScenario 执行场景前置步骤（如只保留 1 个 validator 开 DKG、停掉某节点 TEE）。
	SetupScenario(ctx context.Context, scenarioName string) error
	// TeardownScenario 恢复环境（如重新开启所有 DKG、启动 TEE）。
	TeardownScenario(ctx context.Context, scenarioName string) error
}

// PeriodConfigurator 用于自动化修改 devnet 的 DKG 周期参数（genesis 重启）。
// 若 Driver 同时实现了 PeriodConfigurator，TestDKG_FullSuite 会自动切换周期。
type PeriodConfigurator interface {
	// SetDKGPeriods 修改 devnet genesis 中的 DKG 周期参数（秒）并重启所有 validator。
	// 调用方应在返回后等待链恢复可查询状态。
	SetDKGPeriods(ctx context.Context, registration, dealing, finalization, active uint64) error
}

// NoopDriver 不执行任何操作；未配置 Driver 时使用。
type NoopDriver struct{}

func (NoopDriver) SetupScenario(context.Context, string) error    { return nil }
func (NoopDriver) TeardownScenario(context.Context, string) error { return nil }

// ScriptDriver 通过执行脚本目录下的 <scenario>/pre.sh 与 <scenario>/post.sh 来造景与恢复。
// 环境变量 DKG_SCENARIO_SCRIPTS_DIR 指向该目录（默认 ./scripts/dkg_e2e/scenarios）。
type ScriptDriver struct {
	Dir string
}

func (s *ScriptDriver) SetupScenario(ctx context.Context, scenarioName string) error {
	return s.runScript(ctx, scenarioName, "pre")
}

func (s *ScriptDriver) TeardownScenario(ctx context.Context, scenarioName string) error {
	return s.runScript(ctx, scenarioName, "post")
}

// SetDKGPeriods 执行 configure_periods.sh 脚本修改 devnet genesis 中的 DKG 周期参数并重启。
// 脚本路径: <Dir>/configure_periods.sh <registration> <dealing> <finalization> <active>
func (s *ScriptDriver) SetDKGPeriods(ctx context.Context, registration, dealing, finalization, active uint64) error {
	dir := s.resolveDir()
	path := filepath.Join(dir, "configure_periods.sh")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("configure_periods.sh not found at %s; please create it to automate period changes", path)
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", path,
		fmt.Sprintf("%d", registration),
		fmt.Sprintf("%d", dealing),
		fmt.Sprintf("%d", finalization),
		fmt.Sprintf("%d", active),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("configure_periods.sh: %w", err)
	}
	return nil
}

func (s *ScriptDriver) runScript(ctx context.Context, scenarioName, kind string) error {
	dir := s.resolveDir()
	path := filepath.Join(dir, scenarioName, kind+".sh")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // 脚本不存在则跳过
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(dir, scenarioName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", path, err)
	}
	return nil
}

func (s *ScriptDriver) resolveDir() string {
	dir := s.Dir
	if dir == "" {
		dir = os.Getenv("DKG_SCENARIO_SCRIPTS_DIR")
	}
	if dir == "" {
		dir = "scripts/dkg_e2e/scenarios"
	}
	if filepath.IsAbs(dir) {
		return dir
	}
	if root := repoRoot(); root != "" {
		return filepath.Join(root, dir)
	}
	return dir
}

func repoRoot() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Scenario 描述一个可"造景"的测试场景：名称、对应用例 ID、前置/恢复说明。
type Scenario struct {
	Name        string   // 目录名，与 ScriptDriver 下 <Name>/pre.sh、post.sh 对应
	CaseIDs     []string // 使用该场景的用例 ID
	SetupDesc   string   // 前置步骤说明（给人工或脚本实现参考）
	TeardownDesc string  // 恢复步骤说明
}

// AllScenarios 返回所有需造景的场景；实现 pre/post 脚本后即可自动化。
func AllScenarios() []Scenario {
	return []Scenario{
		{
			Name:        "insufficient_verified",
			CaseIDs:     []string{"IT-E2E-03", "IT-DL-02", "IT-SKIP-01"},
			SetupDesc:   "只保留 1 个 validator 开启 DKG，其余 2 个关闭 DKG（或停 story-kernel），使该轮 Verified < MinReq。",
			TeardownDesc: "恢复 3 个 validator 均开启 DKG 且 TEE 正常。",
		},
		{
			Name:        "insufficient_finalized",
			CaseIDs:     []string{"IT-E2E-04", "IT-ACT-02", "IT-ACT-03", "IT-SKIP-02"},
			SetupDesc:   "Dealing 阶段结束后，只让 1 个 validator 能 finalize（其余关闭 DKG 或停 TEE），使 finalizedCount < MinReq 或 < Threshold。",
			TeardownDesc: "恢复全部 3 个 validator DKG 与 TEE。",
		},
		{
			Name:        "dkg_disabled_one_node",
			CaseIDs:     []string{"IT-REG-04", "IT-DL-03", "IT-FN-02", "IT-ACT-04"},
			SetupDesc:   "关闭其中 1 个 validator 的 DKG（config dkg.enable=false 并重启该节点），其余 2 个保持开启。",
			TeardownDesc: "重新开启该节点的 DKG 并重启。",
		},
		{
			Name:        "tee_down_one_node",
			CaseIDs:     []string{"IT-REG-10", "IT-REG-11", "IT-DL-08", "IT-DL-10", "IT-FN-07", "IT-FN-08", "IT-ACT-12", "IT-RES-02", "IT-RES-03", "IT-RES-04", "IT-RES-05", "IT-EDGE-06"},
			SetupDesc:   "停止其中 1 个 validator 的 story-kernel（TEE），使该节点 CreateSession/GenerateAndSealKey/GenerateDeals 等失败或 Phase=Failed。",
			TeardownDesc: "重新启动该节点的 story-kernel。",
		},
		{
			Name:        "height_below_dkg_start",
			CaseIDs:     []string{"IT-BB-01"},
			SetupDesc:   "在一条新链或快照上，确保当前区块高度 < dkgStartBlock（如 10）；或查询时使用历史高度（若 RPC 支持）。通常需新起链或从创世重放。",
			TeardownDesc: "无需恢复（或切回正常链）。",
		},
		{
			Name:        "upgrade_scheduled",
			CaseIDs:     []string{"IT-BB-03", "IT-BB-09", "IT-REG-03", "IT-DL-09", "IT-ACT-05", "IT-E2E-06", "IT-UPG-01", "IT-UPG-02", "IT-UPG-04", "IT-UPG-05"},
			SetupDesc:   "通过链上或 Keeper 接口调用 UpgradeScheduled(activationHeight, version)；若为测试可设 activationHeight 为即将到达的高度。",
			TeardownDesc: "调用 UpgradeCancelled 或等待升级完成。",
		},
		{
			Name:        "resharing_second_round",
			CaseIDs:     []string{"IT-REG-02", "IT-BB-07", "IT-E2E-02"},
			SetupDesc:   "先跑完第一轮至 Active，再等待 Active 周期结束，或缩短 active_period 后推进区块，触发新轮 IsResharing=true。",
			TeardownDesc: "无需恢复。",
		},
		{
			Name:        "complaint_justification",
			CaseIDs:     []string{"IT-E2E-05"},
			SetupDesc:   "需要 TEE 或数据层面产生无效 deal（如 story-kernel 测试模式、或篡改某笔 deal 数据），使某 recipient 产生 complaint → response → justification，最终无效 dealer 被 invalidated。",
			TeardownDesc: "恢复正常 TEE 与数据。",
		},
		{
			Name:        "cdr_e2e",
			CaseIDs:     []string{"IT-CDR-01", "IT-CDR-08"},
			SetupDesc:   "确保有 Active 轮（GlobalPublicKey 已设置），CDR 合约已部署且可交互（STORY_ETH_RPC_URL + DKG_SIGNER_PRIVATE_KEY 已配置）。",
			TeardownDesc: "无需恢复。",
		},
	}
}

// ─── Mock Kernel 场景（用于 CL-* 系列需故障注入的用例）───

// AllMockKernelScenarios 返回需要 mock kernel 的场景定义。
func AllMockKernelScenarios() []Scenario {
	return []Scenario{
		{
			Name:         "mock_kernel_bad_sig",
			CaseIDs:      []string{"CL-PD-01", "CL-PD-02", "CL-PD-03"},
			SetupDesc:    "在一个 validator 上部署 mock kernel，PartialDecryptTDH2 返回伪造/错误签名的 partial。其余 validator 用真实 TEE。",
			TeardownDesc: "恢复该 validator 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_bad_deal",
			CaseIDs:      []string{"CL-JUST-01", "CL-JUST-02", "CL-JUST-03", "CL-VE-01", "CL-VE-02", "CL-VE-03", "CL-VE-06"},
			SetupDesc:    "在一个 validator 上部署 mock kernel，GenerateDeals 返回无效 share/oversized payload/garbage bytes。",
			TeardownDesc: "恢复该 validator 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_replay",
			CaseIDs:      []string{"CL-REPLAY-01", "CL-REPLAY-02", "CL-REPLAY-03"},
			SetupDesc:    "在一个 validator 上部署 mock kernel，记录 round N 的 deals/justifications/finalization sig，在 round N+1 重放。",
			TeardownDesc: "恢复该 validator 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_non_committee",
			CaseIDs:      []string{"CL-FEE-02"},
			SetupDesc:    "在一个非 DKG committee 的 validator 上部署 mock kernel，令其提交 partial decryption。",
			TeardownDesc: "恢复该 validator。",
		},
		{
			Name:         "mock_kernel_invalid_argument",
			CaseIDs:      []string{"CL-KERR-01"},
			SetupDesc:    "在 Validator 3 上部署 mock kernel，GenerateDeals 返回 gRPC codes.InvalidArgument（非重试错误）。验证 story 是否立即标记 session Failed 而不是无限重试。",
			TeardownDesc: "恢复 Validator 3 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_bad_dealer_finalize",
			CaseIDs:      []string{"CL-BADDEALER-01", "CL-BADDEALER-02", "CL-BADDEALER-03", "CL-BADDEALER-04", "CL-BADDEALER-05", "CL-BADDEALER-06"},
			SetupDesc:    "在 Validator 3 上部署 mock kernel (WithInvalidVSSDeal)，产出有效签名但无效 VSS share 的 deals。B/C 用真实 TEE。等待 complaint → justification → failed verification 流程完成。",
			TeardownDesc: "恢复 Validator 3 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_2_bad_dealers",
			CaseIDs:      []string{"CL-BADDEALER-07"},
			SetupDesc:    "在 Validator 2 和 3 上同时部署 mock kernel (WithInvalidVSSDeal)，仅 Validator 1 用真实 TEE。",
			TeardownDesc: "恢复两个 validator 的真实 story-kernel。",
		},
		{
			Name:         "mock_kernel_adaptive_dealer",
			CaseIDs:      []string{"CL-BADDEALER-08"},
			SetupDesc:    "在 Validator 3 上部署特制 mock kernel：GenerateDeals 返回无效 share，但 FinalizeDKG 计算时排除自己的 deal contribution 以匹配诚实多数 globalPubKey。",
			TeardownDesc: "恢复真实 story-kernel。",
		},
		{
			Name:         "kernel_restart_after_deals",
			CaseIDs:      []string{"CL-RESTART-01"},
			SetupDesc:    "等待 DKG 进入 Dealing 阶段、validator 3 的 ProcessDeals 完成后，SSH 重启其 story-kernel。rebuildInitDKG 应从磁盘 replay deals 恢复状态。",
			TeardownDesc: "确认 story-kernel 已正常运行，必要时手动重启。",
		},
		{
			Name:         "kernel_restart_before_finalize",
			CaseIDs:      []string{"CL-RESTART-02"},
			SetupDesc:    "等待 DKG 进入 Finalization 阶段（ProcessResponses 已完成），SSH 重启 validator 3 的 story-kernel。rebuildInitDKG 应从磁盘 replay deals+responses 恢复状态。此场景直接验证 PrivatePoly 未持久化的 bug fix。",
			TeardownDesc: "确认 story-kernel 已正常运行。",
		},
		{
			Name:         "kernel_restart_after_justification",
			CaseIDs:      []string{"CL-RESTART-03"},
			SetupDesc:    "先在 validator 3 部署 mock kernel 产出 bad deal → 触发 complaint → justification → 然后重启恢复真实 kernel。rebuildInitDKG 应从磁盘 replay deals+responses+justifications。",
			TeardownDesc: "恢复真实 story-kernel。",
		},
	}
}

// ScenarioNameForCase 返回包含该 caseID 的场景名称（若有）；用于调用 Driver.SetupScenario/TeardownScenario。
func ScenarioNameForCase(caseID string) string {
	for _, s := range AllScenarios() {
		for _, id := range s.CaseIDs {
			if id == caseID {
				return s.Name
			}
		}
	}
	for _, s := range AllMockKernelScenarios() {
		for _, id := range s.CaseIDs {
			if id == caseID {
				return s.Name
			}
		}
	}
	return ""
}
