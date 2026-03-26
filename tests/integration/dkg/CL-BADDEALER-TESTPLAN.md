# CL-BADDEALER: Bad Dealer Finalization & Committee Participation Test Plan

## Bug Summary

A validator who sends cryptographically invalid deals (VSS verification fails) can still:
1. Successfully call `finalize()` on DKG.sol
2. Be marked as `DKGRegStatusFinalized` on-chain
3. Be counted in `finalizedCount` toward the Active threshold
4. Receive UBI committee rewards
5. Submit partial decryptions that pass on-chain validation

**Root cause**: `invalidateDealerRegistration()` exists in `dkg_justification.go:206` but is **never called** from `handleDKGProcessJustifications`. The comment states "justification processing must not affect on-chain state". As a result, a bad dealer's on-chain status permanently stays `DKGRegStatusVerified` even after their deal is proven invalid through the justification process.

## Code Evidence (Root Cause Chain)

### Evidence 1: invalidateDealerRegistration 从未被调用

**File**: `client/x/dkg/keeper/dkg_justification.go:198-205`
```go
// invalidateDealerRegistration marks a dealer's DKG registration as Invalidated.
// ...
// NOTE: This function is NOT called from ProcessJustifications because justification
// processing must not affect on-chain state. It is retained as a utility for
// potential future use (e.g., explicit slashing proposals).
func (k *Keeper) invalidateDealerRegistration(ctx context.Context, latestRound *types.DKGNetwork, dealerIndex uint32) error {
```

验证方法：全局搜索 `invalidateDealerRegistration` 调用点 — 仅在函数定义和单元测试中出现，生产代码中**零调用**。

### Evidence 2: VSS 验证失败的 justification 被静默丢弃

**File**: `client/x/dkg/keeper/dkg_svc_dealing.go:453-468`
```go
for _, j := range deduped {
    valid, err := verifyJustification(dkgNetwork, j)
    if err != nil {
        log.Warn(ctx, "Justification VSS verification error, dropping", err,
            "dealer_index", j.Index,
        )
        continue  // ← 丢弃，无后续动作
    }
    if !valid {
        log.Info(ctx, "Justification VSS verification failed (deal was invalid), dropping",
            "dealer_index", j.Index,
        )
        continue  // ← 丢弃，无后续动作。不调用 invalidateDealerRegistration()
    }
    validJustifications = append(validJustifications, j)
}
```

**关键**: `!valid` 分支仅 log + continue，不触发任何链上状态变更。坏 dealer 的 `DKGRegStatusVerified` 保持不变。

### Evidence 3: Finalized() handler 的 Invalidated 检查形同虚设

**File**: `client/x/dkg/keeper/dkg_handler.go:121-124`
```go
// Reject finalization by invalidated dealers (deal complaint found invalid via VSS verification)
if reg.Status == types.DKGRegStatusInvalidated {
    return errors.New("dealer has been invalidated and cannot finalize")
}
```

**关键**: 这个检查是正确的防线，但因为 Evidence 1 & 2，`DKGRegStatusInvalidated` 永远不会被设置，所以这个 `if` 永远不会触发。

### Evidence 4: finalizeDKGRegistration 无条件将 Verified 标记为 Finalized

**File**: `client/x/dkg/keeper/dkg_handler.go:96-161`（`Finalized` 函数完整流程）
```go
func (k *Keeper) Finalized(ctx context.Context, round uint32, msgSender common.Address, ...) error {
    reg, err := k.getDKGRegistration(ctx, round, msgSender)
    // ...
    if reg.Status == types.DKGRegStatusFinalized {
        return errors.New("validator has already finalized for this round")  // 防重复
    }
    if reg.Status == types.DKGRegStatusInvalidated {
        return errors.New("dealer has been invalidated and cannot finalize")  // 形同虚设
    }
    // ↓ 坏 dealer 状态是 Verified，所以通过上面的检查
    if err := k.validateParticipantsRoot(ctx, round, participantsRoot); err != nil { ... }
    if err := verifyFinalizationSignature(reg.CommPubKey, ...); err != nil { ... }
    voteCount, err := k.AddGlobalPubKeyVote(ctx, round, globalPubKey, publicCoeffs)
    // ...
    if err := k.finalizeDKGRegistration(ctx, round, msgSender, pubKeyShare); err != nil { ... }
    // ↑ 坏 dealer 成功被标记为 Finalized
```

**关键**:
- `validateParticipantsRoot` (line 126): 计算 Keccak256(所有 Verified + Finalized 地址)，坏 dealer 仍是 Verified 所以被包含
- `verifyFinalizationSignature` (line 130): 验证 ECDSA 签名，坏 dealer 用自己的 commPubKey 签名，完全合法
- `AddGlobalPubKeyVote` (line 134): 坏 dealer 投票一个不同的 globalPubKey，投票被计数但不会赢（少数派）
- `finalizeDKGRegistration` (line 148): 无条件标记为 Finalized

### Evidence 5: participantsRoot 包含坏 dealer

**File**: `client/x/dkg/keeper/dkg_handler.go:167-213`
```go
func (k *Keeper) validateParticipantsRoot(ctx context.Context, round uint32, participantsRoot [32]byte) error {
    verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, round, types.DKGRegStatusVerified)
    finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round, types.DKGRegStatusFinalized)
    allRegs := append(verifiedRegs, finalizedRegs...)
    // ↑ 坏 dealer 是 Verified，被包含在 allRegs 中
    // ↑ 没有排除 Invalidated，因为根本就没人会变成 Invalidated
```

### Evidence 6: finalizedCount 包含坏 dealer，用于 Active 判定

**File**: `client/x/dkg/keeper/dkg_finalization.go:29-62`
```go
func (k *Keeper) FinalizeDKGRound(ctx context.Context, latestRound *types.DKGNetwork) error {
    finalizedCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusFinalized)
    // ↑ 坏 dealer 已被标记为 Finalized（Evidence 4），被计入
    if finalizedCount < params.MinReqFinalizedParticipants { ... }  // 3 >= 3 通过
    if finalizedCount < latestRound.Threshold { ... }               // 3 >= 2 通过
    // → 轮次进入 Active
```

### Evidence 7: 坏 dealer 领取委员会奖励

**File**: `client/x/dkg/keeper/dkg_rewards.go:58-66`
```go
func (k *Keeper) distributeRewardsFromModule(ctx context.Context, round *types.DKGNetwork, ...) (math.Int, error) {
    finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round.Round, types.DKGRegStatusFinalized)
    // ↑ 坏 dealer 在此列表中
    memberAddrs := make([]string, 0, len(finalizedRegs))
    for _, reg := range finalizedRegs {
        memberAddrs = append(memberAddrs, reg.ValidatorAddr)  // ← 坏 dealer 地址被添加
    }
    // ...
    perMemberReward := dkgReward.Quo(memberCount)  // memberCount=3 而非正确的 2
    for _, evmAddr := range memberAddrs {
        k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, perMemberCoins)
        // ↑ 坏 dealer 收到与诚实 validator 等额的奖励
    }
```

### Evidence 8: 坏 dealer 的 partial decryption 通过链上验证

**File**: `client/x/dkg/keeper/dkg_handler.go:483-578`（`PartialDecryptionSubmitted`）
```go
func (k *Keeper) PartialDecryptionSubmitted(ctx context.Context, validator common.Address, round uint32, pid uint32, ...) error {
    reg, err := k.getDKGRegistration(ctx, req.Round, validator)
    // ↑ 获取坏 dealer 的注册信息

    if !bytes.Equal(pubShare, reg.PubKeyShare) {
        return errors.New("pubShare mismatch")
    }
    // ↑ reg.PubKeyShare 是坏 dealer 在 finalize 时自己提交的值
    // ↑ pubShare 也是坏 dealer 当前提交的值 → 两者一致 → 通过！

    if err := verifyPartialDecryptionSignature(reg.CommPubKey, round, ciphertext, ...); err != nil { ... }
    // ↑ 坏 dealer 用自己的 commKey 签名 → 签名有效 → 通过！

    if err := k.setPartialDecryptionSubmission(ctx, validator, round, pid, ...); err != nil { ... }
    // ↑ 坏 dealer 的无效 partial decryption 成功存储到链上
```

**关键**: `pubShare` 验证比对的是 dealer 自己在 finalize 时提交的值（`reg.PubKeyShare`），而不是从共识 `globalPubKey` 推导的期望值。这意味着坏 dealer 可以提交与自己错误 key share 一致的 pubShare 和 partial decryption，验证全部通过。

### Evidence 9: 全局公钥投票无法阻止坏 dealer

**File**: `client/x/dkg/keeper/dkg_votes.go:15-46`
```go
func (k *Keeper) AddGlobalPubKeyVote(ctx context.Context, round uint32, globalPubKey []byte, publicCoeffs [][]byte) (uint32, error) {
    key := fmt.Sprintf("%d_%s_%s", round, hex.EncodeToString(globalPubKey), coeffHash)
    current, err := k.GlobalPubKeyVotes.Get(ctx, key)
    newCount := current + 1
    k.GlobalPubKeyVotes.Set(ctx, key, newCount)
    return newCount, nil
}
```

**File**: `client/x/dkg/keeper/dkg_handler.go:139-146`
```go
if voteCount >= latest.Threshold && len(latest.GlobalPublicKey) == 0 {
    latest.GlobalPublicKey = globalPubKey
    latest.PublicCoeffs = publicCoeffs
    // ...
}
```

**关键**: 投票机制确实保证了链上 `GlobalPublicKey` 是诚实多数的值（B+C），但它**不阻止坏 dealer finalize** — finalize 和投票是耦合在同一个 `Finalized()` handler 中的。投票失败（少数派）不影响 finalize 成功。

## Affected code paths (summary)

| # | File:Line | 代码行为 | 问题 |
|---|-----------|---------|------|
| E1 | `dkg_justification.go:203-205` | `invalidateDealerRegistration` 注释说明从不调用 | 根因 |
| E2 | `dkg_svc_dealing.go:463-468` | VSS 失败的 justification 静默丢弃 | 无链上后果 |
| E3 | `dkg_handler.go:122-124` | Invalidated 检查存在 | 永远不触发 |
| E4 | `dkg_handler.go:96-161` | `Finalized()` 完整流程 | 坏 dealer 通过所有检查 |
| E5 | `dkg_handler.go:167-213` | `validateParticipantsRoot` 包含 Verified | 坏 dealer 被包含 |
| E6 | `dkg_finalization.go:29-62` | `FinalizeDKGRound` 按 Finalized 计数 | 坏 dealer 被计入 |
| E7 | `dkg_rewards.go:58-66` | 奖励分配给所有 Finalized | 坏 dealer 领奖 |
| E8 | `dkg_handler.go:483-578` | `PartialDecryptionSubmitted` pubShare 自比对 | 坏 dealer 通过验证 |
| E9 | `dkg_votes.go:15-46` + `dkg_handler.go:139-146` | 投票不阻止 finalize | 少数派投票 ≠ 拒绝 finalize |

## Prerequisites

- 3-validator devnet (Validator A, B, C)
- Mock kernel server capability (to inject `WithInvalidVSSDeal` on Validator A)
- On-chain query access (DKG registrations, CDR partials)
- DKG_DRIVER=script enabled
- Long period configuration (deals/finalization need enough time for full cycle)

## Scenario Setup: `mock_kernel_bad_dealer_finalize`

1. Deploy mock kernel on **Validator A** with `WithInvalidVSSDeal()` behavior
   - Mock kernel's `GenerateDeals()` returns deals with valid encryption but invalid VSS shares
   - All other RPCs (`ProcessDeals`, `ProcessResponses`, `FinalizeDKG`, etc.) behave normally
2. Validators B and C run normal story-kernel instances
3. Wait for DKG round to reach Registration stage, all 3 validators register successfully

---

## Test Cases

### CL-BADDEALER-01: Bad dealer's registration not invalidated after failed justification

**Priority**: P1
**Category**: Security / State integrity
**Depends on**: mock_kernel_bad_dealer_finalize scenario

**Steps**:
1. Wait for DKG stage = Dealing
2. Validator A (mock kernel) generates deals with invalid VSS shares
3. Validators B and C process A's deals → `ProcessDeal()` returns `Response.Status=false` (complaint)
4. Complaints broadcast via Vote Extension
5. Validator A receives complaints → generates justification (reveals plaintext deal)
6. Consensus layer verifies justification:
   - Schnorr signature → PASS (A's real key signed it)
   - Pedersen VSS → FAIL (share * G ≠ C_0 + i*C_1 + ...)
7. Justification dropped (logged as "deal was invalid")
8. Wait for DKG stage = Finalization

**Verification**:
```
Query: GetDKGRegistration(round, ValidatorA)
Assert: reg.Status == DKGRegStatusVerified   ← BUG: should be Invalidated
Assert: reg.Status != DKGRegStatusInvalidated ← confirms invalidation never happened
```

**Code Evidence**:
- Bad deal 产生 complaint: `story-kernel/service/dkg_process_deals.go:75` — `distKeyGen.ProcessDeal(deal)` 返回 `Response.Status=false`
- Complaint 产生 justification: `story-kernel/service/dkg_process_responses.go:90-114` — `distKeyGen.ProcessResponse(resp)` 返回 justification
- Justification VSS 验证失败被丢弃: `client/x/dkg/keeper/dkg_svc_dealing.go:463-468` — `!valid → continue`
- **缺失的调用**: `client/x/dkg/keeper/dkg_justification.go:206` — `invalidateDealerRegistration` 存在但未被调用 (**E1, E2**)

**Expected (current behavior, demonstrating the bug)**:
- Validator A's status remains `Verified` despite proven bad deal

**Expected (correct behavior after fix)**:
- Validator A's status should be `Invalidated`

---

### CL-BADDEALER-02: Bad dealer successfully finalizes

**Priority**: P1
**Category**: Security / Finalization integrity
**Depends on**: CL-BADDEALER-01

**Steps**:
1. (Continues from CL-BADDEALER-01) DKG stage = Finalization
2. Validator A's kernel calls `FinalizeDKG()` → computes `DistKeyShare` from its own DKG state
3. Validator A submits `finalize()` tx to DKG.sol
4. Keeper's `Finalized()` handler processes A's finalization:
   - `reg.Status == DKGRegStatusInvalidated` check → **PASS** (status is Verified, not Invalidated)
   - `validateParticipantsRoot()` → **PASS** (A is in Verified set, included in root)
   - `verifyFinalizationSignature()` → **PASS** (A's commPubKey signed correctly)
   - `AddGlobalPubKeyVote()` → A votes with `globalPubKey_A` (different from B,C's `globalPubKey_BC`)
   - `finalizeDKGRegistration()` → A's status updated to **Finalized**
5. Validators B and C also finalize with `globalPubKey_BC`

**Verification**:
```
Query: GetDKGRegistration(round, ValidatorA)
Assert: reg.Status == DKGRegStatusFinalized   ← BUG: bad dealer should not be Finalized

Query: GetFinalizedRegistrations(round)
Assert: len(finalizedRegs) == 3   ← BUG: should be 2 (only B and C)

Query: GetLatestDKGNetwork(round)
Assert: network.GlobalPublicKey == globalPubKey_BC  ← B+C reach threshold without A's vote
```

**Code Evidence**:
- Invalidated 检查形同虚设: `client/x/dkg/keeper/dkg_handler.go:122-124` — `reg.Status` 永远不会是 `Invalidated` (**E3**)
- participantsRoot 包含坏 dealer: `client/x/dkg/keeper/dkg_handler.go:168-178` — 查询 Verified + Finalized，坏 dealer 在 Verified 中 (**E5**)
- 签名验证通过: `client/x/dkg/keeper/dkg_handler.go:130` — `verifyFinalizationSignature` 只验证 ECDSA，不验证 globalPubKey 正确性
- 无条件标记 Finalized: `client/x/dkg/keeper/dkg_handler.go:148` — `finalizeDKGRegistration` 不检查 globalPubKey 是否匹配多数派 (**E4**)
- 投票不阻止 finalize: `client/x/dkg/keeper/dkg_handler.go:134-146` — `AddGlobalPubKeyVote` 返回 count，但即使 count=1（少数派），finalize 仍然继续 (**E9**)

**Expected (current behavior, demonstrating the bug)**:
- Validator A is Finalized, 3 validators in finalized set

**Expected (correct behavior after fix)**:
- Validator A's finalize() tx should be rejected

---

### CL-BADDEALER-03: Bad dealer counted in finalizedCount, inflates committee size

**Priority**: P1
**Category**: Security / Threshold integrity
**Depends on**: CL-BADDEALER-02

**Steps**:
1. (Continues from CL-BADDEALER-02) DKG transitions from Finalization → Active
2. `FinalizeDKGRound()` checks:
   - `finalizedCount >= MinReqFinalizedParticipants` (3 >= 3)
   - `finalizedCount >= Threshold` (3 >= 2)
3. Round becomes Active

**Verification**:
```
Query: GetLatestActiveDKGNetwork()
Assert: network.Stage == DKGStageActive
Assert: network.Round == current_round

# Count finalized registrations
finalizedRegs = GetFinalizedRegistrations(round)
Assert: len(finalizedRegs) == 3   ← BUG: includes bad dealer
```

**Code Evidence**:
- finalizedCount 查询: `client/x/dkg/keeper/dkg_finalization.go:30` — `countDKGRegistrationsByStatus(ctx, round, DKGRegStatusFinalized)` 包含坏 dealer (**E6**)
- Active 判定: `client/x/dkg/keeper/dkg_finalization.go:41-61` — 两个 threshold 检查都用 `finalizedCount`，坏 dealer 使其虚增

**Why this matters**:
- In a 3-validator network with threshold=2, if 2 honest + 1 bad finalize, the system thinks it has 3 functioning members
- If one of the 2 honest validators goes offline, the system expects 2 remaining members to handle decryption, but only 1 can actually produce valid partial decryptions
- Effective committee size is 2, not 3 — the system's fault tolerance is silently degraded

---

### CL-BADDEALER-04: Bad dealer receives committee rewards

**Priority**: P2
**Category**: Economic / Reward integrity
**Depends on**: CL-BADDEALER-03

**Steps**:
1. (Continues from CL-BADDEALER-03) Round is Active
2. UBI withdrawal cycle triggers `DistributeRewardsToActiveCommittee()`
3. Rewards distributed to ALL Finalized registrations

**Verification**:
```
# Before reward distribution, record balances of A, B, C
balanceBefore_A = getBalance(ValidatorA)
balanceBefore_B = getBalance(ValidatorB)

# Wait for UBI distribution event
# DKG committee rewards = portion * totalUBI / memberCount

balanceAfter_A = getBalance(ValidatorA)
Assert: balanceAfter_A > balanceBefore_A   ← BUG: bad dealer receives reward

# Verify equal distribution
reward_A = balanceAfter_A - balanceBefore_A
reward_B = balanceAfter_B - balanceBefore_B
Assert: reward_A == reward_B   ← bad dealer gets same reward as honest validators
```

**Code Evidence**:
- 奖励分配查询 Finalized: `client/x/dkg/keeper/dkg_rewards.go:59` — `getDKGRegistrationsByStatus(ctx, round.Round, DKGRegStatusFinalized)` (**E7**)
- 坏 dealer 地址被加入奖励列表: `client/x/dkg/keeper/dkg_rewards.go:69-72` — 遍历 `finalizedRegs` 无过滤
- 等额分配: `client/x/dkg/keeper/dkg_rewards.go:78-79` — `perMemberReward = dkgReward / memberCount`，memberCount 包含坏 dealer
- 活跃委员会奖励同理: `client/x/dkg/keeper/dkg_rewards.go:24-41` — `DistributeRewardsToActiveCommittee` → `distributeRewardsFromModule` 走相同路径

**Expected (correct behavior after fix)**:
- Validator A should NOT receive committee rewards

---

### CL-BADDEALER-05: Bad dealer's partial decryption accepted on-chain

**Priority**: P1
**Category**: Security / Decryption integrity
**Depends on**: CL-BADDEALER-03

**Steps**:
1. (Continues from CL-BADDEALER-03) Round is Active
2. A CDR `read()` request triggers `ThresholdDecryptRequested` event
3. All 3 validators (A, B, C) receive the decrypt request
4. Validator A's kernel calls `PartialDecryptTDH2()`:
   - Uses its `DistKeyShare` (computed with wrong globalPubKey_A)
   - Produces partial decryption with its own key share
   - Signs with Secp256k1 commKey
5. Validator A submits `submitEncryptedPartialDecryption()` to CDR.sol
6. Keeper's `PartialDecryptionSubmitted()` handler:
   - `pubShare` comparison: `bytes.Equal(pubShare, reg.PubKeyShare)` → **PASS** (A stored its own pubKeyShare during finalize)
   - `verifyPartialDecryptionSignature()` → **PASS** (A's commKey signature is valid)
   - Partial decryption stored on-chain

**Verification**:
```
# After all 3 submit partials
partials = GetCDRPartials(uuid, requesterPubKey)
Assert: len(partials) == 3   ← includes A's invalid partial

# Attempt to combine partials for threshold decryption
# Need threshold=2 partials to decrypt
# If A's partial is selected (any 2 of 3), decryption fails
# Only B+C combination works
```

**Code Evidence**:
- pubShare 自比对漏洞: `client/x/dkg/keeper/dkg_handler.go:546` — `bytes.Equal(pubShare, reg.PubKeyShare)` 比对的是 dealer 自己在 finalize 时提交的值 (**E8**)
- pubKeyShare 存储时机: `client/x/dkg/keeper/dkg_handler.go:148` — `finalizeDKGRegistration(ctx, round, msgSender, pubKeyShare)` 将坏 dealer 自己算的 pubKeyShare 存入链上
- 签名验证通过: `client/x/dkg/keeper/dkg_handler.go:553` — `verifyPartialDecryptionSignature(reg.CommPubKey, ...)` 只验证 ECDSA 签名来源，不验证 partial 的密码学正确性
- 无 committee membership 检查: `PartialDecryptionSubmitted` 中**不检查** validator 的 globalPubKey 投票是否与共识一致
- 存储成功: `client/x/dkg/keeper/dkg_handler.go:557` — `setPartialDecryptionSubmission` 存入链上

**Expected (current behavior, demonstrating the bug)**:
- A's partial decryption is accepted and stored on-chain
- If client picks A's partial + one honest partial → decryption fails
- Only B+C partial combination produces correct decryption

**Expected (correct behavior after fix)**:
- A should not be able to submit partial decryption (not in committee), OR
- A's partial should be rejected during submission validation

---

### CL-BADDEALER-06: Decryption fails when bad dealer's partial is selected

**Priority**: P1
**Category**: Security / CDR availability
**Depends on**: CL-BADDEALER-05

**Steps**:
1. (Continues from CL-BADDEALER-05) 3 partial decryptions on-chain
2. Client fetches partials and attempts to combine any threshold=2 subset
3. Test all 3 combinations:
   - A + B → decrypt → FAIL (A's share is wrong)
   - A + C → decrypt → FAIL (A's share is wrong)
   - B + C → decrypt → SUCCESS

**Verification**:
```
# Fetch all 3 partials
partials = GetCDRPartials(uuid, requesterPubKey)

# Combination 1: A + B
result_AB = TDH2Combine(partials[A], partials[B], globalPubKey_BC)
Assert: result_AB == ERROR or result_AB != original_plaintext

# Combination 2: A + C
result_AC = TDH2Combine(partials[A], partials[C], globalPubKey_BC)
Assert: result_AC == ERROR or result_AC != original_plaintext

# Combination 3: B + C
result_BC = TDH2Combine(partials[B], partials[C], globalPubKey_BC)
Assert: result_BC == original_plaintext   ← only valid combination
```

**Code Evidence**:
- 坏 dealer 的 DistKeyShare 基于错误 globalPubKey: `story-kernel/service/dkg_finalize.go:72` — `distKeyGen.DistKeyShare()` 从各自独立的 DKG 状态计算，A 的包含自己坏 deal 的贡献
- globalPubKey 差异源头: `story-kernel/service/dkg_finalize.go:110` — `distKeyShare.Public()` 返回的公钥取决于哪些 deal 被接受。A 接受了自己的 deal（self-deal），B/C 排除了 A 的 deal → 不同的 public key
- TDH2 部分解密使用 DistKeyShare: `story-kernel/service/dkg_partial_decrypt.go` — `mpc.TDH2PartialDecrypt(ownPID, privShare, pubKey, ct, label)`，A 的 `privShare` 基于错误的 `DistKeyShare` → 产生无效的 partial

**Impact assessment**:
- With 3 validators and threshold=2, theoretical availability = C(3,2)/C(3,2) = 100%
- Actual availability with 1 bad dealer = 1/3 = 33% (only B+C works)
- If 2 out of 3 validators are bad dealers, availability = 0%

---

## Edge Case Tests

### CL-BADDEALER-07: Multiple bad dealers simultaneously

**Priority**: P2
**Category**: Security / Threshold boundary

**Setup**: Validators A and B both have mock kernels with `WithInvalidVSSDeal()`

**Verification**:
```
# After Finalization
finalizedRegs = GetFinalizedRegistrations(round)
Assert: len(finalizedRegs) == 3   ← BUG: both bad dealers finalized

# Threshold = ceil(3 * 667 / 1000) = 2
# Only C has valid key share
# GlobalPubKey vote: A votes pubkey_A, B votes pubkey_B, C votes pubkey_C
# No one reaches threshold=2 → GlobalPublicKey is empty → round should fail

network = GetLatestDKGNetwork()
# If GlobalPublicKey is empty but finalizedCount >= threshold, what happens?
```

**This tests whether the round correctly fails when no globalPubKey reaches threshold despite enough finalizations.**

### CL-BADDEALER-08: Bad dealer with globalPubKey matching honest majority

**Priority**: P2
**Category**: Security / Theoretical attack

**Description**: If Validator A runs a modified kernel that:
1. Sends invalid deals to recipients (so its share contribution is excluded)
2. But computes its `DistKeyShare` as if its deals were rejected (matching B+C's computation)
3. Submits `globalPubKey_BC` (same as honest validators) during finalize

**Question**: Can A obtain a valid key share despite not contributing valid deals?

**This requires deeper kyber analysis**: In Pedersen DKG, the final key share includes contributions from all accepted dealers. If A's deals are excluded from B and C's computation, A's own computation must also exclude its own contribution to produce the same globalPubKey. This is theoretically possible if A manually reconstructs its DKG state without self-dealing.

---

## Summary Table

| Case ID | Description | Priority | Type |
|---------|-------------|----------|------|
| CL-BADDEALER-01 | Bad dealer not invalidated after failed justification | P1 | Bug verification |
| CL-BADDEALER-02 | Bad dealer successfully finalizes | P1 | Bug verification |
| CL-BADDEALER-03 | Bad dealer inflates finalizedCount | P1 | Threshold integrity |
| CL-BADDEALER-04 | Bad dealer receives committee rewards | P2 | Economic impact |
| CL-BADDEALER-05 | Bad dealer's partial decryption accepted | P1 | Decryption integrity |
| CL-BADDEALER-06 | Decryption fails with bad dealer's partial | P1 | CDR availability |
| CL-BADDEALER-07 | Multiple bad dealers vs threshold | P2 | Edge case |
| CL-BADDEALER-08 | Adaptive attack: match honest globalPubKey | P2 | Theoretical attack |

## Recommended Fix Direction

1. **Option A (Minimal)**: In `Finalized()` handler, verify that the submitted `globalPubKey` matches the majority vote before allowing finalization
2. **Option B (Correct)**: Implement on-chain invalidation via a consensus message (`MsgInvalidateDealer`) that can be proposed during Dealing/Finalization stage when justification VSS verification fails
3. **Option C (Defense in depth)**: Both A and B, plus exclude Invalidated registrations from reward distribution and partial decryption acceptance

## Related Existing Cases

- **CL-JUST-01**: Claims "dealer removed from finalized set" but this may not actually happen (needs re-verification)
- **IT-FN-09**: Tests "invalidated dealer cannot finalize" — correct logic exists in code but `Invalidated` status is never set
- **IT-E2E-05**: Complaint/justification path — does not verify finalization outcome
