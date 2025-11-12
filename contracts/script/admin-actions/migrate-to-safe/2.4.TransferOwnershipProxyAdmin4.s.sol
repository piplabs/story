/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";

/// @title TransferOwnershipsProxyAdmin4
/// @notice Generates json files with the timelocked operations to transfer the ownership of proxy admins
/// from index 0 to 256
contract TransferOwnershipsProxyAdmin4 is BaseTransferOwnershipProxyAdmin {
    constructor()
        BaseTransferOwnershipProxyAdmin(
            "2.4-safe-migr-transfer-ownerships-proxy-admin-4",
            1, // fromIndex
            256 // toIndex
        )
    {}
}
