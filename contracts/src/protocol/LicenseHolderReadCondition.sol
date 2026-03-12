// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { ICDRReadCondition } from "../interfaces/ICDRReadCondition.sol";

/// @notice Minimal interface for Story Protocol's LicenseToken (ERC-721)
interface ILicenseToken {
    /// @notice Returns the owner of a license token
    function ownerOf(uint256 tokenId) external view returns (address);

    /// @notice Returns the licensor IP ID for a license token
    function getLicensorIpId(uint256 tokenId) external view returns (address);

    /// @notice Checks whether a license token has been revoked
    function isLicenseTokenRevoked(uint256 tokenId) external view returns (bool);
}

/// @title LicenseHolderReadCondition
/// @notice CDR read condition that gates vault decryption to holders of a valid
///         Story Protocol license token for a specific IP asset.
/// @dev    readConditionData encodes the target IP asset: abi.encode(address ipId)
///         accessAuxData encodes the proof of license: abi.encode(uint256 tokenId)
///         The caller must own the referenced license token, the token must be
///         issued for the target IP, and it must not be revoked.
contract LicenseHolderReadCondition is ICDRReadCondition {
    /// @notice The Story Protocol LicenseToken contract
    ILicenseToken public immutable LICENSE_TOKEN;

    constructor(address licenseToken) {
        require(licenseToken != address(0), "LicenseHolderReadCondition: zero address");
        LICENSE_TOKEN = ILicenseToken(licenseToken);
    }

    /// @inheritdoc ICDRReadCondition
    function checkReadCondition(
        uint32,
        bytes calldata accessAuxData,
        bytes calldata readConditionData,
        address caller
    ) external view returns (bool) {
        address ipId = abi.decode(readConditionData, (address));
        uint256 tokenId = abi.decode(accessAuxData, (uint256));

        // Caller must own the license token
        if (LICENSE_TOKEN.ownerOf(tokenId) != caller) return false;

        // Token must be a license for the target IP
        if (LICENSE_TOKEN.getLicensorIpId(tokenId) != ipId) return false;

        // Token must not be revoked
        if (LICENSE_TOKEN.isLicenseTokenRevoked(tokenId)) return false;

        return true;
    }
}
