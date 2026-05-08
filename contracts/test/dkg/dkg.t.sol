// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { DKG } from "../../src/protocol/DKG.sol";
import { SGXValidationHook } from "../../src/protocol/SGXValidationHook.sol";
import { TDXValidationHook } from "../../src/protocol/TDXValidationHook.sol";
import { IDKG } from "../../src/interfaces/IDKG.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

contract DKGTest is Test {
    DKG internal dkg;

    // Genesis enclave types from GenerateAlloc:
    //   - SGX = bytes32(uint256(1)) via setSGXValidationHook()
    //   - TDX = bytes32(uint256(2)) via setTDXValidationHook()
    bytes32 internal constant SGX_ENCLAVE_TYPE = bytes32(uint256(1));
    bytes32 internal constant TDX_ENCLAVE_TYPE = bytes32(uint256(2));

    // Devnet config: USE_DEPLOYER_AS_OWNER=true makes the foundry test deployer
    // (the standard anvil[0] address) the owner of all admin contracts so that
    // whitelistEnclaveType / scheduleUpgrade can be called without timelock.
    address internal constant DEVNET_OWNER = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;

    // Code commitments planted by GenerateAlloc constants.
    bytes32 internal constant SGX_CODE_COMMITMENT =
        hex"cfac25c990dc7517d9704fc51e65199a379332802f4c15d7fb966cba0813301c";
    bytes32 internal constant TDX_CODE_COMMITMENT =
        hex"0000000000000000000000000000000000000000000000000000000000000002";

    function setUp() public virtual override {
        super.setUp();
        dkg = DKG(Predeploys.DKG);
    }

    function testDKG_Initialize() public view {
        // SGX + TDX 2-validator devnet uses n=2 t=2 (100% threshold).
        assertEq(dkg.minReqRegisteredParticipants(), 2);
        assertEq(dkg.minReqFinalizedParticipants(), 2);
        assertEq(dkg.operationalThreshold(), 1000);
        assertEq(dkg.fee(), 1 ether);
    }

    function testDKG_GenesisSGXEnclaveTypeWhitelisted() public view {
        assertTrue(dkg.isEnclaveTypeWhitelisted(SGX_ENCLAVE_TYPE));

        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(SGX_ENCLAVE_TYPE);
        assertEq(data.codeCommitment, SGX_CODE_COMMITMENT);
        assertTrue(data.validationHookAddr != address(0));
    }

    function testDKG_GenesisTDXEnclaveTypeWhitelisted() public view {
        assertTrue(dkg.isEnclaveTypeWhitelisted(TDX_ENCLAVE_TYPE));

        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(TDX_ENCLAVE_TYPE);
        assertEq(data.codeCommitment, TDX_CODE_COMMITMENT);
        assertTrue(data.validationHookAddr != address(0));
    }

    function testDKG_GenesisSGXValidationHook() public view {
        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(SGX_ENCLAVE_TYPE);
        SGXValidationHook sgxHook = SGXValidationHook(data.validationHookAddr);

        assertEq(sgxHook.owner(), DEVNET_OWNER);
        assertEq(sgxHook.DKG(), Predeploys.DKG);
        assertEq(sgxHook.automataValidationAddr(), address(uint160(1000)));
        assertEq(sgxHook.tcbEvaluationDataNumber(), 0);
    }

    function testDKG_GenesisTDXValidationHook() public view {
        IDKG.EnclaveTypeData memory data = dkg.enclaveTypeData(TDX_ENCLAVE_TYPE);
        TDXValidationHook tdxHook = TDXValidationHook(data.validationHookAddr);

        assertEq(tdxHook.owner(), DEVNET_OWNER);
        assertEq(tdxHook.DKG(), Predeploys.DKG);
        assertEq(tdxHook.automataValidationAddr(), address(uint160(1000)));
        assertEq(tdxHook.tcbEvaluationDataNumber(), 0);
    }

    function testDKG_WhitelistNewEnclaveType() public {
        address sgxHookImpl = address(new SGXValidationHook(Predeploys.DKG));
        address sgxHookProxy = address(
            new ERC1967Proxy(
                sgxHookImpl,
                abi.encodeCall(SGXValidationHook.initialize, (DEVNET_OWNER, address(2000), 1))
            )
        );

        bytes32 newEnclaveType = bytes32(uint256(99));
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(42)),
            validationHookAddr: sgxHookProxy
        });

        // Devnet uses deployer-as-owner; whitelist is a direct call, no timelock.
        vm.prank(DEVNET_OWNER);
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
