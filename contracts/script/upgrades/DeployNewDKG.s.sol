// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */

import { Script } from "forge-std/Script.sol";
import { console2 } from "forge-std/console2.sol";

import { DKG } from "../../src/protocol/DKG.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";

/**
 * @title DeployNewDKG
 * @notice Deploys a new implementation of DKG contract to be used for upgrading the predeploy at
 *         0xCcCcCC0000000000000000000000000000000004
 * @dev This script only deploys the implementation contract, it does not perform the upgrade.
 *
 *      After running this script, the following steps must be completed through governance (timelock):
 *        1. Upgrade DKG proxy — call ProxyAdmin.upgradeAndCall(DKG_PROXY, newImpl, "")
 *        2. Deploy SGXValidationHook — run DeploySGXValidationHook.s.sol
 *        3. Whitelist enclave type — call DKG.whitelistEnclaveType(enclaveType, enclaveTypeData, true)
 *           (see DeploySGXValidationHook.s.sol output for the required arguments)
 */
contract DeployNewDKG is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        Create3 create3 = Create3(Predeploys.Create3);

        // DKG has no constructor args (only _disableInitializers in constructor)
        bytes memory creationCode = type(DKG).creationCode;

        bytes32 salt = keccak256(abi.encodePacked("DKG_Implementation_v1_0_0"));

        // Deploy using Create3 for deterministic address
        address newImplementation = create3.deploy(salt, creationCode);
        if (create3.getDeployed(deployer, salt) != newImplementation) {
            revert("Deployment failed");
        }

        vm.stopBroadcast();

        console2.log("New DKG implementation deployed at:", newImplementation);
        console2.log("DKG proxy address:", Predeploys.DKG);
        console2.log("---");
        console2.log("NOTE: After deployment, complete the following through governance (timelock):");
        console2.log("  1. Upgrade DKG proxy to the new implementation");
        console2.log("  2. Deploy SGXValidationHook (run DeploySGXValidationHook.s.sol)");
        console2.log("  3. Whitelist enclave type on DKG using the SGXValidationHook proxy address");
    }
}
