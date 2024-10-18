package types

// evmstaking module event types.
const (
	EventTypeDelegateFailure = "delegate_failure"

	AttributeKeyBlockHeight          = "block_height"
	AttributeKeyDelegatorUncmpPubKey = "delegator_uncmp_pubkey"
	AttributeKeyValidatorUncmpPubKey = "validator_uncmp_pubkey"
	AttributeKeyAmount               = "amount"
	AttributeKeyPeriodType           = "staking_period"
	AttributeKeyDelegateID           = "delegation_id"
	AttributeKeyOperatorAddress      = "operator_address"
)
