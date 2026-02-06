// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IDKG {
    /// @notice Struct for the enclave type data unique to each enclave type
    /// @param codeCommitment The code commitment
    /// @param validationHookAddr The address of the validation hook
    /// @param enclaveAddr The address of the enclave
    struct EnclaveTypeData {
        bytes32 codeCommitment;
        address validationHookAddr;
        address enclaveAddr;
    }

    /// @notice Struct for the enclave instance data unique to each instance
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param enclaveCommKey The communication key of the enclave
    /// @param dkgPubKey The DKG public key
    struct EnclaveInstanceData {
        uint32 round;
        address validatorAddr;
        bytes32 enclaveType;
        bytes enclaveCommKey;
        bytes dkgPubKey;
    }

    /// @notice Emitted when the minimum required registered participants is set
    /// @param newMinReqRegisteredParticipants The new minimum required registered participants
    event MinReqRegisteredParticipantsSet(uint256 newMinReqRegisteredParticipants);

    /// @notice Emitted when the minimum required finalized participants is set
    /// @param newMinReqFinalizedParticipants The new minimum required finalized participants
    event MinReqFinalizedParticipantsSet(uint256 newMinReqFinalizedParticipants);

    /// @notice Emitted when the operational threshold is set
    /// @param newOperationalThreshold The new operational threshold
    event OperationalThresholdSet(uint256 newOperationalThreshold);

    /// @notice Emitted when the fee is set
    /// @param newFee The new fee
    event FeeSet(uint256 newFee);

    /// @notice Emitted when an enclave type is whitelisted
    /// @param enclaveType The type of the enclave
    /// @param codeCommitment The code commitment
    /// @param validationHookAddr The address of the validation hook
    /// @param enclaveAddr The address of the enclave
    /// @param isWhitelisted Whether the enclave type is whitelisted
    event EnclaveTypeWhitelisted(
        bytes32 enclaveType,
        bytes32 codeCommitment,
        address validationHookAddr,
        address enclaveAddr,
        bool isWhitelisted
    );

    /// @notice Emitted when an enclave instance is registered
    /// @param enclaveReport The enclave report
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param enclaveCommKey The communication key of the enclave
    /// @param dkgPubKey The DKG public key
    /// @param validationContext The validation context
    event Registered(
        bytes enclaveReport,
        uint32 round,
        address indexed validatorAddr,
        bytes32 enclaveType,
        bytes enclaveCommKey,
        bytes dkgPubKey,
        bytes validationContext
    );

    /// @notice Emitted when an enclave instance is finalized
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param enclaveCommKey The communication key of the enclave
    /// @param dkgPubKey The DKG public key
    /// @param valSetRoot The value set root
    /// @param globalPubKey The global public key
    /// @param publicCoeffs The public coefficients
    /// @param signature The signature
    event Finalized(
        uint32 round,
        address indexed validatorAddr,
        bytes32 enclaveType,
        bytes enclaveCommKey,
        bytes dkgPubKey,
        bytes32 valSetRoot,
        bytes globalPubKey,
        bytes[] publicCoeffs,
        bytes signature
	);
}
