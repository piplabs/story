package types

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}
