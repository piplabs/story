package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmtconfig "github.com/cometbft/cometbft/config"
	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	storycfg "github.com/piplabs/story/client/config"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

// InitConfig is the config for the init command.
type InitConfig struct {
	HomeDir         string
	Network         netconf.ID
	TrustedSync     bool
	Force           bool
	Clean           bool
	Cosmos          bool
	ExecutionHash   common.Hash
	RPCLaddr        string
	ExternalAddress string
	Seeds           string
	SeedMode        bool
	Moniker         string
	PersistentPeers string
}

// newInitCmd returns a new cobra command that initializes the files and folders required by story.
func newInitCmd() *cobra.Command {
	// Default config flags
	cfg := InitConfig{
		HomeDir: storycfg.DefaultHomeDir(),
		Force:   false,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes required story files and directories",
		Long: `Initializes required story files and directories.

Ensures all the following files and directories exist:
  <home>/                            # Story home directory
  ├── config                         # Config directory
  │   ├── config.toml                # CometBFT configuration
  │   ├── genesis.json               # Story chain genesis file
  │   ├── story.toml                  # Story configuration
  │   ├── node_key.json              # Node P2P identity key
  │   └── priv_validator_key.json    # CometBFT private validator key (back this up and keep it safe)
  ├── data                           # Data directory
  │   ├── snapshots                  # Snapshot directory
  │   ├── priv_validator_state.json  # CometBFT private validator state (slashing protection)

Existing files are not overwritten, unless --clean is specified.
The home directory should only contain subdirectories, no files, use --force to ignore this check.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := libcmd.LogFlags(ctx, cmd.Flags()); err != nil {
				return err
			}

			return InitFiles(cmd.Context(), cfg)
		},
	}

	bindInitFlags(cmd.Flags(), &cfg)

	return cmd
}

// InitFiles initializes the files and folders required by story.
// It ensures a network and genesis file is generated/downloaded for the provided network.
//
//nolint:gocognit,nestif // This is just many sequential steps.
func InitFiles(ctx context.Context, initCfg InitConfig) error {
	if initCfg.Network == "" {
		return errors.New("required flag --network empty")
	}

	log.Info(ctx, "Initializing story files and directories")
	homeDir := initCfg.HomeDir
	network := initCfg.Network

	if err := prepareHomeDirectory(ctx, initCfg, homeDir); err != nil {
		return err
	}

	// Initialize default configs.
	comet := DefaultCometConfig(homeDir)

	var cfg storycfg.Config

	switch {
	case network == netconf.Iliad:
		cfg = storycfg.IliadConfig
	case network == netconf.Odyssey:
		cfg = storycfg.OdysseyConfig
	case network == netconf.Homer:
		cfg = storycfg.HomerConfig
	case network == netconf.Story:
		cfg = storycfg.StoryConfig
	case network == netconf.Local:
		cfg = storycfg.LocalConfig
	default:
		cfg = storycfg.DefaultConfig()
		cfg.Network = network
	}
	cfg.HomeDir = homeDir

	// Folders
	folders := []struct {
		Name string
		Path string
	}{
		{"home", homeDir},
		{"data", filepath.Join(homeDir, cmtconfig.DefaultDataDir)},
		{"config", filepath.Join(homeDir, cmtconfig.DefaultConfigDir)},
		{"comet db", comet.DBDir()},
		{"snapshot", cfg.SnapshotDir()},
		{"app db", cfg.AppStateDir()},
	}
	for _, folder := range folders {
		if cmtos.FileExists(folder.Path) {
			// Dir exists, just skip
			continue
		} else if err := cmtos.EnsureDir(folder.Path, 0o755); err != nil {
			return errors.Wrap(err, "create folder")
		}
		log.Info(ctx, "Generated folder", "reason", folder.Name, "path", folder.Path)
	}

	if initCfg.Moniker != "" {
		comet.Moniker = initCfg.Moniker
		log.Info(ctx, "Overriding node moniker", "moniker", comet.Moniker)
	}

	if initCfg.RPCLaddr != "" {
		comet.RPC.ListenAddress = initCfg.RPCLaddr
		log.Info(ctx, "Overriding RPC listen address", "address", comet.RPC.ListenAddress)
	}

	if initCfg.ExternalAddress != "" {
		comet.P2P.ExternalAddress = initCfg.ExternalAddress
		log.Info(ctx, "Overriding P2P external address", "address", comet.P2P.ExternalAddress)
	}

	// Handle P2P seeds with prioritization
	if initCfg.Seeds != "" {
		// If seeds are provided via the flag, use them
		seeds := SplitAndTrim(initCfg.Seeds)
		comet.P2P.Seeds = strings.Join(seeds, ",")
		log.Info(ctx, "Overriding P2P seeds with provided flag", "seeds", comet.P2P.Seeds)
	} else if networkSeeds := network.Static().ConsensusSeeds(); len(networkSeeds) > 0 {
		// Otherwise, use the network's default seeds
		comet.P2P.Seeds = strings.Join(networkSeeds, ",")
		log.Info(ctx, "Using network's default P2P seeds", "seeds", comet.P2P.Seeds)
	}

	if initCfg.PersistentPeers != "" {
		persistentPeers := SplitAndTrim(initCfg.PersistentPeers)
		comet.P2P.PersistentPeers = strings.Join(persistentPeers, ",")
		log.Info(ctx, "Overriding P2P persistent peers", "persistent-peers", comet.P2P.PersistentPeers)
	}

	if initCfg.SeedMode {
		comet.P2P.SeedMode = true
		log.Info(ctx, "Seed mode enabled")
	}

	// Setup comet config
	cmtConfigFile := filepath.Join(homeDir, cmtconfig.DefaultConfigDir, cmtconfig.DefaultConfigFileName)
	if cmtos.FileExists(cmtConfigFile) {
		log.Info(ctx, "Found comet config file", "path", cmtConfigFile)
	} else {
		WriteConfigFile(cmtConfigFile, &comet) // This panics on any error :(
		log.Info(ctx, "Generated default comet config file", "path", cmtConfigFile)
	}

	// Setup story config
	storyConfigFile := cfg.ConfigFile()
	if cmtos.FileExists(storyConfigFile) {
		log.Info(ctx, "Found story config file", "path", storyConfigFile)
	} else if err := storycfg.WriteConfigTOML(cfg, log.DefaultConfig()); err != nil {
		return err
	} else {
		log.Info(ctx, "Generated default story config file", "path", storyConfigFile)
	}

	// Setup comet private validator
	var pv *privval.FilePV
	privValKeyFile := comet.PrivValidatorKeyFile()
	privValStateFile := comet.PrivValidatorStateFile()
	if cmtos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile) // This hard exits on any error.
		log.Info(ctx, "Found cometBFT private validator",
			"key_file", privValKeyFile,
			"state_file", privValStateFile,
		)
	} else {
		pv = privval.NewFilePV(k1.GenPrivKey(), privValKeyFile, privValStateFile)
		pv.Save()
		log.Info(ctx, "Generated private validator",
			"key_file", privValKeyFile,
			"state_file", privValStateFile)
	}

	// Setup node key
	nodeKeyFile := comet.NodeKeyFile()
	if cmtos.FileExists(nodeKeyFile) {
		log.Info(ctx, "Found node key", "path", nodeKeyFile)
	} else if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
		return errors.Wrap(err, "load or generate node key")
	} else {
		log.Info(ctx, "Generated node key", "path", nodeKeyFile)
	}

	// Setup genesis file
	genFile := comet.GenesisFile()
	if cmtos.FileExists(genFile) {
		log.Info(ctx, "Found genesis file", "path", genFile)
	} else if len(network.Static().ConsensusGenesisJSON) > 0 {
		if err := os.WriteFile(genFile, network.Static().ConsensusGenesisJSON, 0o644); err != nil {
			return errors.Wrap(err, "failed to write genesis file")
		}
		pubKey, err := pv.GetPubKey()
		if err != nil {
			return errors.Wrap(err, "failed to get public key")
		}

		// Derive the various addresses from the public key
		evmAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKey.Bytes())
		if err != nil {
			return errors.Wrap(err, "failed to convert to evm addr")
		}
		accAddr := sdk.AccAddress(evmAddr.Bytes()).String()
		valAddr := sdk.ValAddress(evmAddr.Bytes()).String()
		pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKey.Bytes())
		fmt.Println("Base64 Encoded Public Key:", pubKeyBase64)

		genesisJSON := string(network.Static().ConsensusGenesisJSON)
		genesisJSON = strings.ReplaceAll(genesisJSON, "{{LOCAL_ACCOUNT_ADDRESS}}", accAddr)
		genesisJSON = strings.ReplaceAll(genesisJSON, "{{LOCAL_VALIDATOR_ADDRESS}}", valAddr)
		genesisJSON = strings.ReplaceAll(genesisJSON, "{{LOCAL_VALIDATOR_KEY}}", pubKeyBase64)

		err = os.WriteFile(genFile, []byte(genesisJSON), 0o644)

		if err != nil {
			return errors.Wrap(err, "save genesis file")
		}
		log.Info(ctx, "Generated well-known network genesis file", "path", genFile)
	} else {
		return errors.New("network genesis file not supported yet", "network", network)
	}

	return nil
}

func checkHomeDir(homeDir string) error {
	files, _ := os.ReadDir(homeDir) // Ignore error, we'll just assume it's empty.
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		return errors.New("home directory contains unexpected file(s), use --force to initialize anyway",
			"home", homeDir, "example_file", file.Name())
	}

	return nil
}

func prepareHomeDirectory(ctx context.Context, initCfg InitConfig, homeDir string) error {
	if !initCfg.Force {
		log.Info(ctx, "Ensuring provided home folder does not contain files, since --force=true")
		if err := checkHomeDir(homeDir); err != nil {
			return err
		}
	}

	if initCfg.Clean {
		log.Info(ctx, "Deleting home directory, since --clean=true")
		if err := os.RemoveAll(homeDir); err != nil {
			return errors.Wrap(err, "remove home dir")
		}
	}

	return nil
}

func SplitAndTrim(input string) []string {
	l := strings.Split(input, ",")
	var ret []string
	for _, r := range l {
		if r = strings.TrimSpace(r); r != "" {
			ret = append(ret, r)
		}
	}

	return ret
}
