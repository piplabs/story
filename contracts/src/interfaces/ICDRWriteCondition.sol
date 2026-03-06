// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface ICDRWriteCondition {
    /// @notice Checks the write condition
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for writing
    /// @param writeConditionData The data of the write condition
    /// @param caller The address of the caller
    /// @return True if the write condition is met, false otherwise
    function checkWriteCondition(
        uint32 uuid,
        bytes calldata accessAuxData,
        bytes calldata writeConditionData,
        address caller
    ) external view returns (bool);
}
