// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Test } from "../utils/Test.sol";
import { IPTokenStakingV2 } from "../../src/upgrades/IPTokenStakingV2.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";

import { console2 } from "forge-std/console2.sol";

/**
 * @title IPTokenStakingV2Test
 * @dev A test for the IPTokenStakingV2 contract
 */
contract IPTokenStakingV2Test is Test {

    IPTokenStakingV2 ipTokenStakingProxy;
    address safeGovernanceMultisig;
    
    function setUp() public override {
        // Fork the desired network where UMA contracts are deployed
        uint256 forkId = vm.createFork("https://mainnet.storyrpc.io/");
        vm.selectFork(forkId);

        // Mainnet related addresses
        ipTokenStakingProxy = IPTokenStakingV2(0xCCcCcC0000000000000000000000000000000001);
        safeGovernanceMultisig = 0xF07cA4b61022F0399C1511E7E668A57567f2138B;
        timelock = TimelockController(payable(0x6c7FA8DF1B8Dc29a7481Bb65ad590D2D16787a82));

        // deploy new implementation
        IPTokenStakingV2 newImpl = new IPTokenStakingV2(1 ether, 256);

        // upgrade the proxy to the new implementation
        ProxyAdmin proxyAdmin = ProxyAdmin(EIP1967Helper.getAdmin(address(ipTokenStakingProxy)));
        console2.log("proxyAdmin", address(proxyAdmin));
        vm.startPrank(safeGovernanceMultisig);
        timelock.schedule(
            address(proxyAdmin),
            0,
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(address(ipTokenStakingProxy)),
                newImpl,
                ""
            ),
            bytes32(0),
            bytes32(0),
            timelock.getMinDelay()
        );

        vm.warp(block.timestamp + timelock.getMinDelay() + 1);

        timelock.execute(
            address(proxyAdmin),
            0,
            abi.encodeWithSelector(
                ProxyAdmin.upgradeAndCall.selector,
                ITransparentUpgradeableProxy(address(ipTokenStakingProxy)),
                newImpl,
                ""
            ),
            bytes32(0),
            bytes32(0)
        );
        vm.stopPrank();
    }

    function test_IPTokenStakingV2_pause() public {
        assertEq(ipTokenStakingProxy.paused(), false);

        vm.startPrank(safeGovernanceMultisig);
        timelock.schedule(
            address(ipTokenStakingProxy),
            0,
            abi.encodeWithSelector(IPTokenStakingV2.pause.selector),
            bytes32(0),
            bytes32(0),
            0
        );
        timelock.execute(
            address(ipTokenStakingProxy),
            0,
            abi.encodeWithSelector(IPTokenStakingV2.pause.selector),
            bytes32(0),
            bytes32(0)
        );
        vm.stopPrank();

        assertEq(ipTokenStakingProxy.paused(), true);
    }
}