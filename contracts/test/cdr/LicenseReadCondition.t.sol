// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";

import { LicenseReadCondition } from "../../src/protocol/cdr/LicenseReadCondition.sol";
import { MockLicenseToken } from "./MockLicenseToken.sol";

contract LicenseReadConditionTest is ForgeTest {
    LicenseReadCondition internal condition;
    MockLicenseToken internal licenseToken;

    address internal ipId = address(0x1111);
    address internal otherIpId = address(0xBEEF);
    address internal licenseTemplate = address(0x7E9B7A7E);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    function setUp() public {
        condition = new LicenseReadCondition();
        licenseToken = new MockLicenseToken();
    }

    // -----------------------------------------------------------------------
    // Helpers
    // -----------------------------------------------------------------------

    function _readConditionData() internal view returns (bytes memory) {
        return abi.encode(address(licenseToken), ipId);
    }

    function _accessAuxData(uint256[] memory tokenIds) internal pure returns (bytes memory) {
        return abi.encode(tokenIds);
    }

    function _singleTokenArray(uint256 tokenId) internal pure returns (uint256[] memory) {
        uint256[] memory ids = new uint256[](1);
        ids[0] = tokenId;
        return ids;
    }

    function _check(address caller, uint256[] memory tokenIds) internal view returns (bool) {
        return
            condition.checkReadCondition(
                0, // uuid (unused)
                _accessAuxData(tokenIds),
                _readConditionData(),
                caller
            );
    }

    // -----------------------------------------------------------------------
    // Tests: valid license holder can read
    // -----------------------------------------------------------------------

    function testCheckReadCondition_ValidLicense() public {
        uint256 tokenId = licenseToken.mint(alice, ipId, licenseTemplate, 1, true, 0);
        assertTrue(_check(alice, _singleTokenArray(tokenId)));
    }

    function testCheckReadCondition_MultipleTokens_OneValid() public {
        // Token for a different IP
        uint256 wrongIp = licenseToken.mint(alice, otherIpId, licenseTemplate, 1, true, 0);
        // Token for the correct IP
        uint256 correctIp = licenseToken.mint(alice, ipId, licenseTemplate, 2, true, 0);

        uint256[] memory ids = new uint256[](2);
        ids[0] = wrongIp;
        ids[1] = correctIp;
        assertTrue(_check(alice, ids));
    }

    // -----------------------------------------------------------------------
    // Tests: no license → denied
    // -----------------------------------------------------------------------

    function testCheckReadCondition_NoTokens() public view {
        uint256[] memory empty = new uint256[](0);
        assertFalse(_check(alice, empty));
    }

    function testCheckReadCondition_NotOwner() public {
        uint256 tokenId = licenseToken.mint(bob, ipId, licenseTemplate, 1, true, 0);
        // alice tries to use bob's token
        assertFalse(_check(alice, _singleTokenArray(tokenId)));
    }

    function testCheckReadCondition_WrongIpId() public {
        uint256 tokenId = licenseToken.mint(alice, otherIpId, licenseTemplate, 1, true, 0);
        assertFalse(_check(alice, _singleTokenArray(tokenId)));
    }

    function testCheckReadCondition_RevokedToken() public {
        uint256 tokenId = licenseToken.mint(alice, ipId, licenseTemplate, 1, true, 0);
        licenseToken.revoke(tokenId);
        assertFalse(_check(alice, _singleTokenArray(tokenId)));
    }

    // -----------------------------------------------------------------------
    // Tests: edge cases
    // -----------------------------------------------------------------------

    function testCheckReadCondition_RevokedAndValid() public {
        uint256 revoked = licenseToken.mint(alice, ipId, licenseTemplate, 1, true, 0);
        licenseToken.revoke(revoked);

        uint256 valid = licenseToken.mint(alice, ipId, licenseTemplate, 2, true, 0);

        uint256[] memory ids = new uint256[](2);
        ids[0] = revoked;
        ids[1] = valid;
        assertTrue(_check(alice, ids));
    }

    function testCheckReadCondition_TransferredAway() public {
        uint256 tokenId = licenseToken.mint(alice, ipId, licenseTemplate, 1, true, 0);
        // alice transfers to bob
        licenseToken.transferFrom(alice, bob, tokenId);

        // alice no longer has access
        assertFalse(_check(alice, _singleTokenArray(tokenId)));
        // bob now has access
        assertTrue(_check(bob, _singleTokenArray(tokenId)));
    }

    function testCheckReadCondition_AllTokensInvalid() public {
        uint256 wrongOwner = licenseToken.mint(bob, ipId, licenseTemplate, 1, true, 0);
        uint256 wrongIp = licenseToken.mint(alice, otherIpId, licenseTemplate, 2, true, 0);
        uint256 revoked = licenseToken.mint(alice, ipId, licenseTemplate, 3, true, 0);
        licenseToken.revoke(revoked);

        uint256[] memory ids = new uint256[](3);
        ids[0] = wrongOwner;
        ids[1] = wrongIp;
        ids[2] = revoked;
        assertFalse(_check(alice, ids));
    }

    function testCheckReadCondition_NonTransferable_StillValid() public {
        // transferable=false doesn't affect read access, only transferability
        uint256 tokenId = licenseToken.mint(alice, ipId, licenseTemplate, 1, false, 0);
        assertTrue(_check(alice, _singleTokenArray(tokenId)));
    }

    function testCheckReadCondition_WithCommercialRevShare() public {
        uint256 tokenId = licenseToken.mint(alice, ipId, licenseTemplate, 1, true, 500);
        assertTrue(_check(alice, _singleTokenArray(tokenId)));
    }
}
