// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.23;

import { CREATE3 } from "solady/src/utils/CREATE3.sol";

/**
 * @title Create3
 * @notice Modified to use Solady. Original: (https://github.com/ZeframLou/create3-factory/blob/main/src/CREATE3Factory.sol)
 */
contract Create3 {
    /**
     * @notice Deploys a contract using CREATE3
     * @dev The provided salt is hashed together with msg.sender to generate the final salt
     * @param salt          The deployer-specific salt for determining the deployed contract's address
     * @param creationCode  The creation code of the contract to deploy
     * @return deployed     The address of the deployed contract
     */
    function deploy(bytes32 salt, bytes memory creationCode) external payable returns (address deployed) {
        return CREATE3.deployDeterministic(0, creationCode, salt);
    }

    /**
     * @notice Predicts the address of a deployed contract
     * @dev The provided salt is hashed together with the deployer address to generate the final salt
     * @param salt      The deployer-specific salt for determining the deployed contract's address
     * @return deployed The address of the contract that will be deployed
     */
    function getDeployed(bytes32 salt) external view returns (address deployed) {
        return CREATE3.predictDeterministicAddress(salt);
    }
}
