package keepers

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// GetSubspace gets existing substore from keeper.
func (keepers *Keepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := keepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}
