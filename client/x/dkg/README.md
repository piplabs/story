# DKG Module

The DKG (Distributed Key Generation) module implements a TEE-based distributed key generation and threshold decryption
protocol for the Story blockchain. It coordinates validators running TEE services (`story-kernel`) to collectively
generate cryptographic key shares, perform periodic resharing, and provide TDH2 (Threshold Decryption based on Hybrid
encryption) partial decryption services.

## Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                         Story Consensus Client                       │
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────────────┐  │
│  │  EVMEngine   │────▶│  DKG Keeper  │────▶│   Kernel Router      │  │
│  │ (EVM events) │     │ (on-chain    │     │ (CC-based routing)   │  │
│  └──────────────┘     │  state mgmt) │     └────┬───────────┬─────┘  │
│                       └──────┬───────┘          │           │        │
│                              │                  │           │        │
│  ┌──────────────┐     ┌──────▼───────┐   ┌──────▼───┐ ┌────▼─────┐   │
│  │ Vote Ext.    │◀────│  DKG Service │   │ Kernel A │ │ Kernel B │   │
│  │ (CometBFT)   │     │ (off-chain   │   │ (gRPC)   │ │ (gRPC)   │   │
│  └──────────────┘     │  TEE calls)  │   └──────────┘ └──────────┘   │
│                       └──────────────┘                               │
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐                               │
│  │  DKG.sol     │     │ State Mgr    │                               │
│  │ (EVM contract│     │ (local JSON  │                               │
│  │  predeploy)  │     │  sessions)   │                               │
│  └──────────────┘     └──────────────┘                               │
└──────────────────────────────────────────────────────────────────────┘
```

**Key components:**

- **EVMEngine**: Listens for DKG.sol contract events (`Registered`, `Finalized`, `UpgradeScheduled`, etc.) and forwards them to the DKG Keeper.
- **DKG Keeper**: Manages on-chain state (DKG rounds, registrations, votes, parameters) and orchestrates stage transitions via `BeginBlocker`.
- **DKG Service**: Off-chain service layer that communicates with `story-kernel` (TEE) via gRPC for key generation, dealing, and finalization.
- **Kernel Router**: Routes gRPC calls to the correct `story-kernel` binary based on code commitment. Supports multiple simultaneous binaries during upgrades.
- **State Manager**: Persists local `DKGSession` state to disk as JSON files, enabling recovery across restarts.
- **Vote Extensions**: Deals, responses, and justifications are propagated between validators using CometBFT vote extensions.

## DKG Lifecycle

Each DKG round progresses through the following on-chain stages:

```
                    ┌─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─┐
                    │           ONE COMPLETE DKG ROUND          │

  ┌──────────────┐   ┌─────────┐   ┌──────────────┐   ┌──────────┐
  │ Registration │──▶│ Dealing │──▶│ Finalization │──▶│  Active  │
  │              │   │         │   │              │   │          │
  │ Validators   │   │ VSS deal│   │ Signature    │   │ TDH2     │
  │ register     │   │ exchange│   │ verification │   │ service  │
  │ via DKG.sol  │   │ via Vote│   │ + pub key    │   │ running  │
  │              │   │ Ext.    │   │ consensus    │   │          │
  └──────────────┘   └─────────┘   └──────────────┘   └──────────┘
       │                 │                │                 │
       │  period ends    │  period ends   │  period ends    │  period ends
       ▼                 ▼                ▼                 ▼
  (next stage)      (next stage)    (next stage)     (new round:
                                                      resharing)
                    └─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─┘
```

### Stage Details

| Stage | Duration (default) | What Happens |
|-------|--------------------|--------------|
| **Registration** | 86,400 blocks (~1 day) | Validators call `story-kernel` to generate and seal a DKG key pair, then submit a `register` transaction to DKG.sol. The TEE attestation report is verified on-chain. |
| **Dealing** | 86,400 blocks (~1 day) | Requires `min_req_registered_participants` (default: 3) verified registrations. Validators generate VSS deals via TEE, exchange them through CometBFT vote extensions, and process received deals to produce responses. |
| **Finalization** | 86,400 blocks (~1 day) | Validators call TEE to finalize the DKG, producing a global public key, public coefficients, and a signature. Results are submitted to DKG.sol. A quorum vote determines the accepted global public key. |
| **Active** | 1,814,400 blocks (~21 days) | The DKG committee is live. TDH2 threshold decryption requests are processed. UBI rewards are distributed to committee members. When the period ends, a new resharing round begins automatically. |

### Upgrade Resharing Priority

> **Note:** During upgrade periods, the upgrade resharing takes priority over the active period lifecycle. If an upgrade
> is scheduled and the activation height is reached, the DKG module will initiate a resharing round regardless of the
> current round's remaining active period.

### Failure Handling

If a round fails to meet minimum participation thresholds when transitioning from Registration to Dealing, or from
Finalization to Active, the round is marked as `Failed` and a new round is initiated immediately. The previous active
round remains in service until a new round successfully reaches the Active stage.

## Module Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `registration_period` | 86,400 | Number of blocks for the registration stage |
| `dealing_period` | 86,400 | Number of blocks for the dealing stage |
| `finalization_period` | 86,400 | Number of blocks for the finalization stage |
| `active_period` | 1,814,400 | Number of blocks for the active stage |
| `dkg_committee_reward_portion` | `"0.10"` (10%) | Fraction of UBI allocated to the DKG committee |
| `min_req_registered_participants` | 3 | Minimum verified registrations required to enter dealing |
| `min_req_finalized_participants` | 3 | Minimum finalized participants required to enter active stage |
| `operational_threshold` | 667 (basis 1000) | Quorum threshold for global public key consensus (66.7%) |

The operational threshold determines the DKG quorum: `threshold = ceil(total * operationalThreshold / 1000)`.

## Configuration

Validator operators configure the DKG service via CLI flags or `config/dkg.toml`:

```toml
[dkg]
enable = true                                          # Enable DKG service
kernel_endpoints = ["127.0.0.1:50051"]                 # story-kernel gRPC endpoints
engine_rpc_endpoint = "http://127.0.0.1:8545"          # Execution layer JSON-RPC
enclave_type = 1                                       # TEE enclave type (1 = SGX)
```

| CLI Flag | Description |
|----------|-------------|
| `--dkg-enable` | Enable the DKG service |
| `--dkg-kernel-endpoints` | Comma-separated list of story-kernel gRPC endpoints |
| `--dkg-engine-rpc-endpoint` | Execution layer JSON-RPC endpoint |
| `--dkg-enc-type` | TEE enclave type identifier |

## On-Chain State

The module stores the following data in the Cosmos SDK KV store:

| Collection | Key Format | Value | Description |
|------------|-----------|-------|-------------|
| `Params` | prefix(0) | `Params` | Module parameters |
| `DKGNetworks` | prefix(1) + `round` | `DKGNetwork` | Per-round network state (stage, threshold, global key, etc.) |
| `LatestDKGNetwork` | prefix(2) | `string` | Pointer to the latest round |
| `DKGRegistrations` | prefix(3) + `round_address` | `DKGRegistration` | Per-validator registration record |
| `LatestActiveRound` | prefix(4) | `string` | Pointer to the latest active round |
| `GlobalPubKeyVotes` | prefix(5) + `round_pubkey_coeffhash` | `uint32` | Vote count for global public key consensus |
| `KernelUpgradeInfos` | prefix(6) + `upgradeVersion` | `KernelUpgradeInfo` | Pending kernel upgrade schedules |
| `SettlementBalance` | prefix(7) | `string` | Accumulated UBI balance for committee rewards |

## Contract Integration (DKG.sol)

The DKG module interacts with a pre-deployed EVM contract (`DKG.sol`) at a predefined address. All validator actions
(register, finalize, submit partial decryption) are executed as EVM transactions. The contract emits events that the
consensus layer processes:

| Contract Event | CL Handler | Description |
|----------------|------------|-------------|
| `Registered` | `Registered()` | Stores verified registration with attestation |
| `Finalized` | `Finalized()` | Verifies signature, tallies global key votes |
| `MinReqRegisteredParticipantsSet` | `SetMinReqRegisteredParticipants()` | Updates parameter |
| `MinReqFinalizedParticipantsSet` | `SetMinReqFinalizedParticipants()` | Updates parameter |
| `OperationalThresholdSet` | `SetOperationalThreshold()` | Updates parameter |
| `UpgradeScheduled` | `UpgradeScheduled()` | Stores pending kernel upgrade info |
| `UpgradeCancelled` | `UpgradeCancelled()` | Deletes pending kernel upgrade info |

## TDH2 Threshold Decryption

Once a DKG round reaches the Active stage, the committee can process threshold decryption requests:

1. An external contract emits a `ThresholdDecryptRequested` event with ciphertext and label.
2. The DKG handler queues the request in the local session.
3. A background worker (3-second polling interval) processes queued requests.
4. Each validator calls `story-kernel`'s `PartialDecryptTDH2` RPC to produce a partial decryption share.
5. The partial decryption is submitted on-chain via `DKG.submitPartialDecryption()`.
6. Once enough partial decryptions are collected (meeting the threshold), the ciphertext can be fully decrypted.

## Vote Extensions

The DKG module uses CometBFT vote extensions to propagate dealing data between validators without requiring on-chain transactions for every deal:

- **`ExtendVote`**: Dequeues up to `maxItemsPerVote` (80) deals, responses, and justifications from the in-memory queue and serializes them into the vote extension payload.
- **`VerifyVoteExtension`**: Validates the structure and enforces size/count limits (256 KB max, 80 items per type) on received vote extension payloads.
- **`PrepareVotes`**: Aggregates votes from `ExtendedCommitInfo`, deduplicates deals/responses/justifications, and produces a `MsgAddDkgVote` transaction included in the block.

## Kernel Upgrade (TEE Binary Upgrade)

SGX enclaves seal key shares using `SealWithUniqueKey` (MRENCLAVE-based), which means a new binary cannot unseal data
sealed by the old one. The kernel upgrade mechanism uses **proactive resharing** to redistribute key shares from the old
binary to the new one.

### How It Works

```
┌──────────────┐       ┌──────────────┐       ┌──────────────────────┐
│   DKG.sol    │       │  EVMEngine   │       │     DKG Keeper       │
│ (EVM Layer)  │──────▶│  (CL Layer)  │──────▶│     BeginBlocker     │
└──────────────┘       └──────────────┘       └──────────┬───────────┘
  scheduleUpgrade()      UpgradeScheduled         checkPendingUpgrade
  cancelUpgrade()        UpgradeCancelled                │
                                               At activation height:
                                               InitiateDKGRound(isUpgrade=true)
                                                         │
                                            ┌────────────▼────────────┐
                                            │     Kernel Router       │
                                            └─────┬──────────────┬────┘
                                                  │              │
                                         ┌────────▼───┐   ┌──────▼────────┐
                                         │ Old Kernel │   │  New Kernel   │
                                         │ dealer role│   │  recipient    │
                                         └────────────┘   └───────────────┘
```

During an upgrade round:
- **Registration**: New members generate keys using the **new** kernel binary.
- **Dealing**: Dealers use the **old** kernel binary (sealed key shares are bound to old code commitment).
- **ProcessDeals/Responses**: Each member uses their respective kernel binary based on `session.CodeCommitment`.
- **Finalization**: New members finalize using the **new** kernel binary.

### Operator Guide

#### Step 1: Deploy the New Kernel Binary

Start the new `story-kernel` binary on a different port while keeping the old one running:

```bash
# Old binary (already running)
./story-kernel --port 50051

# New binary (start alongside)
./story-kernel-v2 --port 50052
```

#### Step 2: Update Consensus Client Configuration

Add both endpoints to the configuration and restart the consensus client:

```toml
[dkg]
enable = true
kernel_endpoints = ["127.0.0.1:50051", "127.0.0.1:50052"]
```

The client calls `GetCodeCommitment` on each endpoint to discover code commitments (MRENCLAVE values for SGX):

```
INFO Connected to kernel endpoint  endpoint=127.0.0.1:50051 code_commitment=<old_cc_hex>
INFO Connected to kernel endpoint  endpoint=127.0.0.1:50052 code_commitment=<new_cc_hex>
```

#### Step 3: Whitelist the New Enclave Type (if needed)

If the new binary has a different code commitment, the DKG contract owner must whitelist it:

```solidity
DKG.whitelistEnclaveType(
    enclaveType,
    EnclaveTypeData({ codeCommitment: <new_cc>, validationHookAddr: <hook_addr> }),
    true
);
```

#### Step 4: Schedule the Upgrade

The DKG contract owner schedules the upgrade:

```solidity
DKG.scheduleUpgrade(activationHeight, "v1.6.0");
```

- `activationHeight` must be in the future.
- Only one pending upgrade is allowed at a time. Cancel first if one already exists.

#### Step 5: Wait for Activation

At `activationHeight`, the consensus layer automatically initiates an upgrade resharing round:

```
INFO Kernel upgrade activated, initiating upgrade resharing round  upgrade_version=v1.6.0
```

#### Step 6: Clean Up

After the upgrade round reaches Active:

1. Stop the old kernel binary.
2. Remove the old endpoint from `kernel_endpoints`.
3. Optionally restart the consensus client.

#### Cancelling a Scheduled Upgrade

Cancel before the activation height, then re-schedule:

```solidity
DKG.cancelUpgrade();
DKG.scheduleUpgrade(newActivationHeight, "v1.6.0");
```

### Troubleshooting

| Error | Cause | Resolution |
|-------|-------|------------|
| `no NEW kernel client found for upgrade registration` | New kernel binary not connected | Ensure new binary is running and listed in `kernel_endpoints`, then restart the consensus client |
| `pending upgrade already exists` | Duplicate scheduling | Call `DKG.cancelUpgrade()` first, then `DKG.scheduleUpgrade()` |
| `no kernel client for code commitment` | Required kernel binary not running | Start the required binary and update `kernel_endpoints` |
| Upgrade round stuck in registration | Validators haven't deployed the new binary | All active validators must deploy the new binary and update their config |

## Directory Structure

```
client/x/dkg/
├── keeper/
│   ├── abci.go                    # BeginBlocker, upgrade activation
│   ├── contract_client.go         # EVM DKG.sol interaction
│   ├── dkg_dealing.go             # Dealing stage orchestration
│   ├── dkg_events.go              # SDK event emission
│   ├── dkg_finalization.go        # Finalization + round completion
│   ├── dkg_handler.go             # Contract event handlers
│   ├── dkg_initialization.go      # Round creation, validator selection
│   ├── dkg_network.go             # DKGNetwork CRUD
│   ├── dkg_registration.go        # Registration CRUD
│   ├── dkg_rewards.go             # UBI committee reward distribution
│   ├── dkg_round.go               # Stage transition timing
│   ├── dkg_svc.go                 # Service orchestrator, decrypt worker
│   ├── dkg_svc_dealing.go         # TEE dealing/processing calls
│   ├── dkg_svc_finalization.go    # TEE finalization calls
│   ├── dkg_svc_registration.go    # TEE key generation + routing
│   ├── dkg_votes.go               # Global pub key vote tallying
│   ├── genesis.go                 # InitGenesis / ExportGenesis
│   ├── keeper.go                  # Keeper struct, collection setup
│   ├── kernel_router.go           # Code commitment-based client routing
│   ├── kernel_upgrade.go          # Upgrade info CRUD
│   ├── msg_server.go              # MsgServer (AddVote)
│   ├── params.go                  # Parameter accessors
│   ├── proposal_server.go         # Block proposal vote aggregation
│   ├── query.go                   # gRPC query handlers
│   ├── queue.go                   # Thread-safe deal/response queue
│   ├── state_manager.go           # Local session persistence (JSON)
│   └── tee_client.go              # gRPC client factory
├── module/
│   ├── depinject.go               # Dependency injection wiring
│   └── module.go                  # AppModule registration
├── types/
│   ├── dkg.go                     # DKGSession, DKGPhase (off-chain state)
│   ├── expected_keepers.go        # External keeper interfaces
│   ├── keys.go                    # Store keys and prefixes
│   └── params.go                  # Parameter defaults and validation
└── testutil/
    ├── expected_keepers_mocks.go  # Generated mocks
    └── tee_client_mock.go         # TEEClient mock
```
