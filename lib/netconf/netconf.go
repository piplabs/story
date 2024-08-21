// Package netconf provides the configuration of an Iliad network, an instance
// of the Iliad cross chain protocol.
package netconf

import (
	"time"

	"github.com/piplabs/story/lib/evmchain"
)

// Network defines a deployment of the Iliad cross chain protocol.
// It spans an iliad chain (both execution and consensus) and a set of
// supported EVMs.
type Network struct {
	ID     ID      `json:"name"`   // ID of the network. e.g. "simnet", "testnet", "staging", "mainnet"
	Chains []Chain `json:"chains"` // Chains that are part of the network
}

// Validate returns an error if the configuration is invalid.
func (n Network) Validate() error {
	if err := n.ID.Verify(); err != nil {
		return err
	}

	// TODO: Validate chains

	return nil
}

// EVMChains returns all evm chains in the network. It excludes the iliad consensus chain.
func (n Network) EVMChains() []Chain {
	resp := make([]Chain, 0, len(n.Chains))
	for _, chain := range n.Chains {
		if IsIliadConsensus(n.ID, chain.ID) {
			continue
		}

		resp = append(resp, chain)
	}

	return resp
}

// ChainIDs returns the all chain IDs in the network.
// This is a convenience method.
func (n Network) ChainIDs() []uint64 {
	resp := make([]uint64, 0, len(n.Chains))
	for _, chain := range n.Chains {
		resp = append(resp, chain.ID)
	}

	return resp
}

// ChainNamesByIDs returns the all chain IDs and names in the network.
// This is a convenience method.
func (n Network) ChainNamesByIDs() map[uint64]string {
	resp := make(map[uint64]string)
	for _, chain := range n.Chains {
		resp[chain.ID] = chain.Name
	}

	return resp
}

// IliadEVMChain returns the Iliad execution chain config or false if it does not exist.
func (n Network) IliadEVMChain() (Chain, bool) {
	for _, chain := range n.Chains {
		if IsIliadExecution(n.ID, chain.ID) {
			return chain, true
		}
	}

	return Chain{}, false
}

// IliadConsensusChain returns the Iliad consensus chain config or false if it does not exist.
func (n Network) IliadConsensusChain() (Chain, bool) {
	for _, chain := range n.Chains {
		if IsIliadConsensus(n.ID, chain.ID) {
			return chain, true
		}
	}

	return Chain{}, false
}

// EthereumChain returns the ethereum Layer1 chain config or false if it does not exist.
func (n Network) EthereumChain() (Chain, bool) {
	for _, chain := range n.Chains {
		switch n.ID {
		case Mainnet:
			if chain.ID == evmchain.IDEthereum {
				return chain, true
			}
		case Local:
			if chain.ID == evmchain.IDLocal {
				return chain, true
			}
		case Iliad:
			if chain.ID == evmchain.IDIliad {
				return chain, true
			}
		case Testnet:
			if chain.ID == evmchain.IDHolesky {
				return chain, true
			}
		default:
			if chain.ID == evmchain.IDMockL1Fast || chain.ID == evmchain.IDMockL1Slow {
				return chain, true
			}
		}
	}

	return Chain{}, false
}

// ChainName returns the chain name for the given ID or an empty string if it does not exist.
func (n Network) ChainName(id uint64) string {
	chain, _ := n.Chain(id)
	return chain.Name
}

// Chain returns the chain config for the given ID or false if it does not exist.
func (n Network) Chain(id uint64) (Chain, bool) {
	for _, chain := range n.Chains {
		if chain.ID == id {
			return chain, true
		}
	}

	return Chain{}, false
}

// Chain defines the configuration of an execution chain that supports
// the Iliad cross chain protocol. This is most supported EVMs, but
// also the Iliad EVM, and the Iliad Consensus chain.
type Chain struct {
	ID   uint64 // Chain ID asa per https://chainlist.org
	Name string // Chain name as per https://chainlist.org
	// RPCURL            string            // RPC URL of the chain
	BlockPeriod time.Duration // Block period of the chain
}
