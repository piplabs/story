// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { ICDRWriteCondition } from "../interfaces/ICDRWriteCondition.sol";

/// @notice Minimal interface for Story Protocol's IPAccount (ERC-6551)
interface IIPAccount {
    /// @notice Returns the owner of the IP Account (the underlying NFT owner)
    function owner() external view returns (address);
}

/// @notice Minimal interface for Story Protocol's IP Asset Registry
interface IIPAssetRegistry {
    /// @notice Checks whether an IP was registered based on its ID
    function isRegistered(address id) external view returns (bool);
}

/// @title IPOwnerWriteCondition
/// @notice CDR write condition that gates vault writes to the owner of a
///         specific Story Protocol IP asset.
/// @dev    writeConditionData encodes the target IP account: abi.encode(address ipAccountAddr)
///         The caller must be the current owner of the IP (i.e., the owner of the
///         underlying NFT bound to the IP Account via ERC-6551).
contract IPOwnerWriteCondition is ICDRWriteCondition {
    /// @notice The Story Protocol IP Asset Registry
    IIPAssetRegistry public immutable IP_ASSET_REGISTRY;

    constructor(address ipAssetRegistry) {
        require(ipAssetRegistry != address(0), "IPOwnerWriteCondition: zero address");
        IP_ASSET_REGISTRY = IIPAssetRegistry(ipAssetRegistry);
    }

    /// @inheritdoc ICDRWriteCondition
    function checkWriteCondition(
        uint32,
        bytes calldata,
        bytes calldata writeConditionData,
        address caller
    ) external view returns (bool) {
        address ipAccountAddr = abi.decode(writeConditionData, (address));

        // IP must be registered
        if (!IP_ASSET_REGISTRY.isRegistered(ipAccountAddr)) return false;

        // Caller must be the IP owner (underlying NFT owner)
        if (IIPAccount(ipAccountAddr).owner() != caller) return false;

        return true;
    }
}
