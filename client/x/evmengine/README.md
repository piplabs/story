# `evmengine`

## Abstract

This document specifies the internal `x/evmengine` module of the Story blockchain.

As Story Network separates the consensus and execution client, like Ethereum, the consensus client (CL) and execution client (EL) needs to communicate to sync to the network, propose proper EVM blocks, and execute EVM-triggered EL actions in CL.

The module exists to facilitate all communications between CL and EL using the Engine API, from staking and upgrades to driving block production and consensus in CL and EL.

## Contents

1. **[State](#state)**
2. **[Prepare Proposal](#prepare-proposal)**
3. **[Process Proposal](#process-proposal)**
4. **[Post Finalize](#post-finalize)**
5. **[Messages](#messages)**
6. **[UBI](#ubi)**
7. **[Upgrades](#upgrades)**

## State

### Build Delay

Type: `time.Duration`

Build delay determines the wait duration from the start of `PrepareProposal` ABCI2 call before fetching the next EVM block data to propose from EL via the Engine API. Applicable to the current proposer only. If the node has a block optimistically built beforehand, the build delay is not used.

### Build Optimistic

Type: `bool`

Enable optimistic building of a block if true. A node will deterministically build the next block if it finds itself as the next proposer in the current block. Optimistic building starts with requesting the next EVM block data (for the next CL block) immediately after the `FinalizeBlock` of ABCI2.

### Head Table

Type: `ExecutionHeadTable`

Head table stores the latest execution head data to be used for partial validation of EVM blocks received from other validators. When the chain initializes, the execution head is populated with the genesis execution hash loaded from `genesis.json`.

The following execution head is stored in the table.

```protobuf
message ExecutionHead {
  option (cosmos.orm.v1.table) = {
    id: 1;
    primary_key: { fields: "id", auto_increment: true }
  };

  uint64 id               = 1; // Auto-incremented ID (always and only 1).
  uint64 created_height   = 2; // Consensus chain height this execution block was created in.
  uint64 block_height     = 3; // Execution block height.
  bytes  block_hash       = 4; // Execution block hash.
  uint64 block_time       = 5; // Execution block time.
}
```

### Upgrade Contract

Type: `*bindings.UpgradeEntrypoint`

Upgrade contract is used to filter and parse upgrade-related events from EL.

### UBI Contract

Type: `*bindings.UBIPool`

UBI contract is used to filter and parse UBI-related events from EL.

### Mutable Payload

Type: struct

Mutable payload stores the optimistic block built, if optimistic building is enabled.

#### Genesis State

The module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height.

```protobuf
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
}

message Params {
  bytes execution_block_hash = 1 [
    (gogoproto.moretags) = "yaml:\"execution_block_hash\""
  ];
}
```

## Prepare Proposal

At each block, if the node is the proposer, ABCI2 triggers `PrepareProposal` which

1. Loads staking & reward withdrawals from the evmstaking module.
2. Builds a valid EVM block.
    - If optimistic building: loads the optimistically built block.
    - Non-optimistic: requests and retrieves an EVM block from EL.
3. Collects the EVM logs of the previous/parent block.
4. Assembles `MsgExecutionPayload` with the built EVM block and previous EVM logs.
5. Returns a transaction containing the assembled `MsgExecutionPayload` data.

This CL block is then propagated to all other validators.

## Process Proposal

At each block, if the node is not a proposer but a validator, ABCI2 triggers `ProcessProposal` with received commits (which should be a transaction of `MsgExecutionPayload` data in the honest case).

The node first validates that the received commit has only one transaction with at least 2/3 of votes committed. Then, the node validates that the one transaction only contains one unmarshalled `MsgExecutionPayload` data. Finally, the node processes the received data and broadcasts its acceptance of the proposal to the network. If any of the validation or processing fails, the node rejects the proposal.

More specifically, the node processes the received `MsgExecutionPayload` data in the following manner:

1. Validates the fields of the received `MsgExecutionPayload` (outlined in [Messages](#msgexecutionpayload)).
2. Compare local stake & reward withdrawals with the received withdrawals data.
3. Push the received execution payload to EL via the Engine API and wait for payload validation.
4. Update the EL forkchoice to the execution payload's block hash.
5. Process staking events using the evmstaking module.
6. Process upgrade events.
7. Update the execution head to the execution payload (finalized block).

## Post Finalize

If optimistic building is enabled, `PostFinalize` is triggered immediately after `FinalizeBlock` set through custom ABCI callback. During this process, the node peeks the staking and reward queues from the evmstaking module, and builds a new execution payload on top of the current execution head. It sets the optimistic block to be used in the next block's `PrepareProposal` phase and returns the response from the forkchoice update.

## Messages

In this section we describe the processing of the evmengine messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the state section.

### MsgExecutionPayload

```protobuf
message MsgExecutionPayload {
  option (cosmos.msg.v1.signer) = "authority";
  string            authority           = 1;
  bytes             execution_payload   = 2;
  repeated EVMEvent prev_payload_events = 3;
}

message EVMEvent {
  bytes          address = 1;
  repeated bytes topics  = 2;
  bytes          data    = 3;
  bytes          tx_hash = 4;
}
```

This message is expected to fail if:

- authority is invalid (not evmengine authority)
- execution payload fails to unmarshal to [ExecutableData](https://github.com/piplabs/story/blob/c38b80c13579d3df7174ea10c3368ef0692f52da/client/x/evmengine/types/executable_data.go#L17-L35) for reasons such as invalid fields
- execution payload's block number does not match CL head's block number + 1
- execution payload's block parent hash does not match CL head's hash
- execution payload's timestamp is invalid
- execution payload's RANDAO does not match CL head's hash (ie. parent hash)
- execution payload's `Withdrawals`, `BlobGasUsed`, and `ExcessBlobGas` fields are nil
- execution payload's `Withdrawals` count does not match local node's sum of dequeued stake & reward withdrawals

The message must contain previous block's events, which gets processed at the current CL block (in other words, execution events from EL block n-1 are processed at CL block n). In the future, the message will remove `prev_payload_events` and rely on Engine API to get the current finalized EL block's events.

Also note that EVM events are processed in CL in the order they are generated in EL.

## UBI

All UBI-related changes must be triggered from the canonical UBI contract in the EVM execution layer. This module handles the execution handling of those triggers in CL. Read more about [UBI for validators](https://docs.story.foundation/docs/tokenomics-staking#ubi-for-validators)

### Set UBI Distribution

The `UBIPool` contract emits the UBI distribution set event, which is parsed by the module to set the UBI percentage in the distribution module.

## Upgrades

All chain upgrade-related logics must be triggered from the canonical upgrade contract in the EVM execution layer. This module handles the execution handling of those triggers in CL.

### Software Upgrade

The `UpgradeEntrypoint` contract emits the software upgrade event, which is parsed by the module to schedule an upgrade at a given height for a given binary name. Currently, all upgrades must either be set via forks or by the software upgrade events; the latter process is a multisig-controlled process, which will transition into a voting-based process in the future.

### Cancel Upgrade

Similar to the software upgrade, the module processes the cancel upgrade event from EVM logs of the previous block, and clears an existing upgrade plan.
