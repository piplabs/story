# DKG 集成测试（Go）

基于 **Go 标准 testing 框架** 的 DKG 集成测试，面向已运行的 devnet（如 5 个 enable DKG 的 validator）。用例与 [docs/DKG_Integration_Test_Cases.md](../../docs/DKG_Integration_Test_Cases.md) 中的 148 条一一对应。

## 架构

- **build tag**：`//go:build integration`，默认 `go test ./...` 不跑集成测试，需显式 `-tags=integration`。
- **TestMain**：建立 gRPC 连接、创建 DKG Query 客户端、注入 Harness（轮询间隔、超时等）。
- **Harness**：封装 `dkgtypes.QueryClient`，提供 `Params`、`GetLatestDKGNetwork`、`GetVerifiedRegistrations`、`WaitForStage`、`WaitForVerifiedCount` 等。
- **TestCase**：每条用例 ID（如 IT-E2E-01）、Priority（P0/P1/P2/P3）、Description、可选 `SkipIfLive`、`Run func(t *testing.T, h *Harness)`。
- **分组**：`P0Cases()`、`P1Cases()`、`P2Cases()`、`P3Cases()`，由 `AllCases()` 汇总为 148 条。

## 环境

- **STORY_GRPC_URL**（必填）：任意一个 validator 的 gRPC 地址，例如 `localhost:9090` 或 `1.2.3.4:9090`。
- **DKG_POLL_INTERVAL**（可选）：轮询间隔，如 `10s`。
- **DKG_MAX_WAIT**（可选）：单步最大等待时间，如 `600s`。
- **链上交互**（可选，用于 IT-BB-03、upgrade_scheduled、IT-CDR-* 等造景）：
  - **STORY_ETH_RPC_URL**：执行层 JSON-RPC，如 `http://127.0.0.1:8545`。
  - **DKG_CONTRACT_ADDRESS**：DKG 合约地址，空则用 predeploy 默认。
  - **DKG_SIGNER_PRIVATE_KEY**：调用 DKG 的私钥（hex，须为合约 owner），用于 `scheduleUpgrade` / `cancelUpgrade`。
  - 未配置时 `Harness.ChainClient` 为 Noop，仅读用例不受影响。

## 运行

在仓库根目录执行：

```bash
# 仅编译/运行集成测试（需能访问 devnet gRPC）
export STORY_GRPC_URL=localhost:9090
go test -tags=integration -v ./tests/integration/dkg/

# 仅跑 P0
go test -tags=integration -v -run TestDKG_P0 ./tests/integration/dkg/

# 仅跑 P1
go test -tags=integration -v -run TestDKG_P1 ./tests/integration/dkg/

# 跑全部 156 条（会跳过需故障注入/特殊拓扑的用例）
go test -tags=integration -v -run TestDKG_All ./tests/integration/dkg/

# 按优先级跑
go test -tags=integration -v -run TestDKG_ByPriority/P0 ./tests/integration/dkg/

# 连通性检查
go test -tags=integration -v -run TestDKG_Params ./tests/integration/dkg/
```

## 长短周期切换跑全部测试

需要**自动切换短周期 → 跑阶段转换类用例 → 再切回默认周期 → 跑快照/造景/upgrade 类**时，用 **TestDKG_FullSuite**。前提：`DKG_DRIVER=script`，且脚本能 SSH 到 devnet 节点执行 `configure_periods.sh` 和 `reset_devnet.sh`。

### 方式一：全自动（推荐）

在 **story 仓库根目录** 下：

```bash
# 1. 加载配置（含 STORY_RPC_URL、STORY_ETH_RPC_URL、DKG_DRIVER=script、DKG_SCENARIO_SCRIPTS_DIR 等）
source tests/integration/dkg/config.env

# 2. 确保脚本目录可访问（默认 scripts/dkg_e2e/scenarios，内含 configure_periods.sh）
#    configure_periods.sh 会 SSH 到各 validator 改 upgrades.go、重编译、调 reset_devnet.sh

# 3. 跑全套（内部：短周期 → 跑非 upgrade 用例 → 恢复默认周期 → 跑长周期/upgrade 用例）
go test -tags=integration -v -run TestDKG_FullSuite -timeout 7200s ./tests/integration/dkg/
```

可选环境变量（周期单位为**区块数**）：

- 短周期：`DKG_SHORT_REGISTRATION`（默认 20）、`DKG_SHORT_DEALING`（80）、`DKG_SHORT_FINALIZATION`（50）、`DKG_SHORT_ACTIVE`（20）
- 默认周期：`DKG_DEFAULT_REGISTRATION`（200）、`DKG_DEFAULT_DEALING`（300）、`DKG_DEFAULT_FINALIZATION`（300）、`DKG_DEFAULT_ACTIVE`（200）

### 方式二：手动切周期再跑

没有 ScriptDriver 或不能跑脚本时，可手动改周期、重置链，再分两次跑：

```bash
source tests/integration/dkg/config.env

# 1. 手动把 devnet 调成短周期（改 upgrades.go + 重编译 + reset_devnet.sh），等 DKG 激活
# 2. 只跑“需要阶段转换”的用例
go test -tags=integration -v -run TestDKG_ShortPeriod -timeout 3600s ./tests/integration/dkg/

# 3. 手动恢复默认周期并重置链，等 DKG 激活
# 4. 只跑“快照/不依赖阶段推进”的用例
go test -tags=integration -v -run TestDKG_DefaultPeriod -timeout 3600s ./tests/integration/dkg/
```

upgrade 相关用例（IT-BB-03、IT-E2E-06、IT-UPG-* 等）在 DefaultPeriod 里会跑，需配置好 `STORY_ETH_RPC_URL` 和 `DKG_SIGNER_PRIVATE_KEY`。

## 用例说明

- **P0（1 条）**：IT-E2E-01 完整快乐路径，轮询直到 Registration → Dealing → Finalization → Active，并校验 GlobalPublicKey。
- **P1（105 条）**：BeginBlocker、Registration、Dealing、Finalization、Active/Skip/E2E。凡仅通过查询链上状态可验证的均已实现 `Run`；需故障注入或单验证者等拓扑的标 `SkipIfLive`，运行时自动跳过。
- **P2（36 条）**：Resume、Upgrade、Encrypt/Decrypt；IT-ENC-01、IT-DEC-01 等在 Active 且存在 GlobalPublicKey 时做存在性校验。
- **P3（6 条）**：边界与压测；IT-EDGE-01 等在可观测时做校验。

未实现 `Run` 或标有 `SkipIfLive` 的用例在 live devnet 上会 `t.Skip(...)`，便于 CI 只跑可自动化部分，其余依赖人工或故障注入环境。
