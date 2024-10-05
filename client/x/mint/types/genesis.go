package types

import (
	"context"

	"cosmossdk.io/math"
)

// InflationCalculationFn defines the function required to calculate inflation amount during
// BeginBlock. It receives the params stored in the keeper, along with the current
// bondedRatio and returns the newly calculated inflation amount.
// It can be used to specify a custom inflation calculation logic, instead of relying on the
// default logic provided by the sdk.
type InflationCalculationFn func(ctx context.Context, params Params, bondedRatio math.LegacyDec) math.LegacyDec

// DefaultInflationCalculationFn is the default function used to calculate inflation.
func DefaultInflationCalculationFn(_ context.Context, params Params, _ math.LegacyDec) math.LegacyDec {
	return math.LegacyNewDec(int64(params.InflationsPerYear)).Quo(math.LegacyNewDec(int64(params.BlocksPerYear)))
}

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
