package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// SetKernelUpgradeInfo stores a kernel upgrade info in the store using upgradeVersion as the key.
func (k *Keeper) SetKernelUpgradeInfo(ctx context.Context, upgradeInfo *types.KernelUpgradeInfo) error {
	key := upgradeInfo.UpgradeVersion
	if err := k.KernelUpgradeInfos.Set(ctx, key, *upgradeInfo); err != nil {
		return errors.Wrap(err, "failed to set kernel upgrade info")
	}

	return nil
}

// SetKernelUpgradeInfos stores multiple kernel upgrade infos (used for genesis initialization).
func (k *Keeper) SetKernelUpgradeInfos(ctx context.Context, upgradeInfos []types.KernelUpgradeInfo) error {
	for _, info := range upgradeInfos {
		if err := k.SetKernelUpgradeInfo(ctx, &info); err != nil {
			return errors.Wrap(err, "failed to set kernel upgrade info during genesis initialization")
		}
	}

	return nil
}

// GetKernelUpgradeInfo retrieves a kernel upgrade info by upgradeVersion.
func (k *Keeper) GetKernelUpgradeInfo(ctx context.Context, upgradeVersion string) (*types.KernelUpgradeInfo, error) {
	info, err := k.KernelUpgradeInfos.Get(ctx, upgradeVersion)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "kernel upgrade info not found")
		}

		return nil, errors.Wrap(err, "failed to get kernel upgrade info")
	}

	return &info, nil
}

// GetAllKernelUpgradeInfos retrieves all kernel upgrade infos.
func (k *Keeper) GetAllKernelUpgradeInfos(ctx context.Context) ([]types.KernelUpgradeInfo, error) {
	var infos []types.KernelUpgradeInfo

	err := k.KernelUpgradeInfos.Walk(ctx, nil, func(_ string, info types.KernelUpgradeInfo) (bool, error) {
		infos = append(infos, info)
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate kernel upgrade infos")
	}

	return infos, nil
}

// GetPendingUpgrade finds a pending kernel upgrade info.
// Activated upgrades are marked with IsActivated=true and filtered out,
// so any remaining non-activated entry is a pending upgrade.
func (k *Keeper) GetPendingUpgrade(ctx context.Context) (*types.KernelUpgradeInfo, error) {
	var pending *types.KernelUpgradeInfo

	err := k.KernelUpgradeInfos.Walk(ctx, nil, func(_ string, info types.KernelUpgradeInfo) (bool, error) {
		if info.IsActivated {
			return false, nil // skip activated entries, continue walking
		}

		pending = &info
		return true, nil // found pending, stop
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate kernel upgrade infos for pending upgrade")
	}

	return pending, nil
}

// DeleteKernelUpgradeInfo removes a kernel upgrade info from the store.
func (k *Keeper) DeleteKernelUpgradeInfo(ctx context.Context, upgradeVersion string) error {
	if err := k.KernelUpgradeInfos.Remove(ctx, upgradeVersion); err != nil {
		return errors.Wrap(err, "failed to delete kernel upgrade info")
	}

	return nil
}

// deleteActivatedUpgradeInfo removes all upgrade infos that have been activated.
// Called after a successful upgrade resharing round to clean up completed upgrades.
func (k *Keeper) deleteActivatedUpgradeInfo(ctx context.Context) error {
	var activatedVersions []string

	err := k.KernelUpgradeInfos.Walk(ctx, nil, func(_ string, info types.KernelUpgradeInfo) (bool, error) {
		if info.IsActivated {
			activatedVersions = append(activatedVersions, info.UpgradeVersion)
		}

		return false, nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to iterate kernel upgrade infos for cleanup")
	}

	for _, version := range activatedVersions {
		if err := k.DeleteKernelUpgradeInfo(ctx, version); err != nil {
			return errors.Wrap(err, "failed to delete activated kernel upgrade info")
		}
	}

	return nil
}
