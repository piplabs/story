// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface ICDR {
    /// @notice Struct for the vault
    /// @param updatable Whether the vault is updatable
    /// @param writeConditionAddr The address of the write condition
    /// @param readConditionAddr The address of the read condition
    /// @param writeConditionData The data of the write condition
    /// @param readConditionData The data of the read condition
    /// @param encryptedData The encrypted data
    struct Vault {
        bool updatable;
        address writeConditionAddr;
        address readConditionAddr;
        bytes writeConditionData;
        bytes readConditionData;
        bytes encryptedData;
    }

    /// @notice Emitted when a vault is allocated
    /// @param uuid The UUID of the vault
    /// @param updatable Whether the vault is updatable
    /// @param writeConditionAddr The address of the write condition
    /// @param readConditionAddr The address of the read condition
    /// @param writeConditionData The data of the write condition
    /// @param readConditionData The data of the read condition
    event VaultAllocated(
        uint32 uuid,
        bool updatable,
        address writeConditionAddr,
        address readConditionAddr,
        bytes writeConditionData,
        bytes readConditionData
    );

    /// @notice Emitted when a vault is written
    /// @param uuid The UUID of the vault
    /// @param encryptedData The encrypted data
    event VaultWritten(uint32 uuid, bytes encryptedData);

    /// @notice Emitted when a vault is read
    /// @param uuid The UUID of the vault
    /// @param encryptedData The encrypted data
    /// @param recipientPublicKey The public key of the recipient
    event VaultRead(uint32 uuid, bytes encryptedData, uint256[2] recipientPublicKey);

    /// @notice Emitted when an encrypted partial decryption is submitted
    /// @param enclaveID The ID of the enclave
    /// @param encryptedPartial The encrypted partial decryption
    /// @param signature The signature of the encrypted partial decryption
    event EncryptedPartialDecryptionSubmitted(address enclaveID, bytes encryptedPartial, bytes signature);

    /// @notice Sets the base fee
    /// @param newBaseFee The base fee
    function setBaseFee(uint256 newBaseFee) external;

    /// @notice Sets the write fee
    /// @param newWriteFee The write fee
    function setWriteFee(uint256 newWriteFee) external;

    /// @notice Sets the read fee
    /// @param newReadFee The read fee
    function setReadFee(uint256 newReadFee) external;

    /// @notice Sets the allocate fee
    /// @param newAllocateFee The allocate fee
    function setAllocateFee(uint256 newAllocateFee) external;

    /// @notice Allocates a new vault
    /// @param updatable Whether the vault is updatable
    /// @param writeConditionAddr The address of the write condition
    /// @param readConditionAddr The address of the read condition
    /// @param writeconditionData The data of the write condition
    /// @param readconditionData The data of the read condition
    /// returns the uuid of the new vault
    function allocate(
        bool updatable,
        address writeConditionAddr,
        address readConditionAddr,
        bytes calldata writeconditionData,
        bytes calldata readconditionData
    ) external payable returns (uint32 newVaultUuid);

    /// @notice Writes data to a vault
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for writing
    /// @param encryptedData The encrypted data to write
    function write(uint32 uuid, bytes calldata accessAuxData, bytes calldata encryptedData) external payable;

    /// @notice Reads data from a vault
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for reading
    /// @param recipientPublicKey The public key of the recipient
    function read(uint32 uuid, bytes memory accessAuxData, uint256[2] calldata recipientPublicKey) external payable;

    /// @notice Submits an encrypted partial decryption
    /// @param enclaveID The ID of the enclave
    /// @param encryptedPartial The encrypted partial decryption
    /// @param signature The signature of the encrypted partial decryption
    function submitEncryptedPartialDecryption(
        address enclaveID,
        bytes calldata encryptedPartial,
        bytes calldata signature
    ) external payable;

    /// @notice Gets the UUID of the vault
    /// @return uuid The UUID of the vault
    function uuid() external view returns (uint32 uuid);

    /// @notice Gets the base fee
    /// @return baseFee The base fee
    function baseFee() external view returns (uint256);

    /// @notice Gets the write fee
    /// @return writeFee The write fee
    function writeFee() external view returns (uint256);

    /// @notice Gets the read fee
    /// @return readFee The read fee
    function readFee() external view returns (uint256);

    /// @notice Gets the allocate fee
    /// @return allocateFee The allocate fee
    function allocateFee() external view returns (uint256);

    /// @notice Gets the vault
    /// @param uuid The UUID of the vault
    /// @return vault The vault
    function vaults(uint32 uuid) external view returns (Vault memory vault);
}
