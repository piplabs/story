// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import { ICDR } from "../interfaces/ICDR.sol";
import { ICDRWriteCondition } from "../interfaces/ICDRWriteCondition.sol";
import { ICDRReadCondition } from "../interfaces/ICDRReadCondition.sol";

contract CDR is ICDR, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    /// @dev Storage structure for the CDR
    /// @param uuid The UUID of the vault
    /// @param baseFee The base fee
    /// @param writeFee The write fee
    /// @param readFee The read fee
    /// @param allocateFee The allocate fee
    /// @param vaults The mapping of the vaults
    /// @custom:storage-location erc7201:story.CDR
    struct CDRStorage {
        uint32 uuid;
        uint256 baseFee;
        uint256 writeFee;
        uint256 readFee;
        uint256 allocateFee;
        mapping(uint32 uuid => Vault vault) vaults;
    }

    // keccak256(abi.encode(uint256(keccak256("story.CDR")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant CDRStorageLocation = 0x38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf600;

    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract
    /// @param owner The address of the owner of the contract
    /// @param baseFee The base fee for partial decryption submissions
    /// @param writeFee The fee for writing data to a vault
    /// @param readFee The fee for reading data from a vault
    /// @param allocateFee The fee for allocating a new vault
    function initialize(
        address owner,
        uint256 baseFee,
        uint256 writeFee,
        uint256 readFee,
        uint256 allocateFee
    ) external initializer {
        __Ownable_init(owner);
        __ReentrancyGuard_init();
        __Pausable_init();
        __UUPSUpgradeable_init();

        _setBaseFee(baseFee);
        _setWriteFee(writeFee);
        _setReadFee(readFee);
        _setAllocateFee(allocateFee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Admin Setters                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the base fee
    /// @param newBaseFee The base fee
    function setBaseFee(uint256 newBaseFee) external onlyOwner {
        _setBaseFee(newBaseFee);
    }

    /// @notice Sets the write fee
    /// @param newWriteFee The write fee
    function setWriteFee(uint256 newWriteFee) external onlyOwner {
        _setWriteFee(newWriteFee);
    }

    /// @notice Sets the read fee
    /// @param newReadFee The read fee
    function setReadFee(uint256 newReadFee) external onlyOwner {
        _setReadFee(newReadFee);
    }

    /// @notice Sets the allocate fee
    /// @param newAllocateFee The allocate fee
    function setAllocateFee(uint256 newAllocateFee) external onlyOwner {
        _setAllocateFee(newAllocateFee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              CDR Operations                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Allocates a new vault
    /// @param updatable Whether the vault is updatable
    /// @param writeConditionAddr The address of the write condition
    /// @param readConditionAddr The address of the read condition
    /// @param writeConditionData The data of the write condition
    /// @param readConditionData The data of the read condition
    /// returns the uuid of the new vault
    function allocate(
        bool updatable,
        address writeConditionAddr,
        address readConditionAddr,
        bytes calldata writeConditionData,
        bytes calldata readConditionData
    ) external payable whenNotPaused returns (uint32 newVaultUuid) {
        require(writeConditionAddr != address(0) || readConditionAddr != address(0), "Invalid condition address");

        CDRStorage storage $ = _getCDRStorage();
        // collect allocation fee and burn it
        _collectFee($.allocateFee, ICDR.FeeType.Allocate);

        newVaultUuid = $.uuid++;
        $.vaults[newVaultUuid] = Vault(
            updatable,
            writeConditionAddr,
            readConditionAddr,
            writeConditionData,
            readConditionData,
            "",
            ""
        );

        emit VaultAllocated(
            newVaultUuid,
            updatable,
            writeConditionAddr,
            readConditionAddr,
            writeConditionData,
            readConditionData
        );
    }

    /// @notice Writes data to a vault
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for writing
    /// @param encryptedData The encrypted data to write
    /// @param label The TDH2 label used as authenticated associated data during encryption
    function write(
        uint32 uuid,
        bytes calldata accessAuxData,
        bytes calldata encryptedData,
        bytes calldata label
    ) external payable nonReentrant whenNotPaused {
        require(encryptedData.length > 0, "CDR: Encrypted data cannot be empty");

        CDRStorage storage $ = _getCDRStorage();
        // check if the vault exists
        Vault memory vault = $.vaults[uuid];
        require(vault.writeConditionAddr != address(0), "CDR: Write condition address not set");

        // check the write condition
        if (msg.sender != vault.writeConditionAddr) {
            require(
                ICDRWriteCondition(vault.writeConditionAddr).checkWriteCondition(
                    uuid,
                    accessAuxData,
                    vault.writeConditionData,
                    msg.sender
                ),
                "CDR: Write condition not met"
            );
        }

        // if the vault has data and is not updatable, revert
        if (vault.encryptedData.length > 0) require(vault.updatable, "CDR: Vault is not updatable");

        // collect the write fee and burn it
        _collectFee($.writeFee, ICDR.FeeType.Write);

        // update the data on the vault
        $.vaults[uuid].encryptedData = encryptedData;
        $.vaults[uuid].label = label;

        emit VaultWritten(uuid, encryptedData, label);
    }

    /// @notice Reads data from a vault
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for reading
    /// @param requesterPubKey The public key of the requester
    function read(
        uint32 uuid,
        bytes memory accessAuxData,
        bytes calldata requesterPubKey
    ) external payable nonReentrant whenNotPaused {
        CDRStorage storage $ = _getCDRStorage();
        // check if the vault has data to read
        Vault memory vault = $.vaults[uuid];
        require(vault.encryptedData.length > 0, "CDR: Vault has no data to read");

        // check the read condition
        if (msg.sender != vault.readConditionAddr) {
            require(
                ICDRReadCondition(vault.readConditionAddr).checkReadCondition(
                    uuid,
                    accessAuxData,
                    vault.readConditionData,
                    msg.sender
                ),
                "CDR: Read condition not met"
            );
        }

        // collect the read fee and burn it
        _collectFee($.readFee, ICDR.FeeType.Read);

        emit VaultRead(uuid, msg.sender, vault.encryptedData, requesterPubKey, vault.label);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              CL Operations                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Submits an encrypted partial decryption
    /// @param round The DKG round number
    /// @param pid The participant index of the validator
    /// @param encryptedPartial The encrypted partial decryption
    /// @param ephemeralPubKey The ephemeral public key used for encryption
    /// @param pubShare The validator's public key share
    /// @param requesterPubKey The public key of the requester
    /// @param ciphertext The ciphertext associated with the request
    /// @param uuid The UUID of the vault
    /// @param signature The signature over the partial decryption payload
    function submitEncryptedPartialDecryption(
        uint32 round,
        uint32 pid,
        bytes calldata encryptedPartial,
        bytes calldata ephemeralPubKey,
        bytes calldata pubShare,
        bytes calldata requesterPubKey,
        bytes calldata ciphertext,
        uint32 uuid,
        bytes calldata signature
    ) external payable whenNotPaused {
        // collect the base fee and burn it
        uint256 fee = _getCDRStorage().baseFee;
        _collectFee(fee, ICDR.FeeType.SubmitPartial);

        emit EncryptedPartialDecryptionSubmitted(
            msg.sender,
            round,
            pid,
            encryptedPartial,
            ephemeralPubKey,
            pubShare,
            requesterPubKey,
            ciphertext,
            uuid,
            signature,
            fee
        );
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Get Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Gets the UUID of the vault
    /// @return uuid The UUID of the vault
    function uuid() external view returns (uint32 uuid) {
        return _getCDRStorage().uuid;
    }

    /// @notice Gets the base fee
    /// @return baseFee The base fee
    function baseFee() external view returns (uint256) {
        return _getCDRStorage().baseFee;
    }

    /// @notice Gets the write fee
    /// @return writeFee The write fee
    function writeFee() external view returns (uint256) {
        return _getCDRStorage().writeFee;
    }

    /// @notice Gets the read fee
    /// @return readFee The read fee
    function readFee() external view returns (uint256) {
        return _getCDRStorage().readFee;
    }

    /// @notice Gets the allocate fee
    /// @return allocateFee The allocate fee
    function allocateFee() external view returns (uint256) {
        return _getCDRStorage().allocateFee;
    }

    /// @notice Gets the vault
    /// @param uuid The UUID of the vault
    /// @return vault The vault
    function vaults(uint32 uuid) external view returns (Vault memory vault) {
        return _getCDRStorage().vaults[uuid];
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the base fee
    /// @param newBaseFee The base fee
    function _setBaseFee(uint256 newBaseFee) internal {
        _getCDRStorage().baseFee = newBaseFee;
    }

    /// @notice Sets the write fee
    /// @param newWriteFee The write fee
    function _setWriteFee(uint256 newWriteFee) internal {
        _getCDRStorage().writeFee = newWriteFee;
    }

    /// @notice Sets the read fee
    /// @param newReadFee The read fee
    function _setReadFee(uint256 newReadFee) internal {
        _getCDRStorage().readFee = newReadFee;
    }

    /// @notice Sets the allocate fee
    /// @param newAllocateFee The allocate fee
    function _setAllocateFee(uint256 newAllocateFee) internal {
        _getCDRStorage().allocateFee = newAllocateFee;
    }

    /// @notice Collects a fee and emits a FeeCollected event
    /// @param feeAmountToCollect The fee amount to collect
    /// @param feeType The type of operation generating the fee
    function _collectFee(uint256 feeAmountToCollect, ICDR.FeeType feeType) internal {
        require(msg.value == feeAmountToCollect, "CDR: Invalid fee amount");
        payable(address(0x0)).transfer(feeAmountToCollect);
        emit FeeCollected(msg.sender, feeAmountToCollect, feeType);
    }

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /// @dev Returns the storage struct of CDR.
    function _getCDRStorage() private pure returns (CDRStorage storage $) {
        assembly {
            $.slot := CDRStorageLocation
        }
    }
}
