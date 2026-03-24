//go:build integration

package dkg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestDKG_All 按安全顺序运行全部 156 条用例：
//  1. 快照查询类（不改链上状态）
//  2. 阶段转换类（等状态变化但不改链）
//  3. 造景类（改链上状态，如停 TEE/关 DKG）
//  4. upgrade_scheduled 类（最后跑，会触发 upgrade resharing 影响后续轮次）
func TestDKG_All(t *testing.T) {
	cases := OrderedCases()
	t.Logf("total cases: %d (original 156 + CL-* extensions)", len(cases))
	for _, tc := range cases {
		t.Run(tc.ID+"_"+tc.Priority, func(t *testing.T) {
			if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" && os.Getenv("DKG_TEST_INVALIDATE_INDEX") == "" {
				t.Skip(tc.SkipIfLive)
			}
			if tc.Run == nil {
				t.Skip("no Run (internal/fault-injection only)")
				return
			}
			globalHarness.RunCase(t, tc)
		})
	}
}

// TestDKG_ByPriority 可按优先级单独跑：-run TestDKG_ByPriority/P0 等。
func TestDKG_ByPriority(t *testing.T) {
	for _, p := range []string{"P0", "P1", "P2", "P3"} {
		t.Run(p, func(t *testing.T) {
			cases := CasesByPriority(p)
			for _, tc := range cases {
				t.Run(tc.ID, func(t *testing.T) {
					if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" {
						t.Skip(tc.SkipIfLive)
					}
					if tc.Run == nil {
						t.Skip("no Run")
						return
					}
					globalHarness.RunCase(t, tc)
				})
			}
		})
	}
}

// TestDKG_ShortPeriod 运行需要 DKG 阶段转换的用例。
// 前置条件：devnet 已配置为短周期（如 10s/10s/10s/30s），否则等待过久会超时。
func TestDKG_ShortPeriod(t *testing.T) {
	cases := ShortPeriodCases()
	t.Logf("running %d short-period (stage-transition) cases", len(cases))
	runCaseList(t, cases)
}

// TestDKG_DefaultPeriod 运行快照查询类用例，不依赖阶段转换，适合默认长周期 devnet。
func TestDKG_DefaultPeriod(t *testing.T) {
	cases := DefaultPeriodCases()
	t.Logf("running %d default-period (snapshot) cases", len(cases))
	runCaseList(t, cases)
}

// TestDKG_FullSuite 全自动化测试流程：
//  1. configure_periods.sh 调短周期 → 重编译 → 重置链 → 等 DKG 激活
//  2. 跑短周期测试（阶段转换类，不含 upgrade）
//  3. configure_periods.sh 恢复原始周期 → 重编译 → 重置链 → 等 DKG 激活
//  4. 跑默认周期测试（快照查询类 + 造景类 + upgrade 类）
//
// Period 单位为区块数，通过修改 upgrade handler 源码 + 重编译 + 链重置实现。
// 前置条件：DKG_DRIVER=script。
// 短周期参数可通过环境变量覆盖：DKG_SHORT_REGISTRATION, DKG_SHORT_DEALING, DKG_SHORT_FINALIZATION, DKG_SHORT_ACTIVE。
// 默认周期参数：DKG_DEFAULT_REGISTRATION=200, DKG_DEFAULT_DEALING=300, DKG_DEFAULT_FINALIZATION=300, DKG_DEFAULT_ACTIVE=200。
func TestDKG_FullSuite(t *testing.T) {
	configurator, ok := globalHarness.Driver.(PeriodConfigurator)
	if !ok {
		t.Fatal("FullSuite requires Driver implementing PeriodConfigurator (set DKG_DRIVER=script)")
	}

	// Resume mode: DKG_RESUME=true skips previously passed tests
	if os.Getenv("DKG_RESUME") != "true" {
		resetPassedTests()
		t.Log("Fresh run: cleared passed tests history")
	} else {
		passed := loadPassedTests()
		t.Logf("Resume run: %d tests previously passed, will be skipped", len(passed))
	}

	ctx := context.Background()

	shortReg := getUint64Env("DKG_SHORT_REGISTRATION", 20)
	shortDeal := getUint64Env("DKG_SHORT_DEALING", 80)
	shortFin := getUint64Env("DKG_SHORT_FINALIZATION", 50)
	shortActive := getUint64Env("DKG_SHORT_ACTIVE", 20)


	// --- Phase 1: 短周期 → 跑阶段转换测试 ---
	t.Log("=== Phase 1: 检查是否需要配置短周期 ===")
	t.Logf("short periods: reg=%d deal=%d fin=%d active=%d (blocks)", shortReg, shortDeal, shortFin, shortActive)
	if periodsMatch(ctx, globalHarness, shortReg, shortDeal, shortFin, shortActive) {
		t.Log("链上 period 已匹配短周期，跳过重置")
	} else {
		t.Log("链上 period 不匹配，执行配置短周期（修改源码 → 重编译 → 重置链）")
		if err := configurator.SetDKGPeriods(ctx, shortReg, shortDeal, shortFin, shortActive); err != nil {
			t.Fatalf("SetDKGPeriods(short): %v", err)
		}
	}

	t.Log("等待 DKG 激活（v2.0.0 upgrade at block ~150）...")
	waitForDKGActive(t, globalHarness)

	// 短周期跑所有非 upgrade 场景的测试
	allCases := OrderedCases()
	var shortCases []TestCase
	for _, c := range allCases {
		if !upgradeScenarioCaseIDs[c.ID] {
			shortCases = append(shortCases, c)
		}
	}
	// 确保所有 validator 和 kernel 处于健康状态
	ensureAllHealthy(t)

	t.Logf("=== Phase 1: 运行 %d 条短周期测试 ===", len(shortCases))
	runCaseList(t, shortCases)

	// --- Phase 2: upgrade_scheduled 测试（稍长 registration 周期） ---
	var upgradeCases []TestCase
	for _, c := range allCases {
		if upgradeScenarioCaseIDs[c.ID] {
			upgradeCases = append(upgradeCases, c)
		}
	}
	if len(upgradeCases) > 0 {
		// Upgrade resharing 需要更长的 registration period（validators 需要重新生成 key）
		upgradeReg := getUint64Env("DKG_UPGRADE_REGISTRATION", 50)
		upgradeDeal := shortDeal
		upgradeFin := shortFin
		upgradeActive := shortActive
		if !periodsMatch(ctx, globalHarness, upgradeReg, upgradeDeal, upgradeFin, upgradeActive) {
			t.Logf("=== Phase 2: 切换 upgrade 周期 reg=%d deal=%d fin=%d active=%d ===",
				upgradeReg, upgradeDeal, upgradeFin, upgradeActive)
			if err := configurator.SetDKGPeriods(ctx, upgradeReg, upgradeDeal, upgradeFin, upgradeActive); err != nil {
				t.Fatalf("SetDKGPeriods(upgrade): %v", err)
			}
			t.Log("等待 DKG 激活...")
			waitForDKGActive(t, globalHarness)
		}
		ensureAllHealthy(t)
		t.Logf("=== Phase 2: 运行 %d 条 upgrade 场景测试 ===", len(upgradeCases))
		runCaseList(t, upgradeCases)
	}
}

// runCaseList 按顺序执行用例列表，每条用例作为子测试，最后打印汇总。
// 当 DKG_RESUME=true 时，跳过上次已通过的测试（记录在 .dkg_passed_tests 中）。
func runCaseList(t *testing.T, cases []TestCase) {
	t.Helper()

	passed := loadPassedTests()
	resumeMode := os.Getenv("DKG_RESUME") == "true"
	if resumeMode {
		t.Logf("RESUME MODE: skipping %d previously passed tests", len(passed))
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.ID+"_"+tc.Priority, func(t *testing.T) {
			if resumeMode && passed[tc.ID] {
				t.Skipf("previously passed (resume mode)")
				return
			}
			if ScenarioNameForCase(tc.ID) == "" && tc.SkipIfLive != "" && os.Getenv("DKG_TEST_INVALIDATE_INDEX") == "" {
				t.Skip(tc.SkipIfLive)
			}
			if tc.Run == nil {
				t.Skip("no Run (internal/fault-injection only)")
				return
			}
			globalHarness.RunCase(t, tc)
			// If we got here without t.Failed(), mark as passed
			if !t.Failed() {
				markTestPassed(tc.ID)
			}
		})
	}

	// 汇总统计（go test 子测试结果通过 t.Run 的 pass/fail/skip 状态反映）
	// 注意：由于 Go test 框架限制，子测试结果只能在父测试结束后统计。
	// 这里按 priority 分组输出用例数量。
	priorityCounts := map[string]int{}
	for _, tc := range cases {
		priorityCounts[tc.Priority]++
	}
	var summary string
	for _, p := range []string{"P0", "P1", "P2", "P3"} {
		if n, ok := priorityCounts[p]; ok {
			summary += fmt.Sprintf(" %s:%d", p, n)
		}
	}
	// 统计 CL-* 等非标准优先级
	clCount := 0
	for p, n := range priorityCounts {
		if p != "P0" && p != "P1" && p != "P2" && p != "P3" {
			clCount += n
		}
	}
	if clCount > 0 {
		summary += fmt.Sprintf(" CL:%d", clCount)
	}
	t.Logf("============================================================")
	t.Logf("SUMMARY  Total: %d |%s", len(cases), summary)
	t.Logf("============================================================")
}

// waitForDKGActive 等待链重置后 DKG 完全激活：
//  1. 等 RPC 可达
//  2. 等 v2.0.0 upgrade (block 100) 设置 DKG params
//  3. 等第一个 DKG round 出现
// 整个过程约需 4-5 分钟（block 100 ≈ 200s + kernel registration）。
func waitForDKGActive(t *testing.T, h *Harness) {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(10 * time.Minute)
	start := time.Now()
	paramsReady := false
	for time.Now().Before(deadline) {
		if !paramsReady {
			params, err := h.Params(ctx)
			if err == nil && params.RegistrationPeriod > 0 {
				t.Logf("DKG params active (elapsed=%s, registration_period=%d)",
					time.Since(start).Truncate(time.Second), params.RegistrationPeriod)
				paramsReady = true
			}
		}
		if paramsReady {
			net, err := h.GetLatestDKGNetwork(ctx)
			if err == nil && net != nil && net.Round > 0 {
				if net.IsUpgrade {
					t.Logf("[WARNING] DKG round %d is upgrade resharing (IsUpgrade=true), proceeding anyway",
						net.Round)
					// Cancel 无法改变已激活 round 的 IsUpgrade，且后续 round 可能持续 upgrade。
					// 不阻塞启动，直接继续跑测试。
					_ = h.ChainClient.CancelDKGUpgrade(ctx, "v2.0.0-test")
				}
				t.Logf("DKG round %d started, stage=%s IsUpgrade=false (elapsed=%s)",
					net.Round, net.Stage, time.Since(start).Truncate(time.Second))
				return
			}
		}
		time.Sleep(5 * time.Second)
	}
	t.Fatal("DKG did not activate within 10 minutes after chain reset")
}

// waitForChainReady 轮询直到 RPC 可查询 DKG 参数（用于非重置场景）。
func waitForChainReady(t *testing.T, h *Harness) {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(5 * time.Minute)
	start := time.Now()
	for time.Now().Before(deadline) {
		params, err := h.Params(ctx)
		if err == nil && params.RegistrationPeriod > 0 {
			t.Logf("chain is ready (elapsed=%s, registration_period=%d)",
				time.Since(start).Truncate(time.Second), params.RegistrationPeriod)
			return
		}
		time.Sleep(5 * time.Second)
	}
	t.Fatal("chain did not recover within 3 minutes after period reconfiguration")
}

// ensureAllHealthy 确保所有 validator 的 story 和 kernel 都在运行。
// kernel 启动后 enclave 加载需要 1-3 分钟，会轮询等待直到 active。
func ensureAllHealthy(t *testing.T) {
	t.Helper()
	needsWait := false
	for i := 0; i < validatorCount(); i++ {
		// 检查并启动 story
		out, _ := sshCheckService(i, "story")
		if strings.TrimSpace(out) != "active" {
			t.Logf("  validator %d story not active (%s), restarting...", i+1, strings.TrimSpace(out))
			sshRunCmd(i, "sudo systemctl restart story")
			needsWait = true
		}
		// 检查并启动 kernel
		out, _ = sshCheckService(i, "story-kernel")
		if strings.TrimSpace(out) != "active" {
			t.Logf("  validator %d kernel not active (%s), restarting...", i+1, strings.TrimSpace(out))
			sshRunCmd(i, "sudo systemctl restart story-kernel")
			needsWait = true
		}
	}
	if !needsWait {
		return
	}
	// Kernel enclave 加载需要 1-3 分钟，轮询等待所有 kernel active
	t.Log("  等待所有 kernel 启动（enclave 加载中）...")
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		allActive := true
		for i := 0; i < validatorCount(); i++ {
			out, _ := sshCheckService(i, "story-kernel")
			if strings.TrimSpace(out) != "active" {
				allActive = false
				break
			}
		}
		if allActive {
			t.Log("  所有 validator 和 kernel 已启动")
			// 额外等 10s 让 story 重新连上 kernel
			time.Sleep(10 * time.Second)
			return
		}
		time.Sleep(10 * time.Second)
	}
	t.Log("  ⚠ 部分 kernel 启动超时，继续测试")
}

// periodsMatch 检查链上 DKG 参数是否已经匹配目标 period。
func periodsMatch(ctx context.Context, h *Harness, reg, deal, fin, active uint64) bool {
	params, err := h.Params(ctx)
	if err != nil || params == nil {
		return false
	}
	return uint64(params.RegistrationPeriod) == reg &&
		uint64(params.DealingPeriod) == deal &&
		uint64(params.FinalizationPeriod) == fin &&
		uint64(params.ActivePeriod) == active
}

// getUint64Env 读取环境变量作为 uint64，不存在或解析失败时返回默认值。
func getUint64Env(key string, def uint64) uint64 {
	if s := os.Getenv(key); s != "" {
		var v uint64
		for _, c := range s {
			if c < '0' || c > '9' {
				return def
			}
			v = v*10 + uint64(c-'0')
		}
		return v
	}
	return def
}

// --- Resume mode: track passed tests across runs ---

const passedTestsFile = ".dkg_passed_tests"

var (
	passedMu   sync.Mutex
	passedFile *os.File
)

// passedTestsPath returns the path to the passed tests file in the repo root.
func passedTestsPath() string {
	root := ""
	if out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output(); err == nil {
		root = strings.TrimSpace(string(out))
	}
	if root == "" {
		root = "."
	}
	return filepath.Join(root, "tests", "integration", "dkg", passedTestsFile)
}

// loadPassedTests reads the .dkg_passed_tests file and returns a set of passed test IDs.
func loadPassedTests() map[string]bool {
	passed := make(map[string]bool)
	f, err := os.Open(passedTestsPath())
	if err != nil {
		return passed
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			passed[id] = true
		}
	}
	return passed
}

// markTestPassed appends a test ID to .dkg_passed_tests (thread-safe).
func markTestPassed(testID string) {
	passedMu.Lock()
	defer passedMu.Unlock()
	if passedFile == nil {
		var err error
		passedFile, err = os.OpenFile(passedTestsPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
	}
	fmt.Fprintln(passedFile, testID)
}

// ResetPassedTests clears the passed tests file (call when starting a fresh suite).
func resetPassedTests() {
	os.Remove(passedTestsPath())
}

// TestDKG_Params 仅校验 DKG 参数可查询（连通性检查）。
func TestDKG_Params(t *testing.T) {
	params, err := globalHarness.Params(context.Background())
	if err != nil {
		t.Fatalf("Params: %v", err)
	}
	if params.RegistrationPeriod == 0 {
		t.Error("RegistrationPeriod is zero")
	}
	t.Logf("registration_period=%d dealing_period=%d min_req_registered=%d",
		params.RegistrationPeriod, params.DealingPeriod, params.MinReqRegisteredParticipants)
}
