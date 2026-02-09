// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { AccessManagedUpgradeable } from "@openzeppelin/contracts-upgradeable/access/manager/AccessManagedUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { IDKG } from "../interfaces/IDKG.sol";
import { IAttestationReportValidator } from "../interfaces/IAttestationReportValidator.sol";

contract DKG is IDKG, PausableUpgradeable, AccessManagedUpgradeable, UUPSUpgradeable {
    /// @dev Storage structure for the DKG
    /// @param automataDcapAttestationFeeAddr The address of the automata dcap attestation fee contract
    /// @param minReqRegisteredParticipants The minimum number of participants needed to be registered for each round
    /// @param minReqFinalizedParticipants The minimum number of participants needed to finish dkg for each round
    /// @param operationalThreshold The operational threshold
    /// @param fee The fee paid to request DKG registration (register and finalize)
    /// @param enclaveTypeData The data of the enclave type
    /// @param isEnclaveTypeWhitelisted The whitelist of enclave types
    /// @custom:storage-location erc7201:story.DKG
    struct DKGStorage {
        address automataDcapAttestationFeeAddr;
        uint256 minReqRegisteredParticipants;
        uint256 minReqFinalizedParticipants;
        uint256 operationalThreshold;
        uint256 fee;
        mapping(bytes32 enclaveType => EnclaveTypeData enclaveTypeData) enclaveTypeData;
        mapping(bytes32 enclaveType => bool isWhitelisted) isEnclaveTypeWhitelisted;
    }

    // basis point for the operational threshold
    uint256 public constant operationalThresholdBasis = 1000;

    // keccak256(abi.encode(uint256(keccak256("story.DKG")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant DKGStorageLocation = 0x12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f00;

    modifier chargesFee() {
        require(msg.value == _getDKGStorage().fee, "DKG: Invalid fee amount");
        payable(address(0x0)).transfer(msg.value);
        _;
    }

    constructor() {
        _disableInitializers();
    }

    /// @notice Initializer for this implementation contract
    /// @param accessManager The address of the access manager of the contract
    function initialize(address accessManager) external initializer {
        __Pausable_init();
        __AccessManaged_init(accessManager);
        __UUPSUpgradeable_init();
        // TODO: add from admin setters the variables to be set at initialization
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Admin Setters                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the minimum number of participants needed to be registered for each round
    /// @param newMinReqRegisteredParticipants The minimum number of participants needed to be registered for each round
    function setMinReqRegisteredParticipants(uint256 newMinReqRegisteredParticipants) external restricted {
        require(newMinReqRegisteredParticipants > 0, "DKG: MinReqRegisteredParticipants cannot be zero");
        _getDKGStorage().minReqRegisteredParticipants = newMinReqRegisteredParticipants;
        emit MinReqRegisteredParticipantsSet(newMinReqRegisteredParticipants);
    }

    /// @notice Sets the minimum number of participants needed to finish dkg for each round
    /// @param newMinReqFinalizedParticipants The minimum number of participants needed to finish dkg for each round
    function setMinReqFinalizedParticipants(uint256 newMinReqFinalizedParticipants) external restricted {
        require(newMinReqFinalizedParticipants > 0, "DKG: MinReqFinalizedParticipants cannot be zero");
        _getDKGStorage().minReqFinalizedParticipants = newMinReqFinalizedParticipants;
        emit MinReqFinalizedParticipantsSet(newMinReqFinalizedParticipants);
    }

    /// @notice Sets the operational threshold
    /// @param newOperationalThreshold The operational threshold
    function setOperationalThreshold(uint256 newOperationalThreshold) external restricted {
        require(newOperationalThreshold > 0, "DKG: Operational threshold cannot be zero");
        require(
            newOperationalThreshold <= operationalThresholdBasis,
            "DKG: Operational threshold cannot be greater than 1000"
        );
        _getDKGStorage().operationalThreshold = newOperationalThreshold;
        emit OperationalThresholdSet(newOperationalThreshold);
    }

    /// @notice Sets the fee paid to request DKG registration (register and finalize)
    /// @param newFee The fee paid to request DKG registration (register and finalize)
    function setFee(uint256 newFee) external restricted {
        _getDKGStorage().fee = newFee;
        emit FeeSet(newFee);
    }

    /// @notice Whitelists an enclave type
    /// @param enclaveType The type of the enclave
    /// @param enclaveTypeData The data of the enclave type
    /// @param isWhitelisted Whether the enclave type is whitelisted
    function whitelistEnclaveType(
        bytes32 enclaveType,
        EnclaveTypeData memory enclaveTypeData,
        bool isWhitelisted
    ) external restricted {
        require(enclaveType != bytes32(0), "DKG: Enclave type cannot be empty");
        require(enclaveTypeData.codeCommitment != bytes32(0), "DKG: Code commitment cannot be empty");
        require(enclaveTypeData.validationHookAddr != address(0), "DKG: Validation hook cannot be empty");
        DKGStorage storage $ = _getDKGStorage();

        $.enclaveTypeData[enclaveType] = enclaveTypeData;
        $.isEnclaveTypeWhitelisted[enclaveType] = isWhitelisted;

        emit EnclaveTypeWhitelisted(
            enclaveType,
            enclaveTypeData.codeCommitment,
            enclaveTypeData.validationHookAddr,
            isWhitelisted
        );
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Authentication Logic                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Authenticates an enclave
    /// @param enclaveReport The enclave report
    /// @param enclaveInstanceData The data of the enclave instance
    /// @param validationContext The validation context
    function authenticateEnclaveReport(
        bytes calldata enclaveReport,
        EnclaveInstanceData calldata enclaveInstanceData,
        bytes calldata validationContext
    ) external {
        _authenticateEnclaveReport(enclaveReport, enclaveInstanceData, validationContext);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              CL Operations                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Registers an enclave instance
    /// @param enclaveReport The enclave report
    /// @param enclaveInstanceData The data of the enclave instance
    /// @param validationContext The validation context
    function register(
        bytes calldata enclaveReport,
        EnclaveInstanceData calldata enclaveInstanceData,
        bytes calldata validationContext
    ) external payable chargesFee {
        require(enclaveReport.length != 0, "DKG: Enclave report cannot be empty");
        require(enclaveInstanceData.round != 0, "DKG: Round cannot be zero");
        require(enclaveInstanceData.validatorAddr != address(0), "DKG: Validator address cannot be empty");
        require(enclaveInstanceData.enclaveType != bytes32(0), "DKG: Enclave type cannot be empty");
        require(enclaveInstanceData.enclaveCommKey.length != 0, "DKG: Enclave communication key cannot be empty");
        require(enclaveInstanceData.dkgPubKey.length != 0, "DKG: DKG public key cannot be empty");

        _authenticateEnclaveReport(enclaveReport, enclaveInstanceData, validationContext);

        emit Registered(
            enclaveReport,
            enclaveInstanceData.round,
            enclaveInstanceData.validatorAddr,
            enclaveInstanceData.enclaveType,
            enclaveInstanceData.enclaveCommKey,
            enclaveInstanceData.dkgPubKey,
            validationContext
        );
    }

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
    ) external payable chargesFee {
        DKGStorage storage $ = _getDKGStorage();
        require(round != 0, "DKG: Round cannot be zero");
        require(validatorAddr != address(0), "DKG: Validator address cannot be empty");
        require($.isEnclaveTypeWhitelisted[enclaveType], "DKG: Enclave type is not whitelisted");
        require(participantsRoot != bytes32(0), "DKG: Participants root cannot be empty");
        require(globalPubKey.length != 0, "DKG: Global public key cannot be empty");
        require(publicCoeffs.length != 0, "DKG: Public coefficients cannot be empty");
        require(signature.length != 0, "DKG: Signature cannot be empty");

        emit Finalized(
            round,
            validatorAddr,
            enclaveType,
            $.enclaveTypeData[enclaveType].codeCommitment,
            participantsRoot,
            globalPubKey,
            publicCoeffs,
            signature
        );
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Authenticates an enclave report
    /// @param enclaveReport The enclave report
    /// @param enclaveInstanceData The data of the enclave instance
    /// @param validationContext The validation context
    function _authenticateEnclaveReport(
        bytes calldata enclaveReport,
        EnclaveInstanceData calldata enclaveInstanceData,
        bytes calldata validationContext
    ) internal {
        DKGStorage storage $ = _getDKGStorage();
        EnclaveTypeData memory enclaveTypeData = $.enclaveTypeData[enclaveInstanceData.enclaveType];
        require($.isEnclaveTypeWhitelisted[enclaveInstanceData.enclaveType], "DKG: Enclave type is not whitelisted");

        bool isValidReport = IAttestationReportValidator(enclaveTypeData.validationHookAddr).validateReport(
            enclaveTypeData.codeCommitment,
            keccak256(abi.encode(enclaveInstanceData)),
            enclaveReport,
            validationContext
        );
        require(isValidReport, "DKG: Enclave authentication failed");
    }

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override restricted {}

    /// @dev Returns the storage struct of DKG.
    function _getDKGStorage() private pure returns (DKGStorage storage $) {
        assembly {
            $.slot := DKGStorageLocation
        }
    }
}
