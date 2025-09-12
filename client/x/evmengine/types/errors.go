package types

import "cosmossdk.io/errors"

var (
	ErrUpgradePending  = errors.Register(ModuleName, 2, "upgrade is already pending")
	ErrUpgradeNotFound = errors.Register(ModuleName, 3, "upgrade not found")
)
