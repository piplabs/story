// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IDKG {
    /// @notice Struct for the enclave type data unique to each enclave type
    /// @param codeCommitment The code commitment
    /// @param validationHookAddr The address of the validation hook
    struct EnclaveTypeData {
        bytes32 codeCommitment;
        address validationHookAddr;
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
    /// @param isWhitelisted Whether the enclave type is whitelisted
    event EnclaveTypeWhitelisted(
        bytes32 enclaveType,
        bytes32 codeCommitment,
        address validationHookAddr,
        bool isWhitelisted
    );

    /// @notice Emitted when an enclave instance is registered
    /// @param enclaveReport The enclave report
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param enclaveCommKey The communication key of the enclave
    /// @param dkgPubKey The DKG public key
    /// @param codeCommitment The code commitment
    /// @param validationContext The validation context
    event Registered(
        bytes enclaveReport,
        uint32 round,
        address indexed validatorAddr,
        bytes32 enclaveType,
        bytes enclaveCommKey,
        bytes dkgPubKey,
        bytes32 codeCommitment,
        bytes validationContext
    );

    /// @notice Emitted when an enclave instance is finalized
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param codeCommitment The code commitment
    /// @param participantsRoot The participants root
    /// @param globalPubKey The global public key
    /// @param publicCoeffs The public coefficients
    /// @param signature The signature
    event Finalized(
        uint32 round,
        address indexed validatorAddr,
        bytes32 enclaveType,
        bytes32 codeCommitment,
        bytes32 participantsRoot,
        bytes globalPubKey,
        bytes[] publicCoeffs,
        bytes signature
    );

    /// @notice Sets the minimum number of participants needed to be registered for each round
    /// @param newMinReqRegisteredParticipants The minimum number of participants needed to be registered for each round
    function setMinReqRegisteredParticipants(uint256 newMinReqRegisteredParticipants) external;

    /// @notice Sets the minimum number of participants needed to finish dkg for each round
    /// @param newMinReqFinalizedParticipants The minimum number of participants needed to finish dkg for each round
    function setMinReqFinalizedParticipants(uint256 newMinReqFinalizedParticipants) external;

    /// @notice Sets the operational threshold
    /// @param newOperationalThreshold The operational threshold
    function setOperationalThreshold(uint256 newOperationalThreshold) external;

    /// @notice Sets the fee paid to request DKG registration (register and finalize)
    /// @param newFee The fee paid to request DKG registration (register and finalize)
    function setFee(uint256 newFee) external;

    /// @notice Whitelists an enclave type
    /// @param enclaveType The type of the enclave
    /// @param enclaveTypeData The data of the enclave type
    /// @param isWhitelisted Whether the enclave type is whitelisted
    function whitelistEnclaveType(
        bytes32 enclaveType,
        EnclaveTypeData memory enclaveTypeData,
        bool isWhitelisted
    ) external;

    /// @notice Authenticates an enclave report
    /// @param enclaveReport The enclave report
    /// @param enclaveInstanceData The data of the enclave instance
    /// @param validationContext The validation context
    function authenticateEnclaveReport(
        bytes calldata enclaveReport,
        EnclaveInstanceData calldata enclaveInstanceData,
        bytes calldata validationContext
    ) external payable;

    /// @notice Registers an enclave instance
    /// @param enclaveReport The enclave report
    /// @param enclaveInstanceData The data of the enclave instance
    /// @param validationContext The validation context
    function register(
        bytes calldata enclaveReport,
        EnclaveInstanceData calldata enclaveInstanceData,
        bytes calldata validationContext
    ) external payable;

    /// @notice Finalizes an enclave instance
    /// @param round The round
    /// @param validatorAddr The address of the validator
    /// @param enclaveType The type of the enclave
    /// @param participantsRoot The participants root
    /// @param globalPubKey The global public key
    /// @param publicCoeffs The public coefficients
    /// @param signature The signature
    function finalize(
        uint32 round,
        address validatorAddr,
        bytes32 enclaveType,
        bytes32 participantsRoot,
        bytes calldata globalPubKey,
        bytes[] calldata publicCoeffs,
        bytes calldata signature
    ) external payable;

    /// @notice Gets the minimum number of participants needed to be registered for each round
    /// @return The minimum number of participants needed to be registered for each round
    function minReqRegisteredParticipants() external view returns (uint256);

    /// @notice Gets the minimum number of participants needed to finish dkg for each round
    /// @return The minimum number of participants needed to finish dkg for each round
    function minReqFinalizedParticipants() external view returns (uint256);

    /// @notice Gets the operational threshold
    /// @return The operational threshold
    function operationalThreshold() external view returns (uint256);

    /// @notice Gets the fee paid to request DKG registration (register and finalize)
    /// @return The fee paid to request DKG registration (register and finalize)
    function fee() external view returns (uint256);

    /// @notice Gets the enclave type data
    /// @param enclaveType The type of the enclave
    function enclaveTypeData(bytes32 enclaveType) external view returns (EnclaveTypeData memory);

    /// @notice Gets the is enclave type whitelisted
    /// @param enclaveType The type of the enclave
    function isEnclaveTypeWhitelisted(bytes32 enclaveType) external view returns (bool);
}
