package types

// evmstaking module event types.
const (
	// Upgrade events (failure).
	EventTypeUpgradeFailure       = "upgrade_failure"
	EventTypeUpdateUbiFailure     = "update_ubi_failure"
	EventTypeCancelUpgradeFailure = "cancel_upgrade_failure"

	// Upgrade events (success).
	EventTypeUpgradeSuccess       = "upgrade_success"
	EventTypeUpdateUbiSuccess     = "update_ubi_success"
	EventTypeCancelUpgradeSuccess = "cancel_upgrade_success"

	// DKG events (success).
	EventTypeDKGInitializedSuccess                       = "dkg_initialized_success"
	EventTypeDKGCommitmentsUpdatedSuccess                = "dkg_commitments_updated_success"
	EventTypeDKGFinalizedSuccess                         = "dkg_finalized_success"
	EventTypeDKGUpgradeScheduledSuccess                  = "dkg_upgrade_scheduled_success"
	EventTypeDKGRegistrationChallengedSuccess            = "dkg_registration_challenged_success"
	EventTypeDKGInvalidDKGInitializationSuccess          = "dkg_invalid_dkg_initialization_success"
	EventTypeDKGRemoteAttestationProcessedOnChainSuccess = "dkg_remote_attestation_processed_on_chain_success"
	EventTypeDKGDealComplaintsSubmittedSuccess           = "dkg_deal_complaints_submitted_success"
	EventTypeDKGDealVerifiedSuccess                      = "dkg_deal_verified_success"
	EventTypeDKGInvalidDealSuccess                       = "dkg_invalid_deal_success"

	// DKG events (failure).
	EventTypeDKGInitializedFailure                       = "dkg_initialized_failure"
	EventTypeDKGCommitmentsUpdatedFailure                = "dkg_commitments_updated_failure"
	EventTypeDKGFinalizedFailure                         = "dkg_finalized_failure"
	EventTypeDKGUpgradeScheduledFailure                  = "dkg_upgrade_scheduled_failure"
	EventTypeDKGRegistrationChallengedFailure            = "dkg_registration_challenged_failure"
	EventTypeDKGInvalidDKGInitializationFailure          = "dkg_invalid_dkg_initialization_failure"
	EventTypeDKGRemoteAttestationProcessedOnChainFailure = "dkg_remote_attestation_processed_on_chain_failure"
	EventTypeDKGDealComplaintsSubmittedFailure           = "dkg_deal_complaints_submitted_failure"
	EventTypeDKGDealVerifiedFailure                      = "dkg_deal_verified_failure"
	EventTypeDKGInvalidDealFailure                       = "dkg_invalid_deal_failure"

	// Common attributes.
	AttributeKeyErrorCode   = "error_code"
	AttributeKeyBlockHeight = "block_height"
	AttributeKeyTxHash      = "tx_hash"
	// Upgrade attributes.
	AttributeKeyUpgradeName   = "upgrade_name"
	AttributeKeyUpgradeHeight = "upgrade_height"
	AttributeKeyUpgradeInfo   = "upgrade_info"
	AttributeKeyUbiPercentage = "ubi_percentage"
	// DKG attributes.
	AttributeKeyDKGRound            = "dkg_round"
	AttributeKeyDKGTotal            = "dkg_total"
	AttributeKeyDKGThreshold        = "dkg_threshold"
	AttributeKeyDKGIndex            = "dkg_index"
	AttributeKeyDKGCommitments      = "dkg_commitments"
	AttributeKeyDKGSignature        = "dkg_signature"
	AttributeKeyDKGMrenclave        = "dkg_mrenclave"
	AttributeKeyDKGDkgPubKey        = "dkg_dkg_pub_key"
	AttributeKeyDKGCommPubKey       = "dkg_comm_pub_key"
	AttributeKeyDKGRawQuote         = "dkg_raw_quote"
	AttributeKeyDKGFinalized        = "dkg_finalized"
	AttributeKeyDKGActivationHeight = "dkg_activation_height"
	AttributeKeyDKGChallenger       = "dkg_challenger"
	AttributeKeyDKGValidator        = "dkg_validator"
	AttributeKeyDKGComplainIndexes  = "dkg_complain_indexes"
	AttributeKeyDKGRecipientIndex   = "dkg_recipient_index"
	AttributeKeyDKGChalStatus       = "dkg_chal_status"
)
