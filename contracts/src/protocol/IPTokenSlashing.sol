// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import { IIPTokenSlashing } from "../interfaces/IIPTokenSlashing.sol";
import { IPTokenStaking } from "./IPTokenStaking.sol";
import { Secp256k1 } from "../libraries/Secp256k1.sol";

/**
 * @title IPTokenSlashing
 * @notice The EVM interface to the consensus chain's x/slashing module. Calls are proxied to the consensus chain, but
 *         not executed synchronously; execution is left to the consensus chain, which may fail.
 */
contract IPTokenSlashing is IIPTokenSlashing, Ownable2StepUpgradeable, UUPSUpgradeable {
    /// @notice IPTokenStaking contract address.
    IPTokenStaking public immutable IP_TOKEN_STAKING;

    /// @notice The fee paid to unjail a validator.
    uint256 public unjailFee;

    constructor(address ipTokenStaking) {
        require(ipTokenStaking != address(0), "IPTokenSlashing: Invalid IPTokenStaking address");
        IP_TOKEN_STAKING = IPTokenStaking(ipTokenStaking);
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(address accessManager, uint256 newUnjailFee) public initializer {
        __UUPSUpgradeable_init();
        __Ownable_init(accessManager);
        require(newUnjailFee > 0, "IPTokenSlashing: Invalid unjail fee");
        unjailFee = newUnjailFee;
        emit UnjailFeeSet(newUnjailFee);
    }

    /// @notice Verifies that the given 65 byte uncompressed secp256k1 public key (with 0x04 prefix) is valid and
    /// matches the expected EVM address.
    modifier verifyUncmpPubkeyWithExpectedAddress(bytes calldata uncmpPubkey, address expectedAddress) {
        require(uncmpPubkey.length == 65, "IPTokenSlashing: Invalid pubkey length");
        require(uncmpPubkey[0] == 0x04, "IPTokenSlashing: Invalid pubkey prefix");
        require(
            _uncmpPubkeyToAddress(uncmpPubkey) == expectedAddress,
            "IPTokenSlashing: Invalid pubkey derived address"
        );
        _;
    }

    /// @notice Sets the unjail fee.
    /// @param newUnjailFee The new unjail fee.
    function setUnjailFee(uint256 newUnjailFee) external onlyOwner {
        unjailFee = newUnjailFee;
        emit UnjailFeeSet(newUnjailFee);
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes calldata uncmpPubkey) internal pure returns (address) {
        return address(uint160(uint256(keccak256(uncmpPubkey[1:]))));
    }

    /// @notice Requests to unjail the validator. Must pay fee on the execution side to prevent spamming.
    function unjail(
        bytes calldata validatorUncmpPubkey
    ) external payable verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) {
        bytes memory validatorCmpPubkey = Secp256k1.compressPublicKey(validatorUncmpPubkey);
        _verifyExistingValidator(validatorCmpPubkey);
        _unjail(validatorCmpPubkey);
    }

    /// @notice Requests to unjail a validator on behalf. Must pay fee on the execution side to prevent spamming.
    /// @param validatorCmpPubkey The validator's 33-byte compressed Secp256k1 public key
    function unjailOnBehalf(bytes calldata validatorCmpPubkey) external payable {
        _verifyExistingValidator(validatorCmpPubkey);
        _unjail(validatorCmpPubkey);
    }

    /// @dev Emits the Unjail event after burning the fee.
    function _unjail(bytes memory validatorCmpPubkey) internal {
        require(msg.value == unjailFee, "IPTokenSlashing: Insufficient fee");
        payable(address(0x0)).transfer(msg.value);
        emit Unjail(msg.sender, validatorCmpPubkey);
    }

    /// @notice Verifies that the validator with the given pubkey exists.
    /// @param validatorCmpPubkey The validator's 33-byte compressed Secp256k1 public key
    function _verifyExistingValidator(bytes memory validatorCmpPubkey) internal view {
        require(validatorCmpPubkey.length == 33, "IPTokenSlashing: Invalid pubkey length");
        require(
            validatorCmpPubkey[0] == 0x02 || validatorCmpPubkey[0] == 0x03,
            "IPTokenSlashing: Invalid pubkey prefix"
        );
        (bool validatorExists, , , , , ) = IP_TOKEN_STAKING.validatorMetadata(validatorCmpPubkey);
        require(validatorExists, "IPTokenSlashing: Validator does not exist");
    }

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
