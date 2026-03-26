# Story Blockchain Architect Memory

## Key Architecture Patterns

### DKG Module
- Module name: `dkg`, store key: `dkg`
- Keeper at `client/x/dkg/keeper/keeper.go` requires: AccountKeeper, StakingKeeper, BankKeeper, DistributionKeeper, EvmStakingKeeper
- Dependencies wired via depinject in `client/x/dkg/module/depinject.go`
- DKG rewards flow: `BeginBlocker -> FinalizeDKGRound -> distributeDKGCommitteeRewards`
- `distributeDKGCommitteeRewards` is called BEFORE `setLatestActiveRound` so it can read the previous active round
- See [dkg-patterns.md](dkg-patterns.md) for details

### Testing Infrastructure
- DKG keeper tests: internal tests in `keeper` package using `setupDKGKeeperWithMocks`
- Mocks generated via: `mockgen -source=client/x/dkg/types/expected_keepers.go -package testutil -destination client/x/dkg/testutil/expected_keepers_mocks.go`
- Evmstaking tests use `createKeeperWithMockStaking(t)` returning 7 values
- Pre-existing test state leak in `TestKeeper_RegistrationInitialized` (subtests share keeper)

### Cross-Module Interactions
- DKG -> Distribution: `GetUbiBalanceByDenom`, `WithdrawUbiByDenomToModule`
- DKG -> Bank: `SendCoinsFromModuleToAccount` (rewards to validators)
- DKG -> EvmStaking: `ProcessUbiWithdrawalFromAmount` (remaining UBI burn+queue)
- EvmStaking `ProcessUbiWithdrawalFromAmount` skips `minPartialWithdrawalAmount` check

### DKG Justification System (implemented Feb 2026)
- Justifications flow: TEE ProcessResponses -> EnqueueJustifications -> VoteExtension -> AddVote -> ProcessJustifications
- VSS verification is deterministic math (Edwards25519), safe for consensus path
- Invalid deal -> `invalidateDealerRegistration` sets `DKGRegStatusInvalidated` (consensus state change)
- Valid deal -> forwards to TEE `ProcessJustification` RPC (async goroutine, non-consensus)
- Finalization handler has explicit guard: invalidated dealers cannot finalize
- Proto changes were additive only (no upgrade handler needed)
- story-kernel: `ConvertToJustification` converts proto -> kyber types with Edwards25519 unmarshal
- story-kernel: `ProcessJustification` handler pattern matches ProcessDeals/ProcessResponses
- Queue pattern: `EnqueueJustifications([]*types.Justification)` takes pointer slice, dereferences to value
- `lib/vss/verify.go`: VerifyPedersenVSS implements share*G == sum(x^k * C_k)

### TEE Kernel Upgrade Architecture (implemented Mar 2026)
- KernelRouter (`client/x/dkg/keeper/kernel_router.go`): map[ccHex] -> TEEClient, strict matching (no fallback)
- Multi-endpoint config: `client/config/dkg_config.go` `KernelEndpoints []string` (was single `TEEEndpoint`)
- Keeper field: `kernelRouter *KernelRouter` (was single `teeClient`)
- Startup: `ConnectAndDiscover` calls `GetIdentity` on each endpoint, registers by code commitment
- DKGNetwork collection key: round only (string), NOT codeCommitment_round
- Registration key: `round_address`, Votes key: `round_pubkey_coeffhash`
- KernelUpgradeInfos collection (prefix 6): key=upgradeVersion, value=KernelUpgradeInfo{version, height, isActivated}
- Upgrade flow: DKG.sol scheduleUpgrade -> UpgradeScheduled event -> store KernelUpgradeInfo -> BeginBlocker checkPendingUpgradeActivation at height -> InitiateDKGRound(isUpgrade=true)
- DKGNetwork.is_upgrade field (proto field 11, additive)
- DKGSession.OldCodeCommitment: old binary CC for dealer routing during upgrade
- Registration: new binary via getRegistrationKernelClient (exclude old CC)
- Dealing: old binary via OldCodeCommitment (sealed keys bound to old MRENCLAVE)
- ProcessDeals/Responses/Finalization: each member uses session.CodeCommitment
- SkipToNextRound always sets isUpgrade=false (upgrade failures require manual re-scheduling)
- KNOWN ISSUES: activationHeight uses uint32 (should be int64), genesis export missing KernelUpgradeInfos
- `shouldDeal` returns inPrevSet (comment says "current set" -- INCORRECT comment)
- `getEventSafe` pattern for pre-wiring events not yet in contract ABI
- Resharing: `shouldDeal` = prevSet validators deal, `shouldProcessResponses` = prevSet OR curSet
- See [tee-upgrade-design.md](tee-upgrade-design.md) for detailed design

## Security Notes
- EVM addr -> bech32 conversion: `sdk.AccAddress(common.HexToAddress(evmAddr).Bytes())` - deterministic but depends on valid hex
- `sort.Strings` on EIP-55 checksummed Hex() is deterministic but case-sensitive (uppercase < lowercase in ASCII)
- `amount.Uint64()` in evmstaking/ubi.go can silently overflow for amounts > uint64 max (pre-existing pattern)
- Partial send failure during reward distribution leaves module account with "stuck" funds (architectural concern)
