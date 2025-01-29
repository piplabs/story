package types

// evmstaking module event types.
const (
	EventTypeUpgradeFailure       = "upgrade_failure"
	EventTypeUpdateUbiFailure     = "update_ubi_failure"
	EventTypeCancelUpgradeFailure = "cancel_upgrade_failure"

	EventTypeUpgradeSuccess       = "upgrade_success"
	EventTypeUpdateUbiSuccess     = "update_ubi_success"
	EventTypeCancelUpgradeSuccess = "cancel_upgrade_success"

	AttributeKeyErrorCode     = "error_code"
	AttributeKeyBlockHeight   = "block_height"
	AttributeKeyTxHash        = "tx_hash"
	AttributeKeyUpgradeName   = "upgrade_name"
	AttributeKeyUpgradeHeight = "upgrade_height"
	AttributeKeyUpgradeInfo   = "upgrade_info"
	AttributeKeyUbiPercentage = "ubi_percentage"
)
