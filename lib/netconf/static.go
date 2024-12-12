package netconf

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/lib/evmchain"

	_ "embed"
)

const consensusID = "story-1"

// Static defines static config and data for a network.
type Static struct {
	Version               string
	StoryExecutionChainID uint64
	MaxValidators         uint32
	ConsensusGenesisJSON  []byte
	ConsensusSeedTXT      []byte
	ExecutionGenesisJSON  []byte
	ExecutionSeedTXT      []byte
}

type Deployment struct {
	ChainID uint64
	Address common.Address
}

// StoryConsensusChainIDStr returns the chain ID string for the Story consensus client.
func (Static) StoryConsensusChainIDStr() string {
	return consensusID
}

func (s Static) ConsensusSeeds() []string {
	var resp []string
	for _, seed := range strings.Split(string(s.ConsensusSeedTXT), "\n") {
		if seed = strings.TrimSpace(seed); seed != "" {
			resp = append(resp, seed)
		}
	}

	return resp
}

//nolint:gochecknoglobals // Static addresses
var (
	//go:embed iliad/genesis.json
	iliadConsensusGenesisJSON []byte

	//go:embed iliad/seeds.txt
	iliadConsensusSeedsTXT []byte
)

//nolint:gochecknoglobals // Static addresses
var (
	//go:embed local/genesis.json
	localConsensusGenesisJSON []byte

	//go:embed local/seeds.txt
	localConsensusSeedsTXT []byte
)

//nolint:gochecknoglobals // Static addresses
var (
	//go:embed odyssey/genesis.json
	odysseyConsensusGenesisJSON []byte

	//go:embed odyssey/seeds.txt
	odysseyConsensusSeedsTXT []byte
)

//nolint:gochecknoglobals // Static addresses
var (
	//go:embed mainnet/genesis.json
	mainnetConsensusGenesisJSON []byte

	//go:embed mainnet/seeds.txt
	mainnetConsensusSeedsTXT []byte
)

//nolint:gochecknoglobals // Static mappings.
var statics = map[ID]Static{
	Iliad: {
		Version:               "v0.0.1",
		StoryExecutionChainID: evmchain.IDIliad,
		ConsensusGenesisJSON:  iliadConsensusGenesisJSON,
		ConsensusSeedTXT:      iliadConsensusSeedsTXT,
	},
	Local: {
		Version:               "v0.0.1",
		StoryExecutionChainID: evmchain.IDLocal,
		ConsensusGenesisJSON:  localConsensusGenesisJSON,
		ConsensusSeedTXT:      localConsensusSeedsTXT,
	},
	Odyssey: {
		Version:               "v0.0.1",
		StoryExecutionChainID: evmchain.IDOdyssey,
		ConsensusGenesisJSON:  odysseyConsensusGenesisJSON,
		ConsensusSeedTXT:      odysseyConsensusSeedsTXT,
	},
	Mainnet: {
		Version:               "v0.0.1",
		StoryExecutionChainID: evmchain.IDStoryMainnet,
		ConsensusGenesisJSON:  mainnetConsensusGenesisJSON,
		ConsensusSeedTXT:      mainnetConsensusSeedsTXT,
	},
}
