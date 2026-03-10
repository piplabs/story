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
	EventTypeDKGInitializedSuccess                     = "dkg_initialized_success"
	EventTypeDKGFinalizedSuccess                       = "dkg_finalized_success"
	EventTypeDKGUpgradeScheduledSuccess                = "dkg_upgrade_scheduled_success"
	EventTypeDKGUpgradeCanceledSuccess                 = "dkg_upgrade_canceled_success"
	EventTypeDKGMinReqRegisteredParticipantsSetSuccess = "dkg_min_req_registered_participants_set_success"
	EventTypeDKGMinReqFinalizedParticipantsSetSuccess  = "dkg_min_req_finalized_participants_set_success"
	EventTypeDKGOperationalThresholdSetSuccess         = "dkg_operational_threshold_set_success"
	EventTypeDKGThresholdDecryptRequestedSuccess       = "dkg_threshold_decrypt_requested_success"
	EventTypeDKGPartialDecryptionSubmittedSuccess      = "dkg_partial_decryption_submitted_success"

	// DKG events (failure).
	EventTypeDKGInitializedFailure                     = "dkg_initialized_failure"
	EventTypeDKGFinalizedFailure                       = "dkg_finalized_failure"
	EventTypeDKGUpgradeScheduledFailure                = "dkg_upgrade_scheduled_failure"
	EventTypeDKGUpgradeCanceledFailure                 = "dkg_upgrade_canceled_failure"
	EventTypeDKGMinReqRegisteredParticipantsSetFailure = "dkg_min_req_registered_participants_set_failure"
	EventTypeDKGMinReqFinalizedParticipantsSetFailure  = "dkg_min_req_finalized_participants_set_failure"
	EventTypeDKGOperationalThresholdSetFailure         = "dkg_operational_threshold_set_failure"
	EventTypeDKGThresholdDecryptRequestedFailure       = "dkg_threshold_decrypt_requested_failure"
	EventTypeDKGPartialDecryptionSubmittedFailure      = "dkg_partial_decryption_submitted_failure"

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
	AttributeKeyDKGRound                        = "dkg_round"
	AttributeKeyDKGTotal                        = "dkg_total"
	AttributeKeyDKGThreshold                    = "dkg_threshold"
	AttributeKeyDKGIndex                        = "dkg_index"
	AttributeKeyDKGCommitments                  = "dkg_commitments"
	AttributeKeyDKGSignature                    = "dkg_signature"
	AttributeKeyDKGCodeCommitment               = "dkg_code_commitment"
	AttributeKeyDKGDkgPubKey                    = "dkg_dkg_pub_key"
	AttributeKeyDKGCommPubKey                   = "dkg_comm_pub_key"
	AttributeKeyDKGEnclaveReport                = "dkg_enclave_report"
	AttributeKeyDKGStartBlockHeight             = "dkg_start_block_height"
	AttributeKeyDKGStartBlockHash               = "dkg_start_block_hash"
	AttributeKeyDKGActivationHeight             = "dkg_activation_height"
	AttributeKeyDKGUpgradeVersion               = "dkg_upgrade_version"
	AttributeKeyDKGParticipantsRoot             = "dkg_participants_root"
	AttributeKeyDKGChallenger                   = "dkg_challenger"
	AttributeKeyDKGValidator                    = "dkg_validator"
	AttributeKeyDKGComplainIndexes              = "dkg_complain_indexes"
	AttributeKeyDKGRecipientIndex               = "dkg_recipient_index"
	AttributeKeyDKGChalStatus                   = "dkg_chal_status"
	AttributeKeyDKGRequester                    = "dkg_requester"
	AttributeKeyDKGRequesterPubKeyLen           = "dkg_requester_pub_key_len"
	AttributeKeyDKGCiphertextLen                = "dkg_ciphertext_len"
	AttributeKeyDKGLabelLen                     = "dkg_label_len"
	AttributeKeyDKGPid                          = "dkg_pid"
	AttributeKeyDKGEncryptedPartLen             = "dkg_encrypted_part_len"
	AttributeKeyDKGEphemeralKeyLen              = "dkg_ephemeral_key_len"
	AttributeKeyDKGPubShareLen                  = "dkg_pub_share_len"
	AttributeKeyDKGMinReqRegisteredParticipants = "dkg_min_req_registered_participants"
	AttributeKeyDKGMinReqFinalizedParticipants  = "dkg_min_req_finalized_participants"
	AttributeKeyDKGOperationalThreshold         = "dkg_operational_threshold"
)
