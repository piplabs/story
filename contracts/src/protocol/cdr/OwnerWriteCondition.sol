// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/// @title OwnerWriteCondition
/// @notice Simplest CDR write condition — only the address encoded in writeConditionData can write.
///         No register() step needed. The authorized writer is set at allocate time via
///         writeConditionData = abi.encode(ownerAddress).
///
/// Usage:
///   CDR.allocate(updatable, ownerWriteConditionAddr, readConditionAddr,
///                abi.encode(yourAddress), readConditionData)
contract OwnerWriteCondition {
    function checkWriteCondition(
        uint32,
        bytes calldata,
        bytes calldata writeConditionData,
        address caller
    ) external pure returns (bool) {
        address owner = abi.decode(writeConditionData, (address));
        return caller == owner;
    }
}
