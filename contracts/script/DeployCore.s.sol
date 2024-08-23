// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../src/protocol/UpgradeEntrypoint.sol";

/**
 * @title DeployCore
 * @dev A script + utilities to deploy the core contracts
 */
contract DeployCore is Script {
    function run() public {
        // TODO: read env
        address protocolAccessManagerAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);
        require(protocolAccessManagerAddr != address(0), "address not set");

        uint256 deployerKey = vm.envUint("IPTOKENSTAKING_DEPLOYER_KEY");

        vm.startBroadcast(deployerKey);

        address impl = address(
            new IPTokenStaking(
                1 gwei, // stakingRounding
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );
        IPTokenStaking ipTokenStaking = IPTokenStaking(address(new ERC1967Proxy(impl, "")));
        ipTokenStaking.initialize(
            protocolAccessManagerAddr,
            1 ether, // minStakeAmount
            1 ether, // minUnstakeAmount
            1 ether, // minRedelegateAmount
            7 days // withdrawalAddressInterval
        );

        impl = address(new IPTokenSlashing(address(ipTokenStaking)));
        IPTokenSlashing ipTokenSlashing = IPTokenSlashing(address(new ERC1967Proxy(impl, "")));
        ipTokenSlashing.initialize(
            protocolAccessManagerAddr,
            1 ether // unjailFee
        );

        impl = address(new UpgradeEntrypoint());
        UpgradeEntrypoint upgradeEntrypoint = UpgradeEntrypoint(address(new ERC1967Proxy(impl, "")));
        upgradeEntrypoint.initialize(protocolAccessManagerAddr);

        vm.stopBroadcast();

        console2.log("IPTokenStaking deployed at:", address(ipTokenStaking));
        console2.log("IPTokenSlashing deployed at:", address(ipTokenSlashing));
        console2.log("UpgradeEntrypoint deployed at:", address(upgradeEntrypoint));
    }
}
