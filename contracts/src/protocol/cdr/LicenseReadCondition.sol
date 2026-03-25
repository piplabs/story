// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { ICDRReadCondition } from "../../interfaces/ICDRReadCondition.sol";

/// @notice Minimal interface for Story Protocol's LicenseToken (ERC721-based).
interface ILicenseToken {
    struct LicenseTokenMetadata {
        address licensorIpId;
        address licenseTemplate;
        uint256 licenseTermsId;
        bool transferable;
        uint32 commercialRevShare;
    }

    function ownerOf(uint256 tokenId) external view returns (address);
    function getLicenseTokenMetadata(uint256 tokenId) external view returns (LicenseTokenMetadata memory);
    function isLicenseTokenRevoked(uint256 tokenId) external view returns (bool);
}

/// @title LicenseReadCondition
/// @notice CDR read condition that gates access to a vault based on ownership of a Story Protocol
///         license token for the vault's IP asset.
///
/// Usage:
///   - When allocating a vault via CDR.allocate(), set readConditionAddr to this contract's address
///     and readConditionData to abi.encode(licenseTokenAddress, ipId).
///   - When reading via CDR.read(), pass accessAuxData as abi.encode(licenseTokenIds) where
///     licenseTokenIds is a uint256[] of token IDs the caller owns.
contract LicenseReadCondition is ICDRReadCondition {
    /// @notice Checks whether the caller holds a valid (non-revoked) license token
    ///         issued for the vault's IP asset.
    /// @param accessAuxData ABI-encoded uint256[] of license token IDs the caller presents as proof
    /// @param readConditionData ABI-encoded (address licenseToken, address ipId)
    /// @param caller The address attempting to read the vault
    /// @return True if the caller owns at least one valid license token for the IP asset
    function checkReadCondition(
        uint32,
        bytes calldata accessAuxData,
        bytes calldata readConditionData,
        address caller
    ) external view override returns (bool) {
        (address licenseToken, address ipId) = abi.decode(readConditionData, (address, address));
        uint256[] memory tokenIds = abi.decode(accessAuxData, (uint256[]));

        ILicenseToken lt = ILicenseToken(licenseToken);

        for (uint256 i = 0; i < tokenIds.length; i++) {
            // Caller must own the token
            if (lt.ownerOf(tokenIds[i]) != caller) continue;

            // Token must not be revoked
            if (lt.isLicenseTokenRevoked(tokenIds[i])) continue;

            // Token must be issued for the target IP asset
            ILicenseToken.LicenseTokenMetadata memory meta = lt.getLicenseTokenMetadata(tokenIds[i]);
            if (meta.licensorIpId == ipId) {
                return true;
            }
        }

        return false;
    }
}
