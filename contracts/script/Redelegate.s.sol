// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";
import { IPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { IIPTokenStaking } from "../src/protocol/IPTokenStaking.sol";
import { Predeploys } from "../src/libraries/Predeploys.sol";

contract Redelegate is Script {

    function run() external {
        // Load the private key from the environment variable
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        console2.log("Deployer:", deployer);
        // Start broadcasting transactions
        vm.startBroadcast(deployerPrivateKey);

        // Create an instance of the Staking precompile interface
        IPTokenStaking staking = IPTokenStaking(Predeploys.Staking);

        // Define redelegate parameters
        bytes memory delegatorUncmpPubkey = hex"04e02e9bde528ab4c1ac9699483107a168af0592da5d09ebdeba96199242aef9abd9e9e86cbfcda373d59404f1deb167ac2b09808e0b140e62c5fff38eacfccce9";
        bytes memory validatorUncmpSrcPubkey = hex"04e02e9bde528ab4c1ac9699483107a168af0592da5d09ebdeba96199242aef9abd9e9e86cbfcda373d59404f1deb167ac2b09808e0b140e62c5fff38eacfccce9";
        bytes memory validatorUncmpDstPubkey = hex"04cc1401aac253c1b81512dd4f3d44205de73643300f24bdb6a7c9c92fa950c93b41aa5ce1e86c7de957d6fbdf18fa7e1ef8b1afcf405b0f2629b31fb1672f0502";
        uint256 stakeAmount = 1025 ether; // Amount to redelegate, adjust as needed
        uint256 delegationId = 0; // Delegation ID, adjust if needed
        console2.log("Redelegate");
        // Call redelegate function
        console2.log("Staking address:", Predeploys.Staking);

        /*staking.stake{ value: stakeAmount }(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            IIPTokenStaking.StakingPeriod.FLEXIBLE,
            ""
        );*/
        staking.redelegate(
            delegatorUncmpPubkey,
            validatorUncmpSrcPubkey,
            validatorUncmpDstPubkey,
            delegationId,
            stakeAmount
        );

        // Stop broadcasting transactions
        vm.stopBroadcast();
    }
}
