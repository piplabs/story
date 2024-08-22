// Package netconf provides the configuration of an Story network, an instance
// of the Story cross chain protocol.
package netconf

import (
	"time"
)

// Network defines a deployment of the Story cross chain protocol.
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
// the Story cross chain protocol. This is most supported EVMs, but
// also the Story EVM, and the Story Consensus chain.
type Chain struct {
	ID   uint64 // Chain ID asa per https://chainlist.org
	Name string // Chain name as per https://chainlist.org
	// RPCURL            string            // RPC URL of the chain
	BlockPeriod time.Duration // Block period of the chain
}
