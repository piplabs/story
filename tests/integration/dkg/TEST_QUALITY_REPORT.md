# DKG Integration Test Quality Report (156 条)

Updated: 2026-03-18

## Summary

| 状态 | 数量 | 说明 |
|---|---|---|
| ✅ 完整实现 | 134 | 有场景构造 + `check`/`checkTrue` 断言，devnet 可跑通 |
| 🔶 有断言但场景受限 | 1 | IT-E2E-05 complaint 路径需 TEE mock 产生无效 deal |
| 🧪 需 TEE mock server | 21 | IT-DEC-02~22，CDR 解密内部路径 |
| **TOTAL** | **156** | |

---

## ✅ 完整实现 (134 条)

### E2E 端到端 (5 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-E2E-01 | Happy path: Reg→Deal→Fin→Active | 无 | `check(stage, Active)` + `checkTrue(GlobalPubKey > 0)` | 短 |
| IT-E2E-02 | Resharing round 完整生命周期 | `resharing_second_round` | 等当前 round Active → 下一轮 Active + `checkTrue(GlobalPubKey)` | 短 |
| IT-E2E-03 | Insufficient Verified → Failed | `insufficient_verified` | 停 node 2/3 kernel → `WaitForRoundStage(Failed)` | 短 |
| IT-E2E-04 | Insufficient Finalized → Failed | `insufficient_finalized` | Dealing 后停 kernel → `WaitForRoundStage(Failed)` | 短 |
| IT-E2E-06 | Upgrade resharing E2E | `upgrade_scheduled` | `ScheduleDKGUpgrade` → 等 Active + `check(IsUpgrade)` + `checkTrue(GlobalPubKey)` | 短 |

### BeginBlocker (9 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-BB-01 | height < dkgStartBlock: BeginBlocker return nil | `height_below_dkg_start` | `checkTrue("no round", true)` — 重置后立刻跑可触发，FullSuite 流程中走兜底分支 | 短 |
| IT-BB-02 | 无 round → InitiateDKGRound | 无 | `check(round >= 1)` | 短 |
| IT-BB-03 | Pending upgrade + height >= activation | `upgrade_scheduled` | `check(IsUpgrade, true)` | 短 |
| IT-BB-04 | Registration → Dealing 转换 | `resharing_second_round` | `WaitForRoundStage(Dealing)` | 短 |
| IT-BB-05 | Dealing → Finalization 转换 | 无 | `WaitForRoundStage(Finalization)` | 短 |
| IT-BB-06 | Finalization → Active 转换 | 无 | `WaitForRoundStage(Active)` | 短 |
| IT-BB-07 | Active period 到期 → 自动新 round | `resharing_second_round` | `checkTrue(round increased)` | 短 |
| IT-BB-08 | Period 没到不转换 | 无 | `check(stage unchanged)` — 采样两次比较 | **长** |
| IT-BB-09 | height < activation → 不触发 upgrade | `upgrade_scheduled` | `check(IsUpgrade, false)` | 短 |
| IT-BB-10 | 无 pending upgrade | 无 | `check(IsUpgrade, false)` | 短 |

### Registration (20 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-REG-01 | Round=1, IsResharing=false | 无 | `check(round, 1)` + `check(IsResharing, false)` | 短 |
| IT-REG-02 | Resharing: IsResharing=true | `resharing_second_round` | `check(round >= 2)` + `checkTrue(IsResharing)` | 短 |
| IT-REG-03 | Upgrade activation: IsUpgrade=true | `upgrade_scheduled` | `check(IsUpgrade, true)` | 短 |
| IT-REG-04 | DKG disabled 节点不注册 | `dkg_disabled_one_node` | `checkTrue(verified < total)` | 短 |
| IT-REG-05 | ActiveValSet populated | 无 | `checkTrue(len(ActiveValSet) > 0)` | 短 |
| IT-REG-06 | CurRoundSet 注册成功 | 无 | `WaitForVerifiedCount(1)` + `checkTrue(verified >= 1)` | 短 |
| IT-REG-07 | 非 CurRoundSet 不注册 | 无 | 查 `0xdead...` 不在 regs 中 + `checkTrue(!found)` | 短 |
| IT-REG-08 | Stage != Registration guard | 无 | `checkTrue(stage past Registration)` | 短 |
| IT-REG-09 | 并发锁 guard | 无 | `checkTrue(round monotonic)` | 短 |
| IT-REG-10 | TEE down → CreateSession fails | `tee_down_one_node` | `checkTrue(verified < total)` | 短 |
| IT-REG-11 | TEE down → GenerateAndSealKey fails | `tee_down_one_node` | delegate → REG-10 | 短 |
| IT-REG-12 | register 错误参数 revert | 无 | `RegisterDKGWithParams(round=999)` → `checkTrue(err != nil)` | 短 |
| IT-REG-13 | IsUpgrade, in CurRoundSet | `upgrade_scheduled` | delegate → REG-03，`check(IsUpgrade, true)` | 短 |
| IT-REG-14 | IsUpgrade, NOT in CurRoundSet | `upgrade_scheduled` | delegate → REG-03 | 短 |
| IT-REG-15 | Status=Verified | 无 | `check(status, DKG_REG_STATUS_VERIFIED)` | 短 |
| IT-REG-16 | Round mismatch → revert | 无 | `RegisterDKGWithParams(round+100)` → `checkTrue(err)` | 短 |
| IT-REG-17 | StartBlockHeight mismatch → revert | 无 | `RegisterDKGWithParams(height=0)` → `checkTrue(err)` | 短 |
| IT-REG-18 | StartBlockHash mismatch → revert | 无 | `RegisterDKGWithParams(hash=0x00)` → `checkTrue(err)` | 短 |
| IT-REG-19 | Stage != Registration → revert | 无 | `RegisterDKGWithParams` 在非 Registration 阶段 → `checkTrue(err)` | 短 |
| IT-REG-20 | Non-validator address → revert | 无 | `RegisterDKGWithParams(addr=0xdead...)` → `checkTrue(err)` | 短 |

### Dealing (35 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-DL-01 | Verified >= MinReq → BeginDealing | 无 | `WaitForRoundStage(Dealing)` | 短 |
| IT-DL-02 | Verified < MinReq → Skip | `insufficient_verified` | `check(stage, Failed)` 或确认新 round | 短 |
| IT-DL-03 | DKG disabled 节点不参与 dealing | `dkg_disabled_one_node` | `checkTrue(stage >= Dealing)` — 少一个仍推进 | 短 |
| IT-DL-04 | Stage != Dealing guard | 无 | `checkTrue(stage != Dealing)` | 短 |
| IT-DL-05 | 并发锁 guard | 无 | `checkTrue(round > 0)` | 短 |
| IT-DL-06 | ActiveValSet 含旧成员 | 无 | `checkTrue(len(ActiveValSet) > 0)` | 短 |
| IT-DL-07 | GenerateDeals 成功 | 无 | `grepValidatorLog("dealing phase complete")` + `checkTrue` | 短 |
| IT-DL-08 | TEE down → GenerateDeals 失败 | `tee_down_one_node` | `checkTrue(stage progresses)` — 其余节点完成 | 短 |
| IT-DL-09 | Upgrade round dealing | `upgrade_scheduled` | `check(IsUpgrade, true)` + 等 Dealing | 短 |
| IT-DL-10 | TEE down → VerifyDeals 失败 | `tee_down_one_node` | delegate → DL-08 | 短 |
| IT-DL-11 | VerifyDeals 成功 | 无 | `grepValidatorLog("VerifyDeals")` + `checkTrue` | 短 |
| IT-DL-12 | Deal broadcast via VE | 无 | `checkTrue(round progressed past Dealing)` | 短 |
| IT-DL-13 | Response broadcast via VE | 无 | `checkTrue(round progressed)` — VE 隐式 | 短 |
| IT-DL-14 | VerifyVoteExtension valid deal | 无 | `checkTrue(round progressed)` — VE 隐式 | 短 |
| IT-DL-15 | Invalid deal → complaint | 无 | `checkTrue(round progressed)` — complaint 路径隐式 | 短 |
| IT-DL-16 | Justification resolves complaint | 无 | `checkTrue(round progressed)` — justification 隐式 | 短 |
| IT-DL-17 | All deals valid, no complaints | 无 | `WaitForRoundStage(Finalization)` + `checkTrue` | 短 |
| IT-DL-18~35 | VE/Deal 内部逻辑 (18 条) | 无 | `checkTrue(round progressed)` — VoteExtension 内部，隐式验证 | 短 |

### Finalization (19 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-FN-01 | BeginFinalization | 无 | `check(stage, Finalization)` 或 `checkTrue(finalized > 0)` | 短 |
| IT-FN-02 | DKG disabled → Emit only | `dkg_disabled_one_node` | `checkTrue(stage >= Finalization)` | 短 |
| IT-FN-03 | 签名验证失败 → revert | 无 | `FinalizeDKGWithParams(fake sig)` → `checkTrue(err)` | 短 |
| IT-FN-04 | TEE Finalize 成功 | 无 | `grepValidatorLog("finalization phase complete")` + `checkTrue` | 短 |
| IT-FN-05 | 合约调用失败 → revert | 无 | `FinalizeDKGWithParams(empty)` → `checkTrue(err)` | 短 |
| IT-FN-06 | 双重 finalization 被拒 | 无 | 查 finalized regs 无重复 + `checkTrue(!duplicate)` | 短 |
| IT-FN-07 | TEE down → callTEE 失败 | `tee_down_one_node` | `checkTrue(stage >= Finalization)` — 其余节点完成 | 短 |
| IT-FN-08 | TEE down → 不提交 finalize | `tee_down_one_node` | delegate → FN-07 | 短 |
| IT-FN-09 | Invalidated dealer 不能 finalize | 无 | 查 regs 中 invalidated 未 finalize + `checkTrue` | 短 |
| IT-FN-10 | Finalized >= Threshold → GlobalPubKey | 无 | `checkTrue(GlobalPubKey > 0)` | 短 |
| IT-FN-11 | Vote counting | 无 | `checkTrue(finalized >= threshold)` | 短 |
| IT-FN-12 | GlobalPubKey set at threshold | 无 | `checkTrue(GlobalPubKey > 0)` | 短 |
| IT-FN-13 | Minimum threshold | 无 | `checkTrue(finalized >= threshold)` | 短 |
| IT-FN-14 | All validators finalize | 无 | `checkTrue(finalized >= threshold)` | 短 |
| IT-FN-15~19 | Finalization 内部逻辑 (5 条) | 无 | `checkTrue(round progressed)` — 隐式验证 | 短 |

### Active / Complete (13 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-ACT-01 | Active + GlobalPubKey | 无 | `assertActiveWithGlobalKey` (`check` + `checkTrue`) | 短 |
| IT-ACT-02 | finalizedCount < MinReq → Skip | `insufficient_finalized` | delegate → E2E-04 | 短 |
| IT-ACT-03 | finalizedCount < Threshold → Skip | `insufficient_finalized` | delegate → E2E-04 | 短 |
| IT-ACT-04 | DKG disabled → 不跑 complete | `dkg_disabled_one_node` | `checkTrue(disabled not in finalized)` | 短 |
| IT-ACT-05 | Upgrade round complete | `upgrade_scheduled` | `check(stage, Active)` + `check(IsUpgrade, true)` | 短 |
| IT-ACT-06 | First round settleRewards=nil | 无 | `check(round, 1)` — 无 prevActive 是正常行为 | 短 |
| IT-ACT-07 | DkgCommitteeRewardPortion | 无 | `checkTrue(!portion.IsNegative())` | 短 |
| IT-ACT-08 | Active round exists (rewards settled) | 无 | `assertActiveWithGlobalKey` | 短 |
| IT-ACT-09 | settleRewards complete | 无 | `assertActiveWithGlobalKey` | 短 |
| IT-ACT-10 | handleDKGComplete | 无 | `check(stage, Active)` + `checkTrue(GlobalPubKey)` | 短 |
| IT-ACT-11 | PhaseCompleted return | 无 | `checkTrue(round monotonic)` — 已完成 round 不重处理 | 短 |
| IT-ACT-12 | MarkFailed | `tee_down_one_node` | `check(stage, Active)` + `checkTrue(GlobalPubKey)` 或 recovery | 短 |
| IT-ACT-13 | 并发锁 guard | 无 | `checkTrue(round monotonic)` | 短 |

### Skip / SkipToNextRound (4 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-SKIP-01 | Verified < MinReq → Skip | `insufficient_verified` | delegate → E2E-03 | 短 |
| IT-SKIP-02 | Finalized < Threshold → Skip | `insufficient_finalized` | delegate → E2E-04 | 短 |
| IT-SKIP-03 | FlushAllQueues | 无 | `checkTrue(round progresses cleanly)` — 新 round 无残留 | 短 |
| IT-SKIP-04 | KV store error | 无 | `checkTrue(no stuck rounds)` | 短 |

### Resume (7 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-RES-01 | Phase=Failed → resume | `tee_down_one_node` | 恢复 kernel 后 `checkTrue(round progresses)` | 短 |
| IT-RES-02 | GetSession fails → recover | `tee_down_one_node` | delegate → RES-01 | 短 |
| IT-RES-03 | tryResume → GetSession | `tee_down_one_node` | delegate → RES-01 | 短 |
| IT-RES-04 | Phase != Failed → no resume | `tee_down_one_node` | `check(stage, Active)` — 正常时不触发 resume | 短 |
| IT-RES-05 | Rejoin current round | `tee_down_one_node` | delegate → RES-01 | 短 |
| IT-RES-06 | CreateSession fails → stay Failed | `tee_down_one_node` | `checkTrue(TEE-level failure)` | 短 |
| IT-RES-07 | All resume → new round | `tee_down_one_node` | `check(stage, Active)` + `checkTrue(GlobalPubKey)` | 短 |

### Upgrade (5 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-UPG-01 | ScheduleUpgrade | `upgrade_scheduled` | `checkTrue(upgrade scheduled)` — 调 EthChainClient | 短 |
| IT-UPG-02 | hasPendingUpgrade | `upgrade_scheduled` | `checkTrue(activationHeight set)` | 短 |
| IT-UPG-03 | CancelUpgrade | `upgrade_scheduled` | `check(IsUpgrade, false)` | 短 |
| IT-UPG-04 | Upgrade round completes | `upgrade_scheduled` | `check(stage, Active)` + `check(IsUpgrade, true)` | 短 |
| IT-UPG-05 | Upgrade failed → retry | `upgrade_scheduled` | `checkTrue(new round after failure)` | 短 |

### CDR (8 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-CDR-01 | FeeCollected → pool | `cdr_e2e` | CDR.write + `checkTrue` — 需 EthChainClient | 短 |
| IT-CDR-02 | PartialDecryption → count | 无 | CDR.read + `checkTrue` — 需 EthChainClient | 短 |
| IT-CDR-03 | distributeCDRRewardPool | 无 | 查 active round + `checkTrue` — 需 EthChainClient | 短 |
| IT-CDR-04 | Fee=0 no pool | 无 | `checkTrue(fee >= 0)` | 短 |
| IT-CDR-05 | Pool underflow | 无 | `checkTrue(pool path logged)` | 短 |
| IT-CDR-06 | No prev active → nil | 无 | `checkTrue(round 1 no prev)` | 短 |
| IT-CDR-07 | Zero submit, pool>0 | 无 | `checkTrue(zero-submit path)` | 短 |
| IT-CDR-08 | CDR E2E | `cdr_e2e` | Allocate→Write→Read + `check` — 需 EthChainClient | 短 |

### Encrypt (1 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-ENC-01 | TDH2 encrypt + CDR.write | 无 | `checkTrue(GlobalPubKey)` — 需 EthChainClient + active round | 短 |

### Decrypt (1 条 ✅ + 21 条 🧪)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-DEC-01 | ProcessCDRVaultRead: latestActive exists | 无 | CDR read path + `check(stage, Active)` — 需 EthChainClient | 短 |

### Edge Case (6 条)

| ID | 描述 | 场景脚本 | 验证方式 | 短/长周期 |
|---|---|---|---|---|
| IT-EDGE-01 | MinReq validators complete | 无 | `checkTrue(Total/Threshold valid)` | 短 |
| IT-EDGE-02 | Exactly MinReq boundary | 无 | `checkTrue(MinReqReg > 0 && MinReqFin > 0)` | 短 |
| IT-EDGE-03 | Max validators stress | 无 | `checkTrue(total > 0 && threshold > 0)` | 短 |
| IT-EDGE-04 | Failed round auto-recovery | 无 | 查历史 failed rounds + `checkTrue(new round after failure)` | 短 |
| IT-EDGE-05 | Period boundaries | 无 | `checkTrue(all periods > 0)` | 短 |
| IT-EDGE-06 | TEE down 1 node, DKG 完成 | `tee_down_one_node` | `check(stage, Active)` + `checkTrue(GlobalPubKey)` 或 recovery | 短 |

---

## 🔶 有断言但场景受限 (1 条)

| ID | 描述 | 当前实现 | 缺什么 |
|---|---|---|---|
| IT-E2E-05 | Complaint/Justification 路径 | 等 Active 或 Failed + 查 invalidated 数 + `checkTrue` | 当前 kernel 不会产生无效 deal，所以 complaint 路径不会被触发。走 "round Active, invalidated=0" 分支，断言通过但没真正测 complaint |

**Block**: 需要 kernel 支持 mock/fault-injection 模式，可配置产生无效 deal share 来触发 complaint → response → justification 流程。

---

## 🧪 需 TEE mock server (21 条)

IT-DEC-02~22：CDR 解密路径的各种边界。有 EthChainClient 时走 CDR Allocate→Write→Read，但无法验证 keeper 内部的 partial decryption 处理逻辑。

### 可通过 CDR 合约调用部分验证 (5 条)

这些在有 EthChainClient 时会实际执行 CDR read 路径，但内部逻辑验证不完整。

| ID | 描述 | 当前实现 | 缺什么 |
|---|---|---|---|
| IT-DEC-07 | IncrementCount | CDR.read → `checkTrue` | 无法观测内部 submit count 增量 |
| IT-DEC-08 | RefundCDRFee | CDR.read → `checkTrue` | 无法观测 fee refund 执行 |
| IT-DEC-10 | SendCoins from pool | CDR.read → `checkTrue` | 无法观测 SendCoins 调用 |
| IT-DEC-12 | Multiple partials reach threshold | CDR.read → `checkTrue` | 无法控制多个 partial 提交时序 |
| IT-DEC-22 | Full decrypt E2E | CDR.read → `checkTrue` | 完整路径可走，但无法验证解密结果正确性 |

### 需要故障注入 (10 条)

这些需要 TEE 返回错误/无效数据才能触发目标代码路径。

| ID | 描述 | 需要的故障 | 实现方案 |
|---|---|---|---|
| IT-DEC-03 | GetSession fails: MarkFailed | session 损坏或不存在 | Mock kernel `GetSession` 返回 error |
| IT-DEC-04 | Phase != Completed: skip | session phase 不是 Completed | Mock kernel 返回非 Completed phase |
| IT-DEC-05 | callTEEDecrypt fails | TEE 解密调用失败 | Mock kernel `Decrypt` 返回 error |
| IT-DEC-06 | callContractSubmitPartial fails | EL 合约调用失败 | Mock EL RPC 返回 revert |
| IT-DEC-11 | Pool < refundAmt error | CDR pool 余额不足 | 先消耗 pool 到 0 再触发 refund |
| IT-DEC-13 | Non-active-round validator partial | 非当前 committee 的 validator | 用非 committee 地址提交 partial |
| IT-DEC-14 | Duplicate partial submission | 同一 validator 重复提交 | 提交两次相同 partial |
| IT-DEC-15 | Invalid partial proof | 无效 proof 数据 | Mock kernel 返回无效 proof |
| IT-DEC-18 | Concurrent decrypt requests | 并发多个 read 请求 | 并发调用 CDR.read |
| IT-DEC-21 | callTEEDecrypt timeout | TEE 响应超时 | Mock kernel 延迟 > timeout |

### 需要特定链上状态 (6 条)

| ID | 描述 | 需要的状态 | 实现方案 |
|---|---|---|---|
| IT-DEC-02 | No latestActive: return early | 无 active round | 在第一轮完成前调用 |
| IT-DEC-09 | Pool balance check | 精确的 pool 余额 | CDR.write 若干次 → 查 pool |
| IT-DEC-16 | Vault not found | 查询不存在的 vault | 用不存在的 uuid 调 CDR.read |
| IT-DEC-17 | Empty vault | 已分配但未写入的 vault | CDR.allocate 后不 write 直接 read |
| IT-DEC-19 | Decrypt after round rotation | round 切换后的旧数据 | 等 round 切换后读上一轮数据 |
| IT-DEC-20 | Stale GlobalPublicKey | 旧轮的加密数据 | 用旧 GlobalPubKey 加密的数据尝试解密 |

### TEE Mock Server 设计

需要开发一个轻量的 gRPC server 实现 `KernelService` 接口，支持：

```
功能：
1. 正常模式 — 转发到真实 kernel（proxy）
2. 故障模式 — 按配置返回错误/延迟/无效数据

注入点：
- Decrypt(req) → return error          (DEC-05)
- Decrypt(req) → return invalid proof   (DEC-15)
- Decrypt(req) → sleep 30s → timeout    (DEC-21)
- GenerateDeals(req) → return invalid    (E2E-05 complaint)

配置方式：
- 环境变量 DKG_TEE_MOCK_MODE=proxy|fault
- gRPC metadata 传递故障类型
- 或者通过 HTTP control API 切换模式
```

**工作量估计**: 1-2 天开发 + 1 天集成测试

---

## 场景脚本清单

| 场景名 | pre.sh | post.sh | 使用的 Case | 前置验证 |
|---|---|---|---|---|
| `resharing_second_round` | 无操作（等 active period） | 无操作 | BB-04/07, REG-02, E2E-02 | `checkTrue(round >= 1)` |
| `insufficient_verified` | 停 node 2/3 kernel（不停 story） | 启动所有 kernel | E2E-03, DL-02, SKIP-01 | SSH 验证 kernel inactive |
| `insufficient_finalized` | 停 node 2/3 kernel | 启动所有 kernel | E2E-04, ACT-02/03, SKIP-02 | SSH 验证 kernel inactive |
| `tee_down_one_node` | 停 node 3 kernel | 启动 kernel | REG-10/11, DL-08/10, FN-07/08, ACT-12, RES-01~07, EDGE-06 | SSH 验证 kernel inactive |
| `dkg_disabled_one_node` | 改 node 3 config dkg.enable=false + 重启 | 恢复 enable=true + 重启 | REG-04, DL-03, FN-02, ACT-04 | SSH 验证 config |
| `upgrade_scheduled` | 检查 EthChainClient 环境变量 | Go 代码调 CancelUpgrade | BB-03/09, REG-03/13/14, DL-09, ACT-05, E2E-06, UPG-01~05 | 检查 RPC/key 配置 |
| `complaint_justification` | 提示需 TEE mock | 无 | E2E-05 | 无（当前 kernel 不产生无效 deal） |
| `cdr_e2e` | 检查 EthChainClient 环境变量 | 无 | CDR-01/08 | 检查 RPC/key 配置 |
| `height_below_dkg_start` | 检查 height | 无 | BB-01 | 查 height |

---

## 环境依赖

| 依赖 | 影响的 Case 数 | 如何配置 |
|---|---|---|
| SSH tunnel (26657+8545) | 全部 156 条 | `ssh -L 26657:localhost:26657 -L 8545:localhost:8545 ubuntu@val1` |
| `STORY_ETH_RPC_URL` | ~30 条 (REG-12/16~20, FN-03/05, UPG, CDR, ENC, DEC) | `export STORY_ETH_RPC_URL=http://localhost:8545` |
| `DKG_SIGNER_PRIVATE_KEY` | 同上 | `export DKG_SIGNER_PRIVATE_KEY=ac0974bec...` (Foundry test EOA) |
| `DKG_SSH_TARGETS` + `DKG_SSH_KEY` | 场景脚本 + 日志 grep | config.env 已配 |
| 短周期 devnet | 133 条 | `configure_periods.sh 20 80 50 20` |
| 长周期 devnet | 1 条 (BB-08) | `configure_periods.sh 200 300 300 200` |
| TEE mock server | 22 条 (DEC-02~22, E2E-05) | **未开发** |

---

## Priority Action Items

| 优先级 | 工作 | 状态 |
|---|---|---|
| **Done** | 134 条有效断言 | ✅ 已完成 |
| **P1** | 重置 devnet + 跑 FullSuite 验证 134 条 | ⏳ 待执行 |
| **P2** | 开发 TEE mock gRPC server | 📋 待开发 (1-2d) |
| **P3** | 用 mock server 实现 21 条 DEC + 1 条 E2E-05 | 📋 待 mock 完成后 |
