// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { ICDR } from "../interfaces/ICDR.sol";
import { ICDRWriteCondition } from "../interfaces/ICDRWriteCondition.sol";
import { ICDRReadCondition } from "../interfaces/ICDRReadCondition.sol";

contract CDR is ICDR, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PausableUpgradeable {
    /// @dev Storage structure for the CDR
    /// @param uuid The UUID of the vault
    /// @param baseFee The base fee
    /// @param writeFee The write fee
    /// @param readFee The read fee
    /// @param allocateFee The allocate fee
    /// @param maxEncryptedDataSize Maximum allowed size for encrypted vault data (bytes)
    /// @param maxEncryptedPartialSize Maximum allowed size for encrypted partial decryptions (bytes)
    /// @param vaults The mapping of the vaults
    /// @custom:storage-location erc7201:story.CDR
    struct CDRStorage {
        uint32 uuid;
        uint256 baseFee;
        uint256 writeFee;
        uint256 readFee;
        uint256 allocateFee;
        uint256 maxEncryptedDataSize;
        uint256 maxEncryptedPartialSize;
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
    /// @param maxEncryptedDataSize Maximum allowed size for encrypted vault data (bytes)
    /// @param maxEncryptedPartialSize Maximum allowed size for encrypted partial decryptions (bytes)
    function initialize(
        address owner,
        uint256 baseFee,
        uint256 writeFee,
        uint256 readFee,
        uint256 allocateFee,
        uint256 maxEncryptedDataSize,
        uint256 maxEncryptedPartialSize
    ) external initializer {
        __Ownable_init(owner);
        __ReentrancyGuard_init();
        __Pausable_init();

        _setBaseFee(baseFee);
        _setWriteFee(writeFee);
        _setReadFee(readFee);
        _setAllocateFee(allocateFee);
        _setMaxEncryptedDataSize(maxEncryptedDataSize);
        _setMaxEncryptedPartialSize(maxEncryptedPartialSize);
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
    /// @dev Zero fees are intentionally allowed. This enables fee-free operation during
    ///      testing and early deployment phases. The owner can set a non-zero fee later.
    /// @param newAllocateFee The allocate fee
    function setAllocateFee(uint256 newAllocateFee) external onlyOwner {
        _setAllocateFee(newAllocateFee);
    }

    /// @notice Sets the maximum allowed size for encrypted vault data
    /// @param newMaxEncryptedDataSize The maximum size in bytes
    function setMaxEncryptedDataSize(uint256 newMaxEncryptedDataSize) external onlyOwner {
        _setMaxEncryptedDataSize(newMaxEncryptedDataSize);
    }

    /// @notice Sets the maximum allowed size for encrypted partial decryptions
    /// @param newMaxEncryptedPartialSize The maximum size in bytes
    function setMaxEncryptedPartialSize(uint256 newMaxEncryptedPartialSize) external onlyOwner {
        _setMaxEncryptedPartialSize(newMaxEncryptedPartialSize);
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
    ) external payable nonReentrant whenNotPaused returns (uint32 newVaultUuid) {
        require(writeConditionAddr != address(0) && readConditionAddr != address(0), "Invalid condition address");

        CDRStorage storage $ = _getCDRStorage();
        // collect allocation fee and burn it
        _collectFee($.allocateFee, ICDR.FeeType.Allocate);

        require($.uuid < type(uint32).max, "CDR: Vault UUID overflow");
        newVaultUuid = $.uuid++;
        $.vaults[newVaultUuid] = Vault(
            updatable,
            writeConditionAddr,
            readConditionAddr,
            writeConditionData,
            readConditionData,
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
    /// @dev If msg.sender is the writeConditionAddr itself, the condition check is
    ///      bypassed. This is intentional because the condition contract has already evaluated
    ///      its own logic before calling write(), so re-checking would be redundant.
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for writing
    /// @param encryptedData The encrypted data to write
    function write(
        uint32 uuid,
        bytes calldata accessAuxData,
        bytes calldata encryptedData
    ) external payable nonReentrant whenNotPaused {
        CDRStorage storage $ = _getCDRStorage();
        require(encryptedData.length > 0, "CDR: Encrypted data cannot be empty");
        require(encryptedData.length <= $.maxEncryptedDataSize, "CDR: Encrypted data exceeds max size");
        // check if the vault exists
        Vault storage vault = $.vaults[uuid];
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

        emit VaultWritten(uuid, encryptedData);
    }

    /// @notice Reads data from a vault
    /// @dev If msg.sender is the readConditionAddr itself, the condition check
    ///      is bypassed. This is intentional because the condition contract has already evaluated
    ///      its own logic before calling read(). Read access control is enforced entirely through
    ///      the readConditionAddr contract — there is no additional caller restriction.
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
        Vault storage vault = $.vaults[uuid];
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

        emit VaultRead(uuid, msg.sender, vault.encryptedData, requesterPubKey);
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
    ) external payable nonReentrant whenNotPaused {
        CDRStorage storage $ = _getCDRStorage();
        require(
            encryptedPartial.length > 0 && encryptedPartial.length <= $.maxEncryptedPartialSize,
            "CDR: Invalid encrypted partial length"
        );

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

    /// @notice Gets the maximum allowed size for encrypted vault data
    /// @return maxEncryptedDataSize The maximum size in bytes
    function maxEncryptedDataSize() external view returns (uint256) {
        return _getCDRStorage().maxEncryptedDataSize;
    }

    /// @notice Gets the maximum allowed size for encrypted partial decryptions
    /// @return maxEncryptedPartialSize The maximum size in bytes
    function maxEncryptedPartialSize() external view returns (uint256) {
        return _getCDRStorage().maxEncryptedPartialSize;
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

    function _setMaxEncryptedDataSize(uint256 newMaxEncryptedDataSize) internal {
        require(newMaxEncryptedDataSize > 0, "CDR: Max encrypted data size must be > 0");
        _getCDRStorage().maxEncryptedDataSize = newMaxEncryptedDataSize;
    }

    function _setMaxEncryptedPartialSize(uint256 newMaxEncryptedPartialSize) internal {
        require(newMaxEncryptedPartialSize > 0, "CDR: Max encrypted partial size must be > 0");
        _getCDRStorage().maxEncryptedPartialSize = newMaxEncryptedPartialSize;
    }

    /// @notice Collects a fee and emits a FeeCollected event
    /// @param feeAmountToCollect The fee amount to collect
    /// @param feeType The type of operation generating the fee
    function _collectFee(uint256 feeAmountToCollect, ICDR.FeeType feeType) internal {
        require(msg.value == feeAmountToCollect, "CDR: Invalid fee amount");
        payable(address(0x0)).transfer(feeAmountToCollect);
        emit FeeCollected(msg.sender, feeAmountToCollect, feeType);
    }

    /// @dev Returns the storage struct of CDR.
    function _getCDRStorage() private pure returns (CDRStorage storage $) {
        assembly {
            $.slot := CDRStorageLocation
        }
    }
}
