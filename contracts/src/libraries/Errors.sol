// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

/**
 * @title Errors
 * @notice Errors for the staking and parameter contracts
 */
library Errors {
    error IPTokenStaking__ZeroStakingRounding();
    error IPTokenStaking__InvalidDefaultMinUnjailFee();
    error IPTokenStaking__CommissionRateUnderMin();
    error IPTokenStaking__CommissionRateOverMax();

    error IPTokenStaking__InvalidPubkeyLength();
    error IPTokenStaking__InvalidPubkeyPrefix();
    error IPTokenStaking__InvalidPubkeyDerivedAddress();

    error IPTokenStaking__InvalidMinUnjailFee();
    error IPTokenStaking__ZeroMinStakeAmount();
    error IPTokenStaking__ZeroMinUnstakeAmount();
    error IPTokenStaking__ZeroMinCommissionRate();

    error IPTokenStaking__ZeroShortPeriodDuration();
    error IPTokenStaking__ShortPeriodLongerThanMedium();
    error IPTokenStaking__MediumLongerThanLong();

    error IPTokenStaking__StakeAmountUnderMin();
    error IPTokenStaking__LowUnstakeAmount();
    error IPTokenStaking__RedelegatingToSameValidator();

    error IPTokenStaking__FailedRemainerRefund();
    error IPTokenStaking__InvalidFeeAmount();
}
