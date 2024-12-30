// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import "forge-std/Script.sol";
import { IPTokenStaking } from "src/protocol/IPTokenStaking.sol";
import { console2 } from "forge-std/console2.sol";

// forge script src/script/DeployIPTokenStaking.s.sol --rpc-url <RPC_URL> --private-key <PRIVATE_KEY> --broadcast

// deploy and verify
// forge script src/script/DeployIPTokenStaking.s.sol --fork-url http://r1-d.odyssey-devnet.storyrpc.io:8545 -vvvv --private-key <PRIVATE_KEY> --priority-gas-price 1 --legacy --verify  --verifier=blockscout --verifier-url=https://devnet.storyscan.xyz/api --broadcast 
contract DeployIPTokenStaking is Script {
    function run() public {
        // Retrieve constructor arguments
        uint256 stakingRounding = 1 gwei;
        uint256 defaultMinFee = 1 ether;

        // Deploy the contract
        vm.startBroadcast();
        
        IPTokenStaking staking = new IPTokenStaking(stakingRounding, defaultMinFee);

        console2.log("IPTokenStaking deployed at:", address(staking));

        vm.stopBroadcast();
    }
}