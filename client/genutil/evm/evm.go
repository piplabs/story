package evm

import (
	"math/big"

	"cosmossdk.io/math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/netconf"
)

//nolint:unused // added for potential future usage
var eth1k = math.NewInt(1000).MulRaw(params.Ether).BigInt()
var eth1m = math.NewInt(1000000).MulRaw(params.Ether).BigInt()

func newUint64(val uint64) *uint64 { return &val }

// MakeGenesis returns a genesis block for a development chain.
// See geth reference: https://github.com/ethereum/go-ethereum/blob/master/core/genesis.go#L564
func MakeGenesis(network netconf.ID) (core.Genesis, error) {
	predeps, err := predeploys.Alloc(network)
	if err != nil {
		return core.Genesis{}, errors.Wrap(err, "predeploys")
	}

	allocs := mergeAllocs(precompilesAlloc(), predeps)

	if network.IsEphemeral() {
		allocs = mergeAllocs(allocs, stagingPrefundAlloc())
	} else if network == netconf.Testnet {
		allocs = mergeAllocs(allocs, testnetPrefundAlloc())
	} else {
		return core.Genesis{}, errors.New("unsupported network", "network", network.String())
	}

	return core.Genesis{
		Config:     defaultChainConfig(network),
		GasLimit:   miner.DefaultConfig.GasCeil,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Difficulty: big.NewInt(0),
		Alloc:      allocs,
	}, nil
}

// defaultChainConfig returns the default chain config for a network.
// See geth reference: https://github.com/ethereum/go-ethereum/blob/master/params/config.go#L65
func defaultChainConfig(network netconf.ID) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                       big.NewInt(int64(network.Static().StoryExecutionChainID)),
		HomesteadBlock:                big.NewInt(0),
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		ShanghaiTime:                  newUint64(0),
		CancunTime:                    newUint64(0),
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
	}
}

// precompilesAlloc returns allocs for precompiled contracts
// TODO: this matches go-ethereum's precompiles, but we should understand why balances are set to 1.
func precompilesAlloc() types.GenesisAlloc {
	return types.GenesisAlloc{
		common.BytesToAddress([]byte{1}): {Balance: big.NewInt(1)}, // ECRecover
		common.BytesToAddress([]byte{2}): {Balance: big.NewInt(1)}, // SHA256
		common.BytesToAddress([]byte{3}): {Balance: big.NewInt(1)}, // RIPEMD
		common.BytesToAddress([]byte{4}): {Balance: big.NewInt(1)}, // Identity
		common.BytesToAddress([]byte{5}): {Balance: big.NewInt(1)}, // ModExp
		common.BytesToAddress([]byte{6}): {Balance: big.NewInt(1)}, // ECAdd
		common.BytesToAddress([]byte{7}): {Balance: big.NewInt(1)}, // ECScalarMul
		common.BytesToAddress([]byte{8}): {Balance: big.NewInt(1)}, // ECPairing
		common.BytesToAddress([]byte{9}): {Balance: big.NewInt(1)}, // BLAKE2b
	}
}

// devPrefundAlloc returns allocs for pre-funded geth dev accounts.
func stagingPrefundAlloc() types.GenesisAlloc {
	return types.GenesisAlloc{
		// TODO: team accounts
		common.HexToAddress("0x0000000000000000000000000000000000000000"): {Balance: eth1m},
	}
}

func testnetPrefundAlloc() types.GenesisAlloc {
	return types.GenesisAlloc{
		// TODO: team accounts
		common.HexToAddress("0x0000000000000000000000000000000000000000"): {Balance: eth1m},

		// TODO: add validators
	}
}

// mergeAllocs merges multiple allocs into one.
func mergeAllocs(allocs ...types.GenesisAlloc) types.GenesisAlloc {
	merged := make(types.GenesisAlloc)
	for _, alloc := range allocs {
		for addr, account := range alloc {
			merged[addr] = account
		}
	}

	return merged
}
