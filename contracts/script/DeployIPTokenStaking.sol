// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";

/**
 * @title DeployIPTokenStaking
 * @dev A script + utilities to deploy the IPTokenStaking contract
 */
contract DeployIPTokenStaking is Script {
    function run() public {
        // TODO: read env
        address protocolAccessManagerAddr = address(0x438d47Bc1e184b976fB9EB06Fa4EA27e1247674E);

        require(block.chainid == 1, "only mainnet deployment");
        require(protocolAccessManagerAddr != address(0), "address not set");

        uint256 deployerKey = vm.envUint("IPTOKENSTAKING_DEPLOYER_KEY");

        vm.startBroadcast(deployerKey);
        IPTokenStaking ipTokenStaking = new IPTokenStaking(
            protocolAccessManagerAddr,
            1 ether, // minStakeAmount
            1 ether, // minUnstakeAmount
            1 ether, // minRedelegateAmount
            1 gwei, // stakingRounding
            7 days, // withdrawalAddressChangeInterval
            1000, // defaultCommissionRate, 10%
            5000, // defaultMaxCommissionRate, 50%
            500 // defaultMaxCommissionChangeRate, 5%
        );
        vm.stopBroadcast();

        console2.log("IPTokenStaking deployed at:", address(ipTokenStaking));
    }
}
