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
	V181   = "v1.8.1"

	Seneca = "seneca"
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
	},
	LocalChainID: {
		V121:    0,
		Terence: 50,
		V142:    50,
		Horace:  100,
	},
	StoryLocalnetID: {
		V121:    0,
		Terence: 0,
		V142:    0,
		Horace:  100,
		// The e2e drives the fork from this height, so the harness and this
		// value have to agree.
		V181: 130,
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
		// TODO(release): set the real v1.8.1 activation height before cutting the
		// release, and set the same value on both release lines. The placeholder
		// is deliberately unreachable so a forgotten value activates nothing: a
		// line left on the placeholder does not halt, and its nodes would run
		// past the real height under the old rules and diverge.
		//
		// This height is what activates the upgrade -- scheduleForkUpgrade
		// schedules the plan when the chain reaches it. Do not also submit an
		// on-chain plan under this name: that would complete the name early, and
		// when the chain later reached this height ScheduleUpgrade would reject
		// the already-completed name and scheduleForkUpgrade would panic.
		V181: 999999999,
	},
	StoryChainID: {
		Virgil:   809988,
		Ovid:     4477880,
		V121:     5084300,
		Polybius: 8270000,
		Terence:  11538000,
		V142:     11784600,
		Horace:   13780500,
		Seneca:   18173000,
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
