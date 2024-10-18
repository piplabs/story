package types

// evmstaking module event types.
const (
	EventTypeDelegateFailure  = "delegate_failure"
	EventTypeAddWithdrawal    = "add_withdrawal"
	EventTypeRemoveWithdrawal = "remove_withdrawal"

	AttributeKeyValidator        = "validator"
	AttributeKeyDelegator        = "delegator"
	AttributeKeyExecutionAddress = "execution_address"
	AttributeKeyCreationHeight   = "creation_height"
)
