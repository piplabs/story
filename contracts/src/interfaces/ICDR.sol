// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface ICDR {
    enum FeeType {
        Allocate,
        Write,
        Read,
        SubmitPartial
    }
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
    /// @param requester The address requesting the read (msg.sender)
    /// @param ciphertext The encrypted data (ciphertext)
    /// @param requesterPubKey The public key of the requester
    event VaultRead(
        uint32 uuid,
        address indexed requester,
        bytes ciphertext,
        bytes requesterPubKey
    );

    /// @notice Emitted when an encrypted partial decryption is submitted
    /// @param validator The address of the submitting validator (msg.sender)
    /// @param round The DKG round number
    /// @param pid The participant index of the validator
    /// @param encryptedPartial The encrypted partial decryption
    /// @param ephemeralPubKey The ephemeral public key used for encryption
    /// @param pubShare The validator's public key share
    /// @param requesterPubKey The public key of the requester
    /// @param uuid The UUID of the vault
    /// @param signature The signature over the partial decryption payload
    /// @param fee The fee collected for the submission
    event EncryptedPartialDecryptionSubmitted(
        address indexed validator,
        uint32 round,
        uint32 pid,
        bytes encryptedPartial,
        bytes ephemeralPubKey,
        bytes pubShare,
        bytes requesterPubKey,
        uint32 uuid,
        bytes signature,
        uint256 fee
    );

    /// @notice Emitted when a CDR fee is collected
    /// @param payer The address that paid the fee (msg.sender)
    /// @param amount The fee amount collected
    /// @param feeType The type of operation: 0=allocate, 1=write, 2=read, 3=submitPartialDecrypt
    event FeeCollected(address indexed payer, uint256 amount, FeeType feeType);

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
    /// @param requesterPubKey The public key of the requester
    function read(
        uint32 uuid,
        bytes memory accessAuxData,
        bytes calldata requesterPubKey
    ) external payable;

    /// @notice Submits an encrypted partial decryption
    /// @param round The DKG round number
    /// @param pid The participant index of the validator
    /// @param encryptedPartial The encrypted partial decryption
    /// @param ephemeralPubKey The ephemeral public key used for encryption
    /// @param pubShare The validator's public key share
    /// @param requesterPubKey The public key of the requester
    /// @param uuid The UUID of the vault
    /// @param signature The signature over the partial decryption payload
    function submitEncryptedPartialDecryption(
        uint32 round,
        uint32 pid,
        bytes calldata encryptedPartial,
        bytes calldata ephemeralPubKey,
        bytes calldata pubShare,
        bytes calldata requesterPubKey,
        uint32 uuid,
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
