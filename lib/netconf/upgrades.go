package netconf

import (
	"errors"
)

const (
	Virgil   = "virgil"
	Ovid     = "v1.2.0"
	Polybius = "polybius"

	V121 = "v1.2.1"

	Terence = "terence"
	V142    = "v1.4.2"

	Horace = "horace"
	V160   = "v1.6.0"

	Seneca = "seneca"

	// V190 gates the DKG consensus-polynomial validation sweep in FinalizeDKGRound.
	V190 = "v1.9.0"
)

var (
	ErrUnknownChainID = errors.New("unknown chain ID")
	ErrUnknownUpgrade = errors.New("unknown upgrade name")
)

type UpgradeMap map[string]int64

// UpgradeHistories are the map of histories for each network.
var UpgradeHistories = map[string]UpgradeMap{
	TestChainID: {
		V121:    10,
		Terence: 50,
		V142:    50,
		Horace:  100,
		Seneca:  300,
		V190:    400,
	},
	LocalChainID: {
		V121:    0,
		Terence: 50,
		V142:    50,
		Horace:  100,
		V190:    100,
	},
	StoryLocalnetID: {
		V121:    0,
		Terence: 0,
		V142:    0,
		Horace:  100,
		V190:    0,
	},
	AeneidChainID: {
		Virgil:   345158,
		Ovid:     4362990,
		V121:     5238000,
		Polybius: 6008000,
		Terence:  10886688,
		V142:     12088950,
		Horace:   14017000,
		Seneca:   18550000,
	},
	StoryChainID: {
		Virgil:   809988,
		Ovid:     4477880,
		V121:     5084300,
		Polybius: 8270000,
		Terence:  11538000,
		V142:     11784600,
		Horace:   13780500,
	},
}

func (um UpgradeMap) GetUpgradeBlock(upgradeName string) (int64, error) {
	upgradeBlock, ok := um[upgradeName]
	if !ok {
		return 0, ErrUnknownUpgrade
	}

	return upgradeBlock, nil
}

func GetUpgradeHistory(chainID string) (UpgradeMap, error) {
	upgradeHistory, ok := UpgradeHistories[chainID]
	if !ok {
		return nil, ErrUnknownChainID
	}

	return upgradeHistory, nil
}

func GetUpgradeHeight(chainID, upgradeName string) (int64, error) {
	upgradeMap, err := GetUpgradeHistory(chainID)
	if err != nil {
		return 0, err
	}

	upgradeBlock, err := upgradeMap.GetUpgradeBlock(upgradeName)
	if err != nil {
		return 0, err
	}

	return upgradeBlock, nil
}

func IsV121(chainID string, blockNumber int64) (bool, error) {
	v121Block, err := GetUpgradeHeight(chainID, V121)
	if err != nil {
		return false, err
	}

	return blockNumber >= v121Block, nil
}

func IsTerence(chainID string, blockNumber int64) (bool, error) {
	terenceBlock, err := GetUpgradeHeight(chainID, Terence)
	if err != nil {
		return false, err
	}

	return blockNumber >= terenceBlock, nil
}

func IsV142(chainID string, blockNumber int64) (bool, error) {
	v142Block, err := GetUpgradeHeight(chainID, V142)
	if err != nil {
		return false, err
	}

	return blockNumber >= v142Block, nil
}

func IsSeneca(chainID string, blockNumber int64) (bool, error) {
	senecaBlock, err := GetUpgradeHeight(chainID, Seneca)
	if err != nil {
		return false, err
	}

	return blockNumber >= senecaBlock, nil
}

// IsV190 reports whether the v1.9.0 upgrade is active at blockNumber on chainID.
// Returns an error when V190 is not registered for the chain; DKG callers treat that
// as "not active" rather than halting (DKG can run on chains not in UpgradeHistories).
func IsV190(chainID string, blockNumber int64) (bool, error) {
	v190Block, err := GetUpgradeHeight(chainID, V190)
	if err != nil {
		return false, err
	}

	return blockNumber >= v190Block, nil
}
