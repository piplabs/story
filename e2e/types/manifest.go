package types

import (
	e2e "github.com/cometbft/cometbft/test/e2e/pkg"

	"github.com/piplabs/story/e2e/app/key"
	"github.com/piplabs/story/lib/netconf"
)

// Mode defines the iliad consensus node mode.
// Nodes are in general full nodes (light nodes are not supported yet).
// In some cases, additional roles are defined: validator, archive, seed.
//
// Note that the execution clients only have two modes: "default" and "archive".
//
// e2e.Mode is extended so ModeArchive can be added transparently.
type Mode = e2e.Mode

const (
	// ModeValidator defines a validator node.
	// It's validator key has staked tokens and it actively participates in consensus and is subject to rewards and penalties.
	// It must always be online, otherwise it will get stashed/jailed.
	// [genesis_validator_set=true,pruning=default,consensus=default,special_p2p=false].
	// Note technically a validator node is also a "full node".
	ModeValidator = e2e.ModeValidator

	// ModeArchive defines an archive node.
	// It stores all historical blocks and state, it doesn't delete anything ever. It will require TBs of disk.
	// [genesis_validator_set=false,pruning=none,consensus=default,special_p2p=false].
	// Note technically an archive node is also a "full node".
	ModeArchive Mode = "archive"

	// ModeSeed defines a seed node. It must have a long-lived p2p pubkey and address (encoded in repo).
	// It acts as notice board for external nodes to learn about the network and connect to publicly available nodes.
	// It crawls the network regularly, making it available to new nodes.
	// [genesis_validator_set=false,pruning=default,consensus=default,special_p2p=true].
	// Note technically a seed node is also a "full node".
	ModeSeed = e2e.ModeSeed

	// ModeFull defines a full node. A full node a normal node without a special role.
	// [genesis_validator_set=false,pruning=default,consensus=default,special_p2p=false].
	ModeFull = e2e.ModeFull

	// ModeLight defines a light node. This isn't used yet.
	// [genesis_validator_set=false,pruning=no_data,consensus=light,special_p2p=false]
	// Only light nodes are not also full nodes.
	ModeLight = e2e.ModeLight
)

// Perturb defines non-cometBFT perturbations of components like iliad_evm.
type Perturb string

const (
	// PerturbRestart defines a perturbation that restarts a docker container.
	PerturbRestart Perturb = "restart"
	// PerturbStopStart defines a perturbation that stops and then starts a docker container.
	PerturbStopStart Perturb = "stopstart"
)

// Manifest wraps e2e.Manifest with additional iliad-specific fields.
//

type Manifest struct {
	e2e.Manifest

	Network netconf.ID `toml:"network"`

	// AnvilChains defines the anvil chains to deploy; mock_l1, mock_l2, etc.
	AnvilChains []string `toml:"anvil_chains"`

	// PublicChains defines the public chains to connect to; arb_sepolia, etc.
	PublicChains []string `toml:"public_chains"`

	// MultiIliadEVMs defines whether to deploy one or multiple Iliad EVMs.
	MultiIliadEVMs bool `toml:"multi_iliad_evms"`

	// Keys contains long-lived private keys (address by type) by node name.
	Keys map[string]map[key.Type]string `toml:"keys"`

	// Perturb defines additional (non-cometBFT) perturbations by service name.
	Perturb map[string][]Perturb `json:"perturb"`
}

// Seeds returns a map of seed nodes by name.
func (m Manifest) Seeds() map[string]bool {
	resp := make(map[string]bool)
	for name, node := range m.Nodes {
		if Mode(node.Mode) == ModeSeed {
			resp[name] = true
		}
	}

	return resp
}

// IliadEVMs returns a map of iliad evm instances names by <IsArchive> to deploy.
// If only a single Iliad EVM is to be deployed, the name is "iliad_evm".
// Otherwise, the names are "<node>_evm".
func (m Manifest) IliadEVMs() map[string]bool {
	if !m.MultiIliadEVMs {
		return map[string]bool{
			"iliad_evm": false,
		}
	}

	resp := make(map[string]bool)
	for name, node := range m.Nodes {
		resp[name+"_evm"] = Mode(node.Mode) == ModeArchive
	}

	return resp
}
