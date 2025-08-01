/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";

/// @title TransferOwnershipsProxyAdmin1
/// @notice Generates json files with the timelocked operations to transfer the ownership of the
/// last 256 proxy admins to the new timelock
contract TransferOwnershipsProxyAdmin1 is BaseTransferOwnershipProxyAdmin {
    constructor()
        BaseTransferOwnershipProxyAdmin(
            "2.1-safe-migr-transfer-ownerships-proxy-admin-1",
            769, // fromIndex
            1024 // toIndex
        )
    {}
}
