# Apex Report - Story network / Scan #1

## Table of contents

- [High](about:blank#high)
    - [STOR-15 — CDR fee bridging mints wei-denominated fees as gwei-denominated stake, creating a 1e9 reward over-credit](about:blank#finding-stor-15)
    - [STOR-6 — Validators can submit arbitrary or replayed partial decryptions and still collect CDR refunds and rewards](about:blank#finding-stor-6)
    - [STOR-5 — Unauthenticated PartialDecryptTDH2 turns default kernels into a remote threshold-decryption oracle](about:blank#finding-stor-5)
    - [STOR-4 — Unknown or expired partial submissions are rewarded as valid work](about:blank#finding-stor-4)
    - [STOR-3 — Duplicate partial replays accrue fresh refunds and reward weight every time](about:blank#finding-stor-3)
- [Medium](about:blank#medium)
    - [STOR-29 — Rebooting during an active round permanently disables threshold decryption for that round](about:blank#finding-stor-29)
    - [STOR-28 — Decrypt worker dies after one minute while CDR keeps charging 21-day read fees](about:blank#finding-stor-28)
    - [STOR-27 — Final EL-block partial submissions are settled in the wrong committee epoch](about:blank#finding-stor-27)
    - [STOR-26 — Permissionless pre-v2 CDR fees mutate the future DKG store before activation](about:blank#finding-stor-26)
    - [STOR-25 — Single proposer can permanently censor committed DKG complaints by substituting MsgAddDkgVote](about:blank#finding-stor-25)
    - [STOR-24 — Global first-10 truncation permanently discards valid dealer-invalidating justifications](about:blank#finding-stor-24)
    - [STOR-23 — Dealing-to-finalization transitions deterministically drop the last block’s DKG complaints and preserve bad committee members](about:blank#finding-stor-23)
    - [STOR-22 — The first active block rejects the last finalization block’s DKG finalize events, excluding honest validators from the paid committee](about:blank#finding-stor-22)
    - [STOR-21 — Resharing mutates the serving-round pointer into a non-active round and opens a multi-day read-fee capture window](about:blank#finding-stor-21)
    - [STOR-20 — PartialDecryptTDH2 protobuf mismatch makes every CL-to-kernel decrypt RPC invalid on the wire](about:blank#finding-stor-20)
    - [STOR-19 — The first dealing block rejects the last registration block’s DKG register events, shrinking the committee before rewards are split](about:blank#finding-stor-19)
    - [STOR-18 — VerifiedQueryClient accepts proofs for arbitrary keys and stores, letting hostile RPCs feed fake DKG state](about:blank#finding-stor-18)
    - [STOR-16 — Expired DKG rounds keep collecting UBI and blocking validator exits because `LatestActiveRound` is never cleared at rollover](about:blank#finding-stor-16)
    - [STOR-13 — Finalization stores raw pubKeyShare while partial decryptions submit prefixed pubShare, so honest reads cannot be completed](about:blank#finding-stor-13)
    - [STOR-11 — Reusing one TEE keypair across validator addresses can finalize an undecryptable DKG committee](about:blank#finding-stor-11)
    - [STOR-10 — Early-sorted validator can suppress later DKG complaints by spoofing response dedup keys](about:blank#finding-stor-10)
    - [STOR-9 — The CDR read path is dead in every DKG phase because active rounds kill the decrypt worker and resharing rounds reject request queueing](about:blank#finding-stor-9)
    - [STOR-8 — Resharing rounds never populate PIDCache, so every post-resharing active committee permanently fails CDR decryptions](about:blank#finding-stor-8)
    - [STOR-2 — Upgrade resharing rejects every validator that was not in the previous committee](about:blank#finding-stor-2)
    - [STOR-1 — Partial-decrypt signatures omit requester and UUID binding, so a hostile keeper can redirect valid shares while CL accepts the forged request context](about:blank#finding-stor-1)
- [Low](about:blank#low)
    - [STOR-17 — Sealed light-client state is rollbackable, so a hostile host can pin the enclave to stale trust roots](about:blank#finding-stor-17)
    - [STOR-14 — Honest signed DKG deals can be turned into undeliverable shares because nonce and recipient routing are outside the dealer signature](about:blank#finding-stor-14)
    - [STOR-12 — Whitelisting a new enclave type lets validators bypass scheduleUpgrade and rotate binaries early](about:blank#finding-stor-12)
    - [STOR-7 — First-seen vote dedup lets an earlier validator erase later signed DKG deals and responses](about:blank#finding-stor-7)

## High

### STOR-15 — CDR fee bridging mints wei-denominated fees as gwei-denominated stake, creating a 1e9 reward over-credit

### CDR fee bridging mints wei-denominated fees as gwei-denominated stake, creating a 1e9 reward over-credit

### Executive Summary

The CDR fee bridge breaks Story’s core value-conservation invariant by minting raw EL `wei` fee amounts directly into CL `stake`, even though the CL↔︎EL withdrawal path is denominated in `gwei`. The staking bridge explicitly normalizes execution-layer stake amounts by dividing event amounts by `1 gwei` before minting Cosmos-side stake. The CDR fee path does not: `FeeCollected(amount)` is emitted in raw wei, `ProcessCDRFeeCollected()` forwards that raw `uint256` into `AddCDRFeeToPool()`, and the DKG module mints the same numeric amount of CL `stake` into `cdr-fee-pool`.

That over-credit is immediately exploitable. On every valid partial submission, `RefundCDRFee()` pays the validator the raw `ev.Fee` as ordinary CL `stake`. Those coins can then be staked and later withdrawn through the execution withdrawal queue, whose amount field is explicitly interpreted in gwei. The result is a deterministic `1e9` amplification: 1 wei burned on EL becomes 1 gwei-worth of claimable `stake` on CL. A permissionless attacker who controls one or more DKG validators can repeatedly trigger CDR reads and harvest vastly more value than was ever paid into the system.

### Details

Story’s staking bridge clearly treats CL `stake` units as `gwei`-scaled execution units. The execution-layer staking contract rounds stake amounts to `1 gwei`, and the CL normalizer divides `CreateValidator`, `Deposit`, `Redelegate`, and `Withdraw` event amounts by `1e9` before minting or burning Cosmos-side stake:

```solidity
/// @notice Stake amount increments. Consensus Layer requires staking in increments of 1 gwei.
uint256 public constant STAKE_ROUNDING = 1 gwei;
```

```go
gwei, exp := big.NewInt(10), big.NewInt(9)
gwei.Exp(gwei, exp, nil)
...
ev.StakeAmount.Div(ev.StakeAmount, gwei)
...
ev.Amount.Div(ev.Amount, gwei)
```

The execution-withdrawal type then uses `Amount` in gwei:

```go
// The withdrawn amount in Gwei
Amount uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
```

The CDR fee bridge omits that normalization entirely. The CDR contract emits `FeeCollected` using the raw payable fee amount, which is standard EVM `wei`:

```solidity
function _collectFee(uint256 feeAmountToCollect, ICDR.FeeType feeType) internal {
    require(msg.value == feeAmountToCollect, "CDR: Invalid fee amount");
    payable(address(0x0)).transfer(feeAmountToCollect);
    emit FeeCollected(msg.sender, feeAmountToCollect, feeType);
}
```

On the CL side, `ProcessCDRFeeCollected()` forwards that raw `ev.Amount` into `AddCDRFeeToPool()`, which mints the exact numeric value as `stake` without dividing by `1e9`:

```go
if err = k.dkgKeeper.AddCDRFeeToPool(cachedCtx, ev.Amount); err != nil {
    return errors.Wrap(err, "add CDR fee to pool")
}
```

```go
feeAmount, ok, err := parseCDRFeeAmount(amount)
...
coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, feeAmount))
if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
    return errors.Wrap(err, "mint CDR fee coins")
}
```

The same raw-wei amount is then used again for immediate validator refunds on valid partials:

```go
} else if ev.Fee != nil && ev.Fee.Sign() > 0 {
    if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
        partialErr = errors.Wrap(err, "refund CDR fee")
    }
}
```

```go
coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, feeAmount))
if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins); err != nil {
    return errors.Wrap(err, "refund CDR fee coins")
}
```

Finally, those CL `stake` balances are real economic value. Reward withdrawals enqueue the raw account balance into the execution withdrawal queue, and that queue’s amount is interpreted as gwei:

```go
delRewardUint64 := delRewards.AmountOf(sdk.DefaultBondDenom).Uint64()
...
types.NewWithdrawal(..., delRewardUint64, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, ...)
```

This produces the following concrete amplification:
1. A requester pays `F` wei to `CDR.read()` or a validator pays `F` wei to `submitEncryptedPartialDecryption()`.
2. EL burns exactly `F` wei.
3. CL mints `F` units of `stake` into `cdr-fee-pool`, even though CL staking units are gwei-scaled.
4. A valid partial immediately refunds `F` CL `stake` to the validator account.
5. When that balance is later bridged back through the reward/withdrawal queue, it becomes `F gwei = F * 1e9 wei` on EL.

So a single valid partial converts an EL fee payment of `F wei` into a CL/EL credit worth `F * 1e9 wei`.

### Impact Cascade

- Every valid partial submission prints `1e9x` more `stake` on CL than the fee actually burned on EL.
- The attacker receives that amplified value immediately via `RefundCDRFee`, without waiting for end-of-round pool distribution.
- Any remaining over-credited pool balance is distributed again in `distributeCDRRewardPool()`, compounding the inflation.
- The resulting CL `stake` is ordinary user balance, so it can be staked, transferred, and eventually withdrawn to EL through the gwei-denominated withdrawal queue.
- The protocol’s accounting no longer conserves value across EL burns, CL minting, and EL withdrawals; attackers can mint unbacked `stake` supply and drain economic value from the bridge model.

### Assumptions and Uncertainties

1. CL `stake` units are the same gwei-scaled units used everywhere else in the CL↔︎EL staking bridge; this is strongly supported by the explicit `1 gwei` normalization in `ProcessStakingEvents()` and the execution-withdrawal type comment.
2. A malicious actor can control at least one active DKG validator. This is a permissionless role in the protocol, not an admin-only capability.
3. The attacker can eventually realize the over-credited CL `stake` value either directly on CL or by bridging it back to EL via normal staking/withdrawal flows.

### How will the bug recipient respond?

“CDR fees are meant to be raw `stake` values, so no conversion is needed.”

That response conflicts with the surrounding bridge logic. The same codebase explicitly divides EL staking event amounts by `1 gwei` before minting Cosmos `stake`, and the execution withdrawal type explicitly documents that withdrawal amounts are interpreted in gwei. Without the same normalization on the CDR fee bridge, the protocol credits nine extra decimal places of value that were never burned on EL.

### Why did tests miss this issue?

The existing tests in `dkg_cdr_fees_internal_test.go` only check that numeric values round-trip unchanged (`100 -> 100`, `42 -> 42`). They never model the EL↔︎CL unit boundary, never compare CDR fee handling to the staking bridge’s explicit `gwei` division, and never trace the minted fee value through the gwei-denominated withdrawal queue.

### Recommendation

The bridge must normalize CDR fee amounts exactly once at the EL→CL boundary, just like staking events do.

Primary fix:

```go
var gwei = big.NewInt(1_000_000_000)

func normalizeELWeiToCLStakeUnits(amount *big.Int) (math.Int, bool, error) {
    if amount == nil || amount.Sign() == 0 {
        return math.ZeroInt(), false, nil
    }

    normalized := new(big.Int).Div(new(big.Int).Set(amount), gwei)
    if normalized.Sign() == 0 {
        return math.ZeroInt(), false, nil
    }

    feeAmount, ok := math.NewIntFromString(normalized.String())
    if !ok {
        return math.ZeroInt(), false, errors.New("invalid normalized CDR fee amount", "value", normalized.String())
    }

    return feeAmount, true, nil
}
```

Then use that normalization in both `AddCDRFeeToPool()` and `RefundCDRFee()` so the pool and refunds stay in the same unit system as the staking bridge.

Additional hardening:
- Add invariant tests that compare EL fee burns, CL fee-pool balance changes, and eventual EL withdrawal amounts end-to-end.
- Reject or explicitly round sub-gwei fee amounts if the protocol intends the bridge unit to remain gwei-scaled.
- Audit all other EL→CL fee/value bridges for the same missing normalization pattern.

### References

1. [story/contracts/src/protocol/IPTokenStaking.sol#L33]
2. [story/client/x/evmstaking/keeper/keeper.go#L116]
3. [story/client/x/evmstaking/keeper/keeper.go#L183]
4. [story/client/x/evmstaking/keeper/keeper.go#L196]
5. [story/client/x/evmstaking/keeper/keeper.go#L209]
6. [story/client/x/evmstaking/keeper/keeper.go#L222]
7. [story/client/x/evmengine/types/tx.pb.go#L348]
8. [story/contracts/src/protocol/CDR.sol#L325]
9. [story/client/x/evmengine/keeper/cdr.go#L231]
10. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L24]
11. [story/client/x/evmengine/keeper/cdr.go#L177]
12. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L55]
13. [story/client/x/evmstaking/keeper/withdraw.go#L388]
14. [story/client/x/evmstaking/keeper/withdraw.go#L422]
15. [story/client/x/dkg/keeper/contract_client.go#L244]

### STOR-6 — Validators can submit arbitrary or replayed partial decryptions and still collect CDR refunds and rewards

### Validators can submit arbitrary or replayed partial decryptions and still collect CDR refunds and rewards

### Executive Summary

The CDR bridge accepts a partial-decryption submission whenever the request exists, the validator’s stored `pubShare` matches, and the validator’s enclave communication key signed a payload. Crucially, the consensus layer never verifies that `encryptedPartial` is a real TDH2 partial for the request, and the signature preimage omits both `requesterPubKey` and `label`. A malicious committee member can therefore submit arbitrary signed garbage, or replay one previously generated partial across different requesters for the same ciphertext/round, and the CL will still mark the submission successful, refund the submission fee from `cdr-fee-pool`, and increment that validator’s reward counter.

This breaks the whitepaper’s core assumption that up to one-third malicious validators should not be able to monetize invalid decryption work. Instead, malicious validators can siphon CDR fees while delivering unusable shares to requesters.

### Details

The acceptance path only enforces request existence, timeout, registration lookup, and a signature over a truncated preimage:

```go
// client/x/dkg/keeper/dkg_handler.go
req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, round, ciphertext)
...
if !bytes.Equal(pubShare, reg.PubKeyShare) {
    return errors.New("pubShare mismatch")
}

if err := verifyPartialDecryptionSignature(
    reg.CommPubKey, round, ciphertext, encryptedPartial, ephemeralPubKey, pubShare, signature,
); err != nil {
    return errors.Wrap(err, "partial decryption signature verification failed")
}
```

The verifier itself signs only `round || ciphertext || encryptedPartial || ephemeralPubKey || pubShare`:

```go
// client/x/dkg/keeper/dkg_handler.go
encoded := make([]byte, 0, 4+len(ciphertext)+len(encryptedPartial)+len(ephemeralPubKey)+len(pubShare))
encoded = append(encoded, roundBytes...)
encoded = append(encoded, ciphertext...)
encoded = append(encoded, encryptedPartial...)
encoded = append(encoded, ephemeralPubKey...)
encoded = append(encoded, pubShare...)
```

`requesterPubKey`, `label`/`uuid`, and `pid` are not authenticated there, and there is no on-chain TDH2 correctness check before the submission is rewarded. Once the bridge processes the event, the validator is immediately credited:

```go
// client/x/evmengine/keeper/cdr.go
if partialErr == nil {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil {
        ...
    } else if ev.Fee != nil && ev.Fee.Sign() > 0 {
        if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
            ...
        }
    }
}
```

The Solidity side does not authenticate semantic correctness either; it only emits whatever the validator submitted:

```solidity
// contracts/src/protocol/CDR.sol
emit EncryptedPartialDecryptionSubmitted(
    msg.sender,
    round,
    pid,
    encryptedPartial,
    ephemeralPubKey,
    pubShare,
    requesterPubKey,
    ciphertext,
    uuid,
    signature,
    fee
);
```

Exploit path:
1. A requester triggers `read()` and the CL stores a decrypt request.
2. A malicious validator crafts arbitrary bytes for `encryptedPartial` / `ephemeralPubKey` (or replays an old signed payload tied to another requester) and signs only the truncated preimage.
3. `PartialDecryptionSubmitted` accepts the event because the request exists and the signature matches the validator’s comm key.
4. The bridge refunds the validator’s fee and increments `CDRPartialSubmitCount` even though the requester cannot use the share.
5. At fee-pool settlement, the malicious validator also receives a larger portion of `cdr-fee-pool` for work never actually performed.

### Impact Cascade

- Malicious committee members can monetize invalid decryption submissions without generating a valid TDH2 share.
- Requesters receive unusable partials and may fail to reconstruct plaintext even though the chain marked submissions successful.
- `cdr-fee-pool` distribution is skewed toward malicious validators because fake submissions increment `CDRPartialSubmitCount`.
- The protocol’s stated security model (tolerating up to one-third malicious validators) is violated at the reward-accounting layer.

### Assumptions and Uncertainties

1. At least one committee validator is malicious and controls its own enclave communication key; this is within the whitepaper’s stated threat model.
2. Requesters rely on CL-accepted submissions and later settlement accounting, which is exactly how `IncrementCDRPartialSubmitCount` / `RefundCDRFee` are wired today.
3. I did not find any later on-chain validation step that checks TDH2 correctness before fees or rewards are granted.

### How will the bug recipient respond?

“Correctness is checked off-chain by the requester, so on-chain validation is unnecessary.”

That does not address the bug. Even if the requester rejects the bad share later, the CL has already refunded the validator and increased its reward count. The protocol is paying for invalid work.

### Why did tests miss this issue?

Existing tests focus on signature-shape validation and happy-path fee accounting. They do not assert that a submitted partial is cryptographically valid for the specific `(requesterPubKey, label, pid)` request before rewards/refunds are granted.

### Recommendation

Bind the full request identity into the signed payload and refuse to reward a submission until it is request-specific and structurally valid.

A minimum hardening set is:

```go
// include every field that the requester needs to interpret the share
encoded := round || pid || requesterPubKey || label || ciphertext || encryptedPartial || ephemeralPubKey || pubShare
```

In addition:
- reject events whose `pid` does not match the validator’s registered index,
- require the request key `(requesterPubKey, label, round, ciphertext)` to be covered by the signature,
- and move fee refunds / reward counting behind a stronger validity gate (for example, a verifiable proof or a delayed payout once enough consistent shares are observed).

### References

1. [client/x/dkg/keeper/dkg_handler.go#L362-L409]
2. [client/x/dkg/keeper/dkg_handler.go#L491-L576]
3. [client/x/evmengine/keeper/cdr.go#L120-L185]
4. [contracts/src/protocol/CDR.sol#L223-L250]

### STOR-5 — Unauthenticated PartialDecryptTDH2 turns default kernels into a remote threshold-decryption oracle

### Unauthenticated PartialDecryptTDH2 turns default kernels into a remote threshold-decryption oracle

### Executive Summary

The `story-kernel` gRPC server exposes `PartialDecryptTDH2` without any built-in caller authentication, and the default server configuration binds that API to `:50051` on all interfaces while explicitly allowing plaintext transport when TLS is not configured. An attacker who can reach the port first calls `GetCodeCommitment()` to learn the enclave’s expected `code_commitment`, then invokes `PartialDecryptTDH2()` with any current-round ciphertext, the corresponding UUID-derived label, the current global public key, and an attacker-controlled requester public key. Inside the enclave, the method only checks that the supplied round matches the latest active network and that the supplied `code_commitment` equals the enclave’s own MRENCLAVE; it never proves that the request came from a canonical `VaultRead` event or that the requester/key/label tuple was ever authorized on-chain. The method then loads the sealed `DistKeyShare` and emits a real TDH2 partial decryption. Repeating this against a threshold of exposed validators lets a remote attacker decrypt protected vault contents without satisfying the read condition or leaving any on-chain request trail.

### Details

The network boundary is open by default. `story-kernel` enables no authentication interceptor at all, only panic recovery; if TLS is not configured it logs a warning and still starts serving. The default config listens on `:50051`, i.e. all interfaces, while the consensus client assumes the kernel is reachable at `127.0.0.1:50051`.

```go
func Serve(cfg *config.Config) (*grpc.Server, chan error) {
    serverOpts := []grpc.ServerOption{
        grpc.UnaryInterceptor(recoveryInterceptor()),
    }
    if cfg.GRPC.TLSEnabled() {
        serverOpts = append(serverOpts, grpc.Creds(tlsCreds))
    } else {
        log.Warn("gRPC server running without TLS. Set tls_cert_file and tls_key_file to enable.")
    }
    svr := grpc.NewServer(serverOpts...)
}
```

```go
GRPC: GRPCConfig{
    ListenAddr: ":50051",
}
```

```go
func DefaultDKGConfig() DKGConfig {
    return DKGConfig{
        KernelEndpoints: []string{"127.0.0.1:50051"},
    }
}
```

The API needed to satisfy the enclave’s `code_commitment` precondition is also openly exposed:

```go
func (s *DKGServer) GetCodeCommitment(...) (*pb.GetCodeCommitmentResponse, error) {
    codeCommitment, err := enclave.GetSelfCodeCommitment()
    return &pb.GetCodeCommitmentResponse{CodeCommitment: codeCommitment}, nil
}
```

Once the attacker knows the commitment, `PartialDecryptTDH2` performs only syntactic validation plus a latest-round equality check before using the sealed share. The method’s own TODO states that canonical-chain request verification is missing.

```go
//TODO: TEE should verify if the request transaction was indeed submitted to the canonical chain and the unique ID
// and round match to prevent any leakage of data by off-chain collusion.
func (s *DKGServer) PartialDecryptTDH2(ctx context.Context, req *pb.PartialDecryptTDH2Request) (*pb.PartialDecryptTDH2Response, error) {
    latestNetwork, err := s.verifyRoundMatchesLatestNetwork(ctx, req.GetRound())
    if err := enclave.ValidateCodeCommitment(req.GetCodeCommitment()); err != nil { ... }
    ownPID, ok := s.PIDCache.Get(req.GetRound())
    ...
    distKeyShare = load sealed share for (code_commitment, round)
    pd, err := mpc.TDH2PartialDecrypt(int(ownPID), privShare, pubKey, ct, req.GetLabel())
    encryptedPartial, ephPubKey, err := encryptPartialToRequester(req.GetRequesterPubKey(), pd.Bytes)
    signature, err := s.signPartialDecryptResponse(...)
}
```

The whitepaper says the kernel is supposed to “only act when smart contract events are provided with proofs” and should prevent committee members from “trigger[ing] unauthorized malicious actions”, but this path never asks the enclave to verify any proof of a `VaultRead` event or of a matching `DecryptRequestRegistry` entry before touching the sealed share.

Attack path:
1. Read the target ciphertext and UUID from chain data for the active DKG round.
2. Call `GetCodeCommitment()` on exposed validators to learn each enclave’s MRENCLAVE.
3. Call `PartialDecryptTDH2(round, ciphertext, label(uuid), globalPubKey, attackerRequesterPubKey)` directly on those kernels.
4. Each enclave returns a genuine partial encrypted to the attacker’s secp256k1 key.
5. After collecting threshold-many partials, combine them off-chain and recover the plaintext without ever satisfying the on-chain read condition.

### Impact Cascade

- Any reachable validator kernel becomes a decryption oracle for arbitrary current-round vault ciphertexts.
- The attacker does not need to call `CDR.read()` or satisfy any read condition, so access-control failures are silent.
- The TEE boundary no longer protects against a malicious network caller or host path; secrecy reduces to network exposure rather than enclave policy.
- With threshold-many exposed validators, protected vault plaintext can be recovered off-chain with no on-chain evidence beyond normal ciphertext availability.

### Assumptions and Uncertainties

1. The validator exposes the kernel gRPC port beyond a strictly local trust boundary; this is the out-of-the-box bind behavior because the server listens on `:50051` and plaintext is allowed when TLS is unset.
2. The target ciphertext belongs to the currently active DKG round, because `PartialDecryptTDH2` checks only `round == latest active round`.
3. Threshold-many validators are reachable through the same exposure pattern to fully decrypt the target ciphertext.

### How will the bug recipient respond?

“Operators are expected to keep the kernel on localhost or behind a firewall, and TLS can be enabled if they want.”

That response does not hold at the enclave boundary. The server itself chooses an all-interface default bind, accepts unauthenticated requests when TLS is absent, and exposes `GetCodeCommitment()` specifically so callers can satisfy `PartialDecryptTDH2`’s only enclave-identity check. More importantly, the whitepaper’s security claim is that kernels only act on proved smart-contract events; that guarantee must be enforced by the enclave itself, not outsourced to deployment hygiene.

### Why did tests miss this issue?

Existing tests only cover local cryptographic helpers, signature recovery, and PID bounds. They do not exercise the gRPC server with default networking/TLS settings, and they do not assert that `PartialDecryptTDH2` refuses requests lacking a light-client-verifiable `VaultRead` proof.

### Recommendation

The enclave should refuse to partial-decrypt unless the caller supplies a proof of a canonical, still-active `VaultRead` request that binds at minimum `round`, `ciphertext`, `uuid/label`, and `requester_pub_key` to the canonical chain. The proof should be validated inside `PartialDecryptTDH2` using the existing verified query client before any sealed key material is loaded.

A defense-in-depth hardening should also remove the remote attack surface entirely:

```go
GRPC: GRPCConfig{
    ListenAddr: "127.0.0.1:50051",
}
```

and require mTLS or another enclave-authenticated local transport for every gRPC method that can touch sealed state.

### References

1. [story-kernel/service/dkg_partial_decrypt.go#L40]
2. [story-kernel/service/dkg_get_code_commitment.go#L15]
3. [story-kernel/server/server.go#L32]
4. [story-kernel/server/server.go#L79]
5. [story-kernel/config/config.go#L220]
6. [story/client/config/dkg_config.go#L37]
7. [story-kernel/proto/tee.proto#L148]
8. [ai-docs/core/whitepaper/confidentialdatarails.md#L385]

### STOR-4 — Unknown or expired partial submissions are rewarded as valid work

### Unknown or expired partial submissions are rewarded as valid work

### Executive Summary

The confidential-read payout path treats several non-valid outcomes as success. `ProcessDKGPartialDecryptionSubmitted` refunds the submit fee and increments the sender’s reward weight whenever `dkgKeeper.PartialDecryptionSubmitted` returns `nil`. But `PartialDecryptionSubmitted` returns `nil` not only for a valid, stored partial — it also returns `nil` when the decrypt request does not exist anymore, or when the request has already timed out. Because `CDR.submitEncryptedPartialDecryption` is permissionless, any EVM address can emit arbitrary `EncryptedPartialDecryptionSubmitted` events with junk payloads. On the CL side, the not-found / timed-out branches exit before any DKG registration lookup, pubshare check, or signature verification runs, yet the sender still receives reward credit and a fee refund. Later, `distributeCDRRewardPool` pays the entire accumulated CDR fee pool pro-rata to those forged counts, and it does so for arbitrary addresses, not only committee members. The result is a gas-only drain of the user-funded CDR reward pool.

### Details

The entrypoint is fully permissionless. Any account can call `CDR.submitEncryptedPartialDecryption`; the contract only checks that `msg.value == baseFee`, burns the fee, and emits `EncryptedPartialDecryptionSubmitted` with `msg.sender` recorded as the validator.

```solidity
function submitEncryptedPartialDecryption(
    uint32 round,
    uint32 pid,
    bytes calldata encryptedPartial,
    bytes calldata ephemeralPubKey,
    bytes calldata pubShare,
    bytes calldata requesterPubKey,
    bytes calldata ciphertext,
    uint32 uuid,
    bytes calldata signature
) external payable whenNotPaused {
    uint256 fee = _getCDRStorage().baseFee;
    _collectFee(fee, ICDR.FeeType.SubmitPartial);

    emit EncryptedPartialDecryptionSubmitted(
        msg.sender,
        round,
        pid,
        encryptedPartial,
        ephemeralPubKey,
        pubShare,
        requesterPubKey,
        ciphertext,
        uuid,
        signature,
        fee
    );
}
```

The CL then processes that event. First it routes the burned fee into the CDR fee pool, and then it calls `PartialDecryptionSubmitted`. If that function returns `nil`, the CL increments `CDRPartialSubmitCount` and refunds the submit fee to `ev.Validator`.

```go
partialErr := k.dkgKeeper.PartialDecryptionSubmitted(...)

if partialErr == nil {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil {
        partialErr = errors.Wrap(err, "increment CDR submit count")
    } else if ev.Fee != nil && ev.Fee.Sign() > 0 {
        if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
            partialErr = errors.Wrap(err, "refund CDR fee")
        }
    }
}
```

The problem is that `PartialDecryptionSubmitted` uses `nil` for ignored / invalid branches. If the request is missing, it logs and returns `nil`. If the request has already timed out, it deletes the registry entry and returns `nil` again. Both branches execute before any validator registration lookup, pubshare comparison, or signature verification.

```go
req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, round, ciphertext)
if !found {
    log.Info(ctx, "Partial decryption submitted for unknown or cleaned-up request", ...)
    return nil
}

if currentHeight-req.Height > types.PartialDecryptionTimeoutBlocks {
    ...
    if err := k.deleteDecryptRequest(ctx, requesterPubKey, label, round, ciphertext); err != nil {
        return errors.Wrap(err, "failed to delete expired decrypt request registry entry")
    }
    return nil
}
```

This creates a complete exploit path for an unprivileged external account:
1. Call `CDR.submitEncryptedPartialDecryption` from any EVM address with arbitrary round, pid, ciphertext, requester key, uuid, pubshare, and signature.
2. The contract emits a valid Story event and a fee-collection event.
3. `ProcessCDRFeeCollected` mints the burned base fee into the CDR fee pool.
4. `PartialDecryptionSubmitted` hits the `!found` branch (or the timed-out branch) and returns `nil` without checking any validator registration or cryptographic proof.
5. `ProcessDKGPartialDecryptionSubmitted` interprets that `nil` as success, increments `CDRPartialSubmitCount` for the attacker-controlled address, and refunds the fee from the pool.
6. At `FinalizeDKGRound`, `distributeCDRRewardPool` pays the pool balance proportionally to those forged counts. Since the attacker can repeat the transaction arbitrarily many times at near-zero cost, they can dominate `totalCount` and drain essentially the whole pool.

Two implementation details make the theft complete rather than theoretical:
- Reward recipients are not constrained to committee members; `EvmAddressToBech32AccAddress` just reinterprets any 20-byte EVM address as an account address.
- Timed-out requests are only pruned every 1000 blocks while the timeout is 200 blocks, so even legitimate requests remain exploitable for hundreds of blocks after expiry.

### Impact Cascade

- Any external account can mint arbitrary `CDRPartialSubmitCount` entries without being a validator or owning a TEE.
- The submit fee is immediately refunded, so the attack is limited to gas costs.
- The attacker can force `totalCount` to be attacker-dominated before the next DKG finalization.
- `distributeCDRRewardPool` then transfers the accumulated user-funded CDR fee pool to the attacker’s address.
- Legitimate committee members are diluted out of the distribution even though they performed the actual decryption work.

### Assumptions and Uncertainties

1. The CDR fee pool contains non-zero user fees before `FinalizeDKGRound`; this is the intended steady-state because all CDR fees are routed into that pool.
2. The attacker can submit normal EVM transactions to the CDR predeploy; no special privileges are required by the contract.
3. Gas costs do not outweigh the fee-pool balance, which is realistic whenever the pool has meaningful accumulated fees.

### How will the bug recipient respond?

“Those branches intentionally ignore stale or missing submissions; no state is updated.”

That misses the exploit boundary. The ignored branches still return `nil`, and `ProcessDKGPartialDecryptionSubmitted` uses `nil` as the sole signal for refunding fees and incrementing reward counts. The dangerous state transition is not inside `PartialDecryptionSubmitted`; it happens immediately afterwards in the evmengine keeper.

### Why did tests miss this issue?

The unit tests cover fee accounting helpers in isolation, but there is no end-to-end test for `ProcessDKGPartialDecryptionSubmitted` asserting that ignored / missing / timed-out submissions must not trigger `IncrementCDRPartialSubmitCount` or `RefundCDRFee`. The available `cdr_internal_test.go` only checks `uuidToLabel`, so the critical read-path accounting branch is effectively untested.

### Recommendation

The success condition must be explicit instead of inferred from `err == nil`.

Primary fix:

```go
accepted, err := k.dkgKeeper.PartialDecryptionSubmitted(...)
if err != nil {
    ...
}
if accepted {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(...); err != nil { ... }
    if err := k.dkgKeeper.RefundCDRFee(...); err != nil { ... }
}
```

Additional hardening:
- Return a concrete error or `accepted=false` for unknown, expired, and duplicate submissions.
- Reject non-committee / non-registered senders before any payout logic.
- In `distributeCDRRewardPool`, ignore counts for addresses that were not members of the rewarded committee.

### References

1. [story/contracts/src/protocol/CDR.sol#L223-L250]
2. [story/contracts/src/protocol/CDR.sol#L325-L328]
3. [story/client/x/evmengine/keeper/cdr.go#L160-L180]
4. [story/client/x/evmengine/keeper/cdr.go#L193-L236]
5. [story/client/x/dkg/keeper/dkg_handler.go#L491-L547]
6. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L104-L120]
7. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L123-L218]
8. [story/client/x/dkg/keeper/dkg_finalization.go#L68-L76]
9. [story/client/x/dkg/keeper/abci.go#L72-L78]
10. [story/client/x/dkg/types/keys.go#L16-L23]
11. [story/client/server/utils/address.go#L12-L18]

### STOR-3 — Duplicate partial replays accrue fresh refunds and reward weight every time

### Duplicate partial replays accrue fresh refunds and reward weight every time

### Executive Summary

A single valid partial decryption can be replayed indefinitely for profit. The keeper correctly detects that a submission with the same `(requesterPubKey, label, ciphertext, round, validator)` already exists, but it treats that duplicate as a successful outcome by returning `nil`. The evmengine keeper then increments `CDRPartialSubmitCount` and refunds the submit fee again. Because the contract entrypoint is permissionless, a committee validator only needs one legitimate kernel-produced partial for a request; after that, the same signed payload can be resubmitted over and over. Every replay is counted when `distributeCDRRewardPool` divides the user-funded CDR fee pool, even though no additional threshold-decryption work was performed. This lets one malicious validator monopolize the reward pool with repeated copies of a single valid response.

### Details

Duplicate detection exists in storage, not in payout logic. `setPartialDecryptionSubmission` returns `ErrDuplicatePartialDecryptionSubmission` once the `(requesterPubKey, label, ciphertext, round, validator)` key already exists.

```go
key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, validator)
exists, err := k.DKGPartialDecrypt.Has(ctx, key)
if exists {
    return ErrDuplicatePartialDecryptionSubmission
}
```

`PartialDecryptionSubmitted` catches that duplicate sentinel, logs it, and returns `nil`.

```go
if err := k.setPartialDecryptionSubmission(...); err != nil {
    if errors.Is(err, ErrDuplicatePartialDecryptionSubmission) {
        log.Info(ctx, "Duplicate partial decryption submission received; ignoring", ...)
        return nil
    }
    return errors.Wrap(err, "failed to store partial decryption submission")
}
```

But the caller interprets `nil` as success and immediately pays out:

```go
if partialErr == nil {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil {
        ...
    } else if ev.Fee != nil && ev.Fee.Sign() > 0 {
        if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
            ...
        }
    }
}
```

Exploit path:
1. The attacker controls one registered validator in the active DKG committee.
2. For any legitimate read request, the validator lets the kernel produce one valid partial and submits it once.
3. The validator then resubmits the exact same event payload repeatedly through `CDR.submitEncryptedPartialDecryption`.
4. Each replay reaches the duplicate branch, which returns `nil`.
5. Each replay still increments the validator’s `CDRPartialSubmitCount` and refunds the base fee.
6. At the next `FinalizeDKGRound`, `distributeCDRRewardPool` allocates the pool by raw count, so the attacker can make their share arbitrarily close to 100% by spamming duplicates.

This is not just an accounting quirk. The pool being divided contains real user fees, and the attack reuses one genuine partial to claim an unbounded amount of that pool. The fee refund neutralizes the contract-side base fee, so the only marginal cost is gas.

### Impact Cascade

- One valid partial submission is enough to bootstrap the exploit.
- Every duplicate replay adds new reward weight despite zero new cryptographic work.
- The validator’s fee payment is refunded on every replay, making the strategy gas-only.
- Honest validators are diluted out of the fee-pool distribution.
- The malicious validator can drain essentially the entire CDR reward pool before the next round transition.

### Assumptions and Uncertainties

1. The attacker controls one validator that can produce at least one valid partial; this is realistic in a committee-adversary model.
2. The reward pool holds meaningful user fees before the next finalization.
3. Gas costs are lower than the attacker’s eventual share of the pool.

### How will the bug recipient respond?

“Duplicates are intentionally ignored, so they should be harmless.”

Ignoring duplicates is fine only if duplicates are also excluded from every downstream side effect. Here they are not ignored economically: the duplicate branch returns `nil`, and the caller uses `nil` to trigger both the refund and the reward-count increment.

### Why did tests miss this issue?

The duplicate sentinel is only exercised at the storage layer. There is no integration test that feeds a duplicate through `ProcessDKGPartialDecryptionSubmitted` and asserts that duplicate submissions must not call `IncrementCDRPartialSubmitCount` or `RefundCDRFee`.

### Recommendation

Treat duplicates as economically unsuccessful.

Primary fix:

```go
if errors.Is(err, ErrDuplicatePartialDecryptionSubmission) {
    return false, nil // or a typed non-success result
}
```

Then gate the refund and count increment on that explicit success flag. As a defense in depth, the event processor can also check whether a new row was actually inserted before mutating fee accounting.

### References

1. [story/contracts/src/protocol/CDR.sol#L223-L250]
2. [story/client/x/dkg/keeper/dkg_partial_decryption.go#L35-L75]
3. [story/client/x/dkg/keeper/dkg_handler.go#L549-L599]
4. [story/client/x/evmengine/keeper/cdr.go#L160-L180]
5. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L104-L120]
6. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L123-L218]
7. [story/client/x/dkg/keeper/dkg_finalization.go#L68-L76]

## Medium

### STOR-29 — Rebooting during an active round permanently disables threshold decryption for that round

### Rebooting during an active round permanently disables threshold decryption for that round

### Executive Summary

The DKG service persists completed sessions and queued decrypt requests to disk, but its recovery logic never restarts the decrypt worker for an already-completed active round. After any validator reboot during `DKGStageActive`, the node reloads the completed session, keeps accepting and charging new `VaultRead` requests, and stores them in the session JSON, yet no worker is re-armed to service those requests.

This turns a routine restart into fee-for-no-service for the rest of the round. Because the next automatic worker start only happens when some future DKG round completes, the pending requests from the current round miss the 200-block timeout long before another worker can exist. Users therefore pay CDR read fees during the active round even though the validator can no longer produce threshold decryptions until a later round.

### Details

Completed sessions and pending decrypt requests are explicitly persisted:

```go
type DKGSession struct {
    ...
    DecryptRequests []DecryptRequest `json:"decrypt_requests,omitempty"`
}
```

```go
if err := sm.loadSessions(); err != nil {
    return nil, errors.Wrap(err, "failed to load existing sessions")
}
```

`ThresholdDecryptRequested` appends the request to the session and writes it back to disk:

```go
session.AddDecryptRequest(types.DecryptRequest{ ... })
if err := k.stateManager.UpdateSession(ctx, session); err != nil {
    return errors.Wrap(err, "failed to persist decrypt request")
}
```

But on restart, recovery ignores already-completed sessions. `ResumeDKGService` only revives failed or intermediate phases; `PhaseCompleted` in `DKGStageActive` is treated as healthy and nothing calls `StartDecryptWorker` again:

```go
if session.Phase == types.PhaseFailed {
    k.resumeFailedSession(ctx, session, dkgNetwork)
    return
}

if isSessionStuckForStage(session.Phase, dkgNetwork.Stage) {
    ...
}
```

```go
case types.DKGStageActive:
    return phase != types.PhaseCompleted
```

The worker is only started from DKG completion:

```go
func (k *Keeper) handleDKGComplete(ctx context.Context, dkgNetwork *types.DKGNetwork) {
    ...
    k.StartDecryptWorker(ctx)
}
```

So a restart during an active round leaves the node with persisted decrypt requests but no execution path that can ever process them before `PartialDecryptionTimeoutBlocks = 200` expires.

### Impact Cascade

- A normal validator reboot during `DKGStageActive` kills threshold decryption for that validator for the rest of the round.
- The node still accepts `VaultRead` events, records decrypt requests, and charges read fees on-chain.
- Those queued requests are never serviced because recovery never re-arms the worker for `PhaseCompleted` sessions.
- After 200 blocks, the requests expire and users are left with paid-but-unfulfilled confidential reads.
- The economic damage persists until the next round, because the only automatic worker restart occurs on a future `handleDKGComplete`.

### Assumptions and Uncertainties

1. Validators rely on the built-in recovery logic and do not manually invoke an out-of-band worker restart after every reboot.
2. Reboots during a 21-day active round are realistic operational events.
3. User-facing `VaultRead` requests are expected to remain serviceable across validator restarts during the same round.

### How will the bug recipient respond?

“Operators should not reboot during the active period, or they can manually recover the service.”

That is not a sufficient defense because the code explicitly persists decrypt requests for crash recovery, yet the built-in recovery path refuses to restart the worker for the very completed session that owns those persisted requests. The protocol keeps charging users even after such a routine restart.

### Why did tests miss this issue?

The existing tests cover DKG phase transitions and some persistence helpers, but they do not simulate a restart in `DKGStageActive`, reload a completed session from disk, and verify that queued decrypt requests continue to be serviced afterward.

### Recommendation

On startup and in `ResumeDKGService`, explicitly detect a completed session for the current active round and restart the decrypt worker if it is not running.

Primary fix:

```go
func (k *Keeper) ResumeDKGService(ctx context.Context, dkgNetwork *types.DKGNetwork) {
    session, err := k.stateManager.GetSession(dkgNetwork.Round)
    if err != nil {
        ...
    }
    if dkgNetwork.Stage == types.DKGStageActive && session.Phase == types.PhaseCompleted {
        k.StartDecryptWorker(context.Background())
        return
    }
    ...
}
```

Also add a restart test that persists a completed session with pending decrypt requests, reinitializes `StateManager`, runs `BeginBlocker`, and asserts that partial decrypts are still produced.

### References

1. [story/client/x/dkg/types/dkg.go#L73-L81]
2. [story/client/x/dkg/keeper/state_manager.go#L41-L43]
3. [story/client/x/dkg/keeper/state_manager.go#L213-L246]
4. [story/client/x/dkg/keeper/dkg_handler.go#L463-L479]
5. [story/client/x/dkg/keeper/dkg_svc.go#L74-L112]
6. [story/client/x/dkg/keeper/dkg_svc.go#L118-L134]
7. [story/client/x/dkg/keeper/dkg_svc_complete.go#L45-L60]
8. [story/client/x/dkg/types/keys.go#L16-L23]

### STOR-28 — Decrypt worker dies after one minute while CDR keeps charging 21-day read fees

### Decrypt worker dies after one minute while CDR keeps charging 21-day read fees

### Executive Summary

The off-chain decrypt worker is started with a one-minute timeout context even though a DKG round stays active for 21 days by default. Once that minute elapses, every validator’s local decrypt worker exits permanently for the rest of the round. The chain still accepts `VaultRead` events, stores decrypt requests, and routes the associated CDR read fees into the reward pool, but no worker remains alive to submit any threshold partial decryptions before the 200-block request timeout expires.

This is not a mere availability issue. Users continue paying CDR read fees for requests that cannot complete, and those fees remain in the CDR fee pool for later validator distribution. In production defaults, the service window is roughly one minute out of a 21-day active period, so almost every paid confidential read after round activation becomes fee-for-no-service.

### Details

The async helper hardcodes a one-minute deadline:

```go
const dkgAsyncTimeout = 1 * time.Minute

func dkgAsyncContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), dkgAsyncTimeout)
}
```

`handleDKGComplete` starts the decrypt worker with exactly that timed context:

```go
func (k *Keeper) handleDKGComplete(ctx context.Context, dkgNetwork *types.DKGNetwork) {
    ...
    k.StartDecryptWorker(ctx)
}
```

The worker exits forever when the context is canceled:

```go
func (k *Keeper) StartDecryptWorker(ctx context.Context) {
    ...
    go func() {
        defer decryptWorkerRunning.Store(false)
        ticker := time.NewTicker(3 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                k.processDecryptQueue(ctx)
            }
        }
    }()
}
```

Meanwhile, the chain continues accepting CDR reads and queuing decrypt requests during the 21-day active round:

```go
latestRound, err := k.dkgKeeper.GetLatestActiveRound(cachedCtx)
...
round := latestRound.Round
...
if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, round, ev.RequesterPubKey, ev.Ciphertext, label[:], uint64(sdkCtx.BlockHeight())); ...
```

Those requests are rejected only after `PartialDecryptionTimeoutBlocks = 200`, while the active round lasts `DefaultDkgActivePeriod = 21 * 24 * 60 * 60` blocks. So once the worker dies, later requests will time out long before the protocol rotates to another round.

### Impact Cascade

- Every validator’s decrypt worker stops about one minute after round completion.
- Subsequent `VaultRead` calls still charge read fees and persist decrypt requests on-chain.
- No partial decryptions are produced before the 200-block timeout, so requesters pay but receive no confidential data.
- The collected fees remain in `cdr-fee-pool` and are later distributed to validators despite no decryption service having been delivered.
- The issue affects almost the entire 21-day active window, making CDR reads economically unsafe in normal operation.

### Assumptions and Uncertainties

1. Validators use the shipped DKG service path and do not run an out-of-band supervisor that restarts the worker every minute.
2. Production DKG params remain near the documented defaults, where the active period is far longer than the 200-block decrypt timeout.
3. Users expect paid `VaultRead` requests to receive threshold decrypt service throughout the active round, not only during its first minute.

### How will the bug recipient respond?

“This is just liveness; users can retry later or operators can restart the service manually.”

That response does not hold here because the protocol continues charging read fees and retaining those fees for later validator distribution while no built-in path restarts the worker during the same round. The economic harm is fee-for-no-service, not just temporary unavailability.

### Why did tests miss this issue?

Existing DKG lifecycle tests focus on stage transitions and registration/finalization logic. They do not exercise `StartDecryptWorker`, do not simulate real-time worker lifetime, and do not assert that threshold decrypt service remains available throughout an active round.

### Recommendation

Keep the decrypt worker on a long-lived lifecycle that matches the active DKG session rather than the one-minute async RPC helper.

Primary fix:

```go
func (k *Keeper) StartDecryptWorker() {
    if !decryptWorkerRunning.CompareAndSwap(false, true) {
        return
    }
    go func() {
        defer decryptWorkerRunning.Store(false)
        ticker := time.NewTicker(3 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            k.processDecryptQueue(context.Background())
        }
    }()
}
```

Also gate worker shutdown explicitly on round transitions / app shutdown, and add an integration test that submits `VaultRead` requests well after round activation and verifies they are still serviced.

### References

1. [story/client/x/dkg/keeper/dkg_svc.go#L60-L70]
2. [story/client/x/dkg/keeper/dkg_svc.go#L238-L260]
3. [story/client/x/dkg/keeper/dkg_svc_complete.go#L45-L60]
4. [story/client/x/evmengine/keeper/cdr.go#L57-L105]
5. [story/client/x/dkg/keeper/dkg_handler.go#L414-L486]
6. [story/client/x/dkg/types/keys.go#L16-L23]
7. [story/client/x/dkg/types/params.go#L12-L15]

### STOR-27 — Final EL-block partial submissions are settled in the wrong committee epoch

### Final EL-block partial submissions are settled in the wrong committee epoch

### Executive Summary

The CDR fee-pool settlement for a DKG committee runs in `dkg.BeginBlocker`, before the bridge has processed the previous EL block’s `EncryptedPartialDecryptionSubmitted` events. As a result, any partial submissions emitted in the final EL block of committee `R` are absent from `CDRPartialSubmitCount` when `distributeCDRRewardPool()` pays out committee `R`.

Those late events are then bridged during `DeliverTx` after the count map has already been cleared. They repopulate `CDRPartialSubmitCount` with stale submissions from committee `R`, and because reward settlement is not round-scoped, those counts are paid during the next committee settlement instead. This lets validators shift a chosen amount of submission credit across epochs by concentrating activity into the last EL block before a resharing boundary.

### Details

`FinalizeBlock` runs `BeginBlock -> DeliverTx -> EndBlock`, so DKG lifecycle transitions and reward settlement happen before `MsgExecutionPayload` processes the previous EL block’s logs.

```go
// FinalizeBlock calls BeginBlock -> DeliverTx (for all txs) -> EndBlock.
```

When a new round becomes active, `FinalizeDKGRound()` distributes the CDR fee pool and clears `CDRPartialSubmitCount`:

```go
if err := k.distributeCDRRewardPool(ctx); err != nil {
    return errors.Wrap(err, "failed to distribute CDR fee pool")
}
```

```go
if err := k.CDRPartialSubmitCount.Clear(ctx, nil); err != nil {
    return errors.Wrap(err, "clear CDR submit count")
}
```

But the previous EL block’s partial-submission events are only processed later in `DeliverTx`:

```go
if partialErr == nil {
    if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil { ... }
    ...
}
```

Because those late events arrive after the clear, the outgoing committee’s final-block submissions are not included in its own payout window. They instead sit in the global count map until the next `distributeCDRRewardPool()` call, where they are mixed into the next epoch.

Exploit path:
1. The attacker controls validators in committee `R`.
2. They delay valid partial submissions until the last EL block before `R` ends.
3. On the next CL block, `FinalizeDKGRound()` settles and clears committee `R` before those last EL-block events are bridged.
4. `DeliverTx` then increments `CDRPartialSubmitCount` for committee `R` after the clear.
5. Those stale counts are paid during the next epoch’s settlement, letting the attacker move submission credit across committee boundaries.

### Impact Cascade

- Partial-submission accounting is off by one epoch at every committee boundary.
- Outgoing committees can force their last-block work to be paid in the next settlement window.
- Next-epoch reward accounting is diluted by stale counts that do not belong to that service window.
- Committee compensation stops matching the EL block range in which the work actually occurred.

### Assumptions and Uncertainties

1. The attacker can time submissions into the final EL block before a known resharing boundary.
2. Committee compensation is intended to follow the committee epoch that actually served the request.
3. The protocol treats the bridge’s one-block delay as an implementation detail, not as an intended reward-shifting mechanism.

### How will the bug recipient respond?

“The one-block lag is expected, so shifting the final block into the next settlement is harmless.”

That does not hold for rewards, because settlement is explicitly keyed to committee epochs, not to arbitrary bridge lag. Clearing counts before bridging the last EL block makes compensation depend on implementation timing instead of on when the validator actually served the request.

### Why did tests miss this issue?

Tests cover `distributeCDRRewardPool()` and partial-submission processing independently, but they do not simulate a real `FinalizeBlock` boundary where `BeginBlocker` settles rewards before the final EL block’s logs are delivered.

### Recommendation

Settle CDR submission counts only after processing the previous EL block’s logs, or persist counts with an explicit round/epoch identifier and distribute only the counts that belong to the epoch being closed.

### References

1. [story/client/app/abci.go#L98-L105]
2. [story/client/x/dkg/keeper/abci.go#L89-L109]
3. [story/client/x/dkg/keeper/dkg_finalization.go#L68-L85]
4. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L123-L205]
5. [story/client/x/evmengine/keeper/cdr.go#L115-L190]

### STOR-26 — Permissionless pre-v2 CDR fees mutate the future DKG store before activation

### Permissionless pre-v2 CDR fees mutate the future DKG store before activation

### Executive Summary

`v2.0.0` is supposed to be the point where the DKG module becomes live, but the current rollout mounts the DKG store *before* that height and still routes pre-v2 CDR `FeeCollected` logs into `dkgKeeper.AddCDRFeeToPool()`. Any external user can trigger those logs by calling permissionless CDR entrypoints such as `allocate`, `write`, or `read`, all of which collect a fee and emit `FeeCollected`. On a node that has already switched to the `v2` binary, that pre-v2 event mints `sdk.DefaultBondDenom` into the DKG module account and updates `CDRFeePoolBalance`. Pre-v2 binaries do not have that module/store at all, because `v2.0.0` adds the DKG store via `StoreUpgrades`. The result is an attacker-triggerable state split across binaries exactly during the upgrade window: early-upgraded validators compute a different app state than old validators, and any validator replaying history with the latest binary after the upgrade also derives non-canonical pre-v2 DKG state.

### Details

`UpgradeStoreLoader` explicitly pre-adds future stores before their upgrade height so a new binary can start early, and `v2.0.0` declares the DKG KV store as a newly added store. That means a `v2` node can mount the DKG store before the chain has actually reached the activation height.

```go
// story/client/app/upgrades.go
// For future upgrades ... pre-add new stores
// covers rolling upgrades and late-joining validators.
if height > nextVersion {
    for _, key := range su.Added {
        if !addedSet[key] && !mountedStores[key] {
            merged.Added = append(merged.Added, key)
        }
    }
}
```

```go
// story/client/app/upgrades/v_2_0_0/constants.go
StoreUpgrades: storetypes.StoreUpgrades{
    Added: []string{dkgtypes.StoreKey},
}
```

The DKG module itself claims pre-v2 behavior should remain identical to the old binary, but only `BeginBlocker` is gated. EVM log delivery is not. `MsgExecutionPayload` finalization always dispatches DKG and CDR logs, regardless of height.

```go
// story/client/x/evmengine/keeper/msg_server.go
if err := s.ProcessDKGEvents(ctx, payload.Number-1, ethLogs); err != nil { ... }
if err := s.ProcessCDREvents(ctx, payload.Number-1, ethLogs); err != nil { ... }
```

CDR fee collection is permissionless. `allocate`, `write`, and `read` all burn ETH on the EVM side and emit `FeeCollected`.

```solidity
// story/contracts/src/protocol/CDR.sol
function allocate(...) external payable whenNotPaused returns (uint32 newVaultUuid) {
    _collectFee($.allocateFee, ICDR.FeeType.Allocate);
    ...
}

function write(...) external payable nonReentrant whenNotPaused {
    _collectFee($.writeFee, ICDR.FeeType.Write);
    ...
}

function read(...) external payable nonReentrant whenNotPaused {
    _collectFee($.readFee, ICDR.FeeType.Read);
    emit VaultRead(...);
}

function _collectFee(uint256 feeAmountToCollect, ICDR.FeeType feeType) internal {
    require(msg.value == feeAmountToCollect, "CDR: Invalid fee amount");
    payable(address(0x0)).transfer(feeAmountToCollect);
    emit FeeCollected(msg.sender, feeAmountToCollect, feeType);
}
```

When that log reaches the CL, `ProcessCDRFeeCollected()` unconditionally calls `AddCDRFeeToPool()`, which mints new bond-denom coins into the DKG module account and updates the DKG fee-pool state.

```go
// story/client/x/evmengine/keeper/cdr.go
if err = k.dkgKeeper.AddCDRFeeToPool(cachedCtx, ev.Amount); err != nil {
    return errors.Wrap(err, "add CDR fee to pool")
}
```

```go
// story/client/x/dkg/keeper/dkg_cdr_fees.go
coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, feeAmount))
if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil { ... }
if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.CDRFeePoolName, coins); err != nil { ... }
if err := k.CDRFeePoolBalance.Set(ctx, newBalance.String()); err != nil { ... }
```

This contradicts the intended rollout boundary:

```go
// story/client/x/dkg/keeper/abci.go
// DKG module activates at the v2.0.0 upgrade height. Before that,
// BeginBlocker is a complete no-op to ensure identical behavior to
// the pre-upgrade binary during rolling upgrades.
```

A concrete exploit path is straightforward:
1. Wait until some validators have restarted on the `v2` binary but the chain is still below the `v2.0.0` activation height.
2. Submit any permissionless CDR call (`allocate`, `write`, or `read`) with the required fee.
3. The EVM emits `FeeCollected` in an otherwise valid execution payload.
4. `v2` validators mint DKG fee-pool coins and update DKG store state; old validators do not have that module/store at all.
5. The same block now has two different state-transition functions depending on binary version.

Because the current binary also replays historic blocks through the same ungated path, any late-joining validator that syncs with the `v2` binary after the upgrade will derive the same non-canonical pre-v2 DKG writes for every historical pre-v2 `FeeCollected` event.

### Impact Cascade

- Any external user can trigger a pre-v2 block whose app-state transition differs across old and new binaries.
- Early-upgraded validators mint DKG fee-pool balances that the canonical pre-v2 chain never had.
- Mixed-version rollout can be bricked by a single permissionless CDR fee event.
- Late joiners replaying history with the latest binary derive non-canonical pre-v2 DKG state and cannot safely sync.
- The split is asset-bearing, not just metadata drift, because the upgraded path mints `sdk.DefaultBondDenom` and mutates module balances.

### Assumptions and Uncertainties

1. At least one permissionless CDR fee-bearing transaction lands before the `v2.0.0` activation height.
2. Either some validators upgrade before the height, or a node later replays/syncs historical pre-v2 blocks with the `v2` binary.
3. The pre-v2 production binary does not mount the DKG store, which is consistent with `v2.0.0` adding `dkgtypes.StoreKey` as a new store.

### How will the bug recipient respond?

“Validators are expected to coordinate the rollout, and the comment already says everyone must switch together.”

That does not solve the bug. The code explicitly claims to support early startup for rolling upgrades and late-joining validators, yet a normal permissionless CDR transaction makes the pre-v2 state transition depend on which binary is running. That is a protocol-level safety failure, not an operator convenience issue.

### Why did tests miss this issue?

The existing `v2.0.0` tests focus on vote-extension activation and migration idempotence; they do not exercise historical pre-v2 EVM log replay into a future-mounted DKG store, nor do they compare pre-v2 state transitions across old/new binaries.

### Recommendation

Gate *all* DKG and CDR-to-DKG state mutations on actual DKG activation, not just `BeginBlocker`. A safe fix is to reject or ignore `ProcessDKGEvents` / `ProcessCDREvents` until the DKG module is live at the current height, or until migration has explicitly marked the module active. As defense in depth, the migration should also assert the new store is empty before calling `RunMigrations`, and fail fast if any pre-activation writes already occurred.

### References

1. [story/client/app/upgrades.go#L82-L156]
2. [story/client/app/upgrades/v_2_0_0/constants.go#L17-L23]
3. [story/client/app/upgrades/v_2_0_0/upgrades.go#L27-L40]
4. [story/client/x/dkg/keeper/abci.go#L18-L27]
5. [story/client/x/evmengine/keeper/msg_server.go#L157-L174]
6. [story/client/x/evmengine/keeper/cdr.go#L193-L235]
7. [story/client/x/dkg/keeper/dkg_cdr_fees.go#L23-L48]
8. [story/contracts/src/protocol/CDR.sol#L108-L113]
9. [story/contracts/src/protocol/CDR.sol#L143-L168]
10. [story/contracts/src/protocol/CDR.sol#L180-L204]
11. [story/contracts/src/protocol/CDR.sol#L322-L328]

### STOR-25 — Single proposer can permanently censor committed DKG complaints by substituting MsgAddDkgVote

### Single proposer can permanently censor committed DKG complaints by substituting MsgAddDkgVote

### Executive Summary

The DKG transport assumes that the `MsgAddDkgVote` inside the proposal is a faithful aggregation of the previous block’s committed vote extensions. That invariant is never enforced. During `ExtendVote`, each validator destructively dequeues its local deals, responses, and justifications; once a message is placed into a vote extension it is removed from the local retry path. However, during proposal processing the chain only checks that one `MsgAddDkgVote` exists and that its `authority` equals the module address. It never verifies that the vote body matches `ProposedLastCommit` / `LocalLastCommit`, nor that every committed DKG item was carried forward.

A Byzantine proposer can therefore build a valid proposal that omits selected complaint responses or dealer justifications from the previous commit while still passing `ProcessProposal`. Because the omitted items were already dequeued from honest validators, they are lost forever. This lets the proposer suppress the only consensus path that invalidates malicious dealers, allowing them to survive finalization, remain in the rewarded committee, and continue serving as threshold-decryption participants.

### Details

`ExtendVote` removes locally generated DKG messages from in-memory queues before any proposal-level acknowledgment exists:

```go
dequeuedDeals := k.DequeueDeals(maxItemsPerVote)
dequeuedResponses := k.DequeueResponses(maxItemsPerVote)
dequeuedJustifications := k.DequeueJustifications(maxItemsPerVote)
```

The process-proposal router only checks message counts/types, then hands the DKG message to a proposal server that performs an authority check and nothing else:

```go
if req.Height > v200Height+1 {
    expectedMsgCounts[sdk.MsgTypeURL(&dkgtypes.MsgAddDkgVote{})] = 1
}
...
if _, err := handler(ctx, msg); err != nil {
    return rejectProposal(ctx, errors.Wrap(err, "execute message"))
}
```

```go
func (s proposalServer) AddVote(ctx context.Context, msg *types.MsgAddDkgVote) (*types.AddDkgVoteResponse, error) {
    if msg.Authority != s.Keeper.GetAuthority() {
        return nil, errors.New("unauthorized")
    }
    return &types.AddDkgVoteResponse{}, nil
}
```

Finalization then executes whatever vote body the proposer supplied:

```go
if len(msg.Vote.Responses) > 0 {
    if err := s.ProcessResponses(ctx, latestRound, msg.Vote.Responses); err != nil { ... }
}
if len(msg.Vote.Justifications) > 0 {
    if err := s.ProcessJustifications(ctx, latestRound, msg.Vote.Justifications); err != nil { ... }
}
```

The omitted-message effect is permanent because there is no “proposal rejected / message omitted” requeue path anywhere after `Dequeue*`. For justifications, that directly suppresses the only consensus-side invalidation hook:

```go
if !valid {
    if err := k.invalidateDealerRegistration(ctx, latestRound, j.Index); err != nil { ... }
    continue
}
```

Once the round finalizes, rewards are paid to whichever registrations remain finalized:

```go
finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round.Round, types.DKGRegStatusFinalized)
...
if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, perMemberCoins); err != nil { ... }
```

Exploit path:
1. A malicious dealer sends invalid or selectively inconsistent deals during the dealing phase.
2. Honest recipients generate complaint responses / justifications and broadcast them via vote extensions.
3. The next proposer is Byzantine and includes a handcrafted `MsgAddDkgVote` that omits exactly those complaint artifacts while still keeping one syntactically valid DKG message in the block.
4. `ProcessProposal` accepts because it only checks message type/count and `proposalServer.AddVote` only checks `authority`.
5. Honest complaints never reach `FinalizeBlock`; the malicious dealer is not invalidated and can still be part of the finalized / rewarded committee.

### Impact Cascade

- Complaint responses and dealer justifications can be censored by a single proposer after they already reached committed vote extensions.
- Malicious dealers avoid `invalidateDealerRegistration`, even though honest committee members produced the evidence needed to eject them.
- Finalized committee composition becomes proposer-controlled rather than vote-extension-controlled.
- UBI / CDR rewards are paid to the surviving finalized set, so malicious committee members continue receiving funds from the reward pool.
- The active threshold-decryption committee can retain malicious participants that the complaint flow was supposed to remove.

### Assumptions and Uncertainties

1. The attacker controls at least one proposer during the DKG dealing window; this is realistic in the stated Byzantine-validator threat model.
2. At least one malicious dealer has generated complaint-worthy DKG traffic that honest recipients tried to transport via vote extensions.
3. Honest validators do not locally re-enqueue omitted messages, which matches the current queue implementation.

### How will the bug recipient respond?

“The proposer always uses `PrepareVotes`, so this body cannot differ from committed vote extensions.”

That is only true for an honest proposer binary. `ProcessProposal` never enforces that invariant against a Byzantine proposer, which is exactly the actor CometBFT expects the application to defend against.

### Why did tests miss this issue?

Existing tests focus on message presence/count and authority routing. They do not construct a proposal whose `MsgAddDkgVote` body diverges from `ProposedLastCommit`, nor do they model the destructive dequeue behavior that makes omissions permanent.

### Recommendation

Reject proposals unless the DKG vote body is a deterministic re-aggregation of `req.ProposedLastCommit` / committed vote extensions. The check must run inside proposal validation, not only inside the proposer helper. In addition, make queue removal acknowledgement-based: either re-enqueue items not carried into the accepted proposal, or persist them until a finalized block proves inclusion.

### References

1. [story/client/x/dkg/keeper/vote.go#L28-L52]
2. [story/client/x/dkg/keeper/queue.go#L15-L33]
3. [story/client/x/dkg/keeper/queue.go#L44-L61]
4. [story/client/x/dkg/keeper/queue.go#L75-L93]
5. [story/client/app/prouter.go#L76-L130]
6. [story/client/x/dkg/keeper/proposal_server.go#L15-L23]
7. [story/client/x/dkg/keeper/msg_server.go#L16-L51]
8. [story/client/x/dkg/keeper/dkg_dealing.go#L150-L178]
9. [story/client/x/dkg/keeper/dkg_rewards.go#L58-L124]

### STOR-24 — Global first-10 truncation permanently discards valid dealer-invalidating justifications

### Global first-10 truncation permanently discards valid dealer-invalidating justifications

### Executive Summary

The DKG vote-extension pipeline admits vastly more justifications than finalization will execute. Each validator may include up to 80 justifications in its vote extension, and the proposer may aggregate all of them into one `MsgAddDkgVote`. But during `FinalizeBlock`, `ProcessJustifications` silently slices the aggregated list down to the first 10 entries and drops the rest. There is no requeue, persistence, or retry path for the discarded evidence.

This is not just a throughput issue. Justifications are the only consensus path that can invalidate a dealer whose complaints prove its deal was wrong. Once a justification is dropped by the first-10 slice, the related dealer is never invalidated even though valid evidence existed. Because vote-extension items were already dequeued from honest validators, the dropped evidence is gone forever. A malicious proposer or a small colluding set can exploit this by flooding the first 10 slots with harmless signed justifications, thereby protecting targeted malicious dealers from invalidation while preserving their ability to finalize and collect UBI/CDR rewards.

### Details

The vote-extension verifier allows up to 80 justifications per validator:

```go
if len(vote.Justifications) > maxItemsPerVote {
    return nil, false, fmt.Errorf("too many justifications ...")
}
```

The proposer helper aggregates all committed justifications into a single `Vote`:

```go
return &types.Vote{
    Deals:          deduplicateDeals(allDeals),
    Responses:      deduplicateResponses(allResponses),
    Justifications: deduplicateJustifications(allJustifications),
}
```

Finalization then silently truncates to ten:

```go
if len(justifications) > MaxJustificationsPerBlock {
    justifications = justifications[:MaxJustificationsPerBlock]
}
```

Yet invalidation only occurs inside this function:

```go
if !valid {
    if err := k.invalidateDealerRegistration(ctx, latestRound, j.Index); err != nil { ... }
    continue
}
```

Justifications themselves are produced by response processing and pushed into the same destructive queue used by vote extensions:

```go
if processResp != nil && len(processResp.GetJustifications()) > 0 {
    k.EnqueueJustifications(processResp.GetJustifications())
}
```

```go
out := make([]types.Justification, count)
copy(out, justifications[:count])
justifications = justifications[count:]
```

Exploit path:
1. Colluding validators create more than 10 signed complaint / justification pairs (easy because the transport budget is 80 per validator).
2. The malicious proposer orders those harmless justifications first inside `MsgAddDkgVote`.
3. Honest justifications that would invalidate malicious dealers appear after slot 10.
4. `ProcessJustifications` slices them away, and the omitted evidence is never retried because the queue entry was already dequeued during `ExtendVote`.
5. Malicious dealers avoid invalidation and remain eligible to finalize and receive committee rewards.

### Impact Cascade

- Honest dealer-invalidating evidence can be dropped even when it was validly transported through vote extensions.
- The first 10 justifications in proposer-controlled order become the only ones that affect consensus state.
- Small colluding sets can crowd out evidence against targeted malicious dealers with inexpensive filler justifications.
- Dealers who should have been invalidated can remain finalized committee members.
- Reward payouts derived from finalized registrations are redirected toward the attacker-controlled surviving set.

### Assumptions and Uncertainties

1. Attackers can generate more than 10 valid signed justifications during the dealing window; this is realistic because the upstream transport allows 80 per validator and justification generation is automatic once complaint responses are processed.
2. The attacker controls the proposer that chooses the ordering inside `MsgAddDkgVote`, or otherwise can rely on the honest aggregation order being unfavorable.
3. Discarded justifications are not reintroduced later, which matches the current dequeue-and-drop implementation.

### How will the bug recipient respond?

“The 10-item cap is only a DoS safeguard; remaining justifications can be carried in later blocks.”

They cannot. The system removes justifications from the local queue when `ExtendVote` runs, and there is no mechanism that re-enqueues the overflow discarded by `ProcessJustifications`.

### Why did tests miss this issue?

The tests explicitly assert that truncation to `MaxJustificationsPerBlock` works, but they treat it as benign slicing. They do not combine that slice with the one-shot dequeue semantics and the fact that invalidation only happens inside the truncated execution window.

### Recommendation

Make justification processing complete and retryable. Either (a) reject any proposal whose aggregated justifications exceed the finalize-time budget, (b) persist overflow justifications and carry them into subsequent blocks deterministically, or (c) move the cap to the vote-extension aggregation stage so that proposers cannot cause silent finalize-time evidence loss.

### References

1. [story/client/x/dkg/keeper/vote.go#L92-L104]
2. [story/client/x/dkg/keeper/vote.go#L159-L173]
3. [story/client/x/dkg/keeper/dkg_dealing.go#L116-L124]
4. [story/client/x/dkg/keeper/dkg_dealing.go#L150-L178]
5. [story/client/x/dkg/keeper/dkg_svc_dealing.go#L334-L344]
6. [story/client/x/dkg/keeper/queue.go#L63-L93]
7. [story/client/x/dkg/keeper/dkg_rewards.go#L58-L124]

### STOR-23 — Dealing-to-finalization transitions deterministically drop the last block’s DKG complaints and preserve bad committee members

### Dealing-to-finalization transitions deterministically drop the last block’s DKG complaints and preserve bad committee members

### Executive Summary

`MsgAddDkgVote` carries vote-extension data from the previous block, but `msgServer.AddVote` decides whether to process that data by looking at the **current** round stage after `BeginBlocker` has already run. That creates a deterministic stage-boundary loss: on the first block of `DKGStageFinalization`, the proposal still carries the final vote extensions from the previous dealing block, but `msgServer.AddVote` silently skips them because the round has already transitioned out of `DKGStageDealing`.

This is not proposer censorship and does not require any malformed tx. A malicious dealer can wait until the last dealing block to trigger complaint justifications about its invalid deal; honest validators will include those complaints in vote extensions, yet the next block will discard them automatically at the Dealing→Finalization transition. Since invalidation is the only on-chain path that prevents a bad dealer from finalizing, the attacker keeps its finalized registration and continues collecting DKG-linked UBI/CDR rewards.

### Details

DKG stages transition in `BeginBlocker`, before the block’s transactions are executed. As soon as the dealing period elapses, the round is moved into finalization:

```go
case types.DKGStageDealing:
    if elapsed >= dealingEnd {
        return types.DKGStageFinalization, true
    }
```

```go
if shouldTransition {
    latestRound.Stage = nextStage
    if err := k.setDKGNetwork(ctx, latestRound); err != nil {
        return err
    }

    switch nextStage {
    case types.DKGStageFinalization:
        return k.BeginFinalization(ctx, latestRound)
    }
}
```

But the `MsgAddDkgVote` included in that same block still represents `LocalLastCommit`, i.e. the **previous** block’s vote extensions. `PrepareProposal` explicitly builds it from `req.LocalLastCommit`:

```go
voteMsg, err := k.voteProvider.PrepareVotes(ctx, req.LocalLastCommit, uint64(req.Height-1))
```

During finalize, the DKG msg server processes those previous-block votes only if the **current** round stage is still dealing:

```go
latestRound, err := s.GetLatestDKGRound(ctx)
...
if latestRound != nil && latestRound.Stage == types.DKGStageDealing {
    if len(msg.Vote.Deals) > 0 { ... }
    if len(msg.Vote.Responses) > 0 { ... }
    if len(msg.Vote.Justifications) > 0 {
        if err := s.ProcessJustifications(ctx, latestRound, msg.Vote.Justifications); err != nil { ... }
    }
}
```

So on the first finalization block:
1. `BeginBlocker` advances the round from `DKGStageDealing` to `DKGStageFinalization`.
2. The proposal still carries the final dealing block’s vote extensions via `LocalLastCommit`.
3. `msgServer.AddVote` checks the new current stage, sees `Finalization`, and discards the vote payload without processing it.
4. Any complaint justifications from the last dealing block are lost forever.

This is especially dangerous because complaint data is ephemeral and single-use. Honest validators dequeue justifications out of local queues when producing vote extensions:

```go
dequeuedJustifications := k.DequeueJustifications(maxItemsPerVote)
```

And `ProcessJustifications` is the only path that can invalidate a malicious dealer:

```go
if !valid {
    if err := k.invalidateDealerRegistration(ctx, latestRound, j.Index); err != nil { ... }
    continue
}
```

Once those last-block complaints are skipped, the malicious dealer can still finalize because finalization only blocks registrations that were actually marked invalidated:

```go
if reg.Status == types.DKGRegStatusInvalidated {
    return errors.New("dealer has been invalidated and cannot finalize")
}
```

The dealer then remains eligible for the same committee reward paths used throughout the round lifecycle:

```go
if err := k.settleRewardsForPreviousCommittee(ctx); err != nil { ... }
if err := k.distributeCDRRewardPool(ctx); err != nil { ... }
```

```go
return k.distributeRewardsFromModule(ctx, activeRound, senderModule, totalAmount)
```

### Impact Cascade

- Deterministically drops honest complaints from the final dealing block even with an honest proposer.
- Lets a malicious dealer place invalid deals near the stage boundary and avoid on-chain invalidation.
- Preserves the attacker’s ability to finalize and remain inside the reward-paying DKG committee.
- Directly increases the attacker’s share of UBI and CDR fee-pool distributions.

### Assumptions and Uncertainties

1. The attacker is a DKG participant and can cause complaint justifications to be generated near the end of the dealing period.
2. Complaint justifications are not replayed later; the current implementation eagerly dequeues them for vote extensions.
3. The malicious dealer benefits economically from keeping finalized status, which is exactly how DKG committee rewards are distributed.

### How will the bug recipient respond?

“Those vote extensions belong to the previous block, and the current stage is finalization, so dropping them is harmless.”

It is not harmless because those previous-block vote extensions are the canonical transport for the **last** dealing block’s complaints. The stage transition happens before tx execution, so the system is deterministically discarding valid dealing-stage evidence before it can ever reach the invalidation logic.

### Why did tests miss this issue?

The lifecycle tests verify that stage transitions happen at the configured heights, but they do not exercise `MsgAddDkgVote` on the exact boundary block where `LocalLastCommit` still contains dealing-stage data and `BeginBlocker` has already moved the round into finalization.

### Recommendation

Bind DKG vote processing to the stage/round that produced the vote extensions, not the current post-`BeginBlock` stage. For example, include the source round or commit height in `MsgAddDkgVote` and accept the final dealing block’s votes while transitioning:

```go
if msg.SourceRound != latestRound.Round || msg.SourceStage != types.DKGStageDealing {
    return nil, errors.New("unexpected dkg vote source")
}
```

Alternatively, delay the Dealing→Finalization transition until after the previous block’s `MsgAddDkgVote` has been processed, or locally recompute/process the final dealing vote extensions before flipping the stage.

### References

1. [story/client/x/dkg/keeper/dkg_round.go#L10]
2. [story/client/x/dkg/keeper/abci.go#L89]
3. [story/client/x/evmengine/keeper/abci.go#L193]
4. [story/client/x/dkg/keeper/msg_server.go#L16]
5. [story/client/x/dkg/keeper/vote.go#L28]
6. [story/client/x/dkg/keeper/dkg_dealing.go#L104]
7. [story/client/x/dkg/keeper/dkg_handler.go#L129]
8. [story/client/x/dkg/keeper/dkg_finalization.go#L68]
9. [story/client/x/dkg/keeper/dkg_rewards.go#L18]

### STOR-22 — The first active block rejects the last finalization block’s DKG finalize events, excluding honest validators from the paid committee

### The first active block rejects the last finalization block’s DKG finalize events, excluding honest validators from the paid committee

### Executive Summary

Story imports EVM `DKG.Finalized` logs one block late through `MsgExecutionPayload`, but `BeginBlocker` advances the round from `DKGStageFinalization` to `DKGStageActive` **before** those previous-block logs are processed. `dkgKeeper.Finalized` then refuses to accept any `Finalized` event unless the current round stage is still `DKGStageFinalization`. The result is deterministic: all finalize transactions emitted in the final allowed EVM finalization block are dropped on the floor when the next CL block begins.

This is exploitable economically. A malicious validator can finalize early while honest validators finalize near the end of the finalization window; the boundary logic strips those late honest finalizations from the committee snapshot that `FinalizeDKGRound` uses for activation and reward accounting. The attacker ends up inside a smaller paid committee and captures a larger share of DKG-linked UBI/CDR distributions.

### Details

Stage transitions happen in `BeginBlocker`. Once the finalization period has elapsed, the round is immediately moved into `Active`, and `FinalizeDKGRound` runs right there in begin block:

```go
case types.DKGStageFinalization:
    if elapsed >= finalizationEnd {
        return types.DKGStageActive, true
    }
```

```go
case types.DKGStageActive:
    return k.FinalizeDKGRound(ctx, latestRound)
```

The previous execution payload’s logs are only imported later, inside the block transaction executed by `MsgExecutionPayload`:

```go
if err := s.ProcessDKGEvents(ctx, payload.Number-1, ethLogs); err != nil {
    return nil, errors.Wrap(err, "deliver dkg-related event logs")
}
```

But `dkgKeeper.Finalized` rejects those imported events unless the round is still in finalization:

```go
if latest.Stage != types.DKGStageFinalization {
    return errors.New("round is not in network set stage")
}
```

Therefore, for the first block after the finalization window expires:
1. `BeginBlocker` flips the round to `Active` and calls `FinalizeDKGRound` before tx execution.
2. The proposal’s `MsgExecutionPayload` still carries the previous block’s EVM `DKG.Finalized` logs.
3. `ProcessDKGFinalized` forwards those logs into `dkgKeeper.Finalized`.
4. `dkgKeeper.Finalized` rejects them because the current stage is already `Active`.
5. Honest validators who finalized in the last EVM finalization block are omitted from the finalized set used for activation and reward distribution.

That omission directly affects who gets paid, because both per-round settlement and ongoing committee rewards iterate the finalized registrations of the activated round:

```go
finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round.Round, types.DKGRegStatusFinalized)
```

```go
if err := k.settleRewardsForPreviousCommittee(ctx); err != nil { ... }
if err := k.distributeCDRRewardPool(ctx); err != nil { ... }
```

```go
return k.distributeRewardsFromModule(ctx, activeRound, senderModule, totalAmount)
```

### Impact Cascade

- Deterministically excludes honest validators who finalize in the last permitted EL block.
- Lets early-finalizing attackers remain in a smaller activated committee.
- Increases each attacker’s share of DKG-linked UBI and CDR fee-pool rewards.
- Can also push finalized count below threshold, forcing a skip/retry and stripping honest finalizers from the intended round.

### Assumptions and Uncertainties

1. Honest validators may finalize near the end of the finalization window, which is realistic because kernel finalization depends on off-chain TEE work and transaction inclusion timing.
2. The attacker is a validator who finalized earlier and therefore benefits from a smaller finalized set.
3. No compensating replay path re-applies the dropped `Finalized` events after activation; the current code rejects them outright.

### How will the bug recipient respond?

“Finalization had already ended, so those late events are out of window.”

The contract events are not actually late relative to the protocol’s own EL→CL bridge: they were emitted in the final allowed EL block and delivered through the normal `MsgExecutionPayload` path. The CL simply advances the stage one step too early and then rejects its own canonical bridge input.

### Why did tests miss this issue?

The lifecycle tests assert stage transitions, but they do not simulate `DKG.Finalized` logs landing in the final EL finalization block and being imported one block later through `MsgExecutionPayload` after `BeginBlocker` has already activated the round.

### Recommendation

Consume the previous block’s finalization events before switching to `Active`, or bind each imported finalization event to the stage/window of the block that emitted it instead of the current post-begin-block stage.

```go
if msg.SourceHeight == ctx.BlockHeight()-1 && latest.Stage == types.DKGStageActive {
    allowLateFinalizationForPreviousWindow = true
}
```

The safer fix is to move the Finalization→Active transition until after `MsgExecutionPayload` has processed the previous block’s DKG logs.

### References

1. [story/client/x/dkg/keeper/dkg_round.go#L14]
2. [story/client/x/dkg/keeper/abci.go#L89]
3. [story/client/x/dkg/keeper/dkg_finalization.go#L29]
4. [story/client/x/evmengine/keeper/msg_server.go#L170]
5. [story/client/x/evmengine/keeper/dkg.go#L126]
6. [story/client/x/dkg/keeper/dkg_handler.go#L103]
7. [story/client/x/dkg/keeper/dkg_rewards.go#L18]

### STOR-21 — Resharing mutates the serving-round pointer into a non-active round and opens a multi-day read-fee capture window

### Resharing mutates the serving-round pointer into a non-active round and opens a multi-day read-fee capture window

### Executive Summary

When a DKG round’s active period expires, `BeginBlocker()` rewrites that very round from `Active` to `Registration` before it starts the successor resharing round. The `LatestActiveRound` pointer is not updated until the successor later finalizes, so all `CDR.read()` requests during resharing are still keyed to the stale old round even though that round is no longer active. `ProcessCDRVaultRead()` accepts the read and `ThresholdDecryptRequested()` records the paid request in consensus state first, but then silently skips the honest validator auto-queue because the round stage is not `Active`. With the default parameters this mismatch lasts roughly three days, while decrypt requests expire after only 200 blocks. Honest default operators therefore never service these reads before timeout. A malicious validator can still use the manual `PartialDecryptTDH2` path against the stale round, become the only credited submitter during the resharing window, and capture the read-fee-backed CDR reward flow from users transacting in that period.

### Details

The stale-pointer window starts in `BeginBlocker()`. When the active period elapses, the keeper mutates the current round in place and immediately starts a new resharing round:

```go
nextStage, shouldTransition := k.shouldTransitionStage(currentHeight, latestRound, params)
if shouldTransition {
    latestRound.Stage = nextStage
    if err := k.setDKGNetwork(ctx, latestRound); err != nil {
        return err
    }
    switch nextStage {
    case types.DKGStageRegistration:
        return k.InitiateDKGRound(ctx, false)
    }
}
```

`LatestActiveRound` is only a pointer to a stored round number, so after that mutation it still resolves to the same round object, now with `Stage = Registration`:

```go
func (k *Keeper) getLatestActiveDKGNetwork(ctx context.Context) (*types.DKGNetwork, error) {
    key, err := k.LatestActiveRound.Get(ctx)
    ...
    dkgNetwork, err := k.DKGNetworks.Get(ctx, key)
    ...
    return &dkgNetwork, nil
}
```

`ProcessCDRVaultRead()` trusts this stale pointer as the serving committee and forwards the user-paid read into `ThresholdDecryptRequested()`:

```go
latestRound, err := k.dkgKeeper.GetLatestActiveRound(cachedCtx)
...
round := latestRound.Round
...
k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, round, ev.RequesterPubKey, ev.Ciphertext, label[:], uint64(sdkCtx.BlockHeight()))
```

`ThresholdDecryptRequested()` stores the request first, then notices the round is no longer active and skips the honest queueing path entirely:

```go
if err := k.setDecryptRequest(ctx, requesterPubKey, label, types.DecryptRequest{ ... }); err != nil {
    return errors.Wrap(err, "failed to register decrypt request")
}
...
if dkgNetwork.Stage != types.DKGStageActive {
    log.Info(ctx, "Skipping threshold decrypt request; DKG round is not active", ...)
    return nil
}
```

The timeout mismatch makes this exploitable at scale. Requests expire after 200 blocks, but the default resharing stages are each 86,400 blocks (`Registration`, `Dealing`, `Finalization`). The old round also never gets marked `Ended` in the normal cleanup path because `endPreviousActiveRound()` only changes the stage when `prevActive.Stage == Active`, which is already false by the time the successor finalizes.

A malicious validator can still monetize the window because the manual TEE path only checks that the requested round matches the stale `LatestActiveRound` pointer, not that the pointed round is actually active:

```go
latestNetwork, err := s.verifyRoundMatchesLatestNetwork(ctx, req.GetRound())
if err != nil {
    return nil, status.Errorf(codes.FailedPrecondition, "round does not match latest active network")
}
```

On the CL side, `PartialDecryptionSubmitted()` never checks the round stage before accepting, refunding, and crediting the submission:

```go
req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, round, ciphertext)
...
if currentHeight-req.Height > types.PartialDecryptionTimeoutBlocks {
    ...
    return nil
}
...
if err := k.setPartialDecryptionSubmission(...); err != nil {
    ...
}
```

Combined, these paths let a custom validator bot service stale-round reads that the honest default service never even queues.

### Impact Cascade

- Every `CDR.read()` during resharing is accepted and charged even though the default honest service path drops it immediately after registry insertion.
- The request expires after 200 blocks, long before the successor round can finish its three-stage resharing process, so honest automation cannot recover the paid read.
- A validator running a custom manual submitter can still produce stale-round partials during this window and become the only credited submitter for those reads.
- Because `ProcessDKGPartialDecryptionSubmitted()` increments `CDRPartialSubmitCount` and refunds successful submissions, the attacker captures the read-fee-backed reward distribution generated by users who transact during every resharing window.

### Assumptions and Uncertainties

1. At least one validator is willing to use a custom manual submission path instead of the default auto-queue.
2. Users continue calling `CDR.read()` during resharing windows because the EL contract does not expose the CL-side stage mismatch.
3. The network uses the default timing scale, or any configuration where resharing duration greatly exceeds `PartialDecryptionTimeoutBlocks`.

### How will the bug recipient respond?

“Reads are only expected during active rounds, and validators can always submit manually if needed.”

That is precisely the vulnerability. The EL still accepts and charges reads during non-active rounds because `ProcessCDRVaultRead()` trusts a stale serving-round pointer, while the honest default service path silently refuses to queue them. Manual submission is not a mitigation; it gives custom validators an exclusive fee-capture channel that default honest operators do not have.

### Why did tests miss this issue?

The existing registry and timeout tests exercise CRUD and expiry in isolation, but they do not simulate the lifecycle boundary where the current active round is rewritten to `Registration` while `LatestActiveRound` still points at it. There is no integration coverage for `CDR.read()` during resharing or for the mismatch between the three resharing stages and the 200-block decrypt-request timeout.

### Recommendation

Keep the serving committee pointer bound only to a truly `Active` round. Two safe patterns are:
1. Do not mutate the current serving round out of `Active` until the successor round has finalized.
2. If a separate lifecycle tracker is required, decouple it from `LatestActiveRound` so `GetLatestActiveRound()` never resolves to a non-active network.

In addition:
- Reject `CDR.read()` processing unless the pointed round is truly `Active` before storing the decrypt request.
- Add a round-stage check to `PartialDecryptionSubmitted()` so stale-round manual submissions cannot be accepted or rewarded.
- Align decrypt-request expiry with resharing duration, or explicitly suspend EL read acceptance whenever no active committee is serving requests.

### References

1. [story/client/x/dkg/keeper/abci.go#L89-L109]
2. [story/client/x/dkg/keeper/dkg_network.go#L179-L236]
3. [story/client/x/evmengine/keeper/cdr.go#L57-L99]
4. [story/client/x/dkg/keeper/dkg_handler.go#L414-L476]
5. [story/client/x/dkg/keeper/dkg_handler.go#L491-L599]
6. [story/client/x/dkg/types/keys.go#L16-L18]
7. [story/client/x/dkg/types/params.go#L10-L15]
8. [story-kernel/service/dkg_partial_decrypt.go#L40-L60]
9. [story/docs/design/DKG.md#L633-L640]
10. [story/docs/design/CDR.md#L73-L76]

### STOR-20 — PartialDecryptTDH2 protobuf mismatch makes every CL-to-kernel decrypt RPC invalid on the wire

### PartialDecryptTDH2 protobuf mismatch makes every CL-to-kernel decrypt RPC invalid on the wire

### Executive Summary

The consensus client and `story-kernel` compile different protobuf schemas for the same gRPC method `KernelService.PartialDecryptTDH2`. The CL request assigns field 5 to `pid`, field 6 to `global_pub_key`, and field 8 to `requester_pub_key`, while the kernel request expects field 5 to be `global_pub_key` bytes and field 6 to be `requester_pub_key` bytes, with no `pid` field at all. The CL always populates `Pid` before calling the kernel. On the wire that emits tag `0x28` (field 5, varint), but the kernel decodes field 5 as a length-delimited bytes field. That is a protobuf wire-type mismatch, so the request cannot be decoded as the kernel’s request message. As a result, CL-to-kernel partial decryption fails before the handler logic runs, independently of the already-reported PIDCache issue. Users can pay `CDR.read()` fees even in round 1 while every decrypt RPC is malformed.

### Details

The CL protobuf for `PartialDecryptTDH2Request` is:

```protobuf
bytes code_commitment = 1;
uint32 round = 2;
bytes ciphertext = 3;
bytes label = 4;
uint32 pid = 5;
bytes global_pub_key = 6;
string sealed_share_id = 7;
bytes requester_pub_key = 8;
```

The kernel protobuf for the same RPC method is:

```protobuf
bytes code_commitment = 1;
uint32 round = 2;
bytes ciphertext = 3;
bytes label = 4;
bytes global_pub_key = 5;
bytes requester_pub_key = 6;
```

The CL call site always sets `Pid`, `GlobalPubKey`, and `RequesterPubKey`:

```go
client.PartialDecryptTDH2(ctx, &types.PartialDecryptTDH2Request{
    CodeCommitment:  session.CodeCommitment,
    Round:           session.Round,
    Ciphertext:      req.Ciphertext,
    Label:           req.Label,
    Pid:             pid,
    GlobalPubKey:    session.GlobalPubKey,
    RequesterPubKey: req.RequesterPubKey,
})
```

The generated CL marshal code emits:
- field 5 / wire type varint for `Pid` (`0x28`),
- field 6 / wire type bytes for `GlobalPubKey` (`0x32`),
- field 8 / wire type bytes for `RequesterPubKey`.

But the kernel-generated request type defines:
- field 5 as `bytes global_pub_key`,
- field 6 as `bytes requester_pub_key`,
- no field 8 at all.

So the first nontrivial partial decrypt request sent by the CL is malformed from the kernel’s perspective:
- field 5 arrives with wire type `varint`, while the kernel expects `bytes` for `global_pub_key`;
- field 6 is interpreted as `requester_pub_key`, not `global_pub_key`;
- the actual requester key in field 8 is discarded as unknown.

This is not a semantic bug inside the handler — it is a schema incompatibility at the RPC boundary. The kernel cannot reliably reconstruct the request the CL thinks it sent.

### Impact Cascade

- Every `handleDecryptRequest()` call sends a malformed protobuf for `PartialDecryptTDH2()`.
- Threshold decryptions fail even before considering PIDCache or other handler-side logic.
- `CDR.read()` still burns the read fee before any validator attempts the malformed decrypt RPC.
- The confidential-read feature is therefore broken from the CL↔︎TEE boundary itself, not only after resharing.

### Assumptions and Uncertainties

1. Validators use the shipped `story` client protobufs and the shipped `story-kernel` protobufs, which is the normal deployment model.
2. No custom compatibility shim rewrites the request message in transit.

### How will the bug recipient respond?

“The kernel ignores unknown fields, so this is harmless.”

Unknown-field tolerance does not save a known-field wire-type conflict. The CL sends field 5 as a varint because it believes field 5 is `pid`; the kernel declares field 5 as a bytes field (`global_pub_key`). That is an incompatible wire type for a known field, not a benign extra field.

### Why did tests miss this issue?

The test suites are split across repos and validate each side against its own generated types. There is no compatibility test that serializes the CL request type and deserializes it with the kernel request type for `PartialDecryptTDH2`.

### Recommendation

Unify the protobuf source of truth for `KernelService` across the CL and `story-kernel`, regenerate both clients/servers from the same schema, and add a cross-repo compatibility test for every RPC. For `PartialDecryptTDH2`, field numbers must be aligned before the method can be considered safe to use.

### References

1. [story/client/proto/story/dkg/v1/types/kernel.proto#L231]
2. [story-kernel/proto/tee.proto#L148]
3. [story/client/x/dkg/keeper/dkg_svc.go#L319]
4. [story/client/x/dkg/types/kernel.pb.go#L785]
5. [story-kernel/types/pb/v0/tee.pb.go#L1225]
6. [story/contracts/src/protocol/CDR.sol#L203]

### STOR-19 — The first dealing block rejects the last registration block’s DKG register events, shrinking the committee before rewards are split

### The first dealing block rejects the last registration block’s DKG register events, shrinking the committee before rewards are split

### Executive Summary

DKG registrations are emitted on the EVM side and imported into CL one block later through `MsgExecutionPayload`, but `BeginBlocker` advances the round from `Registration` to `Dealing` before those previous-block registration logs are processed. `dkgKeeper.Registered` then refuses to accept any imported registration unless the current round stage is still `DKGStageRegistration`. As a result, every registration submitted in the final allowed EL registration block is deterministically discarded.

This is economically exploitable because the dealing transition computes the round’s `Total` and `Threshold` from the already-imported registration set only. A malicious validator can register early while honest validators who land in the last registration block are excluded by the bridge timing bug, producing a smaller committee and a lower threshold. The attacker then competes for the same UBI/CDR reward pool with fewer peers.

### Details

`BeginBlocker` moves a round out of registration as soon as the registration window elapses and immediately starts dealing using the registrations currently stored in CL:

```go
case types.DKGStageRegistration:
    if elapsed >= registrationEnd {
        return types.DKGStageDealing, true
    }
}
```

```go
case types.DKGStageDealing:
    return k.BeginDealing(ctx, latestRound)
```

`BeginDealing` snapshots the committee size and threshold from the CL registrations that exist at that moment:

```go
verifiedRegCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusVerified)
...
latestRound.Total = verifiedRegCount
latestRound.Threshold = types.CalculateThreshold(verifiedRegCount, params.OperationalThreshold)
```

But the previous EL block’s `DKG.Registered` logs are only imported later when `MsgExecutionPayload` executes:

```go
if err := s.ProcessDKGEvents(ctx, payload.Number-1, ethLogs); err != nil {
    return nil, errors.Wrap(err, "deliver dkg-related event logs")
}
```

`dkgKeeper.Registered` rejects those imported events unless the round is still in registration:

```go
if latest.Stage != types.DKGStageRegistration {
    return errors.New("round is not in registration stage")
}
```

Therefore, on the first dealing block:
1. `BeginBlocker` advances the round into `DKGStageDealing` and snapshots `Total`/`Threshold`.
2. The block proposal’s `MsgExecutionPayload` still carries the final registration block’s `DKG.Registered` logs.
3. `ProcessDKGRegistered` forwards them to `dkgKeeper.Registered`.
4. `dkgKeeper.Registered` rejects them because the stage is already `Dealing`.
5. Honest validators whose registration landed in the last EL registration block never join the round.

This smaller committee directly affects economic outcomes because finalized committee rewards are divided by the set that survives registration and finalization:

```go
finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round.Round, types.DKGRegStatusFinalized)
...
perMemberReward := dkgReward.Quo(memberCount)
```

```go
if err := k.distributeCDRRewardPool(ctx); err != nil { ... }
```

### Impact Cascade

- Deterministically discards valid registrations from the final allowed EL registration block.
- Shrinks the dealing/finalization committee and lowers the effective threshold.
- Increases each surviving attacker’s share of UBI and CDR fee-pool rewards.
- Makes committee composition depend on bridge timing instead of the advertised registration window.

### Assumptions and Uncertainties

1. Honest validators may register near the end of the registration window, which is realistic because registration depends on off-chain TEE setup and transaction inclusion timing.
2. The attacker registers earlier and benefits from fewer competing committee members.
3. No compensating replay path re-applies the dropped registration events after dealing begins.

### How will the bug recipient respond?

“Those registrations were too late.”

They were not too late for the protocol’s own EL window: they were emitted in the final permitted registration block and delivered through the standard `MsgExecutionPayload` bridge. The CL stage machine simply closes registration one step before it processes the previous EL block’s canonical logs.

### Why did tests miss this issue?

The lifecycle tests advance the round stage directly from CL state, but they do not simulate `DKG.Registered` events arriving in the last EL registration block and then being imported one block later after the round has already moved into dealing.

### Recommendation

Do not begin dealing until the previous EL block’s registration logs have been consumed, or explicitly accept registrations sourced from the immediately preceding registration block while processing the first dealing block.

### References

1. [story/client/x/dkg/keeper/dkg_round.go#L14]
2. [story/client/x/dkg/keeper/abci.go#L89]
3. [story/client/x/dkg/keeper/dkg_dealing.go#L15]
4. [story/client/x/evmengine/keeper/msg_server.go#L170]
5. [story/client/x/evmengine/keeper/dkg.go#L72]
6. [story/client/x/dkg/keeper/dkg_handler.go#L22]
7. [story/client/x/dkg/keeper/dkg_rewards.go#L58]

### STOR-18 — VerifiedQueryClient accepts proofs for arbitrary keys and stores, letting hostile RPCs feed fake DKG state

### VerifiedQueryClient accepts proofs for arbitrary keys and stores, letting hostile RPCs feed fake DKG state

### Executive Summary

The kernel’s “verified” query path does not actually bind an ICS23 proof to the key or store that the enclave asked for. `queryWithProof()` accepts any non-empty `resp.Key`, but never checks it equals the requested key, and `verifyProof()`/`verifyMultiStoreProof()` then trust the proof-operation keys supplied by the RPC response instead of the caller’s query key or store name. A hostile RPC can therefore answer `GetDKGNetwork`, `getDKGRegistration`, or `GetLatestActiveDKGNetwork` with a proof for some other key and the enclave will still treat the bytes as authenticated Story state.

This breaks the core host-boundary invariant described in the whitepaper: the enclave is supposed to consume only canonical, proof-verified DKG state. In practice, a malicious host/RPC can suppress real registrations with arbitrary non-membership proofs, swap in historical `DKGNetwork` records for a different round, or redirect the “latest active round” lookup to an unrelated numeric-string value. The resulting fake round context is then used for key generation, deal processing, finalization, and decryption gating.

### Details

`getStoreData()` forwards the caller’s requested key to `queryWithProof()`, but the response is accepted as long as `resp.Key` is merely non-empty.

```go
func (q *VerifiedQueryClient) queryWithProof(ctx context.Context, storeKey string, key []byte, height int64) (*ctypes.ResultABCIQuery, error) {
    ...
    if len(resp.Key) == 0 {
        return nil, errors.New("empty key in response")
    }
    ...
}
```

`verifyProof()` then parses the proof ops and passes the proof-supplied keys (`op.Key`) into `verifyMultiStoreProof()`.

```go
if err := verifyMultiStoreProof(
    proofs[moduleProofIdx], proofs[simpleProofIdx],
    proofTypes[moduleProofIdx], proofTypes[simpleProofIdx],
    keys[moduleProofIdx], keys[simpleProofIdx],
    key, value, appHash,
); err != nil {
    ...
}
```

Inside `verifyMultiStoreProof()`, those proof-supplied keys become authoritative. For non-membership, `verifyModuleNonExistenceProof()` proves absence of `moduleKey` when present; for membership, `verifyModuleProof()` proves membership of `moduleKey`; and `verifySimpleProof()` proves whatever `simpleKey` the RPC supplied. None of these are compared against the caller’s expected module key or the requested store key.

```go
proofKey := moduleKey
if len(proofKey) == 0 {
    proofKey = queryKey
}
...
if !ics23.VerifyNonMembership(spec, expectedModuleRoot, moduleProof, proofKey) {
    return errors.New("module non-existence proof verification failed")
}
```

```go
if err := verifySimpleProof(simpleProof, simpleType, simpleKey, expectedModuleRoot, appHash); err != nil {
    return err
}
```

The downstream consumers assume these queries are trustworthy. `GetLatestActiveDKGNetwork()` parses whatever bytes came back as a decimal round number and immediately calls `GetDKGNetwork()` for that round. `PartialDecryptTDH2()` uses this path in `verifyRoundMatchesLatestNetwork()`, and `GetOrLoadRoundContext()` uses the queried network/registrations to build the DKG state machine. `validateRegistrations()` only checks count, duplicate indices, and `reg.Round`; it never verifies that each returned registration actually corresponds to the queried validator or belongs to `network.ActiveValSet`.

A hostile RPC can therefore execute the following concrete attacks:
1. For a real validator registration query, return a valid non-membership proof for any unrelated absent key. The kernel interprets the target validator as “not registered”.
2. For `GetDKGNetwork(round=X)`, return a proof for some other `DKGNetwork` record. The decoded protobuf is trusted without checking `network.Round == X`.
3. For `GetLatestActiveDKGNetwork()`, return a proof for any numeric-string value in state, steering the kernel to the attacker-chosen round.

Once that forged state is accepted, the enclave proceeds to keygen/finalization/decryption decisions under attacker-selected committee data.

### Impact Cascade

- A hostile RPC can make the enclave treat arbitrary registration absence as authentic, removing honest participants from its internal committee view.
- A hostile RPC can swap in stale `DKGNetwork` data for the wrong round, defeating the quote-to-round/start-block trust chain.
- `PartialDecryptTDH2()` can be gated on a forged “latest active round”, allowing decryption decisions against stale or fake committee state.
- The enclave’s most important security claim — “only canonical Story state enters the TEE” — no longer holds.

### Assumptions and Uncertainties

1. The attacker controls, or can fully spoof, the RPC endpoint(s) feeding the enclave. This is explicitly the review threat model for this focus area.
2. The attacker does not need to break CometBFT or ICS23; they only need to replay valid proofs for different keys because the enclave never binds proofs to the requested key/store.
3. Some downstream CL handlers reject certain malformed outputs, but the enclave still acts on fake state before those later checks.

### How will the bug recipient respond?

“The proof still verifies against a real AppHash, so this is only an RPC integrity issue, not an enclave bug.”

That response misses the root cause: the enclave is supposed to verify *the queried key*, not merely *some key*. Returning a valid proof for a different key is enough to violate the enclave’s trust boundary, and several production paths (`GetLatestActiveDKGNetwork`, registration loading, round-context construction) consume the forged value immediately.

### Why did tests miss this issue?

Existing tests exercise successful proof verification and malformed-proof rejection, but they do not adversarially swap `resp.Key`, `moduleKey`, or `simpleKey` to prove that a response for one key can satisfy a query for another key.

### Recommendation

Bind the proof chain to the caller’s expectation at every layer:
- Reject any ABCI response where `resp.Key` does not exactly equal the requested key.
- Pass the expected store key into `verifyProof()` and require the multistore proof to prove that exact store.
- Require `moduleKey` (if present) to equal the requested key; otherwise reject.
- After decoding, validate semantic invariants such as `network.Round == requestedRound` and `registration.ValidatorAddr == queriedValidatorAddr`.

A minimal shape is:

```go
if !bytes.Equal(resp.Key, key) {
    return nil, fmt.Errorf("response key mismatch")
}
if len(moduleKey) > 0 && !bytes.Equal(moduleKey, queryKey) {
    return fmt.Errorf("module proof key mismatch")
}
if !bytes.Equal(simpleKey, []byte(expectedStoreKey)) {
    return fmt.Errorf("multistore proof key mismatch")
}
```

### References

1. [story-kernel/story/query_client.go#L303-L348]
2. [story-kernel/story/query_client.go#L357-L449]
3. [story-kernel/story/query_client.go#L494-L620]
4. [story-kernel/service/round_context.go#L78-L125]
5. [story-kernel/service/dkg_partial_decrypt.go#L192-L205]

### STOR-16 — Expired DKG rounds keep collecting UBI and blocking validator exits because `LatestActiveRound` is never cleared at rollover

### Expired DKG rounds keep collecting UBI and blocking validator exits because `LatestActiveRound` is never cleared at rollover

### Executive Summary

When a DKG round reaches the end of its configured `ActivePeriod`, `BeginBlocker` immediately mutates that round’s stored stage to `Registration` and starts the next resharing round. However, the module does not clear or advance `LatestActiveRound` until some later round successfully completes `FinalizeDKGRound`. Every reward-distribution and self-unstake guard in the rest of the system trusts `LatestActiveRound` alone, not the round’s actual stage. The result is that an already-expired committee remains the authoritative “active” committee for accounting purposes for as long as resharing can be kept incomplete. A malicious committee coalition can therefore keep an obsolete round collecting the DKG share of UBI withdrawals indefinitely, while honest former committee members remain unable to self-unstake because the chain still treats them as members of the active committee.

### Details

`BeginBlocker` computes the next stage from elapsed height, writes that stage back into the current round, and then starts a new round when `Active -> Registration` is reached:

```go
nextStage, shouldTransition := k.shouldTransitionStage(currentHeight, latestRound, params)
if shouldTransition {
    latestRound.Stage = nextStage
    if err := k.setDKGNetwork(ctx, latestRound); err != nil {
        return err
    }

    switch nextStage {
    case types.DKGStageRegistration:
        return k.InitiateDKGRound(ctx, false)
    }
}
```

That transition does **not** touch `LatestActiveRound`; only `FinalizeDKGRound` updates the pointer after a new round succeeds:

```go
if err := k.setLatestActiveRound(ctx, latestRound); err != nil {
    return errors.Wrap(err, "failed to set the latest active round of DKG")
}
```

The module even documents that previous rounds are only ended if their stored stage is still `Active`:

```go
if prevActive.Stage == types.DKGStageActive {
    prevActive.Stage = types.DKGStageEnded
    if err := k.setDKGNetwork(ctx, prevActive); err != nil {
        return errors.Wrap(err, "failed to set previous active round to ended")
    }
}
```

But by the time the next round finalizes, the previous active round has already been overwritten to `Registration`, so `endPreviousActiveRound` does nothing. The stale pointer then drives downstream economic logic. UBI rewards are sent to whatever `LatestActiveRound` returns, with no stage check:

```go
activeRound, err := k.getLatestActiveDKGNetwork(ctx)
if activeRound == nil {
    return math.ZeroInt(), nil
}

return k.distributeRewardsFromModule(ctx, activeRound, senderModule, totalAmount)
```

Self-unstaking is also blocked solely by finalized membership in `LatestActiveRound`:

```go
activeRound, err := k.dkgKeeper.GetLatestActiveRound(cachedCtx)
if err == nil && activeRound != nil {
    hasFinalized, err := k.dkgKeeper.HasFinalizedRegistration(cachedCtx, activeRound.Round, valEvmAddr)
    if err == nil && hasFinalized {
        return errors.WrapErrWithCode(errors.ActiveDKGMemberSelfUnstake, ...)
    }
}
```

An attacker who controls enough of the current committee to keep subsequent resharing rounds from finalizing never needs any privileged role. They simply let `ActivePeriod` expire, keep the next round(s) from reaching `FinalizeDKGRound`, and the chain will continue treating the old committee as active forever for reward and exit checks.

### Impact Cascade

- The stale committee continues receiving the configured DKG share of every UBI withdrawal even after its active period has already ended.
- Honest former committee members cannot self-unstake, because `withdraw.go` still treats the obsolete round as the active committee.
- A malicious committee can monetize failed resharing by farming UBI distributions from `UBIPool` reserves without ever allowing authority to rotate.
- Because no timeout or “clear active round” fallback exists, this stale-accounting state persists until a later DKG round successfully finalizes.

### Assumptions and Uncertainties

1. The attacker controls enough members of the currently active committee to keep follow-on rounds from finalizing.
2. UBI withdrawals continue while resharing is stalled, so `ProcessUbiWithdrawal` keeps invoking DKG reward distribution.
3. The stale committee members have finalized registrations in the expired round, which is exactly the membership set rewarded and exit-blocked by the code.

### How will the bug recipient respond?

“This is intended: the old committee stays active until a new one finalizes so service remains continuous.”

That explanation does not match the implementation. The old round is explicitly moved out of `Active` into `Registration`, yet rewards and self-unstake checks continue to treat it as active because they trust `LatestActiveRound` only. If the design intended continuous authority, the stage should remain active or there should be an explicit fallback state; instead, the code creates a split-brain state where the round is expired for lifecycle purposes but still active for money and validator lockups.

### Why did tests miss this issue?

The lifecycle test suite explicitly observes the bad state (`Round 1 stage was overwritten to Registration during Active→Registration transition`) but only asserts the stored value instead of checking downstream consumers such as reward distribution or self-unstake guards. The tests therefore codify the state drift without validating its economic consequences.

### Recommendation

Make rollover atomic for lifecycle and accounting state.

Primary fix:

```go
case types.DKGStageRegistration:
    if err := k.clearLatestActiveRound(ctx); err != nil {
        return err
    }
    latestRound.Stage = types.DKGStageEnded
    if err := k.setDKGNetwork(ctx, latestRound); err != nil {
        return err
    }
    return k.InitiateDKGRound(ctx, false)
```

Also harden all downstream consumers:
- `DistributeRewardsToActiveCommittee` should require `activeRound.Stage == DKGStageActive`.
- The self-unstake guard should reject only if the returned round is still actually active.
- If continuous service from the old committee is desired, represent that with an explicit state instead of leaving `LatestActiveRound` pointing at an expired round.

### References

1. [story/client/x/dkg/keeper/abci.go#L81-L109]
2. [story/client/x/dkg/keeper/dkg_network.go#L188-L236]
3. [story/client/x/dkg/keeper/dkg_finalization.go#L78-L86]
4. [story/client/x/dkg/keeper/dkg_rewards.go#L18-L40]
5. [story/client/x/evmstaking/keeper/withdraw.go#L532-L542]

### STOR-13 — Finalization stores raw pubKeyShare while partial decryptions submit prefixed pubShare, so honest reads cannot be completed

### Finalization stores raw pubKeyShare while partial decryptions submit prefixed pubShare, so honest reads cannot be completed

### Executive Summary

The stock Story DKG finalization path and the stock kernel partial-decryption path serialize the validator public share in two incompatible byte formats. During finalization, the kernel returns `pubKeyShare` as raw Edwards25519 point bytes, and the consensus layer stores that value verbatim in `DKGRegistration.PubKeyShare`. During partial decryption, the same kernel derives `pubShare` from the same scalar but prepends the TDH2 prefix bytes `[0x04, 0x3f]` before returning it. The CL submission handler then enforces exact byte equality between submitted `pubShare` and stored `PubKeyShare`. Because no normalization occurs in the contract call, event bridge, or keeper logic, an honest validator following the production flow will always hit `pubShare mismatch` and have its partial rejected. Users are still charged `readFee` on `CDR.read`, but legitimate threshold decryptions cannot complete through the shipped path.

### Details

The finalization path stores the kernel’s `pubKeyShare` output unchanged:

```go
// story-kernel/service/dkg_finalize.go
pubKeyShare, err := s.Suite.Point().Mul(priShare.V, nil).MarshalBinary()
```

```go
// story/client/x/dkg/keeper/dkg_registration.go
dkgReg.PubKeyShare = pubKeyShare
dkgReg.Status = types.DKGRegStatusFinalized
```

That value is passed through unchanged by the client and contract/event bridge:

```go
// story/client/x/dkg/keeper/dkg_svc_finalization.go
session.PubKeyShare = resp.GetPubKeyShare()
```

```solidity
// story/contracts/src/protocol/DKG.sol
emit Finalized(..., pubKeyShare, signature);
```

During partial decryption, the kernel computes the same share point but intentionally prepends a TDH2 prefix so `TDH2Combine` can deserialize it:

```go
// story-kernel/service/dkg_partial_decrypt.go
func marshalPubShare(scalar kyber.Scalar) ([]byte, error) {
    pubSharePoint := suite.Point().Mul(scalar, nil)
    pointBz, err := pubSharePoint.MarshalBinary()
    ...
    return append([]byte{sec1UncompressedPrefix, tdh2Edwards25519CurveID}, pointBz...), nil
}
```

The kernel test suite explicitly asserts that the prefix is present:

```go
// story-kernel/service/dkg_partial_decrypt_test.go
assert.Equal(t, byte(sec1UncompressedPrefix), result[0])
assert.Equal(t, byte(tdh2Edwards25519CurveID), result[1])
```

The CL partial-submission verifier rejects unless the bytes match exactly:

```go
// story/client/x/dkg/keeper/dkg_handler.go
if !bytes.Equal(pubShare, reg.PubKeyShare) {
    return errors.New("pubShare mismatch: submitted pubShare does not match stored pubKeyShare")
}
```

The production worker forwards the kernel’s prefixed `resp.PubShare` directly into the contract call; it does not strip or normalize the prefix first. Therefore a validator using the shipped `handleDecryptRequest -> PartialDecryptTDH2 -> SubmitEncryptedPartialDecryption -> ProcessDKGPartialDecryptionSubmitted` flow can never satisfy the equality check.

### Impact Cascade

- Every honest kernel-generated partial decryption is rejected before storage, so legitimate requesters cannot gather an accepted threshold set.
- `CDR.read` still charges the user `readFee` before any decryption work is proven possible, so users can pay for a service that the shipped implementation cannot complete.
- Honest validators never receive fee refunds or reward credit for real work because their submissions fail at `pubShare mismatch`.
- The only submissions that can accrue CDR reward weight are non-standard or already-invalid paths, worsening the previously validated reward-accounting bugs.
- A protocol-wide read outage can persist until one side of the serialization mismatch is patched and all validators upgrade.

### Assumptions and Uncertainties

1. Validators are running the published `story-kernel` implementation and the published `story/client` submission flow without an out-of-band byte rewrite.
2. No hidden normalization exists in off-repo requester tooling or validator wrappers; I did not find any normalization in the audited repositories.
3. The intended wire format for stored `PubKeyShare` is the same value later compared against submitted `pubShare`, as enforced by `bytes.Equal` in the keeper.

### How will the bug recipient respond?

“The formats are intentionally different; operators can translate between them before submitting.”

That is not what the shipped code does. The stock worker forwards `resp.PubShare` exactly as returned by the kernel, and the stock verifier compares it byte-for-byte against the raw `PubKeyShare` stored during finalization. If a translation layer were intended, it is missing from the production path.

### Why did tests miss this issue?

The current tests cover the two sides in isolation: kernel tests verify that `marshalPubShare` adds the `[0x04, 0x3f]` prefix, and DKG keeper tests exercise `Finalized` / signature helpers with synthetic byte slices. I did not find an end-to-end test that finalizes a registration using the real kernel output and then feeds a real kernel partial-decryption response into `PartialDecryptionSubmitted`.

### Recommendation

Normalize the representation on one side and enforce that representation everywhere.

Preferred fix:

```go
// Store the same prefixed format that partial decryptions submit.
prefixedPubShare, err := marshalPubShare(priShare.V)
if err != nil { ... }
pubKeyShare = prefixedPubShare
```

Alternative fix:
- Keep storing raw point bytes, but strip the `[0x04, 0x3f]` prefix from `resp.PubShare` before contract submission and before any equality checks.
- Add an integration test covering `FinalizeDKG -> PartialDecryptTDH2 -> SubmitEncryptedPartialDecryption -> PartialDecryptionSubmitted` with the real kernel serialization.

### References

1. [story-kernel/service/dkg_finalize.go#L92]
2. [story/client/x/dkg/keeper/dkg_registration.go#L94]
3. [story/client/x/dkg/keeper/dkg_svc_finalization.go#L121]
4. [story/contracts/src/protocol/DKG.sol#L229]
5. [story-kernel/service/dkg_partial_decrypt.go#L253]
6. [story-kernel/service/dkg_partial_decrypt_test.go#L15]
7. [story/client/x/dkg/keeper/dkg_handler.go#L554]
8. [story/client/x/dkg/keeper/dkg_svc.go#L319]

### STOR-11 — Reusing one TEE keypair across validator addresses can finalize an undecryptable DKG committee

### Reusing one TEE keypair across validator addresses can finalize an undecryptable DKG committee

### Executive Summary

The DKG stack lets a multi-validator operator reuse the exact same TEE long-term keys for multiple validator addresses in the same round, and the CL then counts those duplicate identities as independent finalized participants. `story-kernel` seals both the Ed25519 DKG key and secp256k1 communication key only by `(round, codeCommitment)`, not by validator address, so a node that is reconfigured for another validator in the same round will reload the same keys and produce a fresh attestation quote for the new address. The on-chain registration path accepts duplicate `DkgPubKey` / `CommPubKey` values, and the finalization signature does not bind `validatorAddr`; `DKG.sol.finalize()` also does not require `validatorAddr == msg.sender`. As a result, the same finalization output and signature can be replayed across every duplicated registration, causing `AddGlobalPubKeyVote` and `FinalizeDKGRound` to reach threshold on registration count even though the committee has fewer unique key shares than the configured threshold. The round becomes active and starts enforcing committee lockups/reward flows, but threshold decryptions are no longer guaranteed to work.

### Details

`story-kernel` persists round keys without the validator address in the key path. Both long-term key loaders use only `round` and `codeCommitmentHex`:

```go
func (s *DKGStore) ed25519Path(codeCommitmentHex string, round uint32) string {
    return filepath.Join(s.keyDir, strconv.FormatUint(uint64(round), 10), codeCommitmentHex, KeyEd25519File)
}

func (s *DKGStore) secp256k1Path(codeCommitmentHex string, round uint32) string {
    return filepath.Join(s.keyDir, strconv.FormatUint(uint64(round), 10), codeCommitmentHex, KeySecp256k1File)
}
```

`GenerateAndSealKey` then reuses those keys while recomputing report data with whatever validator address is supplied in the request:

```go
_, edPub, err := s.DKGStore.LoadOrGenerateEd25519Key(codeCommitmentHex, req.GetRound())
_, secpPub, err := s.DKGStore.LoadOrGenerateSecp256k1Key(codeCommitmentHex, req.GetRound())

reportData, err := calculateReportData(
    req.Address,
    req.Round,
    edPubBz,
    ecrypto.FromECDSAPub(secpPub)[1:],
    network.StartBlockHeight,
    network.StartBlockHash,
)
```

So an operator controlling validators `A` and `B` can reuse one enclave/state directory, request `GenerateAndSealKey` once as `A`, then request it again as `B`; the quote for `B` is valid because the report binds `B`, but the underlying `dkgPubKey` and `commPubKey` stay identical.

The CL registration path only enforces uniqueness per validator address, not per key material:

```go
exists, err := k.hasDKGRegistration(ctx, round, validator)
if exists {
    return errors.New("validator already registered for this round")
}

dkgReg := &types.DKGRegistration{
    Round:         round,
    ValidatorAddr: validator.Hex(),
    Index:         uint32(index),
    DkgPubKey:     dkgPubKey,
    CommPubKey:    commPubKey,
}
```

Finalization then compounds the problem. The signed message omits the validator address entirely:

```go
// codeCommitment || round || participantsRoot || globalPubKey || publicCoeffs... || pubKeyShare
encoded = append(encoded, codeCommitment[:]...)
binary.BigEndian.PutUint32(roundBytes, round)
encoded = append(encoded, roundBytes...)
encoded = append(encoded, participantsRoot[:]...)
encoded = append(encoded, globalPubKey...)
for _, coeff := range publicCoeffs {
    encoded = append(encoded, coeff...)
}
encoded = append(encoded, pubKeyShare...)
```

And the EL contract does not bind the event’s `validatorAddr` to `msg.sender`:

```solidity
function finalize(
    uint32 round,
    address validatorAddr,
    ...
) external payable chargesFee whenNotPaused {
    require(validatorAddr != address(0), "DKG: Validator address cannot be empty");
    emit Finalized(round, validatorAddr, enclaveType, ...);
}
```

CL-side verification recovers the signer from the signature and compares it only to `keccak256(commPubKey)` loaded from the registration selected by `validatorAddr`:

```go
reg, err := k.getDKGRegistration(ctx, round, msgSender)
if err := verifyFinalizationSignature(reg.CommPubKey, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature); err != nil {
    return errors.Wrap(err, "finalization signature verification failed")
}
voteCount, err := k.AddGlobalPubKeyVote(ctx, round, globalPubKey, publicCoeffs)
if err := k.finalizeDKGRegistration(ctx, round, msgSender, pubKeyShare); err != nil {
    return errors.Wrap(err, "failed to update dkg registration status")
}
```

Therefore, once multiple validator registrations share the same `CommPubKey`, the exact same `globalPubKey/publicCoeffs/pubKeyShare/signature` tuple can be replayed for each validator address. `AddGlobalPubKeyVote` stores only `(round, globalPubKey, coeffHash)` and increments blindly, and `FinalizeDKGRound` checks only the number of finalized registrations:

```go
newCount := current + 1
if err := k.GlobalPubKeyVotes.Set(ctx, key, newCount); err != nil { ... }
```

```go
finalizedCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusFinalized)
if finalizedCount < latestRound.Threshold {
    return k.SkipToNextRound(ctx, latestRound)
}
```

That reaches an “active” round even though the finalized set may contain fewer unique TEE identities / key shares than `Threshold`. The kernel code itself assumes one unique PID per key and simply takes the first registration whose `DkgPubKey` matches its own key:

```go
for _, reg := range regs {
    if bytes.Equal(reg.GetDkgPubKey(), ownPubKeyBytes) {
        ownPID = reg.GetIndex()
        break
    }
}
```

This confirms the duplicated registrations do not create duplicated private shares inside the TEE; they only create duplicated on-chain identities that are counted as if they were independent participants.

### Impact Cascade

- A staking operator with multiple validator addresses can make one enclave/keypair occupy multiple DKG seats in the same round.
- The chain can finalize and activate a committee whose number of unique key shares is below the configured threshold.
- `FinalizeDKGRound` still advances the round to `Active`, applies committee rewards, and keeps the self-unstake lockup logic keyed to this broken committee.
- CDR readers may now observe a “healthy” active round on-chain while threshold decryptions cannot reliably reach the cryptographic threshold, causing long-lived confidentiality-service failure until a later resharing succeeds.

### Assumptions and Uncertainties

1. The attacker controls multiple validator addresses in the active validator snapshot for the round.
2. The attacker can reuse the same `story-kernel` state/key directory across those validators or otherwise preserve the same `(DkgPubKey, CommPubKey)` pair.
3. The threshold cryptography requires distinct private shares / participant identities; duplicated registrations backed by one keypair do not increase the number of unique shares. This is the standard threshold assumption and matches the kernel’s single-PID-per-key behavior.

### How will the bug recipient respond?

“A large operator who controls multiple validators can already refuse to participate, so this is not a new attack.”

That response misses the actual exploit. Without this bug, a malformed or under-participating round should fail finalization and be skipped. Here, the attacker can make the round *pass* the protocol’s finalized-count and vote thresholds and become the active committee even though the committee is cryptographically incapable of meeting the advertised threshold. The issue is not mere non-participation; it is the ability to satisfy the protocol’s independence checks with duplicated identities and thereby activate a broken committee.

### Why did tests miss this issue?

Existing tests cover per-validator registration existence, signature verification, and threshold counting, but they do not model duplicate `DkgPubKey` / `CommPubKey` reuse across multiple validator addresses in the same round. The kernel-side tests also validate PID bounds, yet do not assert uniqueness of registration public keys before caching the first matching PID.

### Recommendation

Bind TEE identity and finalization to the validator address end to end:

```go
// key material must be namespaced by validator address as well as round
filepath.Join(keyDir, strconv.FormatUint(uint64(round), 10), codeCommitmentHex, strings.ToLower(validatorAddr), KeyEd25519File)
```

And reject duplicate key material on registration/finalization:
- Reject any registration whose `DkgPubKey` or `CommPubKey` already exists in the round.
- Include `validatorAddr` in the finalization hash.
- Enforce `validatorAddr == msg.sender` in `DKG.sol.finalize()`.
- Track finalization votes per validator address as well as per `(globalPubKey, coeffHash)`.

### References

1. [story-kernel/store/key_store.go#L17]
2. [story-kernel/service/dkg_generate_key.go#L18]
3. [story/client/x/dkg/keeper/dkg_handler.go#L24]
4. [story/client/x/dkg/keeper/dkg_handler.go#L288]
5. [story/contracts/src/protocol/DKG.sol#L229]
6. [story/client/x/dkg/keeper/dkg_votes.go#L14]
7. [story/client/x/dkg/keeper/dkg_finalization.go#L29]
8. [story-kernel/service/dkg_generate_deals.go#L126]
9. [story/docs/design/DKG.md#L368]

### STOR-10 — Early-sorted validator can suppress later DKG complaints by spoofing response dedup keys

### Early-sorted validator can suppress later DKG complaints by spoofing response dedup keys

### Executive Summary

The DKG vote-extension transport treats each validator’s vote extension as authenticated, but it does not authenticate the nested `Response` items before globally deduplicating them. `PrepareVotes` simply concatenates all committed vote extensions in deterministic commit order, then `deduplicateResponses` keeps the first `(dealer, responder)` pair it sees. Because Comet/Cosmos sorts `ExtendedCommitInfo.Votes` by voting power descending and validator address ascending, any earlier-sorted Byzantine validator can place forged `Response` objects for later validators inside its own valid vote extension and permanently suppress the real complaint responses those later validators actually signed.

This does not require proposer control. An honest proposer will still aggregate the attacker’s forged entries first, drop the honest entries as duplicates, and forward only the attacker-selected set into `FinalizeBlock`. When the forged complaints later fail kernel signature checks, no justification is emitted for the targeted dealer. That lets a malicious dealer suppress the complaint → justification → invalidation path that is supposed to exclude bad dealers before finalization and committee reward distribution.

### Details

`parseAndVerifyVoteExtension` only enforces raw size and per-extension item counts; it does not authenticate any nested `Deal`, `Response`, or `Justification` objects carried inside the vote extension:

```go
func (*Keeper) parseAndVerifyVoteExtension(voteExt []byte) ([]*types.Vote, bool, error) {
    if len(vote.Responses) > maxItemsPerVote {
        return nil, false, fmt.Errorf("too many responses in vote extension: %d exceeds max %d", len(vote.Responses), maxItemsPerVote)
    }
    return []*types.Vote{vote}, true, nil
}
```

`PrepareVotes` then processes committed vote extensions in commit order and calls `aggregateVotes`, which performs a first-wins deduplication across all nested response items:

```go
for _, vote := range commit.Votes {
    selected, _, err := k.parseAndVerifyVoteExtension(vote.VoteExtension)
    ...
    allVotes = append(allVotes, selected...)
}

votes := aggregateVotes(allVotes)
```

```go
func deduplicateResponses(responses []types.Response) []types.Response {
    seen := make(map[dedupKey]struct{})
    for _, r := range responses {
        var dealerIdx uint32
        if r.VssResponse != nil {
            dealerIdx = r.VssResponse.Index
        }
        key := dedupKey{responderIndex: r.Index, dealerIndex: dealerIdx}
        if _, exists := seen[key]; exists {
            continue
        }
        seen[key] = struct{}{}
        result = append(result, r)
    }
    return result
}
```

The naming in `deduplicateResponses` is misleading, but the actual pair is still the unique `(dealer, responder)` tuple: kyber defines the outer `dkg.Response.Index` as the target dealer index, while the inner `vss.Response.Index` is the responder/verifier index, and the Schnorr signature is checked only against the inner `vss.Response` fields.

```go
type Response struct {
    // Index of the Dealer for which this response is for
    Index uint32
    // Response issued from another participant
    Response *vss.Response
}
```

```go
func (a *Aggregator) verifyResponse(r *Response) error {
    pub, ok := findPub(a.verifiers, r.Index)
    if !ok { return errors.New("vss: index out of bounds in response") }
    if err := schnorr.Verify(a.suite, pub, r.Hash(a.suite), r.Signature); err != nil {
        return err
    }
    return a.addResponse(r)
}
```

This creates an exploitable provenance gap:

1. A malicious validator with an earlier commit-order position constructs a valid vote extension signed by itself.
2. Inside that vote extension, it inserts forged `Response` entries whose `(dealerIndex, responderIndex)` pairs match later validators’ real complaints against a targeted dealer. The outer vote extension is valid; only the nested complaint signatures are fake.
3. `PrepareVotes` appends the attacker’s vote first because `ExtendedCommitInfo` is deterministically ordered by voting power descending and address ascending.
4. `deduplicateResponses` keeps the attacker’s forged `(dealer,responder)` pairs and drops the later honest complaints as duplicates before any nested signature verification runs.
5. `msgServer.AddVote` forwards the reduced response set into `ProcessResponses`, which hands it to `story-kernel`.
6. `story-kernel` rejects the forged complaints during `distKeyGen.ProcessResponse`, but by then the real complaints are already gone and no justifications are produced.
7. Without those justifications, the dealer never reaches `invalidateDealerRegistration`, so it remains eligible to finalize and share in committee rewards.

The exploit is practical within the protocol’s own bounds. A single vote extension can carry up to `maxItemsPerVote = 80` responses, which is enough for one malicious committee member to spoof every other validator’s complaint pair in the intended committee size.

### Probing Questions

1. Are nested response items authenticated before `deduplicateResponses` runs?
2. Can an honest proposer remove this attack by recomputing the aggregate more carefully?
3. Does kernel-side signature verification recover the dropped honest complaints later?
4. Is commit ordering attacker-influenceable or deterministic enough for first-wins censorship?
5. Does suppressing complaint responses only affect liveness, or does it change consensus-visible DKG eligibility and rewards?

### Valid Argument

1. Nested items are not authenticated before dedup: `parseAndVerifyVoteExtension` checks only size/count limits, not nested response provenance.
2. Honest proposers still apply the vulnerable logic: `PrepareVotes` iterates committed vote extensions in order and then calls `aggregateVotes`/`deduplicateResponses`.
3. First-wins ordering is deterministic: Cosmos validates `ExtendedCommitInfo` ordering by voting power descending, address ascending.
4. Kernel verification happens too late: `ProcessResponses` forwards only the already-deduplicated set, and `story-kernel` merely logs and skips forged responses that fail `distKeyGen.ProcessResponse`.
5. Dealer invalidation is consensus-critical: invalidated dealers cannot finalize, and reward distribution iterates finalized registrations only.

### Invalid Counter-Argument

“This is harmless because forged complaints fail signature verification in the kernel, and honest complaints will still be present somewhere else.”

That does not hold. Kernel verification runs after `deduplicateResponses` has already removed the honest complaint entries. The honest duplicates never reach `ProcessResponses` or the kernel at all, so there is nothing left to recover once the forged complaints are rejected. This is also not proposer-only censorship: the honest proposer applies the same first-wins aggregation logic to the committed vote extensions it receives.

### Impact Cascade

- A single earlier-sorted Byzantine validator can suppress later validators’ complaint responses without proposer control.
- A malicious dealer can spoof complaint pairs targeting itself, preventing the honest complaints that should trigger dealer justifications from ever reaching kernel processing.
- Missing justifications prevent on-chain dealer invalidation, so the bad dealer remains eligible to finalize.
- Finalized registration status feeds directly into DKG committee reward distribution, so undeserving members can continue drawing rewards from the UBI-funded committee pool.

### Assumptions and Uncertainties

1. At least one Byzantine validator has an earlier `ExtendedCommitInfo` position than the honest complainants it wants to suppress.
2. The targeted dealer relies on suppressing complaint processing to avoid later invalidation.
3. I do not rely on direct theft from arbitrary user accounts; the concrete impact proven from code is committee-composition corruption and reward misallocation.

### How will the bug recipient respond?

“Those forged responses are invalid anyway, so the kernel drops them; this is only redundant gossip.”

The kernel indeed drops the forged responses, but that is the bug’s mechanism, not a defense. Because deduplication runs before nested signature verification, the forged entries consume the unique `(dealer,responder)` slots and erase the honest complaints before the kernel ever sees them.

### Why did tests miss this issue?

The existing dedup test explicitly treats “first duplicate wins” as benign and only checks slice length. It never models an earlier-sorted invalid response suppressing a later valid one, and the proposal-router tests never exercise honest aggregation of adversarially crafted nested response items.

### Recommendation

The aggregate step must preserve provenance until after nested semantic verification.

Primary fix:

```go
type attributedResponse struct {
    validatorAddr []byte
    response      types.Response
}
```

- Carry the enclosing validator identity alongside every nested response when aggregating vote extensions.
- Before deduplicating by `(dealer,responder)`, verify that the nested response’s claimed responder index is compatible with the enclosing validator and/or verify the nested signature against the registered DKG public key for that responder.
- Only deduplicate responses after provenance and signature checks pass.

Alternative fix:
- Remove first-wins dedup at proposal time and let finalize-time processing verify all candidate responses while discarding only semantically valid duplicates.
- Apply the same provenance hardening to `Deal` and `Justification`, which currently share the same “trust the enclosing vote extension, then first-win dedup” pattern.

### References

1. [/home/user/workspace/story/client/x/dkg/keeper/vote.go#L112-L227]
2. [/home/user/workspace/story/client/x/dkg/keeper/msg_server.go#L16-L48]
3. [/home/user/workspace/story/client/x/dkg/keeper/dkg_svc_dealing.go#L276-L338]
4. [/home/user/workspace/story/client/x/dkg/keeper/dkg_handler.go#L124-L132]
5. [/home/user/workspace/story/client/x/dkg/keeper/dkg_rewards.go#L259-L276]
6. [/home/vercel-sandbox/go/pkg/mod/github.com/piplabs/cosmos-sdk@v0.50.14-piplabs-v1.1/baseapp/abci_utils.go#L152-L190]
7. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/dkg/pedersen/structs.go#L61-L68]
8. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/vss/pedersen/vss.go#L609-L659]
9. [/home/user/workspace/story-kernel/service/dkg_process_responses.go#L84-L131]
10. [/home/user/workspace/story/client/x/dkg/keeper/vote_justification_test.go#L307-L319]

### STOR-9 — The CDR read path is dead in every DKG phase because active rounds kill the decrypt worker and resharing rounds reject request queueing

### The CDR read path is dead in every DKG phase because active rounds kill the decrypt worker and resharing rounds reject request queueing

### Executive Summary

Story’s threshold-decryption service is broken in both halves of the DKG lifecycle.

During `Active`, `FinalizeDKGRound` launches `handleDKGComplete` on a timeout-scoped async context and the wrapper goroutine defers `cancel()`. `handleDKGComplete` immediately calls `StartDecryptWorker(ctx)` and returns, so the wrapper goroutine exits and cancels the worker context before the worker’s first 3-second tick. `ResumeDKGService` treats `PhaseCompleted` as healthy during `Active`, so the dead worker is never restarted.

During `Registration` / `Dealing` / `Finalization`, `BeginBlocker` overwrites the previous active round’s stage to `Registration` before starting the resharing round. `ProcessCDRVaultRead` still routes reads to `LatestActiveRound`, but `ThresholdDecryptRequested` refuses to queue them because that round’s stage is no longer `Active`.

Together these bugs mean CDR reads are never automatically serviced: active-round requests have no worker, and resharing-window requests are skipped until timeout.

### Details

Active-round failure:

```go
// client/x/dkg/keeper/dkg_finalization.go
asyncCtx, cancel := dkgAsyncContext()
go func() {
    defer cancel()
    k.handleDKGComplete(asyncCtx, latestRound)
}()
```

`handleDKGComplete` updates the session to `PhaseCompleted` and starts the background decrypt worker on that same context:

```go
// client/x/dkg/keeper/dkg_svc_complete.go
session.UpdatePhase(types.PhaseCompleted)
session.IsFinalized = true
...
k.StartDecryptWorker(ctx)
```

But `StartDecryptWorker` runs a nested goroutine whose first useful action occurs only on a 3-second ticker:

```go
// client/x/dkg/keeper/dkg_svc.go
go func() {
    defer decryptWorkerRunning.Store(false)
    ticker := time.NewTicker(3 * time.Second)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            k.processDecryptQueue(ctx)
        }
    }
}()
```

Because `handleDKGComplete` returns immediately after `StartDecryptWorker`, the outer goroutine runs `defer cancel()` and kills the worker before its first tick. The session still sits in `PhaseCompleted`, and `ResumeDKGService` explicitly treats `PhaseCompleted` as the correct terminal state for `Active`, so no later block restarts the worker:

```go
// client/x/dkg/keeper/dkg_svc.go
case types.DKGStageActive:
    return phase != types.PhaseCompleted
```

Resharing-window failure:

```go
// client/x/dkg/keeper/abci.go
nextStage, shouldTransition := k.shouldTransitionStage(currentHeight, latestRound, params)
if shouldTransition {
    latestRound.Stage = nextStage
    ...
    case types.DKGStageRegistration:
        return k.InitiateDKGRound(ctx, false)
}
```

`ProcessCDRVaultRead` still binds new reads to `LatestActiveRound`:

```go
// client/x/evmengine/keeper/cdr.go
latestRound, err := k.dkgKeeper.GetLatestActiveRound(cachedCtx)
round := latestRound.Round
return k.dkgKeeper.ThresholdDecryptRequested(..., round, ...)
```

But `ThresholdDecryptRequested` rejects queueing unless that round’s stored stage is `Active`:

```go
// client/x/dkg/keeper/dkg_handler.go
dkgNetwork, err := k.getDKGNetwork(ctx, round)
if dkgNetwork.Stage != types.DKGStageActive {
    log.Info(ctx, "Skipping threshold decrypt request; DKG round is not active")
    return nil
}
session.AddDecryptRequest(...)
```

So:
1. Active rounds have an immediately canceled worker.
2. Resharing rounds route reads to a round whose stage has already been rewritten away from `Active`.
3. `PartialDecryptionTimeoutBlocks` is only 200, so skipped requests expire long before the next round can become active.

### Impact Cascade

- Every `CDR.read()` call can burn user fees without any automatic threshold-decryption service ever running.
- Active rounds fail because the decrypt worker dies before processing even a single queue tick.
- Registration / dealing / finalization fail because reads are routed to a non-active round and never enqueued.
- Timed-out requests are pruned, making the failure permanent rather than delayed.

### Assumptions and Uncertainties

1. Validators rely on the built-in decrypt worker started by `handleDKGComplete`; that is the documented production path.
2. I did not find any alternate automatic worker bootstrap after `PhaseCompleted` is reached.
3. Manual out-of-band validator intervention does not negate the bug; the protocol’s automatic service path is still dead by construction.

### How will the bug recipient respond?

“The worker uses a one-minute async context, so this only shortens service, not kills it.”

The code is stricter than that: the context is canceled immediately by the wrapper goroutine’s `defer cancel()` as soon as `handleDKGComplete` returns. Because the worker waits for its first ticker event, it exits before processing any request at all.

### Why did tests miss this issue?

The existing tests cover stage transitions and helper behavior in isolation, but they do not exercise the production composition of:
- `FinalizeDKGRound -> go handleDKGComplete -> StartDecryptWorker`,
- `Active -> Registration` stage mutation of the previous round,
- and a real `VaultRead` path across those phases.

### Recommendation

Fix both halves of the pipeline:

1. Decouple the decrypt worker from the short-lived async context. It should run on its own long-lived process context and be restarted if it dies.
2. Stop rewriting the previous active round’s stage away from `Active` until the successor round is ready to take over request servicing, or route reads using a separate “serving round” pointer instead of `LatestActiveRound` + current stage.
3. Add an integration test that covers `VaultRead` during `Active`, the `Active -> Registration` boundary, and the full resharing window.

### References

1. [client/x/dkg/keeper/dkg_finalization.go#L99-L106]
2. [client/x/dkg/keeper/dkg_svc_complete.go#L45-L60]
3. [client/x/dkg/keeper/dkg_svc.go#L238-L260]
4. [client/x/dkg/keeper/dkg_svc.go#L118-L133]
5. [client/x/dkg/keeper/abci.go#L89-L109]
6. [client/x/evmengine/keeper/cdr.go#L47-L105]
7. [client/x/dkg/keeper/dkg_handler.go#L414-L486]
8. [client/x/dkg/types/keys.go#L16-L23]

### STOR-8 — Resharing rounds never populate PIDCache, so every post-resharing active committee permanently fails CDR decryptions

### Resharing rounds never populate PIDCache, so every post-resharing active committee permanently fails CDR decryptions

### Executive Summary

`story-kernel.PartialDecryptTDH2()` requires a per-round PID entry from `PIDCache`, but the only place that populates that cache is `GenerateDeals()` and it explicitly skips `CachePID()` whenever `isResharing == true`. Story’s lifecycle makes every round after genesis a resharing round, including every kernel-upgrade round. Once such a round becomes active, validators start receiving `CDR.read()` requests, the CL routes them into `PartialDecryptTDH2()`, and the kernel deterministically fails with `PID not found`. There is no fallback path because the request’s `pid` field is ignored and no later stage backfills the cache. The result is that the first periodic resharing or upgrade permanently bricks threshold decryption for that round’s committee. Since `CDR.read()` burns the read fee before any partial decryptions are produced, users can keep paying for reads that can never succeed.

### Details

`PartialDecryptTDH2()` does not use the PID supplied by the CL request. Instead it loads `ownPID` only from `PIDCache`:

```go
ownPID, ok := s.PIDCache.Get(req.GetRound())
if !ok {
    return nil, status.Errorf(codes.FailedPrecondition,
        "PID not found: SetupDKGNetwork may not have been called for this round")
}
```

The sole population site for that cache is `GenerateDeals()`:

```go
if !req.GetIsResharing() {
    if err := s.CachePID(codeCommitmentHex, req.Round, rc.Registrations); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to cache PID")
    }
}
```

The code explicitly skips `CachePID()` for all resharing rounds. That is fatal because Story makes every round after the first a resharing round:

```go
func (k *Keeper) shouldReshare(ctx context.Context) (bool, error) {
    activeNetwork, err := k.getLatestActiveDKGNetwork(ctx)
    ...
    return activeNetwork != nil, nil
}
```

So after round 1, every subsequent round (periodic resharing and kernel upgrade resharing) is created with `IsResharing=true` and therefore never seeds PIDCache.

When that reshared round becomes active, decrypt requests flow normally:
- `CDR.read()` burns `readFee` and emits `VaultRead`.
- CL records the request and eventually calls `handleDecryptRequest()`.
- `handleDecryptRequest()` sends `PartialDecryptTDH2Request{Pid: session.Index, ...}`.
- `PartialDecryptTDH2()` ignores `req.Pid`, looks up `PIDCache`, and fails.

Because there is no alternate PID initialization after resharing finalization, the active committee for that round can never serve TDH2 partial decryptions.

### Impact Cascade

- The first periodic resharing or upgrade resharing round that becomes active disables threshold decryption for the new committee.
- All later `CDR.read()` calls still pass access control and burn `readFee`, but the committee can never produce partial decryptions for that round.
- Confidential data becomes unreadable exactly when the system rotates committees, which is a core-function failure in the primary production path.
- Users can repeatedly lose funds on read attempts with no possibility of successful decryption until the software is fixed.

### Assumptions and Uncertainties

1. The chain has progressed beyond the genesis DKG round, which is the normal steady state because every later round is resharing by design.
2. Validators run the shipped kernel/consensus code; no custom out-of-band patch backfills PIDCache.

### How will the bug recipient respond?

“The request already carries `pid`, so the cache miss is harmless or only affects tests.”

That is contradicted by the production implementation: `PartialDecryptTDH2()` never reads `req.GetPid()`. The only live PID source is `PIDCache`, and it is never populated for resharing rounds.

### Why did tests miss this issue?

The existing PID tests exercise bounds checks on an already-populated cache, but they do not simulate the real lifecycle where resharing rounds skip `CachePID()` and later become active decryption rounds.

### Recommendation

Use the request PID as the source of truth for partial decryption, or populate `PIDCache` for resharing rounds as soon as the validator’s round registration/index is known. At minimum:

```go
ownPID := req.GetPid()
if ownPID == 0 {
    ownPID, ok = s.PIDCache.Get(req.GetRound())
    ...
}
```

Also add an end-to-end regression covering: round 1 active -> round 2 resharing -> round 2 active -> `CDR.read()` -> successful `PartialDecryptTDH2()`.

### References

1. [story-kernel/service/dkg_partial_decrypt.go#L43]
2. [story-kernel/service/dkg_generate_deals.go#L49]
3. [story/client/x/dkg/keeper/dkg_initialization.go#L41]
4. [story/client/x/dkg/keeper/dkg_svc.go#L299]
5. [story/contracts/src/protocol/CDR.sol#L176]
6. [story/docs/design/DKG.md#L363]
7. [story/docs/design/DKG.md#L475]

### STOR-2 — Upgrade resharing rejects every validator that was not in the previous committee

### Upgrade resharing rejects every validator that was not in the previous committee

### Executive Summary

The upgrade resharing path hard-requires every current-round validator to already have a code commitment from the previous active round. That assumption is false for any validator that joined the active set after the current committee was formed, which is exactly the case resharing is supposed to handle. `InitiateDKGRound` precomputes `oldCC` from the validator’s previous-round registration and drops the lookup error. `handleDKGRegistration` only records `session.OldCodeCommitment` when that lookup succeeds, and `getRegistrationKernelClient` refuses to register an upgrade-round validator unless `session.OldCodeCommitment` is non-empty. The result is that every newly activated validator deterministically fails `GenerateAndSealKey` during a legitimate kernel upgrade and never reaches on-chain registration.

This is not an operator-only mistake. It triggers under normal validator-set churn once governance legitimately schedules an upgrade. The protocol therefore cannot rotate binaries and committee membership at the same time, even though the round is explicitly initialized from the current bonded validator set.

### Details

Upgrade rounds are initialized from the current bonded validator set, not the previous committee:

```go
activeValidators, err := k.GetActiveValidators(ctx)
...
dkgNetwork := types.DKGNetwork{
    ActiveValSet: activeValidators,
    IsResharing:  isResharing,
    IsUpgrade:    isUpgrade,
}
```

Yet the off-chain registration handler immediately derives the routing commitment from the *previous* active round only:

```go
oldCC, _ := k.getOldCodeCommitment(ctx)
go func() {
    k.handleDKGRegistration(asyncCtx, &dkgNetwork, oldCC, alreadyRegistered)
}()
```

`getOldCodeCommitment` reads the validator’s registration in the previous active round and returns an error when that validator did not previously belong to the committee:

```go
prevReg, err := k.getDKGRegistration(ctx, prevActive.Round, common.HexToAddress(k.validatorEVMAddr))
if err != nil {
    return nil, err
}
return prevReg.CodeCommitment, nil
```

The upgrade registration path then becomes impossible for new validators. `handleDKGRegistration` only fills `session.OldCodeCommitment` when `oldCC` is present, and `getRegistrationKernelClient` hard-fails otherwise:

```go
if dkgNetwork.IsUpgrade && len(oldCC) > 0 {
    session.OldCodeCommitment = oldCC
}
...
client, clientCC, cErr := k.getRegistrationKernelClient(isUpgrade, oldCC, session.OldCodeCommitment)
...
if len(upgradeOldCC) == 0 {
    return nil, nil, errors.New("old code commitment required for upgrade registration")
}
```

Because `Registered` only accepts validators from the current round’s `ActiveValSet`, these new validators have no alternate path to obtain a current-round registration later:

```go
if !slices.Contains(latest.ActiveValSet, strings.ToLower(validator.Hex())) {
    return errors.New("msg sender is not in the active validator set")
}
```

A validator that just joined the bonded set is therefore permanently excluded from the upgrade resharing round.

### Impact Cascade

- Newly active validators cannot generate a new enclave keypair or register during a scheduled upgrade.
- Binary rotation only works for validators that were already in the previous committee, defeating the intended “new membership + new binary” transition.
- Any upgrade that coincides with validator churn becomes structurally biased toward stale committee members.
- The subsequent threshold calculation and reward logic can operate on the surviving overlap subset rather than the scheduled validator set.

### Assumptions and Uncertainties

1. Governance/owner legitimately schedules a kernel upgrade.
2. At least one validator has joined the active bonded set since the previous DKG committee became active.
3. The joining validator is otherwise honest and runs the expected new kernel binary.

### How will the bug recipient respond?

“Upgrade resharing only supports the same validator set, so new validators are intentionally excluded.”

That explanation is contradicted by the implementation itself: `InitiateDKGRound` explicitly rebuilds `ActiveValSet` from the current bonded validator set, and `Registered` enforces membership in that new set. If unchanged membership were intended, the round should fail closed before registration starts instead of silently initializing an upgrade round that some active validators can never join.

### Why did tests miss this issue?

The upgrade lifecycle tests reuse the exact same mock validator set across both rounds and then call `env.registerValidators(t, 2)` for every validator in the second round. They never model a validator that appears only in the upgrade round, so the `old code commitment required for upgrade registration` path is never exercised.

### Recommendation

The upgrade registration path must support first-time committee members explicitly.

Primary fix:

```go
if isUpgrade && len(session.OldCodeCommitment) == 0 {
    // validator is new to the committee; register directly against the new binary
    return newClient, newCC, nil
}
```

Complementary hardening:
- Abort upgrade-round initialization up front if validator-set churn is unsupported.
- Store previous-round membership separately from current-round registration state so the router can distinguish “new member” from “misconfigured old member”.
- Add an integration test where round N+1 contains at least one validator absent from round N.

### References

1. [story/client/x/dkg/keeper/dkg_initialization.go#L41-L100]
2. [story/client/x/dkg/keeper/dkg_svc_registration.go#L49-L60]
3. [story/client/x/dkg/keeper/dkg_svc_registration.go#L118-L145]
4. [story/client/x/dkg/keeper/dkg_svc_registration.go#L192-L229]
5. [story/client/x/dkg/keeper/dkg_svc_registration.go#L232-L249]
6. [story/client/x/dkg/keeper/dkg_handler.go#L24-L49]
7. [story/client/x/dkg/keeper/dkg_lifecycle_internal_test.go#L352-L432]

### STOR-1 — Partial-decrypt signatures omit requester and UUID binding, so a hostile keeper can redirect valid shares while CL accepts the forged request context

### Partial-decrypt signatures omit requester and UUID binding, so a hostile keeper can redirect valid shares while CL accepts the forged request context

### Executive Summary

The kernel encrypts each TDH2 partial to whatever `requester_pub_key` the caller supplies, but the response signature omits every field that tells the chain who that partial was actually generated for. `signPartialDecryptResponse()` covers only `round || ciphertext || encryptedPartial || ephemeralPubKey || pubShare`. On the consensus side, `verifyPartialDecryptionSignature()` recreates the same truncated hash, while `handleDecryptRequest()` submits `requesterPubKey` and `uuid` as separate event fields after the signature has already been produced. As a result, a hostile keeper/host can take a legitimate queued decrypt request, swap in an attacker-controlled requester key when talking to the enclave, receive a valid partial encrypted to the attacker, and then submit that same signed payload to the contract with the original queued requester key and UUID. The CL accepts the submission because the forged context is outside the signed domain. Repeating this across enough malicious validators leaks the plaintext to the attacker while the victim only sees apparently valid but undecryptable partials on-chain.

### Details

The kernel’s signature domain does not include the requester public key, UUID/label, PID, or even the code commitment. Only the round, ciphertext, encrypted partial, ephemeral pubkey, and pubshare are signed.

```go
func (s *DKGServer) signPartialDecryptResponse(codeCommitmentHex string, round uint32, ciphertext []byte, encryptedPartial []byte, ephPubKey []byte, pubShareBz []byte) ([]byte, error) {
    encoded := make([]byte, 0, 4+len(ciphertext)+len(encryptedPartial)+len(ephPubKey)+len(pubShareBz))
    encoded = append(encoded, uint32ToBytes(round)...)
    encoded = append(encoded, ciphertext...)
    encoded = append(encoded, encryptedPartial...)
    encoded = append(encoded, ephPubKey...)
    encoded = append(encoded, pubShareBz...)
    respHash := ecrypto.Keccak256(encoded)
    ...
}
```

The CL verifier reproduces exactly the same domain, so it cannot detect requester or UUID substitution.

```go
// encoded = round(4B big-endian) || ciphertext || encryptedPartial || ephPubKey || pubShare
func verifyPartialDecryptionSignature(commPubKey []byte, round uint32, ciphertext []byte, encryptedPartial, ephemeralPubKey, pubShare, signature []byte) error {
    encoded = append(encoded, roundBytes...)
    encoded = append(encoded, ciphertext...)
    encoded = append(encoded, encryptedPartial...)
    encoded = append(encoded, ephemeralPubKey...)
    encoded = append(encoded, pubShare...)
    respHash := crypto.Keccak256(encoded)
    ...
}
```

Meanwhile the keeper sends `requesterPubKey` to the enclave, but then independently submits `requesterPubKey` and `uuid` to the contract after the response is signed.

```go
resp, err := client.PartialDecryptTDH2(ctx, &types.PartialDecryptTDH2Request{
    Round:           session.Round,
    Ciphertext:      req.Ciphertext,
    Label:           req.Label,
    GlobalPubKey:    session.GlobalPubKey,
    RequesterPubKey: req.RequesterPubKey,
})

k.contractClient.SubmitEncryptedPartialDecryption(
    ctx,
    session.Round,
    pid,
    resp.EncryptedPartialDecryption,
    resp.EphemeralPubKey,
    resp.PubShare,
    req.RequesterPubKey,
    req.Ciphertext,
    uuid,
    resp.Signature,
)
```

The contract event also treats `requesterPubKey` and `uuid` as unsigned payload fields.

```solidity
emit EncryptedPartialDecryptionSubmitted(
    msg.sender,
    round,
    pid,
    encryptedPartial,
    ephemeralPubKey,
    pubShare,
    requesterPubKey,
    ciphertext,
    uuid,
    signature,
    fee
);
```

Finally, CL acceptance checks only that the request registry matches the submitted `(requesterPubKey, label, round, ciphertext)` tuple and that the signature matches the truncated domain above. It never proves that the enclave signed over the same requester or UUID.

Exploit path:
1. The keeper queue contains a legitimate request `(round R, ciphertext C, label L(uuid), requester key K_victim)`.
2. A hostile keeper/host calls the enclave with the same `(R, C, L)` but swaps `requester_pub_key` to `K_attacker`.
3. The enclave returns `(encryptedPartial_to_K_attacker, ephPub, pubShare, sig)`; `sig` is valid because it does not cover `K_attacker` or `uuid`.
4. The same hostile keeper submits that response to `submitEncryptedPartialDecryption()` using the original `K_victim` and `uuid` from the queue.
5. `PartialDecryptionSubmitted()` accepts the event, because the registry lookup uses the forged `(K_victim, uuid)` values and the signature verifier ignores both of them.
6. The attacker keeps the usable encrypted partial off-chain, while the victim later fetches an apparently valid on-chain partial that cannot be decrypted with `K_victim`.

The protocol’s own whitepaper says kernels should only act on proved smart-contract events and that committee members should not be able to trigger unauthorized malicious actions. This bug leaves the most important request-identity fields entirely outside the enclave-authenticated transcript.

### Impact Cascade

- A malicious validator host can silently redirect its kernel’s partial decryption to an attacker-controlled requester key.
- The forged submission still passes CL verification, so the attack is camouflaged as a legitimate partial-decryption submission.
- Victims receive unusable partials and cannot detect the substitution from signature verification alone.
- Across threshold-many malicious validators, plaintext is disclosed to the attacker even though the on-chain artifacts claim the decryption was produced for the victim requester.

### Assumptions and Uncertainties

1. The attacker controls the validator-side caller that talks to the enclave (host, keeper, or another component on that path). This is the exact trust boundary the kernel is supposed to defend.
2. Enough such validators misbehave to reach the TDH2 threshold for the target ciphertext.
3. The victim later relies on CL-stored partials and signature verification to infer that the shares were produced for their request context.

### How will the bug recipient respond?

“A hostile keeper can always refuse to submit or can already exfiltrate whatever it sees, so rebinding the requester does not change the security model.”

That response contradicts the protocol’s stated TEE goal. The whitepaper explicitly claims kernels only act on proved smart-contract events and prevent committee members from triggering unauthorized malicious actions. Here the enclave is doing cryptographic work for the attacker’s requester key, and the CL cannot detect the substitution because requester identity and UUID are missing from the signed transcript. This is a confidentiality break at the exact oracle-containment boundary the TEE is supposed to enforce.

### Why did tests miss this issue?

The unit tests for `verifyPartialDecryptionSignature()` only mutate fields that are already in the signed message (`ciphertext`, `encryptedPartial`, `ephemeralPubKey`, `pubShare`). They never test that changing `requesterPubKey`, `uuid/label`, or the request context should invalidate the signature, because those fields were omitted from both the implementation and the mirrored tests.

### Recommendation

Extend the kernel signature domain to cover the entire request/response context that the chain relies on, at minimum:

```
keccak256(
    code_commitment || round || ciphertext || label || requester_pub_key ||
    encrypted_partial || ephemeral_pub_key || pub_share || pid
)
```

Then update `verifyPartialDecryptionSignature()` and the on-chain/CL submission path to reject any event whose signed requester key, label/UUID, PID, or code commitment differs from the canonical request being processed. As a second line of defense, `PartialDecryptTDH2()` should validate a canonical request proof before generating any partial at all.

### References

1. [story-kernel/service/dkg_partial_decrypt.go#L142]
2. [story-kernel/service/dkg_partial_decrypt.go#L207]
3. [story/client/x/dkg/keeper/dkg_svc.go#L299]
4. [story/client/x/dkg/keeper/dkg_handler.go#L362]
5. [story/client/x/dkg/keeper/dkg_handler.go#L491]
6. [story/contracts/src/protocol/CDR.sol#L223]
7. [story/docs/design/CDR.md#L119]
8. [ai-docs/core/whitepaper/confidentialdatarails.md#L385]

## Low

### STOR-17 — Sealed light-client state is rollbackable, so a hostile host can pin the enclave to stale trust roots

### Sealed light-client state is rollbackable, so a hostile host can pin the enclave to stale trust roots

### Executive Summary

The kernel stores its trusted Comet light-client state in a sealed LevelDB, but sealing is used only for confidentiality/integrity of individual values — not for freshness. Every record is sealed with `SealWithUniqueKey(..., nil)` and later resumed solely based on the presence of LevelDB entries. There is no monotonic counter, no authenticated “latest seen height”, and no rollback marker that would let the enclave detect that the host restored an older snapshot of `light_client.db`.

Under this focus area’s threat model (host OS and RPC endpoints may be hostile), that is enough to break canonical-state guarantees. The host can snapshot a previously trusted database, restore it later while it is still within the trusted period, and reboot the kernel into that stale trust root. If the primary/witness RPC set is also attacker-controlled, the enclave can be kept on that stale view indefinitely and will verify future state queries against the rolled-back trust root. The TEE then makes key-generation or decryption decisions against outdated Story state even though its “verified query client” appears healthy.

### Details

The generic sealing helper uses a unique SGX key with `nil` additional authenticated data. Nothing binds the sealed blob to freshness, height, or a monotonic version.

```go
sealedData, err := ecrypto.SealWithUniqueKey(data, nil)
...
key, err := ecrypto.Unseal(sealed, nil)
```

The light-client database applies the same pattern to every LevelDB value:

```go
sealedValue, err := ecrypto.SealWithUniqueKey(value, nil)
...
unsealedValue, err := ecrypto.Unseal(sealedValue, nil)
```

On startup, `initializeQueryClient()` simply checks whether any trusted light-client state exists in the sealed DB and, if so, resumes from it.

```go
hasExistingState, err := story.HasTrustedState(db, cfg.LightClient.ChainID)
...
if hasExistingState {
    queryClient, err := story.LoadVerifiedQueryClient(ctx, cfg, db)
    ...
    return queryClient, nil
}
```

There is no attempt to authenticate that this is the *latest* state ever seen by the enclave. A host can therefore:
1. Snapshot `light_client/light_client.db` at height `H`.
2. Wait until the chain advances and committee state changes.
3. Restore the old sealed snapshot while `H` is still inside the configured trusted period.
4. Restart the kernel, which resumes from the rolled-back store.
5. Feed matching stale headers/app hashes from attacker-controlled RPC/witness endpoints.

From the enclave’s perspective, all subsequent ICS23 queries are now being verified against an attacker-selected stale trust root.

### Impact Cascade

- The enclave can be pinned to an outdated committee view even though it is using a “verified” light client.
- Key generation and finalization can be driven from stale DKG rounds after the chain has already rotated committees.
- Start-block validation and latest-active-round selection inherit the rolled-back trust root, undermining report-data anchoring and decryption gating.
- Because the rollback happens inside the sealed persistence layer, ordinary operator checks may not notice that the enclave resumed from an old snapshot.

### Assumptions and Uncertainties

1. The attacker has host/filesystem control over the validator machine, which is explicitly in scope for this focus area.
2. To keep the enclave on stale state rather than letting it catch up, the attacker also controls the RPC/witness providers the light client talks to.
3. If any independent witness is honest and reachable, the stale view may eventually be detected or superseded.

### How will the bug recipient respond?

“Sealing already protects the database, and the Comet light client enforces a trusted period, so rollback is not exploitable.”

Sealing only protects *contents* from tampering; it does not prove *freshness*. The trusted period merely says how long an old header remains acceptable as a root of trust. That is exactly why snapshot rollback matters: the host can restore an old but still-valid trust root and the enclave has no monotonic evidence that it is stale.

### Why did tests miss this issue?

The existing logic tests fresh startup vs. resume-from-existing-state, but there are no adversarial tests that restore an older sealed snapshot and assert the enclave rejects it as stale.

### Recommendation

Add rollback protection to the trusted-state store:
- Seal an authenticated metadata record containing the highest trusted height and reject any restored DB whose metadata regresses.
- Bind sealed records to a monotonic version (TPM/SGX monotonic counter, remote monotonic service, or an append-only host-independent checkpoint).
- On resume, require the restored height to be at least as new as the last sealed checkpoint before accepting the DB.

Without a freshness primitive, sealing alone is not enough for trusted-state persistence.

### References

1. [story-kernel/enclave/seal.go#L11-L36]
2. [story-kernel/enclave/sealed_leveldb.go#L78-L145]
3. [story-kernel/server/server.go#L90-L131]

### STOR-14 — Honest signed DKG deals can be turned into undeliverable shares because nonce and recipient routing are outside the dealer signature

### Honest signed DKG deals can be turned into undeliverable shares because nonce and recipient routing are outside the dealer signature

### Executive Summary

Story treats the outer `Deal` protobuf as an authenticated dealer message, but the dealer’s Schnorr signature does not cover the fields that actually determine whether the intended recipient can decrypt and route the share. In kyber, `dkg.Deal.MarshalBinary()` signs only `(dealerIndex, cipher)`. Story then routes the message to a recipient using the unsigned outer `recipient_index`, while decryption uses the unsigned `nonce`. A malicious relay/proposer can therefore take an honest kernel-produced deal, keep both existing signatures valid, and alter either `recipient_index` or `nonce` before the deal reaches consensus.

The result is not a harmless drop: the rightful recipient never produces a response, the dealer’s deal is never certified, and Story permanently excludes that honest dealer from `QUAL()` because vote-extension delivery is one-shot and Story never enables kyber’s timeout path. Since finalization only requires a threshold of certified deals, the round can still finalize on an attacker-selected subset of dealers even though the excluded dealer never sent an invalid share.

### Details

Kyber’s dealer signature is too narrow. `dkg.Deal.MarshalBinary()` signs only the outer dealer index and ciphertext bytes:

```go
func (d *Deal) MarshalBinary() ([]byte, error) {
    var b bytes.Buffer
    binary.Write(&b, binary.LittleEndian, d.Index)
    b.Write(d.Deal.Cipher)
    return b.Bytes(), nil
}
```

The ephemeral `EncryptedDeal.Signature` also covers only `DHKey`, not `Nonce` or the outer routing field:

```go
dhPublicBuff, _ := dhPublic.MarshalBinary()
signature, err := schnorr.Sign(d.suite, d.long, dhPublicBuff)
```

Story preserves those unauthenticated fields as consensus inputs. The protobuf deal contains an unsigned outer `recipient_index`, and `story-kernel` reconstructs the kyber deal from `dh_key`, `signature`, `nonce`, and `cipher` exactly as received:

```go
message Deal {
  uint32 index = 1;
  uint32 recipient_index = 2;
  EncryptedDeal deal = 3;
  bytes signature = 4;
}
```

```go
func ConvertToDeal(deal *pb.Deal) *dkg.Deal {
    return &dkg.Deal{
        Index: deal.GetIndex(),
        Deal: &vss.EncryptedDeal{
            DHKey:     deal.GetDeal().GetDhKey(),
            Signature: deal.GetDeal().GetSignature(),
            Nonce:     deal.GetDeal().GetNonce(),
            Cipher:    deal.GetDeal().GetCipher(),
        },
        Signature: deal.GetSignature(),
    }
}
```

The CL routes by the unsigned outer `recipient_index` *before* any cryptographic validation:

```go
for _, deal := range deals {
    if session.Index > 0 && deal.RecipientIndex == session.Index-1 {
        req.Deals = append(req.Deals, deal)
    }
}
```

And the kernel explicitly assumes the routing metadata is already trustworthy:

```go
// ProcessDeals process the deals. It is assumed that the deal has been correctly delivered to the corresponding recipient index.
```

Once the wrong recipient (or the right recipient with a tampered nonce) processes the deal, kyber fails in `decryptDeal()` before generating any complaint response:

```go
decrypted, err := gcm.Open(nil, e.Nonce, e.Cipher, v.hkdfContext)
if err != nil {
    return nil, err
}
```

So the attack is silent: no justification path is created, no dealer-invalidating evidence is emitted, and the only observable effect is an absent response from the rightful recipient.

That absent response is enough to remove an honest dealer from the transcript. Story documents that vote extensions deliver each deal exactly once, so the omitted/corrupted copy is permanently lost. Kyber then includes only certified deals in `QUAL()`, while Story finalizes once merely a threshold subset is certified:

```go
if !v.DealCertified() {
    return nil
}
```

```go
func (d *DistKeyGenerator) ThresholdCertified() bool {
    return len(d.QUAL()) >= d.c.Threshold
}
```

This lets a malicious proposer (who can choose the `MsgAddDkgVote` contents that honest validators do not cross-check) transform an honest dealer’s signed deal into a silent non-delivery and thereby exclude that honest dealer from the final key transcript.

### Impact Cascade

- An attacker does not need the dealer’s private key; they only need one honest signed deal and the ability to relay it into consensus.
- By altering unsigned `nonce` or `recipient_index`, the attacker turns a valid share into a silent non-delivery rather than a complaint that could be adjudicated.
- The targeted honest dealer loses the missing response forever because DKG deal delivery is one-shot.
- `QUAL()` excludes that honest dealer, while `DistKeyShare()` still finalizes if the attacker leaves a threshold subset of dealers intact.
- The active DKG key can therefore be derived from an attacker-selected subset even though the excluded dealer never sent an invalid share and never triggered a complaint/justification trail.

### Assumptions and Uncertainties

1. The attacker must be able to relay the mutated deal into the consensus transcript. A malicious block proposer can do this because `MsgAddDkgVote` is not checked against `ProposedLastCommit` during proposal validation.
2. The excluded dealer impact is transcript exclusion / cryptographic committee skew. Full secret recovery additionally requires the attacker to keep the certified subset within a coalition they control.

### How will the bug recipient respond?

“Changing `recipient_index` or `nonce` only causes a missing response, so this is just liveness.”

That understates the effect. In Story’s implementation, missing responses are not temporary transport noise: vote-extension delivery is one-shot, `DealCertified()` gates `QUAL()`, and finalization proceeds on only-threshold certification. The attacker is not merely delaying a message; they are deterministically removing an honest dealer’s contribution from the final transcript.

### Why did tests miss this issue?

The tests exercise justification helpers and signature checks with synthetic in-process messages, but they do not cover the signed-field boundary of `dkg.Deal.MarshalBinary()` or an end-to-end relay that mutates unsigned `nonce` / `recipient_index` while preserving both existing signatures.

### Recommendation

Bind all deal-delivery-critical fields into the dealer-authenticated transcript.

At minimum:
- include `RecipientIndex`, `Deal.Nonce`, and `Deal.DHKey` in the outer dealer signature;
- or eliminate the unsigned outer `recipient_index` and derive routing from authenticated state;
- reject any proposal-time DKG aggregate that does not match the authenticated prior vote extensions.

### References

1. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/dkg/pedersen/structs.go#L52-L58]
2. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/vss/pedersen/vss.go#L178-L199]
3. [story/client/proto/story/dkg/v1/types/types.proto#L145-L152]
4. [story-kernel/types/types.go#L15-L25]
5. [story/client/x/dkg/keeper/dkg_svc_dealing.go#L166-L171]
6. [story-kernel/service/dkg_process_deals.go#L19-L24]
7. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/vss/pedersen/vss.go#L390-L414]
8. [story/client/x/dkg/keeper/dkg_svc_dealing.go#L137-L143]
9. [/home/vercel-sandbox/go/pkg/mod/go.dedis.ch/kyber/v4@v4.0.0-pre2/share/dkg/pedersen/dkg.go#L500-L510]

### STOR-12 — Whitelisting a new enclave type lets validators bypass scheduleUpgrade and rotate binaries early

### Whitelisting a new enclave type lets validators bypass scheduleUpgrade and rotate binaries early

### Executive Summary

`scheduleUpgrade()` / `cancelUpgrade()` are supposed to be the authoritative control plane for when a new story-kernel measurement becomes active. In practice, the consensus layer never enforces that relationship. Once governance/admin has merely whitelisted another enclave type, current validators can directly call `DKG.register()` and `DKG.finalize()` with that alternative whitelisted measurement in an ordinary round, and the CL will accept those events without consulting pending/activated upgrade state at all. Because dealing/finalization thresholds are computed from accepted registrations rather than an expected round code commitment, the round can become active under an unscheduled binary or under a binary whose scheduled upgrade was later cancelled. This lets a validator coalition move committee authority across enclave measurements outside the EL-controlled activation/cancellation workflow.

### Details

The contract clearly separates two admin operations:
1. `whitelistEnclaveType(...)` approves an enclave measurement.
2. `scheduleUpgrade(activationHeight, upgradeVersion)` is supposed to decide when a new binary takes over.

The docs mirror that sequence: first whitelist the new type, then schedule activation at a future height.

But the production CL path never checks whether the code commitment presented in `Registered` / `Finalized` is actually authorized for the current round. `Registered()` accepts any event whose validator is in `ActiveValSet`; `Finalized()` verifies only the registrant’s signature and participants root; `BeginDealing()` and `FinalizeDKGRound()` count registrations/finalizations without any code-commitment invariant.

Therefore, as soon as a second enclave type is whitelisted, validators can bypass `scheduleUpgrade()` entirely:
- run the alternate binary,
- call `DKG.register()` directly with that enclave type in the next ordinary round,
- finalize under that same code commitment,
- have the round activated by CL as if governance had approved the timing.

The same bypass survives `cancelUpgrade()`: cancellation only deletes stored upgrade metadata, but registration/finalization still accept any whitelisted enclave type.

### Impact Cascade

- A validator coalition can activate a new enclave measurement before its scheduled activation height, or after the scheduled upgrade was cancelled.
- EL-side schedule/cancel events stop being the effective authority over binary rotation timing.
- The active committee can silently shift to a different enclave measurement without entering `IsUpgrade=true` resharing mode.
- That undermines the protocol’s cross-layer trust boundary: the CL treats whitelisting alone as sufficient for committee authority, even though the EL API exposes a separate explicit upgrade scheduler.

### Assumptions and Uncertainties

1. At least two enclave types are simultaneously whitelisted, which the documented upgrade flow explicitly requires before activation.
2. Enough current validators cooperate to finalize the round under the alternate measurement if different binaries produce incompatible transcripts; if transcripts match, even a smaller subset can remain on the alternate measurement inside the active committee.

### How will the bug recipient respond?

“Whitelisting a type is already admin approval, so validators using it early is acceptable.”

That is inconsistent with the public interface and docs. The protocol exposes a dedicated `scheduleUpgrade()` / `cancelUpgrade()` mechanism precisely to control when a new measurement becomes active. If whitelisting alone were sufficient, the separate upgrade control plane would be meaningless.

### Why did tests miss this issue?

Tests cover storing upgrade metadata and activating upgrade rounds, but they do not model validators directly submitting `register` / `finalize` events with an alternate whitelisted enclave type during a non-upgrade round, nor do they assert that ordinary rounds must inherit the previous active code commitment unless `IsUpgrade` is set.

### Recommendation

Make the CL enforce a round-level expected code commitment:
- ordinary rounds must reject registrations/finalizations whose code commitment differs from the previous active round’s commitment;
- upgrade rounds must reject current-round members using the old commitment;
- `cancelUpgrade()` should restore the prior expected commitment immediately.

### References

1. [story/contracts/src/protocol/DKG.sol#L99]
2. [story/contracts/src/protocol/DKG.sol#L129]
3. [story/contracts/src/protocol/DKG.sol#L140]
4. [story/client/x/dkg/keeper/dkg_handler.go#L24]
5. [story/client/x/dkg/keeper/dkg_handler.go#L103]
6. [story/client/x/dkg/keeper/dkg_dealing.go#L15]
7. [story/client/x/dkg/keeper/dkg_finalization.go#L29]
8. [story/client/x/dkg/README.md#L237]

### STOR-7 — First-seen vote dedup lets an earlier validator erase later signed DKG deals and responses

### Executive Summary

`PrepareVotes()` aggregates all validators’ signed vote-extension payloads and immediately deduplicates raw DKG `Deal` and `Response` entries **before** any DKG cryptographic verification. Because the dedup keys are only `(dealerIndex, recipientIndex)` for deals and `(responderIndex, dealerIndex)` for responses, a malicious validator can place a fake colliding entry in its own vote extension and cause the later honest entry to be discarded. The surviving fake is only rejected much later inside the kernel, after the real message has already been dropped forever.

This does **not** require proposer control. `ExtendedCommitInfo.Votes` is processed in deterministic validator order, so any validator that sorts before a target validator can win the “first seen” slot. The result is consensus-wide censorship of honest DKG traffic: the attacker can erase later validators’ real deals or responses, prune selected honest dealers from `QUAL()`, or force round retries by creating permanent absent-response conditions.

### Details

`VerifyVoteExtension()` only checks vote-extension size/count and protobuf parseability; it does not authenticate nested DKG objects.

`PrepareVotes()` then iterates `commit.Votes`, appends every parsed payload, and calls `aggregateVotes()`:

```go
for _, vote := range commit.Votes {
    selected, _, err := k.parseAndVerifyVoteExtension(vote.VoteExtension)
    ...
    allVotes = append(allVotes, selected...)
}
votes := aggregateVotes(allVotes)
```

`aggregateVotes()` performs first-wins dedup on raw fields only:

```go
// deals
key := dedupKey{dealerIndex: d.Index, recipientIndex: d.RecipientIndex}
// responses
key := dedupKey{responderIndex: r.Index, dealerIndex: dealerIdx}
```

No signature or session validation happens before this elimination step.

That ordering is attacker-controlled enough to exploit. Cosmos SDK requires `ExtendedCommitInfo.Votes` to be sorted deterministically by descending voting power, with address as the tiebreaker. An attacker that sorts before an honest validator therefore does not need proposer control: an honest proposer calling `PrepareVotes()` will still see the attacker’s fake colliding entry first and drop the later honest one.

The dropped message is not recovered later:
- `ProcessDeals()` forwards the deduped list directly to the kernel.
- `ProcessResponses()` forwards the deduped list directly to the kernel.
- The kernel verifies each item individually and silently skips invalid ones.

For deals:

```go
resp, err := distKeyGen.ProcessDeal(deal)
if err != nil {
    ...
    continue
}
```

For responses:

```go
j, err := distKeyGen.ProcessResponse(resp)
if err != nil {
    ...
    continue
}
```

Vote-extension delivery is one-shot because `ExtendVote()` dequeues messages from in-memory queues. Once the honest later message is dropped during aggregation, there is no authenticated retransmission path.

A concrete exploit against responses is:
1. Pick an honest validator `R` that sorts after the attacker.
2. The attacker includes a fake `Response` in its own vote extension with outer `Index = R` and inner `VssResponse.Index = D` for any target honest dealer `D`.
3. The real signed response from `R` about dealer `D` appears later in `commit.Votes`, but `deduplicateResponses()` discards it because the `(R, D)` slot is already occupied.
4. The kernel rejects the attacker’s fake because it is not signed by `R`, yet the honest response is already gone.
5. Kyber now sees dealer `D` as missing one verifier response, so `DealCertified()` fails before timeout and `D` is excluded from `QUAL()`.

The same primitive exists for deals: a fake colliding `(dealerIndex, recipientIndex)` deal erases the later honest deal before kernel verification, causing the intended recipient to miss that share and produce no response.

### Impact

The practical impact is DKG transcript censorship without proposer control:
- a single earlier-sorted validator can erase later validators’ honest DKG messages;
- selected honest dealers can be pruned from `QUAL()` even though they sent correct messages;
- enough such suppressions can force DKG round failure/retry.

This is stronger than ordinary withholding because the attacker is not limited to omitting its own vote-extension contents; it can delete other validators’ honest DKG traffic during consensus aggregation.

I am **not** claiming direct secret-key recovery or immediate protocol-funds theft from this bug alone. The substantiated effects are transcript-integrity corruption and round-level censorship/availability degradation.

### Why this was missed

Existing tests cover structural vote-extension parsing and helper-level dedup, but do not exercise the adversarial ordering case where an invalid duplicate appears before the valid message carrying the same dedup key. They also do not combine that first-wins behavior with the kernel’s “log and continue” handling for invalid deals/responses.

### Recommendation

Authenticate first, deduplicate second.

Concretely:
1. Verify DKG message authenticity before first-wins dedup (or keep all colliding candidates until authenticity is checked).
2. Bind proposal-time DKG aggregation to authenticated vote-extension contents and reject invalid colliding placeholders instead of silently preferring the first raw entry.
3. Treat invalid dedup-colliding DKG entries as slashable / attributable faults rather than dropping the honest later message.