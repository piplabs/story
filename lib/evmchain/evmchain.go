// Package evmchain provides static metadata about supported evm chains.
package evmchain

import (
	"time"
)

type Token string

const (
	// Story Mainnet.
	IDStory uint64 = 1514

	// Local Devnet.
	IDLocal uint64 = 1511

	// Iliad Testnets.
	IDIliad uint64 = 1513

	// Odyssey Testnet.
	IDOdyssey uint64 = 1516

	// Aeneid Testnet.
	IDAeneid uint64 = 1315

	storyEVMName        = "story_evm"
	storyEVMBlockPeriod = time.Second * 2

	IP  Token = "IP"
	ETH Token = "ETH"
)

type Metadata struct {
	ChainID     uint64
	Name        string
	BlockPeriod time.Duration
	NativeToken Token
}

func MetadataByID(chainID uint64) (Metadata, bool) {
	resp, ok := static[chainID]
	return resp, ok
}

func MetadataByName(name string) (Metadata, bool) {
	for _, metadata := range static {
		if metadata.Name == name {
			return metadata, true
		}
	}

	return Metadata{}, false
}

var static = map[uint64]Metadata{
	IDIliad: {
		ChainID:     IDIliad,
		Name:        storyEVMName,
		BlockPeriod: storyEVMBlockPeriod,
		NativeToken: IP,
	},
}
