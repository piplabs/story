i# `evmstaking`

## Abstract

This document specifies the internal `x/evmstaking` module of the Story blockchain.

In Story blockchain, the gas token resides on the execution layer (EL) to pay for transactions and interact with smart contracts. However, the consensus layer (CL) manages the consensus staking, slashing, and rewarding. This module exists to facilitate CL-level staking-related logic, such as delegating to validators with custom lock periods.

## Contents

1. **[State](#state)**
2. **[Two Queue System](#two-queue-system)**
3. **[Withdrawal Queue Content](#withdrawal-queue-content)**
4. **[End Block](#end-block)**
5. **[Processing Staking Events](#processing-staking-events)**
6. **[Withdrawing Delegations](#withdrawing-delegations)**
7. **[Withdrawing Rewards](#withdrawing-rewards)**
8. **[Withdrawing UBI](#withdrawing-ubi)**

## State

### Withdrawal Queue

Type: `Queue[types.Withdrawal]`

The (stake) withdrawal queue stores the pending unbonded stakes to be burned on CL and minted on EL. Stakes that are unbonded after 14 days of unstaking period are added to the queue to be processed.

### Reward Withdrawal Queue

Type: `Queue[types.Withdrawal]`

The reward withdrawal queue stores the pending rewards from stakes to be burned on CL and minted on EL. All rewards above a threshold are eligible to be queued in this queue, but there exists a parameter of maximum additions per block.

### Parameters

```protobuf
message Params {
  uint32 max_withdrawal_per_block = 1 [
    (gogoproto.moretags) = "yaml:\"max_withdrawal_per_block\""
  ];
  uint32 max_sweep_per_block = 2 [
    (gogoproto.moretags) = "yaml:\"max_sweep_per_block\""
  ];
  uint64 min_partial_withdrawal_amount = 3 [
    (gogoproto.moretags) = "yaml:\"min_partial_withdrawal_amount\""
  ];
  string ubi_withdraw_address = 4 [
    (gogoproto.moretags) = "yaml:\"ubi_withdraw_address\""
  ];
}
```

- `max_withdrawal_per_block` is the maximum number of withdrawals (reward and unstakes, each) to process per block. This parameter prevents nodes from processing a large amount of withdrawals at once, which could exceed the max chain timeout.
- `max_sweep_per_block` is the maximum number of validator-delegator delegations to sweep per block. This parameter prevents nodes from processing a large amount of delegations at once.
- `min_partial_withdrawal_amount` is the minimum amount required for rewards to get added to the reward withdrawal queue.
- `ubi_withdrawal_address` is the UBI contract address to which UBI withdrawals should be deposited.

### Delegator Withdraw Address

Type: `Map[string, string]`

The delegator-withdraw address mapping tracks the address to which a delegator receives their withdrawn stakes. The (stake) withdrawal queue uses this map to determine the `execution_address` in the `Withdrawal` struct used in building an EVM block payload.

While the delegator can change the withdraw address at any time, existing stake withdraw requests in the (stake) withdrawal queue will maintain their original values.

### Delegator Reward Address

The delegator-reward address mapping tracks the address to which a delegator receives their reward stakes, similar to the delegator-withdraw mapping.

While the delegator can change the reward address at any time, existing reward withdraw requests in the reward withdrawal queue will maintain their original values.

Type: `Map[string, string]`

### Delegator Operator Address

Type: `Map[string, string]`

The delegator-operator address mapping tracks the address to which a delegator has given the privilege to delegate (stake), undelegate (unstake), and redelegate on behalf of themselves.

### IP Token Staking Contract

Type: `*bindings.IPTokenStaking`

IPTokenStaking contract is used to filter and parse staking-related events from EL.

## Two Queue System

The module departs from traditional Cosmos SDK staking module's unstaking system, where all unbonded entries (stakes that have unbonded after 14 days of unbonding period) are immediately distributed into delegators account. Instead, Story's unstaking system assimilates Ethereum 2.0's unstaking system, where 16 full or partial (reward) withdrawals are processed per slot.

In a single queue of withdrawals, reward withdrawals can significantly delay stake withdrawals. Hence, Story blockchain implements a two-queue system where a max amount to process per block is enforced per queue. In other words, the stake/ubi withdrawal and reward withdrawal queues can each process the max parameter per block.

## Withdrawal Queue Content

Since the module only processes unstakes/rewards/ubi and stores them in queues, the actual dequeueing for withdrawal to the execution layer is carried out in the evmengine module. More specifically, a proposer dequeues the max number of withdrawals from each queue and adds them to the EVM block payload, which gets executed by EL via the Engine API. When validators receive proposed block payload from the proposer, they individually peek the local queues and compare them against the received block's withdrawals. Mismatching withdrawals indicate non-determinism in staking logics and should result in chain halt.

In other words, the `evmstaking` module is in charge of parsing, processing, and inserting withdrawal requests to two queues, while the `evmengine` module is in charge of validating and dequeuing withdrawal requests, as well as depositing them to corresponding withdrawal addresses in EL.

## End Block

The `EndBlock` ABCI2 call is responsible for fetching the unbonded entries (stakes that have unbonded after 14 days) from the staking module and inserting them into the (stake) withdrawal queue. Furthermore, it processes stake reward withdrawals into the reward withdrawal queue and UBI withdrawals into the (stake) withdrawal queue.

If the network is in the [Singularity period](https://docs.story.foundation/docs/tokenomics-staking#singularity), the End Block is skipped as there are no staking rewards and withdrawals available during this period. Otherwise, refer to [Withdrawing Delegations](#withdrawing-delegations), [Withdrawing Rewards](#withdrawing-rewards), and [Withdrawing UBI](#withdrawing-ubi) for detailed withdrawal processes.

## Processing Staking Events

The module parses and processes staking events emitted from the [IPTokenStaking contract](https://github.com/piplabs/story/blob/main/contracts/src/protocol/IPTokenStaking.sol), which are collected by the evmengine module. The list of events are:

### Staking events

- Create Validator
- Deposit (delegate)
- Withdraw (undelegate)
- Redelegate
- Unjail: anyone can request to unjail a jailed validator by paying the unjail fee in the contract.

These operations incur a fixed gas cost to prevent spam.

### Parameter events

- Update Validator Commission: update the validator commission.
- Set Withdrawal Address: delegator can modify their withdrawal address for future unstakes/undelegations.
- Set Reward Address: delegator can modify their withdrawl address for future reward emissions.
- Set Operator: delegator can modify their operator with privileges of delegation, undelegation, and redelgation.
- Unset Operator: delegator can remove operator.

These operations incur a fixed gas cost to prevent spam.

## Withdrawing Delegations

## Withdrawing Rewards

## Withdrawing UBI
