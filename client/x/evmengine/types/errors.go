package types

import "cosmossdk.io/errors"

var (
	ErrUpgradePending  = errors.Register(ModuleName, 1, "upgrade is already pending")
	ErrUpgradeNotFound = errors.Register(ModuleName, 2, "upgrade not found")
)
