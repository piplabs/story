// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";

import { InitializableHelper } from "../../script/utils/InitializableHelper.sol";
import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";

contract Test is ForgeTest {
    address private admin = address(this);

    IPTokenStaking internal ipTokenStaking = IPTokenStaking(Predeploys.Staking);
    IPTokenSlashing internal ipTokenSlashing = IPTokenSlashing(Predeploys.Slashing);
    UpgradeEntrypoint internal upgradeEntrypoint = UpgradeEntrypoint(Predeploys.Upgrade);

    function setUp() public virtual {
        // setProxies(Predeploys.Namespace);
        // setStaking();
        // setSlashing();
        // setUpgrade();
    }

    function setStaking() internal {
        address impl = Predeploys.impl(Predeploys.Staking);
        bytes memory stakingImplCreation = abi.encodePacked(
            type(IPTokenStaking).creationCode,
            abi.encode(
                1 ether, // minStakeAmount
                1 ether, // minUnstakeAmount
                1 ether, // minRedelegateAmount
                1 gwei, // stakingRounding
                7 days, // withdrawalAddressChangeInterval
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );
        vm.etch(impl, stakingImplCreation);
        InitializableHelper.disableInitializers(impl);

        vm.etch(Predeploys.Staking, address(new ERC1967Proxy(impl, "")).code);
        EIP1967Helper.setImplementation(Predeploys.Staking, impl);
        EIP1967Helper.setAdmin(Predeploys.Staking, admin);
        IPTokenStaking(Predeploys.Staking).initialize(admin);
    }

    function setSlashing() internal {
        address impl = Predeploys.impl(Predeploys.Slashing);
        bytes memory slashingImplCreation = abi.encodePacked(
            type(IPTokenSlashing).creationCode,
            abi.encode(
                Predeploys.Staking,
                1 ether // unjailFee
            )
        );

        vm.etch(impl, slashingImplCreation);
        InitializableHelper.disableInitializers(impl);

        vm.etch(Predeploys.Slashing, address(new ERC1967Proxy(impl, "")).code);
        EIP1967Helper.setImplementation(Predeploys.Slashing, impl);
        EIP1967Helper.setAdmin(Predeploys.Slashing, admin);
        IPTokenStaking(Predeploys.Slashing).initialize(admin);
    }

    function setUpgrade() internal {
        address impl = Predeploys.impl(Predeploys.Upgrade);
        // UpgradeEntrypoint doesn't have a constructor, so we can use getDeployedCode
        vm.etch(impl, vm.getDeployedCode("UpgradeEntrypoint.sol:UpgradeEntrypoint"));
        InitializableHelper.disableInitializers(impl);

        vm.etch(Predeploys.Upgrade, address(new ERC1967Proxy(impl, "")).code);
        EIP1967Helper.setImplementation(Predeploys.Upgrade, impl);
        EIP1967Helper.setAdmin(Predeploys.Upgrade, admin);
        UpgradeEntrypoint(Predeploys.Upgrade).initialize(admin);
    }

    function setProxies(address ns) internal {
        require(uint32(uint160(ns)) == 0, "invalid namespace");

        for (uint160 i = 1; i <= Predeploys.NamespaceSize; i++) {
            address addr = address(uint160(ns) + i);
            if (Predeploys.notProxied(addr)) {
                continue;
            }

            // For tests, only set proxies for active predeploys
            if (!Predeploys.isActivePredeploy(addr)) {
                return;
            }

            address impl = Predeploys.impl(addr);

            // set impl code to non-zero length, so it passes constructor check
            // assert it is not already set
            require(impl.code.length == 0, "impl already set");
            vm.etch(impl, "00");

            address tmp = address(new ERC1967Proxy(impl, ""));
            vm.etch(addr, tmp.code);

            // set implempentation storage manually
            EIP1967Helper.setImplementation(addr, impl);

            // set admin storage, to follow EIP1967 standard
            // EIP1967Helper.setAdmin(addr, EIP1967Helper.getAdmin(tmp));
            EIP1967Helper.setAdmin(addr, admin);

            // reset impl & tmp
            vm.etch(impl, "");
            vm.etch(tmp, "");
            vm.resetNonce(tmp);
        }

        vm.etch(address(0), "");
    }
}
