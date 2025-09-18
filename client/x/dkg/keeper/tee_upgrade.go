package keeper

import (
	"context"
	"slices"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// SetTEEUpgradeInfo stores a TEE upgrade info in the store using mrenclave as the key.
func (k *Keeper) SetTEEUpgradeInfo(ctx context.Context, teeUpgradeInfo *types.TEEUpgradeInfo) error {
	key := string(teeUpgradeInfo.Mrenclave)
	if err := k.TEEUpgradeInfos.Set(ctx, key, *teeUpgradeInfo); err != nil {
		return errors.Wrap(err, "failed to set tee upgrade info")
	}

	return nil
}

// SetTEEUpgradeInfos stores multiple TEE upgrade infos (used for genesis initialization).
func (k *Keeper) SetTEEUpgradeInfos(ctx context.Context, teeUpgradeInfos []types.TEEUpgradeInfo) error {
	for _, teeInfo := range teeUpgradeInfos {
		if err := k.SetTEEUpgradeInfo(ctx, &teeInfo); err != nil {
			return errors.Wrap(err, "failed to set tee upgrade info during genesis initialization")
		}
	}

	return nil
}

// GetTEEUpgradeInfo retrieves a TEE upgrade info by unique ID.
func (k *Keeper) GetTEEUpgradeInfo(ctx context.Context, mrenclave []byte) (*types.TEEUpgradeInfo, error) {
	key := string(mrenclave)
	teeInfo, err := k.TEEUpgradeInfos.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "tee upgrade info not found")
		}

		return nil, errors.Wrap(err, "failed to get tee upgrade info")
	}

	return &teeInfo, nil
}

// GetAllTEEUpgradeInfos retrieves all TEE upgrade infos.
func (k *Keeper) GetAllTEEUpgradeInfos(ctx context.Context) ([]types.TEEUpgradeInfo, error) {
	var teeInfos []types.TEEUpgradeInfo

	err := k.TEEUpgradeInfos.Walk(ctx, nil, func(_ string, info types.TEEUpgradeInfo) (bool, error) {
		teeInfos = append(teeInfos, info)
		return false, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate tee upgrade infos")
	}

	return teeInfos, nil
}

// GetTEEUpgradeInfoByMrenclave retrieves a TEE upgrade info with a specific mrenclave.
func (k *Keeper) GetTEEUpgradeInfoByMrenclave(ctx context.Context, mrenclave []byte) (types.TEEUpgradeInfo, error) {
	var upgradeInfo types.TEEUpgradeInfo

	err := k.TEEUpgradeInfos.Walk(ctx, nil, func(_ string, info types.TEEUpgradeInfo) (bool, error) {
		if slices.Equal(info.Mrenclave, mrenclave) {
			upgradeInfo = info
			return true, nil // stop iteration
		}

		return false, nil
	})

	if err != nil {
		return upgradeInfo, errors.Wrap(err, "failed to iterate tee upgrade infos by mrenclave")
	}

	return upgradeInfo, nil
}

// DeleteTEEUpgradeInfo removes a TEE upgrade info from the store.
func (k *Keeper) DeleteTEEUpgradeInfo(ctx context.Context, mrenclave []byte) error {
	key := string(mrenclave)
	if err := k.TEEUpgradeInfos.Remove(ctx, key); err != nil {
		return errors.Wrap(err, "failed to delete tee upgrade info")
	}

	return nil
}
