package types

// evmstaking module event types.
const (
	EventTypeSetWithdrawalAddressFailure = "set_withdrawal_address_failure"
	EventTypeSetRewardAddressFailure     = "set_reward_address_failure"
	EventTypeAddOperatorFailure          = "add_operator_failure"
	EventTypeRemoveOperatorFailure       = "remove_operator_failure"
	EventTypeCreateValidatorFailure      = "create_validator_failure"
	EventTypeDelegateFailure             = "delegate_failure"
	EventTypeRedelegateFailure           = "redelegate_failure"
	EventTypeUndelegateFailure           = "undelegate_failure"
	EventTypeUnjailFailure               = "unjail_failure"

	AttributeKeyBlockHeight             = "block_height"
	AttributeKeyDelegatorUncmpPubKey    = "delegator_uncmp_pubkey"
	AttributeKeyValidatorUncmpPubKey    = "validator_uncmp_pubkey"
	AttributeKeySrcValidatorUncmpPubKey = "src_validator_uncmp_pubkey"
	AttributeKeyDstValidatorUncmpPubKey = "dst_validator_uncmp_pubkey"
	AttributeKeyDelegateID              = "delegation_id"
	AttributeKeyPeriodType              = "staking_period"
	AttributeKeyAmount                  = "amount"
	AttributeKeySenderAddress           = "sender_address"
	AttributeKeyWithdrawalAddress       = "withdrawal_address"
	AttributeKeyRewardAddress           = "reward_address"
	AttributeKeyOperatorAddress         = "operator_address"
	AttributeKeyMoniker                 = "moniker"
	AttributeKeyCommissionRate          = "commission_rate"
	AttributeKeyMaxCommissionRate       = "max_commission_rate"
	AttributeKeyMaxCommissionChangeRate = "max_commission_change_rate"
	AttributeKeyTokenType               = "token_type"
)
