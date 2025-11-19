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

	TestV142 = "test-v1.4.2"
	TestV143 = "test-v1.4.3"
)

var (
	ErrUnknownChainID = errors.New("unknown chain ID")
	ErrUnknownUpgrade = errors.New("unknown upgrade name")
)

type UpgradeMap map[string]int64

// UpgradeHistories are the map of histories for each network.
var UpgradeHistories = map[string]UpgradeMap{
	TestChainID: {
		V121:     10,
		Terence:  50,     // internal-devnet test
		TestV142: 111900, // internal-devnet v1.4.2
		TestV143: 112600, // TODO: set actual upgrade height for next test upgrade
	},
	LocalChainID: {
		V121:     0,
		Terence:  50,
		TestV142: 100, // TODO: set actual upgrade height for local chain
		TestV143: 0,
	},
	StoryLocalnetID: {
		V121:     0,
		Terence:  0,
		TestV142: 0,
		TestV143: 0,
	},
	AeneidChainID: {
		Virgil:   345158,
		Ovid:     4362990,
		V121:     5238000,
		Polybius: 6008000,
		Terence:  10886688,
		TestV142: 0, // Not applicable for Aeneid
		TestV143: 0,
	},
	StoryChainID: {
		Virgil:   809988,
		Ovid:     4477880,
		V121:     5084300,
		Polybius: 8270000,
		Terence:  100000000, // TODO: need to set actual upgrade height for story mainnet
		TestV142: 0,         // Not applicable for mainnet
		TestV143: 0,
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

func IsTestV142(chainID string, blockNumber int64) (bool, error) {
	testV142Block, err := GetUpgradeHeight(chainID, TestV142)
	if err != nil {
		return false, err
	}

	return blockNumber >= testV142Block, nil
}

func IsTestV143(chainID string, blockNumber int64) (bool, error) {
	testV143Block, err := GetUpgradeHeight(chainID, TestV143)
	if err != nil {
		return false, err
	}

	return blockNumber >= testV143Block, nil
}
