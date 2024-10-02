// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../src/protocol/UpgradeEntrypoint.sol";

import { EIP1967Helper } from "./utils/EIP1967Helper.sol";
import { InitializableHelper } from "./utils/InitializableHelper.sol";
import { Predeploys } from "../src/libraries/Predeploys.sol";

/**
 * @title EtchInitialState
 * @dev A script + utilities to etch the core contracts
 */
contract EtchInitialState is Script {
    /**
     * @notice Predeploy deployer address, used for each `new` call in this script
     */
    address internal deployer = 0xDDdDddDdDdddDDddDDddDDDDdDdDDdDDdDDDDDDd;


    address internal upgradeAdmin = vm.envAddress("UPGRADE_ADMIN_ADDRESS");
    address internal protocolAdmin = vm.envAddress("ADMIN_ADDRESS");
    string internal dumpPath = getDumpPath();
    bool public saveState = true;

    function disableStateDump() external {
        require(block.chainid == 31337, "Only for local tests");
        saveState = false;
    }

    function getDumpPath() internal view returns (string memory) {
        if (block.chainid == 1513) {
            return "./iliad-state.json";
        } else if (block.chainid == 31337) {
            return "./local-state.json";
        } else {
            revert("Unsupported chain id");
        }
    }

    function run() public {
        require(upgradeAdmin != address(0), "upgradeAdmin not set");
        require(protocolAdmin != address(0), "protocolAdmin not set");

        vm.startPrank(deployer);
        setPredeploys();

        // Reset so its not included state dump
        vm.etch(msg.sender, "");
        vm.resetNonce(msg.sender);
        vm.deal(msg.sender, 0);

        vm.stopPrank();
        if (saveState) {
            vm.dumpState(dumpPath);
        }
    }

    function setProxy(address proxyAddr) internal {
        address impl = Predeploys.getImplAddress(proxyAddr);

        // set impl code to non-zero length, so it passes TransparentUpgradeableProxy constructor check
        // assert it is not already set
        require(impl.code.length == 0, "impl already set");
        vm.etch(impl, "00");

        // new use new, so that the immutable variable the holds the ProxyAdmin proxyAddr is set in properly in bytecode
        address tmp = address(new TransparentUpgradeableProxy(impl, upgradeAdmin, ""));
        vm.etch(proxyAddr, tmp.code);

        // set implempentation storage manually
        EIP1967Helper.setImplementation(proxyAddr, impl);

        // set admin storage, to follow EIP1967 standard
        EIP1967Helper.setAdmin(proxyAddr, EIP1967Helper.getAdmin(tmp));

        // reset impl & tmp
        vm.etch(impl, "");
        vm.etch(tmp, "");

        // can we reset nonce here? we are using "deployer" proxyAddr
        vm.resetNonce(tmp);
    }

    function setPredeploys() internal {
        setProxy(Predeploys.Staking);
        setProxy(Predeploys.Slashing);
        setProxy(Predeploys.Upgrades);

        setStaking();
        setSlashing();
        setUpgrade();
    }

    /**
     * @notice Setup Staking predeploy
     */
    function setStaking() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Staking);

        address tmp = address(new IPTokenStaking(
            1 gwei, // stakingRounding
            1000, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        ));
        console2.log("tpm", tmp);
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        IPTokenStaking(Predeploys.Staking).initialize(protocolAdmin, 1 ether, 1 ether, 1 ether, 7 days);

        console2.log("IPTokenStaking proxy deployed at:", Predeploys.Staking);
        console2.log("IPTokenStaking ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Staking));
        console2.log("IPTokenStaking impl at:", EIP1967Helper.getImplementation(Predeploys.Staking));
        console2.log("IPTokenStaking owner:", IPTokenStaking(Predeploys.Staking).owner());
    }

    /**
     * @notice Setup Slashing predeploy
     */
    function setSlashing() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Slashing);
        address tmp = address(new IPTokenSlashing(Predeploys.Staking));

        console2.log("tpm", tmp);
        vm.etch(impl, tmp.code);

        // reset tmp
        vm.etch(tmp, "");
        vm.store(tmp, 0, "0x");
        vm.resetNonce(tmp);

        InitializableHelper.disableInitializers(impl);
        IPTokenSlashing(Predeploys.Slashing).initialize(protocolAdmin, 1 ether);

        console2.log("IPTokenSlashing proxy deployed at:", Predeploys.Slashing);
        console2.log("IPTokenSlashing ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Slashing));
        console2.log("IPTokenSlashing impl at:", EIP1967Helper.getImplementation(Predeploys.Slashing));
    }

    /**
     * @notice Setup Upgrade predeploy
     */
    function setUpgrade() internal {
        address impl = Predeploys.getImplAddress(Predeploys.Upgrades);
        bytes memory bytecode = type(UpgradeEntrypoint).creationCode;

        vm.etch(Predeploys.Upgrades, bytecode);

        InitializableHelper.disableInitializers(impl);
        UpgradeEntrypoint(Predeploys.Upgrades).initialize(protocolAdmin);

        console2.log("UpgradeEntrypoint proxy deployed at:", Predeploys.Upgrades);
        console2.log("UpgradeEntrypoint ProxyAdmin deployed at:", EIP1967Helper.getAdmin(Predeploys.Upgrades));
        console2.log("UpgradeEntrypoint impl at:", EIP1967Helper.getImplementation(Predeploys.Upgrades));
    }
}
