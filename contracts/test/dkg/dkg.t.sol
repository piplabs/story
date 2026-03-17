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

    function setUp() public virtual override {
        super.setUp();
        dkg = DKG(Predeploys.DKG);
    }

    function testDKG_Initialize() public view {
        assertEq(dkg.minReqRegisteredParticipants(), 3);
        assertEq(dkg.minReqFinalizedParticipants(), 3);
        assertEq(dkg.operationalThreshold(), 670);
        assertEq(dkg.fee(), 1 ether);
    }

    function testDKG_GenesisEnclaveTypeWhitelisted() public view {
        assertTrue(dkg.isEnclaveTypeWhitelisted(GENESIS_ENCLAVE_TYPE));

        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(GENESIS_ENCLAVE_TYPE);
        assertEq(data.codeCommitment, bytes32(uint256(1)));
        assertTrue(data.validationHookAddr != address(0));
    }

    function testDKG_GenesisSGXValidationHook() public view {
        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(GENESIS_ENCLAVE_TYPE);
        SGXValidationHook sgxHook = SGXValidationHook(data.validationHookAddr);

        assertEq(sgxHook.owner(), address(timelock));
        assertEq(sgxHook.DKG(), Predeploys.DKG);
        assertEq(sgxHook.automataValidationAddr(), address(uint160(1000)));
        assertEq(sgxHook.tcbEvaluationDataNumber(), 0);
    }

    function testDKG_WhitelistNewEnclaveType() public {
        address sgxHookImpl = address(new SGXValidationHook(Predeploys.DKG));
        address sgxHookProxy = address(
            new ERC1967Proxy(
                sgxHookImpl,
                abi.encodeCall(SGXValidationHook.initialize, (address(timelock), address(2000), 1))
            )
        );

        bytes32 newEnclaveType = bytes32(uint256(2));
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: sgxHookProxy
        });
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(DKG.whitelistEnclaveType.selector, newEnclaveType, enclaveTypeData, true)
        );

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
