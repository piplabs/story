package types

import (
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/evmchain"
	"github.com/piplabs/story/lib/netconf"
)

//nolint:gochecknoglobals // Static mappings
var (
	chainEthereum = EVMChain{
		Metadata: mustMetadata(evmchain.IDEthereum),
		IsPublic: true,
	}

	chainHolesky = EVMChain{
		Metadata: mustMetadata(evmchain.IDHolesky),
		IsPublic: true,
	}

	chainArbSepolia = EVMChain{
		Metadata: mustMetadata(evmchain.IDArbSepolia),
		IsPublic: true,
	}

	chainOpSepolia = EVMChain{
		Metadata: mustMetadata(evmchain.IDOpSepolia),
		IsPublic: true,
	}
)

// IliadEVMByNetwork returns the Iliad evm chain definition by netconf network.
func IliadEVMByNetwork(network netconf.ID) EVMChain {
	return EVMChain{
		Metadata: mustMetadata(network.Static().IliadExecutionChainID),
	}
}

// AnvilChainsByNames returns the Anvil evm chain definitions by names.
func AnvilChainsByNames(names []string) ([]EVMChain, error) {
	var chains []EVMChain
	for _, name := range names {
		meta, ok := evmchain.MetadataByName(name)
		if !ok {
			return nil, errors.New("unknown anvil chain", "name", name)
		}
		chains = append(chains, EVMChain{
			Metadata: meta,
		})
	}

	return chains, nil
}

// PublicChainByName returns the public chain definition by name.
func PublicChainByName(name string) (EVMChain, error) {
	switch name {
	case chainHolesky.Name:
		return chainHolesky, nil
	case chainArbSepolia.Name:
		return chainArbSepolia, nil
	case chainOpSepolia.Name:
		return chainOpSepolia, nil
	case chainEthereum.Name:
		return chainEthereum, nil
	default:
		return EVMChain{}, errors.New("unknown chain name")
	}
}

// PublicRPCByName returns the public chain RPC address by name.
func PublicRPCByName(name string) string {
	switch name {
	case chainHolesky.Name:
		return "https://ethereum-holesky-rpc.publicnode.com"
	case chainArbSepolia.Name:
		return "https://sepolia-rollup.arbitrum.io/rpc"
	case chainOpSepolia.Name:
		return "https://sepolia.optimism.io"
	case chainEthereum.Name:
		return "https://ethereum-rpc.publicnode.com"
	default:
		return ""
	}
}

func mustMetadata(chainID uint64) evmchain.Metadata {
	meta, ok := evmchain.MetadataByID(chainID)
	if !ok {
		panic("unknown chain ID")
	}

	return meta
}
