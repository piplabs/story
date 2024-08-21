package types

// evmstaking module event types.
const (
	EventTypeAddWithdrawal    = "add_withdrawal"
	EventTypeRemoveWithdrawal = "remove_withdrawal"

	AttributeKeyValidator        = "validator"
	AttributeKeyDelegator        = "delegator"
	AttributeKeyExecutionAddress = "execution_address"
	AttributeKeyCreationHeight   = "creation_height"
)
