package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
)

func (k *Keeper) InitGenesis(ctx context.Context, gs *types.GenesisState) error {
	if err := k.ValidateGenesis(gs); err != nil {
		return err
	}
	if err := k.SetParams(ctx, gs.Params); err != nil {
		return err
	}
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params: params,
	}
}

func (*Keeper) ValidateGenesis(gs *types.GenesisState) error {
	return gs.Params.Validate()
}
