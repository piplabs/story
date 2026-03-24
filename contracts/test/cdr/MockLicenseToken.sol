// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { ILicenseToken } from "../../src/protocol/cdr/LicenseReadCondition.sol";

/// @title MockLicenseToken
/// @notice Mock of Story Protocol's LicenseToken for testing LicenseReadCondition.
contract MockLicenseToken {
    struct TokenData {
        address owner;
        address licensorIpId;
        address licenseTemplate;
        uint256 licenseTermsId;
        bool transferable;
        uint32 commercialRevShare;
        bool revoked;
        bool exists;
    }

    mapping(uint256 => TokenData) private _tokens;
    uint256 private _nextTokenId;

    /// @notice Mints a license token for testing
    function mint(
        address to,
        address licensorIpId,
        address licenseTemplate,
        uint256 licenseTermsId,
        bool transferable,
        uint32 commercialRevShare
    ) external returns (uint256 tokenId) {
        tokenId = _nextTokenId++;
        _tokens[tokenId] = TokenData({
            owner: to,
            licensorIpId: licensorIpId,
            licenseTemplate: licenseTemplate,
            licenseTermsId: licenseTermsId,
            transferable: transferable,
            commercialRevShare: commercialRevShare,
            revoked: false,
            exists: true
        });
    }

    /// @notice Revokes a license token for testing
    function revoke(uint256 tokenId) external {
        require(_tokens[tokenId].exists, "Token does not exist");
        _tokens[tokenId].revoked = true;
    }

    /// @notice Transfers a token for testing
    function transferFrom(address, address to, uint256 tokenId) external {
        require(_tokens[tokenId].exists, "Token does not exist");
        _tokens[tokenId].owner = to;
    }

    // --- ILicenseToken interface ---

    function ownerOf(uint256 tokenId) external view returns (address) {
        require(_tokens[tokenId].exists, "Token does not exist");
        return _tokens[tokenId].owner;
    }

    function getLicenseTokenMetadata(
        uint256 tokenId
    ) external view returns (ILicenseToken.LicenseTokenMetadata memory) {
        require(_tokens[tokenId].exists, "Token does not exist");
        TokenData storage t = _tokens[tokenId];
        return
            ILicenseToken.LicenseTokenMetadata({
                licensorIpId: t.licensorIpId,
                licenseTemplate: t.licenseTemplate,
                licenseTermsId: t.licenseTermsId,
                transferable: t.transferable,
                commercialRevShare: t.commercialRevShare
            });
    }

    function isLicenseTokenRevoked(uint256 tokenId) external view returns (bool) {
        require(_tokens[tokenId].exists, "Token does not exist");
        return _tokens[tokenId].revoked;
    }
}
