//nolint:revive // just use interface{}
package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/lib/errors"
)

// NewParams returns Params instance with the given values.
func NewParams(mintDenom string, inflationsPerYear math.LegacyDec, blocksPerYear uint64) Params {
	return Params{
		MintDenom:         mintDenom,
		InflationsPerYear: inflationsPerYear,
		BlocksPerYear:     blocksPerYear,
	}
}

// DefaultParams returns default x/mint module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:         sdk.DefaultBondDenom,
		InflationsPerYear: math.LegacyNewDec(24625000000000000.000000000000000000),
		BlocksPerYear:     uint64(60 * 60 * 8766 / 5), // assuming 5 seconds block times
	}
}

// Validate does the sanity check on the params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationsPerYear(p.InflationsPerYear); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}

	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return errors.Wrap(err, "mint denom is invalid")
	}

	return nil
}

func validateInflationsPerYear(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("inflations per year cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("inflations per year cannot be negative: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}
