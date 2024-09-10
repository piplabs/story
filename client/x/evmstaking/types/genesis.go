package types

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
		ValidatorSweepIndex: &ValidatorSweepIndex{
			NextValIndex:    0,
			NextValDelIndex: 0,
		},
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		ValidatorSweepIndex: &ValidatorSweepIndex{
			NextValIndex:    0,
			NextValDelIndex: 0,
		},
	}
}
