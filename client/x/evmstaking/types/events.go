package types

// evmstaking module event types.
const (
	EventTypeDelegateFailure  = "delegate_failure"
	EventTypeAddWithdrawal    = "add_withdrawal"
	EventTypeRemoveWithdrawal = "remove_withdrawal"

	AttributeKeyDelegatorUncmpPubKey = "delegator"
	AttributeKeyValidatorUncmpPubKey = "validator"
	AttributeKeyAmount               = "amount"
	AttributeKeyPeriodType           = "staking_period"
	AttributeKeyDelegateID           = "delegation_id"
	AttributeKeyOperatorAddress      = "operator_address"
)
