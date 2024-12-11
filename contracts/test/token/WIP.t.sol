// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { Test } from "../utils/Test.sol";
import { WIP } from "../../src/token/WIP.sol";

contract ContractWithoutReceive {}

contract WIPTest is Test {
    function testMetadata() public view {
        assertEq(wip.name(), "Wrapped IP");
        assertEq(wip.symbol(), "WIP");
        assertEq(wip.decimals(), 18);
    }

    function testFallbackDeposit() public {
        assertEq(wip.balanceOf(address(this)), 0);
        assertEq(wip.totalSupply(), 0);

        (bool success, ) = address(wip).call{ value: 1 ether }("");
        assertTrue(success);

        assertEq(wip.balanceOf(address(this)), 1 ether);
        assertEq(wip.totalSupply(), 1 ether);
    }

    function testDeposit() public {
        assertEq(wip.balanceOf(address(this)), 0);
        assertEq(wip.totalSupply(), 0);

        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(address(this)), 1 ether);
        assertEq(wip.totalSupply(), 1 ether);
    }

    function testWithdraw() public {
        uint256 startingBalance = address(this).balance;

        wip.deposit{ value: 1 ether }();

        wip.withdraw(1 ether);

        uint256 balanceAfterWithdraw = address(this).balance;

        assertEq(balanceAfterWithdraw, startingBalance);
        assertEq(wip.balanceOf(address(this)), 0);
        assertEq(wip.totalSupply(), 0);
    }

    function testPartialWithdraw() public {
        wip.deposit{ value: 1 ether }();

        uint256 balanceBeforeWithdraw = address(this).balance;

        wip.withdraw(0.5 ether);

        uint256 balanceAfterWithdraw = address(this).balance;

        assertEq(balanceAfterWithdraw, balanceBeforeWithdraw + 0.5 ether);
        assertEq(wip.balanceOf(address(this)), 0.5 ether);
        assertEq(wip.totalSupply(), 0.5 ether);
    }

    function testWithdrawToContractWithoutReceiveReverts() public {
        address owner = address(new ContractWithoutReceive());

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(WIP.IPTransferFailed.selector);
        vm.prank(owner);
        wip.withdraw(1 ether);
    }

    function testFallbackDeposit(uint256 amount) public {
        amount = _bound(amount, 0, address(this).balance);

        assertEq(wip.balanceOf(address(this)), 0);
        assertEq(wip.totalSupply(), 0);

        (bool success, ) = address(wip).call{ value: amount }("");
        assertTrue(success);

        assertEq(wip.balanceOf(address(this)), amount);
        assertEq(wip.totalSupply(), amount);
    }

    function testDeposit(uint256 amount) public {
        amount = _bound(amount, 0, address(this).balance);

        assertEq(wip.balanceOf(address(this)), 0);
        assertEq(wip.totalSupply(), 0);

        wip.deposit{ value: amount }();

        assertEq(wip.balanceOf(address(this)), amount);
        assertEq(wip.totalSupply(), amount);
    }

    function testWithdraw(uint256 depositAmount, uint256 withdrawAmount) public {
        depositAmount = _bound(depositAmount, 0, address(this).balance);
        withdrawAmount = _bound(withdrawAmount, 0, depositAmount);

        wip.deposit{ value: depositAmount }();

        uint256 balanceBeforeWithdraw = address(this).balance;

        wip.withdraw(withdrawAmount);

        uint256 balanceAfterWithdraw = address(this).balance;

        assertEq(balanceAfterWithdraw, balanceBeforeWithdraw + withdrawAmount);
        assertEq(wip.balanceOf(address(this)), depositAmount - withdrawAmount);
        assertEq(wip.totalSupply(), depositAmount - withdrawAmount);
    }

    function testTransferToZeroAddressReverts() public {
        address owner = address(0x123);

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(abi.encodeWithSelector(WIP.ERC20InvalidReceiver.selector, address(0)));
        vm.prank(owner);
        wip.transfer(address(0), 1 ether);
    }

    function testTransferToWIPContractReverts() public {
        address owner = address(0x123);

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(abi.encodeWithSelector(WIP.ERC20InvalidReceiver.selector, address(wip)));
        vm.prank(owner);
        wip.transfer(address(wip), 1 ether);
    }

    function testTransferFromReceiverIsWIPContractReverts() public {
        address owner = address(0x123);

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(abi.encodeWithSelector(WIP.ERC20InvalidReceiver.selector, address(wip)));
        vm.prank(owner);
        wip.transferFrom(owner, address(wip), 1 ether);
    }

    function testTransferFromReceiverIsZeroAddressReverts() public {
        address owner = address(0x123);

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(abi.encodeWithSelector(WIP.ERC20InvalidReceiver.selector, address(0)));
        vm.prank(owner);
        wip.transferFrom(owner, address(0), 1 ether);
    }

    function testApprovalToSelfReverts() public {
        address owner = address(0x123);

        vm.deal(owner, 1 ether);

        vm.prank(owner);
        wip.deposit{ value: 1 ether }();

        assertEq(wip.balanceOf(owner), 1 ether);

        vm.expectRevert(abi.encodeWithSelector(WIP.ERC20InvalidSpender.selector, owner));
        vm.prank(owner);
        wip.approve(owner, 1 ether);
    }

    receive() external payable {}
}
