// SPDX-License-Identifier: AGPL-3.0
pragma solidity 0.8.23;

import { CREATE3 } from "solmate/src/utils/CREATE3.sol";

/**
 * @title Create3
 * @notice Factory for deploying contracts to deterministic addresses via CREATE3 Enables deploying
 *         contracts using CREATE3.
 *         It is based on Zefram’s implementation and includes support for deployment using only a salt value.
 *         It support two deployment styles:
 *         1. Each deployer (msg.sender) has its own namespace for deployed addresses.
 *         2. Deterministic deployment using a single salt.
 * @custom:attribution zefram.eth (https://github.com/ZeframLou/create3-factory/blob/main/src/CREATE3Factory.sol)
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
        // hash salt with the deployer address to give each deployer its own namespace
        salt = keccak256(abi.encodePacked(msg.sender, salt));
        return CREATE3.deploy(salt, creationCode, msg.value);
    }

    /**
     * @notice Predicts the address of a deployed contract
     * @dev The provided salt is hashed together with the deployer address to generate the final salt
     * @param deployer  The deployer account that will call deploy()
     * @param salt      The deployer-specific salt for determining the deployed contract's address
     * @return deployed The address of the contract that will be deployed
     */
    function getDeployed(address deployer, bytes32 salt) external view returns (address deployed) {
        // hash salt with the deployer address to give each deployer its own namespace
        salt = keccak256(abi.encodePacked(deployer, salt));
        return CREATE3.getDeployed(salt);
    }

    /**
     * @notice Deploys a contract using CREATE3
     * @param salt          The deployer-specific salt for determining the deployed contract's address
     * @param creationCode  The creation code of the contract to deploy
     * @return deployed     The address of the deployed contract
     */
    function deployDeterministic(bytes memory creationCode, bytes32 salt) external payable returns (address deployed) {
        return CREATE3.deploy(salt, creationCode, msg.value);
    }

    /**
     * @notice Predicts the address of a deployed contract
     * @param salt      The deployer-specific salt for determining the deployed contract's address
     * @return deployed The address of the contract that will be deployed
     */
    function predictDeterministicAddress(bytes32 salt) external view returns (address deployed) {
        // hash salt with the deployer address to give each deployer its own namespace
        return CREATE3.getDeployed(salt);
    }
}
