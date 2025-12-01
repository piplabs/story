/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";

/// @title TransferOwnershipsProxyAdmin2
/// @notice Generates json files with the timelocked operations to transfer the ownership of proxy admins
/// from index 512 to 768
contract TransferOwnershipsProxyAdmin2 is BaseTransferOwnershipProxyAdmin {
    constructor()
        BaseTransferOwnershipProxyAdmin(
            "2.2-safe-migr-transfer-ownerships-proxy-admin-2",
            513, // fromIndex
            768 // toIndex
        )
    {}
}
