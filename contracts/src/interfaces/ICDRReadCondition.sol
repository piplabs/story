// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface ICDRReadCondition {
    /// @notice Checks the read condition
    /// @param uuid The UUID of the vault
    /// @param accessAuxData The auxiliary access data for reading
    /// @param readConditionData The data of the read condition
    /// @param caller The address of the caller
    /// @return True if the read condition is met, false otherwise
    function checkReadCondition(
        uint32 uuid,
        bytes calldata accessAuxData,
        bytes calldata readConditionData,
        address caller
    ) external view returns (bool);
}
