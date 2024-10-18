package types

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
		ValidatorSweepIndex: ValidatorSweepIndex{
			NextValIndex:    0,
			NextValDelIndex: 0,
		},
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		ValidatorSweepIndex: DefaultValidatorSweepIndex(),
	}
}

func NewValidatorSweepIndex(nextValIndex, nextValDelIndex uint64) ValidatorSweepIndex {
	return ValidatorSweepIndex{
		NextValIndex:    nextValIndex,
		NextValDelIndex: nextValDelIndex,
	}
}

func DefaultValidatorSweepIndex() ValidatorSweepIndex {
	return ValidatorSweepIndex{
		NextValIndex:    0,
		NextValDelIndex: 0,
	}
}
