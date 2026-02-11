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
    address internal sgxHookProxy;

    function setUp() public virtual override {
        super.setUp();

        dkg = DKG(Predeploys.DKG);

        address automataValidationAddr = address(1000); // TODO: set this to the actual automata validation address
        uint32 tcbEvaluationDataNumber = 0;
        address sgxHookImpl = address(new SGXValidationHook(Predeploys.DKG));
        sgxHookProxy = address(
            new ERC1967Proxy(
                sgxHookImpl,
                abi.encodeCall(
                    SGXValidationHook.initialize,
                    (address(timelock), automataValidationAddr, tcbEvaluationDataNumber)
                )
            )
        );

        // whitelist the SGX validation hook in DKG contract
        IDKG.EnclaveTypeData memory enclaveTypeData = IDKG.EnclaveTypeData({
            codeCommitment: bytes32(uint256(1)), // TODO: set this to the actual code commitment
            validationHookAddr: sgxHookProxy
        });
        performTimelocked(
            address(dkg),
            abi.encodeWithSelector(
                DKG.whitelistEnclaveType.selector,
                bytes32("SGX"),
                enclaveTypeData,
                true
            )
        );
    }

    function testDKG_Initialize() public {
        assertEq(dkg.minReqRegisteredParticipants(), 3);
        assertEq(dkg.minReqFinalizedParticipants(), 3);
        assertEq(dkg.operationalThreshold(), 670);
        assertEq(dkg.fee(), 1 ether);
        assertEq(dkg.enclaveTypeData(bytes32("SGX")).validationHookAddr, sgxHookProxy);
        assertEq(dkg.enclaveTypeData(bytes32("SGX")).codeCommitment, bytes32(uint256(1)));
        assertEq(dkg.isEnclaveTypeWhitelisted(bytes32("SGX")), true);
    }
}
