// SPDX-License-Identifier: GPL-3.0
// OpenZeppelin Contracts (last updated v5.0.0) (utils/ReentrancyGuard.sol)

pragma solidity ^0.8.23;

import { IPTokenSlashing } from "../../src/protocol/IPTokenSlashing.sol";
import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";


contract MockIPTokenStakingV2 is IPTokenStaking {
    constructor(
        uint256 stakingRounding,
        uint32 defaultCommissionRate,
        uint32 defaultMaxCommissionRate,
        uint32 defaultMaxCommissionChangeRate
    ) IPTokenStaking(
        stakingRounding,
        defaultCommissionRate,
        defaultMaxCommissionRate,
        defaultMaxCommissionChangeRate
    ) {}

    function upgraded() external pure returns (bool) {
        return true;
    }
}

contract MockIPTokenSlashingV2 is IPTokenSlashing {
    constructor(address ipTokenStaking) IPTokenSlashing(ipTokenStaking) {}

    function upgraded() external pure returns (bool) {
        return true;
    }
}

contract MockUpgradeEntryPointV2 is UpgradeEntrypoint {
    function upgraded() external pure returns (bool) {
        return true;
    }
}
