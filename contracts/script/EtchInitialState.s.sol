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

/**
 * @title EtchInitialState
 * @dev A script + utilities to etch the core contracts
 */
contract EtchInitialState is Script {
    /**
     * @notice Predeploy deployer address, used for each `new` call in this script
     */
    address internal deployer = 0xDDdDddDdDdddDDddDDddDDDDdDdDDdDDdDDDDDDd;

    address internal constant StakingProxyAddr = 0xCCcCcC0000000000000000000000000000000001;
    address internal constant SlashingProxyAddr = 0xCccCCC0000000000000000000000000000000002;
    address internal constant UpgradeProxyAddr = 0xccCCcc0000000000000000000000000000000003;

    address internal upgradeAdmin = vm.envAddress("UPGRADE_ADMIN_ADDRESS");
    address internal protocolAdmin = vm.envAddress("ADMIN_ADDRESS");
    string internal dumpPath = getDumpPath();

    function getDumpPath() internal view returns (string memory) {
        if (block.chainid == 1513) {
            return "./iliad-state.json";
        } else {
            revert("Unsupported chain id");
        }
    }

    function run() public {
        require(block.chainid == 1513, "Wrong chain id");

        require(upgradeAdmin != address(0), "upgradeAdmin not set");
        require(protocolAdmin != address(0), "protocolAdmin not set");

        vm.startPrank(deployer);
        setPredeploys();

        // Reset so its not included state dump
        vm.etch(msg.sender, "");
        vm.resetNonce(msg.sender);
        vm.deal(msg.sender, 0);

        vm.stopPrank();

        vm.dumpState(dumpPath);
    }

    /**
     * @notice Return implementation address for a proxied predeploy
     */
    function getImplAddress(address addr) internal pure returns (address) {
        // max uint160 is odd, which gives us unique implementation for each predeploy
        return address(type(uint160).max - uint160(addr));
    }

    function setProxy(address proxyAddr) internal {
        address impl = getImplAddress(proxyAddr);

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
        setProxy(StakingProxyAddr);
        setProxy(SlashingProxyAddr);
        setProxy(UpgradeProxyAddr);

        setStaking();
        setSlashing();
        setUpgrade();
    }

    /**
     * @notice Setup Staking predeploy
     */
    function setStaking() internal {
        address impl = getImplAddress(StakingProxyAddr);

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
        IPTokenStaking(StakingProxyAddr).initialize(protocolAdmin, 1 ether, 1 ether, 1 ether, 7 days);

        console2.log("IPTokenStaking proxy deployed at:", StakingProxyAddr);
        console2.log("IPTokenStaking ProxyAdmin deployed at:", EIP1967Helper.getAdmin(StakingProxyAddr));
        console2.log("IPTokenStaking impl at:", EIP1967Helper.getImplementation(StakingProxyAddr));
    }

    /**
     * @notice Setup Slashing predeploy
     */
    function setSlashing() internal {
        address impl = getImplAddress(SlashingProxyAddr);
        bytes memory bytecode = type(IPTokenSlashing).creationCode;
        // set IPTokenStaking address in constructor
        bytes memory constructorArgs = abi.encode(StakingProxyAddr);

        // Combine bytecode and constructor args
        bytes memory deployCode = abi.encodePacked(bytecode, constructorArgs);
        vm.etch(SlashingProxyAddr, deployCode);

        InitializableHelper.disableInitializers(impl);
        IPTokenSlashing(SlashingProxyAddr).initialize(protocolAdmin, 1 ether);

        console2.log("IPTokenSlashing proxy deployed at:", SlashingProxyAddr);
        console2.log("IPTokenSlashing ProxyAdmin deployed at:", EIP1967Helper.getAdmin(SlashingProxyAddr));
        console2.log("IPTokenSlashing impl at:", EIP1967Helper.getImplementation(SlashingProxyAddr));
    }

    /**
     * @notice Setup Upgrade predeploy
     */
    function setUpgrade() internal {
        address impl = getImplAddress(UpgradeProxyAddr);
        bytes memory bytecode = type(UpgradeEntrypoint).creationCode;

        vm.etch(UpgradeProxyAddr, bytecode);

        InitializableHelper.disableInitializers(impl);
        UpgradeEntrypoint(UpgradeProxyAddr).initialize(protocolAdmin);

        console2.log("UpgradeEntrypoint proxy deployed at:", UpgradeProxyAddr);
        console2.log("UpgradeEntrypoint ProxyAdmin deployed at:", EIP1967Helper.getAdmin(UpgradeProxyAddr));
        console2.log("UpgradeEntrypoint impl at:", EIP1967Helper.getImplementation(UpgradeProxyAddr));
    }
}
