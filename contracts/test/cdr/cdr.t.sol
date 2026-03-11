// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { CDR } from "../../src/protocol/CDR.sol";
import { ICDR } from "../../src/interfaces/ICDR.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Test } from "../utils/Test.sol";

contract CDRTest is Test {
    CDR internal cdr;

    function setUp() public virtual override {
        super.setUp();
        cdr = CDR(Predeploys.CDR);
    }

    function testCDR_Initialize() public view {
        assertEq(cdr.owner(), address(timelock));
        assertEq(cdr.baseFee(), 0);
        assertEq(cdr.writeFee(), 0);
        assertEq(cdr.readFee(), 0);
        assertEq(cdr.allocateFee(), 0);
        assertEq(cdr.uuid(), 0);
    }

    function testCDR_SetBaseFee() public {
        uint256 newBaseFee = 0.1 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setBaseFee.selector, newBaseFee));
        assertEq(cdr.baseFee(), newBaseFee);
    }

    function testCDR_SetWriteFee() public {
        uint256 newWriteFee = 0.2 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setWriteFee.selector, newWriteFee));
        assertEq(cdr.writeFee(), newWriteFee);
    }

    function testCDR_SetReadFee() public {
        uint256 newReadFee = 0.3 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setReadFee.selector, newReadFee));
        assertEq(cdr.readFee(), newReadFee);
    }

    function testCDR_SetAllocateFee() public {
        uint256 newAllocateFee = 0.4 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, newAllocateFee));
        assertEq(cdr.allocateFee(), newAllocateFee);
    }

    function testCDR_SetFee_RevertIfNotOwner() public {
        vm.expectRevert();
        cdr.setBaseFee(1 ether);

        vm.expectRevert();
        cdr.setWriteFee(1 ether);

        vm.expectRevert();
        cdr.setReadFee(1 ether);

        vm.expectRevert();
        cdr.setAllocateFee(1 ether);
    }

    function testCDR_Allocate() public {
        address writeCondition = address(0xAAA);
        address readCondition = address(0xBBB);
        bytes memory writeConditionData = abi.encode("write");
        bytes memory readConditionData = abi.encode("read");

        uint32 vaultUuid = cdr.allocate(true, writeCondition, readCondition, writeConditionData, readConditionData);

        assertEq(vaultUuid, 0);

        ICDR.Vault memory vault = cdr.vaults(vaultUuid);
        assertEq(vault.updatable, true);
        assertEq(vault.writeConditionAddr, writeCondition);
        assertEq(vault.readConditionAddr, readCondition);
        assertEq(vault.writeConditionData, writeConditionData);
        assertEq(vault.readConditionData, readConditionData);
        assertEq(vault.encryptedData.length, 0);
    }

    function testCDR_Allocate_RevertIfBothConditionsZero() public {
        vm.expectRevert("Invalid condition address");
        cdr.allocate(true, address(0), address(0), "", "");
    }

    function testCDR_Allocate_WithFee() public {
        uint256 fee = 0.5 ether;
        performTimelocked(address(cdr), abi.encodeWithSelector(CDR.setAllocateFee.selector, fee));

        // Should revert without fee
        vm.expectRevert("CDR: Invalid fee amount");
        cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");

        // Should succeed with correct fee
        vm.deal(address(this), fee);
        cdr.allocate{ value: fee }(true, address(0xAAA), address(0xBBB), "", "");
    }

    function testCDR_Allocate_IncrementingUuid() public {
        uint32 uuid0 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");
        uint32 uuid1 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");
        uint32 uuid2 = cdr.allocate(true, address(0xAAA), address(0xBBB), "", "");

        assertEq(uuid0, 0);
        assertEq(uuid1, 1);
        assertEq(uuid2, 2);
        assertEq(cdr.uuid(), 3);
    }
}
