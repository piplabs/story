package netconf

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/evmchain"

	_ "embed"
)

const consensusIDPrefix = "iliad-"
const consensusIDOffset = 1_000_000
const maxValidators = 10

// Static defines static config and data for a network.
type Static struct {
	Version               string
	IliadExecutionChainID uint64
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

// IliadConsensusChainIDStr returns the chain ID string for the Iliad consensus chain.
// It is calculated as "iliad-<IliadConsensusChainIDUint64>".
func (s Static) IliadConsensusChainIDStr() string {
	return fmt.Sprintf("%s%d", consensusIDPrefix, s.IliadConsensusChainIDUint64())
}

// IliadConsensusChainIDUint64 returns the chain ID uint64 for the Iliad consensus chain.
// It is calculated as 1_000_000 + IliadExecutionChainID.
func (s Static) IliadConsensusChainIDUint64() uint64 {
	return consensusIDOffset + s.IliadExecutionChainID
}

// IliadConsensusChain returns the story consensus Chain struct.
func (s Static) IliadConsensusChain() Chain {
	return Chain{
		ID:          s.IliadConsensusChainIDUint64(),
		Name:        "iliad_consensus",
		BlockPeriod: time.Second * 2,
	}
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

func (s Static) ExecutionSeeds() []string {
	var resp []string
	for _, seed := range strings.Split(string(s.ExecutionSeedTXT), "\n") {
		if seed = strings.TrimSpace(seed); seed != "" {
			resp = append(resp, seed)
		}
	}

	return resp
}

// Use random runid for version in ephemeral networks.
//
//nolint:gochecknoglobals // Static ID
var runid = uuid.New().String()

//nolint:gochecknoglobals // Static addresses
var (
	//go:embed testnet/consensus-genesis.json
	testnetConsensusGenesisJSON []byte

	//go:embed testnet/consensus-seeds.txt
	testnetConsensusSeedsTXT []byte

	//go:embed testnet/execution-genesis.json
	testnetExecutionGenesisJSON []byte

	//go:embed testnet/execution-seeds.txt
	testnetExecutionSeedsTXT []byte
)

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

//nolint:gochecknoglobals // Static mappings.
var statics = map[ID]Static{
	Simnet: {
		Version:               runid,
		IliadExecutionChainID: evmchain.IDIliadEphemeral,
		MaxValidators:         maxValidators,
	},
	Devnet: {
		Version:               runid,
		IliadExecutionChainID: evmchain.IDIliadEphemeral,
		MaxValidators:         maxValidators,
	},
	Staging: {
		Version:               runid,
		IliadExecutionChainID: evmchain.IDIliadEphemeral,
		MaxValidators:         maxValidators,
	},
	Testnet: {
		Version:               "v0.0.1",
		IliadExecutionChainID: evmchain.IDIliadTestnet,
		MaxValidators:         maxValidators,
		ConsensusGenesisJSON:  testnetConsensusGenesisJSON,
		ConsensusSeedTXT:      testnetConsensusSeedsTXT,
		ExecutionGenesisJSON:  testnetExecutionGenesisJSON,
		ExecutionSeedTXT:      testnetExecutionSeedsTXT,
	},
	Iliad: {
		Version:               "v0.0.1",
		IliadExecutionChainID: evmchain.IDIliadTestnet,
		ConsensusGenesisJSON:  iliadConsensusGenesisJSON,
		ConsensusSeedTXT:      iliadConsensusSeedsTXT,
	},
	Local: {
		Version:               "v0.0.1",
		IliadExecutionChainID: evmchain.IDLocal,
		ConsensusGenesisJSON:  localConsensusGenesisJSON,
		ConsensusSeedTXT:      localConsensusSeedsTXT,
	},
	Mainnet: {
		Version:       "v0.0.1",
		MaxValidators: maxValidators,
	},
}

// ConsensusChainIDStr2Uint64 parses the uint suffix from the provided a consensus chain ID string.
func ConsensusChainIDStr2Uint64(id string) (uint64, error) {
	if !strings.HasPrefix(id, consensusIDPrefix) {
		return 0, errors.New("invalid consensus chain ID", "id", id)
	}

	suffix := strings.TrimPrefix(id, consensusIDPrefix)

	resp, err := strconv.ParseUint(suffix, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse consensus chain ID", "id", id)
	}

	return resp, nil
}

// IsIliadConsensus returns true if provided chainID is the iliad consensus chain for the network.
func IsIliadConsensus(network ID, chainID uint64) bool {
	return network.Static().IliadConsensusChainIDUint64() == chainID
}

// IsIliadExecution returns true if provided chainID is the iliad execution chain for the network.
func IsIliadExecution(network ID, chainID uint64) bool {
	return network.Static().IliadExecutionChainID == chainID
}
