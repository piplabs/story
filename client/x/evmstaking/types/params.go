package types

import (
	"fmt"
)

// Staking params default values.
const (
	DefaultMaxWithdrawalPerBlock uint32 = 32

	DefaultMaxSweepPerBlock uint32 = 128

	DefaultMinPartialWithdrawalAmount uint64 = 8_000_000_000

	DefaultRefundFeeBps uint32 = 100 // 100bps, or 1% (= 10.24 IP min fee)
)

// NewParams creates a new Params instance.
func NewParams(maxWithdrawalPerBlock uint32, maxSweepPerBlock uint32, minPartialWithdrawalAmount uint64, refundFeeBps uint32) Params {
	return Params{
		MaxWithdrawalPerBlock:      maxWithdrawalPerBlock,
		MaxSweepPerBlock:           maxSweepPerBlock,
		MinPartialWithdrawalAmount: minPartialWithdrawalAmount,
		RefundFeeBps:               refundFeeBps,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxWithdrawalPerBlock,
		DefaultMaxSweepPerBlock,
		DefaultMinPartialWithdrawalAmount,
		DefaultRefundFeeBps,
	)
}

func (p Params) Validate() error {
	if err := ValidateMaxWithdrawalPerBlock(p.MaxWithdrawalPerBlock); err != nil {
		return err
	}

	if err := ValidateMaxSweepPerBlock(p.MaxSweepPerBlock, p.MaxWithdrawalPerBlock); err != nil {
		return err
	}

	if err := ValidateMinPartialWithdrawalAmount(p.MinPartialWithdrawalAmount); err != nil {
		return err
	}

	return ValidateRefundFeeBps(p.RefundFeeBps)
}

func ValidateMaxWithdrawalPerBlock(v uint32) error {
	if v == 0 {
		return fmt.Errorf("max withdrawal per block must be positive: %d", v)
	}

	return nil
}

func ValidateMaxSweepPerBlock(maxSweepPerBlock uint32, maxWithdrawalPerBlock uint32) error {
	if maxSweepPerBlock == 0 {
		return fmt.Errorf("max sweep per block must be positive: %d", maxSweepPerBlock)
	}

	if maxSweepPerBlock < maxWithdrawalPerBlock {
		return fmt.Errorf("max sweep per block must be greater than or equal to max withdrawal per block: %d < %d", maxSweepPerBlock, maxWithdrawalPerBlock)
	}

	return nil
}

func ValidateMinPartialWithdrawalAmount(v uint64) error {
	if v == 0 {
		return fmt.Errorf("min partial withdrawal amount must be positive: %d", v)
	}

	return nil
}

func ValidateRefundFeeBps(v uint32) error {
	if v > 10000 {
		return fmt.Errorf("refund fee bps must be less than or equal to 10000bps (100%%): %d", v)
	}

	return nil
}
