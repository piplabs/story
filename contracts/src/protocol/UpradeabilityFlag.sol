// SPDX-License-Identifier: GPL-3.0
// OpenZeppelin Contracts (last updated v5.0.0) (utils/ReentrancyGuard.sol)

pragma solidity ^0.8.23;
import { Initializable } from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title UpgradeabilityFlag
 * @notice The UpgradeabilityFlag contract is used to signal the upgradeability of a UUPSUpgradeable contract.
 * @dev To work, it must be:
 * - inherited by a UUPSUpgradeable contract.
 * - modifier `upgradeabilityEnabled` must be used in _authorizeUpgrade functions.
 * - `disableUpgradeability()` must be protected by an access control mechanism.
 * WARNING: Disabling upgradeability is an irreversible action.
 */
abstract contract UpgradeabilityFlag is Initializable {

    event UpgradeabilityDisabled();

    /// @dev Storage structure for the UpgradeabilityFlag
    /// @custom:storage-location erc7201:story.UpgradeabilityFlag
    struct UpgradeabilityFlagStorage {
        bool upgradeabilityDisabled;
    }

    // keccak256(abi.encode(uint256(keccak256("story.UpgradeabilityFlag")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant UpgradeabilityFlagStorageLocation =
        0x75e222a2214bb4b88a6942927541e7b0da442e311bbd530a7a9ccee34f9dee00;


    /// @dev moodifier to revert the transaction if upgradeability is disabled.
    /// place it in _authorizeUpgrade functions.
    modifier upgradeabilityEnabled() {
        require(!_getUpgradeabilityFlagStorage().upgradeabilityDisabled, "UpgradeabilityFlag: disabled");
        _;
    }

    /// @dev Implementation guide:
    /// - set onlyOwner or other access control on this function.
    /// - call _disableUpgradeability()
    function disableUpgradeability() external virtual;

    /// @notice Returns whether the upgradeability is disabled.
    function upgradeabilityDisabled() external view returns (bool) {
        return _getUpgradeabilityFlagStorage().upgradeabilityDisabled;
    }

    /// @notice Disables the upgradeability of the contract.
    /// WARNING: Disabling upgradeability is an irreversible action.
    function _disableUpgradeability() internal {
        _getUpgradeabilityFlagStorage().upgradeabilityDisabled = true;
        emit UpgradeabilityDisabled();
    }

    /// @dev Returns the storage struct of UpgradeabilityFlag.
    function _getUpgradeabilityFlagStorage() private pure returns (UpgradeabilityFlagStorage storage $) {
        assembly {
            $.slot := UpgradeabilityFlagStorageLocation
        }
    }
}
