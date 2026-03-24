//go:build integration

// Package dkg 提供 DKG 集成测试，面向已运行的 devnet（3 个 enable DKG 的 validator）。
// 运行: go test -tags=integration -v ./tests/integration/dkg/...
// 环境: STORY_RPC_URL（必填，CometBFT RPC 端点，默认 http://localhost:26657）
// 链上交互（可选）: STORY_ETH_RPC_URL, DKG_CONTRACT_ADDRESS, DKG_SIGNER_PRIVATE_KEY
package dkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
)

var (
	globalHarness *Harness
)

func TestMain(m *testing.M) {
	rpcURL := os.Getenv("STORY_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:26657"
	}
	conn, err := NewABCIConn(rpcURL)
	if err != nil {
		os.Stderr.WriteString("DKG integration: RPC connect failed: " + err.Error() + "\n")
		os.Exit(1)
	}

	var driver Driver = NoopDriver{}
	if os.Getenv("DKG_DRIVER") == "script" {
		driver = &ScriptDriver{Dir: os.Getenv("DKG_SCENARIO_SCRIPTS_DIR")}
	}

	chainClient, _ := NewEthChainClient()
	if chainClient == nil {
		chainClient = NoopChainClient{}
	}

	globalHarness = &Harness{
		Conn:         conn,
		QueryClient:  dkgtypes.NewQueryClient(conn),
		PollInterval: getDurationEnv("DKG_POLL_INTERVAL", 10*time.Second),
		MaxWait:      getDurationEnv("DKG_MAX_WAIT", 600*time.Second),
		Driver:       driver,
		ChainClient:  chainClient,
	}
	code := m.Run()
	os.Exit(code)
}

func getDurationEnv(key string, def time.Duration) time.Duration {
	if s := os.Getenv(key); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return def
}

// Harness 持有 RPC 连接、DKG Query 客户端、轮询/超时配置、可选的 Driver（造景）与 ChainClient（链上交互）。
type Harness struct {
	Conn        *ABCIConn
	QueryClient dkgtypes.QueryClient
	PollInterval time.Duration
	MaxWait      time.Duration
	Driver       Driver      // 非 nil 时，对需造景的用例会先 SetupScenario 再跑断言、最后 TeardownScenario
	ChainClient  ChainClient // 非 Noop 时可发链上交易（如 DKG.scheduleUpgrade），用于 upgrade_scheduled 等场景
}

// Params 返回当前 DKG 参数。
func (h *Harness) Params(ctx context.Context) (*dkgtypes.Params, error) {
	resp, err := h.QueryClient.Params(ctx, &dkgtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &resp.Params, nil
}

// GetLatestDKGNetwork 返回当前最新 DKG 轮；无轮时返回 nil, nil（需调用方判断 NotFound）。
func (h *Harness) GetLatestDKGNetwork(ctx context.Context) (*dkgtypes.DKGNetwork, error) {
	resp, err := h.QueryClient.GetLatestDKGNetwork(ctx, &dkgtypes.QueryGetLatestDKGNetworkRequest{})
	if err != nil {
		return nil, err
	}
	return &resp.Network, nil
}

// GetLatestActiveDKGNetwork 返回当前 Active 轮；无时返回 nil, error。
func (h *Harness) GetLatestActiveDKGNetwork(ctx context.Context) (*dkgtypes.DKGNetwork, error) {
	resp, err := h.QueryClient.GetLatestActiveDKGNetwork(ctx, &dkgtypes.QueryGetLatestActiveDKGNetworkRequest{})
	if err != nil {
		return nil, err
	}
	return &resp.Network, nil
}

// GetVerifiedRegistrations 返回指定轮的 Verified 注册列表。
func (h *Harness) GetVerifiedRegistrations(ctx context.Context, round uint32) ([]dkgtypes.DKGRegistration, error) {
	resp, err := h.QueryClient.GetAllVerifiedDKGRegistrations(ctx, &dkgtypes.QueryGetAllVerifiedDKGRegistrationsRequest{
		Round:             round,
		CodeCommitmentHex: "",
	})
	if err != nil {
		return nil, err
	}
	return resp.Registrations, nil
}

// WaitForStage 轮询直到当前轮阶段为 target 或超时；返回是否达到目标阶段。
func (h *Harness) WaitForStage(ctx context.Context, target dkgtypes.DKGStage) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	var lastStage dkgtypes.DKGStage
	for time.Now().Before(deadline) {
		net, err := h.GetLatestDKGNetwork(ctx)
		if err != nil || net == nil {
			time.Sleep(h.PollInterval)
			continue
		}
		if net.Stage == target {
			fmt.Printf("  ✓ stage=%s round=%d (%s)\n", target, net.Round, time.Since(start).Truncate(time.Second))
			return true
		}
		if net.Stage != lastStage {
			fmt.Printf("  … round=%d stage=%s (waiting for %s)\n", net.Round, net.Stage, target)
			lastStage = net.Stage
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	fmt.Printf("  ✗ timeout waiting for stage=%s (%s)\n", target, h.MaxWait)
	return false
}

// WaitForRoundStage 轮询直到指定 round 达到目标阶段或超时。
// 如果 round 变成 FAILED/ENDED 且目标不是这两个状态，返回 false（round 不会再推进）。
func (h *Harness) WaitForRoundStage(ctx context.Context, round uint32, target dkgtypes.DKGStage) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	var lastStage dkgtypes.DKGStage
	for time.Now().Before(deadline) {
		net, err := h.GetDKGNetwork(ctx, round)
		if err != nil || net == nil {
			time.Sleep(h.PollInterval)
			continue
		}
		// 精确匹配目标 stage，或者 stage 已经超过目标（但排除 FAILED/ENDED 误匹配）
		if net.Stage == target {
			fmt.Printf("  ✓ round=%d stage=%s (%s)\n", round, net.Stage, time.Since(start).Truncate(time.Second))
			return true
		}
		// 正常推进：Registration < Dealing < Finalization < Active，允许跳过
		if net.Stage > target && net.Stage < dkgtypes.DKGStageFailed {
			fmt.Printf("  ✓ round=%d stage=%s (past target %s) (%s)\n", round, net.Stage, target, time.Since(start).Truncate(time.Second))
			return true
		}
		// Round 已经 FAILED/ENDED 但目标不是 FAILED/ENDED → 不会再推进了
		if net.Stage >= dkgtypes.DKGStageFailed && target < dkgtypes.DKGStageFailed {
			fmt.Printf("  ✗ round=%d stage=%s (round failed, cannot reach %s)\n", round, net.Stage, target)
			return false
		}
		if net.Stage != lastStage {
			fmt.Printf("  … round=%d stage=%s → waiting for %s\n", round, net.Stage, target)
			lastStage = net.Stage
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	fmt.Printf("  ✗ round=%d timeout waiting for %s (%s)\n", round, target, time.Since(start).Truncate(time.Second))
	return false
}

// WaitForVerifiedCount 轮询直到当前轮 Verified 数 >= minCount 或超时。
func (h *Harness) WaitForVerifiedCount(ctx context.Context, round uint32, minCount int) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	lastCount := -1
	for time.Now().Before(deadline) {
		regs, err := h.GetVerifiedRegistrations(ctx, round)
		if err == nil && len(regs) >= minCount {
			fmt.Printf("  ✓ round=%d verified=%d (%s)\n", round, len(regs), time.Since(start).Truncate(time.Second))
			return true
		}
		count := 0
		if err == nil {
			count = len(regs)
		}
		if count != lastCount {
			fmt.Printf("  … round=%d verified=%d (need %d)\n", round, count, minCount)
			lastCount = count
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	fmt.Printf("  ✗ round=%d verified count timeout (%s)\n", round, time.Since(start).Truncate(time.Second))
	return false
}

// GetDKGNetwork 返回指定轮的 DKG 网络信息。
func (h *Harness) GetDKGNetwork(ctx context.Context, round uint32) (*dkgtypes.DKGNetwork, error) {
	resp, err := h.QueryClient.GetDKGNetwork(ctx, &dkgtypes.QueryGetDKGNetworkRequest{Round: round})
	if err != nil {
		return nil, err
	}
	return &resp.Network, nil
}

// GetAllDKGRegistrations 返回指定轮的所有注册。
func (h *Harness) GetAllDKGRegistrations(ctx context.Context, round uint32) ([]dkgtypes.DKGRegistration, error) {
	resp, err := h.QueryClient.GetAllDKGRegistrations(ctx, &dkgtypes.QueryGetAllDKGRegistrationsRequest{
		Round: round,
	})
	if err != nil {
		return nil, err
	}
	return resp.Registrations, nil
}

// GetFinalizedRegistrations 返回指定轮中 Status=Finalized 的注册。
func (h *Harness) GetFinalizedRegistrations(ctx context.Context, round uint32) ([]dkgtypes.DKGRegistration, error) {
	regs, err := h.GetAllDKGRegistrations(ctx, round)
	if err != nil {
		return nil, err
	}
	var out []dkgtypes.DKGRegistration
	for _, r := range regs {
		if r.Status == dkgtypes.DKGRegStatusFinalized {
			out = append(out, r)
		}
	}
	return out, nil
}

// WaitForRound 轮询直到最新轮 >= targetRound 或超时。
func (h *Harness) WaitForRound(ctx context.Context, targetRound uint32) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	var lastRound uint32
	for time.Now().Before(deadline) {
		net, err := h.GetLatestDKGNetwork(ctx)
		if err == nil && net != nil && net.Round >= targetRound {
			fmt.Printf("  ✓ round=%d (%s)\n", net.Round, time.Since(start).Truncate(time.Second))
			return true
		}
		if net != nil && net.Round != lastRound {
			fmt.Printf("  … round=%d (waiting for %d)\n", net.Round, targetRound)
			lastRound = net.Round
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	fmt.Printf("  ✗ timeout waiting for round=%d (%s)\n", targetRound, time.Since(start).Truncate(time.Second))
	return false
}

// WaitForFinalizedCount 轮询直到指定轮 Finalized 数 >= minCount 或超时。
func (h *Harness) WaitForFinalizedCount(ctx context.Context, round uint32, minCount int) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	fmt.Printf("[WaitForFinalizedCount] round=%d waiting for finalized>=%d (timeout=%s)\n", round, minCount, h.MaxWait)
	for time.Now().Before(deadline) {
		regs, err := h.GetFinalizedRegistrations(ctx, round)
		if err == nil && len(regs) >= minCount {
			fmt.Printf("[WaitForFinalizedCount] round=%d finalized=%d (elapsed=%s)\n", round, len(regs), time.Since(start).Truncate(time.Second))
			return true
		}
		count := 0
		if err == nil {
			count = len(regs)
		}
		fmt.Printf("[WaitForFinalizedCount] round=%d finalized=%d target=%d (elapsed=%s)\n", round, count, minCount, time.Since(start).Truncate(time.Second))
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	return false
}

// WaitForBlockHeight 轮询直到执行层区块高度 >= targetHeight 或超时。
func (h *Harness) WaitForBlockHeight(ctx context.Context, targetHeight uint64) bool {
	deadline := time.Now().Add(h.MaxWait)
	start := time.Now()
	fmt.Printf("[WaitForBlockHeight] waiting for height>=%d (timeout=%s)\n", targetHeight, h.MaxWait)
	for time.Now().Before(deadline) {
		bn, err := h.ChainClient.BlockNumber(ctx)
		if err == nil && bn >= targetHeight {
			fmt.Printf("[WaitForBlockHeight] reached height=%d (elapsed=%s)\n", bn, time.Since(start).Truncate(time.Second))
			return true
		}
		if err == nil {
			fmt.Printf("[WaitForBlockHeight] current=%d target=%d (elapsed=%s)\n", bn, targetHeight, time.Since(start).Truncate(time.Second))
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(h.PollInterval):
		}
	}
	return false
}

// IsNoopChainClient 判断 ChainClient 是否为 NoopChainClient。
func (h *Harness) IsNoopChainClient() bool {
	_, ok := h.ChainClient.(NoopChainClient)
	return ok
}

// GetCDRPartials 查询指定 uuid + requesterPubKey 的 partial decryption 提交列表。
func (h *Harness) GetCDRPartials(ctx context.Context, uuid uint32, requesterPubKeyHex string) (*dkgtypes.QueryGetCDRPartialsResponse, error) {
	return h.QueryClient.GetCDRPartials(ctx, &dkgtypes.QueryGetCDRPartialsRequest{
		Uuid:                uuid,
		RequesterPubKeyHex: requesterPubKeyHex,
	})
}

// GetAllDKGNetworks 返回所有 DKG 网络（轮次）信息。
func (h *Harness) GetAllDKGNetworks(ctx context.Context) ([]dkgtypes.DKGNetwork, error) {
	resp, err := h.QueryClient.GetAllDKGNetworks(ctx, &dkgtypes.QueryGetAllDKGNetworksRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Networks, nil
}

// WaitForActiveRoundGlobalKey 轮询直到存在 Active 轮且 GlobalPublicKey 非空。
func (h *Harness) WaitForActiveRoundGlobalKey(ctx context.Context) (*dkgtypes.DKGNetwork, bool) {
	deadline := time.Now().Add(h.MaxWait)
	for time.Now().Before(deadline) {
		net, err := h.GetLatestActiveDKGNetwork(ctx)
		if err == nil && net != nil && len(net.GlobalPublicKey) > 0 {
			return net, true
		}
		select {
		case <-ctx.Done():
			return nil, false
		case <-time.After(h.PollInterval):
		}
	}
	return nil, false
}

// grepValidatorLog SSH 到验证器并查询其 story journal 日志。
// 使用 DKG_SSH_TARGETS（逗号分隔的 user@ip:port）和 DKG_SSH_KEY。
func grepValidatorLog(pattern string, nodeIndex int) (string, error) {
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
		return "", fmt.Errorf("nodeIndex %d out of range (have %d targets)", nodeIndex, len(targetList))
	}
	target := strings.TrimSpace(targetList[nodeIndex])

	cmd := exec.Command(
		"ssh", "-i", keyPath,
		"-o", "StrictHostKeyChecking=no",
		target,
		fmt.Sprintf("journalctl -u story --no-pager --since \"30 min ago\" | grep -i %q | tail -5", pattern),
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
