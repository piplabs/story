/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";

/// @title TransferOwnershipsProxyAdmin3
/// @notice Generates json files with the timelocked operations to transfer the ownership of proxy admins
/// from index 256 to 512
contract TransferOwnershipsProxyAdmin3 is BaseTransferOwnershipProxyAdmin {
    constructor()
        BaseTransferOwnershipProxyAdmin(
            "2.3-safe-migr-transfer-ownerships-proxy-admin-3",
            257, // fromIndex
            512 // toIndex
        )
    {}
}
