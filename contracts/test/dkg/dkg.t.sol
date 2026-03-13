// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { DKG } from "../../src/protocol/DKG.sol";
import { SGXValidationHook } from "../../src/protocol/SGXValidationHook.sol";
import { IDKG } from "../../src/interfaces/IDKG.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG internal dkg;

    // Genesis enclave type set by GenerateAlloc.setSGXValidationHook()
    bytes32 internal constant GENESIS_ENCLAVE_TYPE = bytes32(uint256(1));

    // NOTE: Must match TEST_DKG_OWNER in GenerateAlloc.s.sol. NEVER use in production.
    address internal constant TEST_DKG_OWNER = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;

    function setUp() public virtual override {
        super.setUp();
        dkg = DKG(Predeploys.DKG);
    }

    function testDKG_Initialize() public view {
        assertEq(dkg.minReqRegisteredParticipants(), 3);
        assertEq(dkg.minReqFinalizedParticipants(), 3);
        assertEq(dkg.operationalThreshold(), 500);
        assertEq(dkg.fee(), 1 ether);
    }

    function testDKG_GenesisEnclaveTypeWhitelisted() public view {
        assertTrue(dkg.isEnclaveTypeWhitelisted(GENESIS_ENCLAVE_TYPE));

        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(GENESIS_ENCLAVE_TYPE);
        assertEq(data.codeCommitment, hex"8518404aed711077ddf6738f51c76a87c0258f158e2e7b7ee6f000369f9394d3");
        assertTrue(data.validationHookAddr != address(0));
    }

    function testDKG_WhitelistNewEnclaveType() public {
        address sgxHookImpl = address(new SGXValidationHook(Predeploys.DKG));
        address sgxHookProxy = address(
            new ERC1967Proxy(
                sgxHookImpl,
                abi.encodeCall(SGXValidationHook.initialize, (TEST_DKG_OWNER, address(2000), 1))
            )
        );

        bytes32 newEnclaveType = bytes32(uint256(2));
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: sgxHookProxy
        });
        vm.prank(TEST_DKG_OWNER);
        dkg.whitelistEnclaveType(newEnclaveType, enclaveTypeData, true);

        assertTrue(dkg.isEnclaveTypeWhitelisted(newEnclaveType));
        assertEq(dkg.enclaveTypeData(newEnclaveType).validationHookAddr, sgxHookProxy);
        assertEq(dkg.enclaveTypeData(newEnclaveType).codeCommitment, bytes32(uint256(42)));
    }

    function testDKG_WhitelistEnclaveType_RevertIfNotOwner() public {
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)),
            validationHookAddr: address(1)
        });

        vm.expectRevert();
        dkg.whitelistEnclaveType(bytes32(uint256(99)), enclaveTypeData, true);
    }
}
