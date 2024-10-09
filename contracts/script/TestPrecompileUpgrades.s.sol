// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../src/protocol/UpgradeEntrypoint.sol";
import { ITransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { EIP1967Helper } from "./utils/EIP1967Helper.sol";
import { Predeploys } from "../src/libraries/Predeploys.sol";

abstract contract MockNewFeatures {
    function foo() external pure returns(string memory) {
        return "bar";
    }
}

contract IPTokenStakingV2 is IPTokenStaking, MockNewFeatures {
    constructor(
        uint256 stakingRounding,
        uint32 defaultCommissionRate,
        uint32 defaultMaxCommissionRate,
        uint32 defaultMaxCommissionChangeRate
    ) IPTokenStaking(stakingRounding, defaultCommissionRate, defaultMaxCommissionRate, defaultMaxCommissionChangeRate) {

    }
}

contract IPTokenSlashingV2 is IPTokenSlashing, MockNewFeatures {
    constructor(address ipTokenStaking) IPTokenSlashing(ipTokenStaking) {}
}

contract UpgradeEntrypointV2 is UpgradeEntrypoint, MockNewFeatures {
}

/**
 * @title TestPrecompileUpgrades
 * @dev A script to test upgrading the precompile contracts
 */
contract TestPrecompileUpgrades is Script {
    // To run the script:
    // - Dry run
    // forge script script/DeployIPTokenSlashing.s.sol --fork-url <fork-url>
    //
    // - Deploy (OK for devnet)
    // forge script script/DeployIPTokenSlashing.s.sol --fork-url <fork-url> --broadcast
    //
    // - Deploy and Verify (for testnet)
    function run() public {
        // Read env for admin address
        uint256 upgradeKey = vm.envUint("UPGRADE_ADMIN_KEY");
        address upgrader = vm.addr(upgradeKey);
        console2.log("upgrader", upgrader);
        vm.startBroadcast(upgradeKey);

        // ---- Staking
        address newImpl = address(new IPTokenStakingV2(
            1 gwei, // stakingRounding
            1000, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        ));
        ProxyAdmin proxyAdmin = ProxyAdmin(
            EIP1967Helper.getAdmin(Predeploys.Staking)
        );
        console2.log("staking proxy admin", address(proxyAdmin));
        console2.log("staking proxy admin owner", proxyAdmin.owner());
        proxyAdmin.upgradeAndCall(
            ITransparentUpgradeableProxy(Predeploys.Staking),
            newImpl,
            ""
        );
        if (EIP1967Helper.getImplementation(Predeploys.Staking) != newImpl) {
            revert("Staking not upgraded");
        }
        if (keccak256(abi.encode(IPTokenStakingV2(Predeploys.Staking).foo())) != keccak256(abi.encode("bar"))) {
            revert("Upgraded to wrong iface");
        }

        // ---- Slashing
        newImpl = address(new IPTokenSlashingV2(
            Predeploys.Staking
        ));
        proxyAdmin = ProxyAdmin(
            EIP1967Helper.getAdmin(Predeploys.Slashing)
        );
        console2.log("slashing proxy admin", address(proxyAdmin));
        console2.log("slashing proxy admin owner", proxyAdmin.owner());
        proxyAdmin.upgradeAndCall(
            ITransparentUpgradeableProxy(Predeploys.Slashing),
            newImpl,
            ""
        );
        if (EIP1967Helper.getImplementation(Predeploys.Slashing) != newImpl) {
            revert("Slashing not upgraded");
        }
        if (keccak256(abi.encode(IPTokenSlashingV2(Predeploys.Slashing).foo())) != keccak256(abi.encode("bar"))) {
            revert("Upgraded to wrong iface");
        }

        // ---- Upgrades
        newImpl = address(new UpgradeEntrypointV2());
        proxyAdmin = ProxyAdmin(
            EIP1967Helper.getAdmin(Predeploys.Upgrades)
        );
        console2.log("upgrades proxy admin", address(proxyAdmin));
        console2.log("upgrades proxy admin owner", proxyAdmin.owner());

        proxyAdmin.upgradeAndCall(
            ITransparentUpgradeableProxy(Predeploys.Upgrades),
            newImpl,
            ""
        );
        if (keccak256(abi.encode(UpgradeEntrypointV2(Predeploys.Upgrades).foo())) != keccak256(abi.encode("bar"))) {
            revert("Upgraded to wrong iface");
        }
        if (EIP1967Helper.getImplementation(Predeploys.Upgrades) != newImpl) {
            revert("UpgradeEntrypoint not upgraded");
        }
        vm.stopBroadcast();
    }
}