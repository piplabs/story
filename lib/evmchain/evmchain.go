// Package evmchain provides static metadata about supported evm chains.
package evmchain

import (
	"time"
)

type Token string

const (
	// Mainnets.
	IDStoryMainnet uint64 = 1514

	// Local Testet.
	IDLocal uint64 = 1511

	// Testnets.
	IDStoryTestnet uint64 = 1511
	IDIliad        uint64 = 1513

	iliadEVMName        = "story_evm"
	iliadEVMBlockPeriod = time.Second * 2

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
	IDStoryTestnet: {
		ChainID:     IDStoryTestnet,
		Name:        iliadEVMName,
		BlockPeriod: iliadEVMBlockPeriod,
		NativeToken: IP,
	},
	IDIliad: {
		ChainID:     IDIliad,
		Name:        iliadEVMName,
		BlockPeriod: iliadEVMBlockPeriod,
		NativeToken: IP,
	},
}
