// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import "forge-std/Script.sol";
import { MyToken } from "src/protocol/MyToken.sol";
import { console2 } from "forge-std/console2.sol";

// forge script src/script/DeployIPTokenStaking.s.sol --rpc-url <RPC_URL> --private-key <PRIVATE_KEY> --broadcast

// deploy and verify

//  forge script script/foundry/deployment/reserved-proxy/DeployMyToken.s.sol  --fork-url http://r1-d.odyssey-devnet.storyrpc.io:8545 -vvvv --private-key < > --priority-gas-price 1 --legacy --verify  --verifier=blockscout --verifier-url=https://devnet.storyscan.xyz/api --broadcast
contract DeployMyToken is Script {
    function run() public {

        // Deploy the contract
        vm.startBroadcast();
        
        MyToken mytoken = new MyToken();

        console2.log("MyToken deployed at:", address(mytoken));

        vm.stopBroadcast();
    }
}