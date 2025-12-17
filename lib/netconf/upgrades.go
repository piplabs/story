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
)

var (
	ErrUnknownChainID = errors.New("unknown chain ID")
	ErrUnknownUpgrade = errors.New("unknown upgrade name")
)

type UpgradeMap map[string]int64

// UpgradeHistories are the map of histories for each network.
var UpgradeHistories = map[string]UpgradeMap{
	DKGTestChainID: {
		V121:    0,
		Terence: 0,
		V142:    0,
	},
	TestChainID: {
		V121:    10,
		Terence: 50,
		V142:    50,
	},
	LocalChainID: {
		V121:    0,
		Terence: 50,
		V142:    50,
	},
	StoryLocalnetID: {
		V121:    0,
		Terence: 0,
		V142:    0,
	},
	AeneidChainID: {
		Virgil:   345158,
		Ovid:     4362990,
		V121:     5238000,
		Polybius: 6008000,
		Terence:  10886688,
		V142:     12088950,
	},
	StoryChainID: {
		Virgil:   809988,
		Ovid:     4477880,
		V121:     5084300,
		Polybius: 8270000,
		Terence:  11538000,
		V142:     11784600,
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
