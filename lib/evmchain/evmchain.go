// Package evmchain provides static metadata about supported evm chains.
package evmchain

import (
	"time"
)

type Token string

const (
	// Mainnets.
	IDEthereum     uint64 = 1
	IDIliadMainnet uint64 = 1514

	// Local Testet.
	IDLocal uint64 = 1511

	// Testnets.
	IDIliadTestnet uint64 = 1511
	IDIliad        uint64 = 1513
	IDHolesky      uint64 = 17000
	IDArbSepolia   uint64 = 421614
	IDOpSepolia    uint64 = 11155420

	// Ephemeral.
	IDIliadEphemeral uint64 = 1651
	IDMockL1Fast     uint64 = 1652
	IDMockL1Slow     uint64 = 1653
	IDMockL2         uint64 = 1654
	IDMockOp         uint64 = 1655
	IDMockArb        uint64 = 1656

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
	IDEthereum: {
		ChainID:     IDEthereum,
		Name:        "ethereum",
		BlockPeriod: 12 * time.Second,
		NativeToken: ETH,
	},
	IDIliadMainnet: {
		ChainID:     IDIliadMainnet,
		Name:        iliadEVMName,
		BlockPeriod: iliadEVMBlockPeriod,
		NativeToken: IP,
	},
	IDIliadTestnet: {
		ChainID:     IDIliadTestnet,
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
	IDHolesky: {
		ChainID:     IDHolesky,
		Name:        "holesky",
		BlockPeriod: 12 * time.Second,
		NativeToken: ETH,
	},
	IDArbSepolia: {
		ChainID:     IDArbSepolia,
		Name:        "arb_sepolia",
		BlockPeriod: 300 * time.Millisecond,
		NativeToken: ETH,
	},
	IDOpSepolia: {
		ChainID:     IDOpSepolia,
		Name:        "op_sepolia",
		BlockPeriod: 2 * time.Second,
		NativeToken: ETH,
	},
	IDIliadEphemeral: {
		ChainID:     IDIliadEphemeral,
		Name:        iliadEVMName,
		BlockPeriod: iliadEVMBlockPeriod,
		NativeToken: IP,
	},
	IDMockL1Fast: {
		ChainID:     IDMockL1Fast,
		Name:        "mock_l1",
		BlockPeriod: time.Second,
		NativeToken: ETH,
	},
	IDMockL1Slow: {
		ChainID:     IDMockL1Slow,
		Name:        "slow_l1",
		BlockPeriod: time.Second * 12,
		NativeToken: ETH,
	},
	IDMockL2: {
		ChainID:     IDMockL2,
		Name:        "mock_l2",
		BlockPeriod: time.Second,
		NativeToken: ETH,
	},
	IDMockOp: {
		ChainID:     IDMockOp,
		Name:        "mock_op",
		BlockPeriod: time.Second * 2,
		NativeToken: ETH,
	},
	IDMockArb: {
		ChainID:     IDMockArb,
		Name:        "mock_arb",
		BlockPeriod: time.Second / 4,
		NativeToken: ETH,
	},
}
