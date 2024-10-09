package config

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	pruningtypes "cosmossdk.io/store/pruning/types"

	cmtos "github.com/cometbft/cometbft/libs/os"
	db "github.com/cosmos/cosmos-db"

	"github.com/piplabs/story/lib/buildinfo"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tracer"

	_ "embed"
)

const (
	configFile      = "story.toml"
	dataDir         = "data"
	configDir       = "config"
	snapshotDataDir = "snapshots"
	networkFile     = "network.json"

	DefaultEngineEndpoint     = "http://localhost:8551" // Default host endpoint for the Engine API
	defaultSnapshotInterval   = 1000                    // Roughly once an hour (given 3s blocks)
	defaultSnapshotKeepRecent = 2
	defaultMinRetainBlocks    = 0 // Retain all blocks

	defaultPruningOption      = pruningtypes.PruningOptionNothing // Prune nothing
	defaultDBBackend          = db.GoLevelDBBackend
	defaultEVMBuildDelay      = time.Millisecond * 600 // 100ms longer than geth's --miner.recommit=500ms.
	defaultEVMBuildOptimistic = true
)

var (
	IliadConfig = Config{
		HomeDir:            DefaultHomeDir(),
		Network:            "iliad",
		EngineEndpoint:     DefaultEngineEndpoint,
		EngineJWTFile:      DefaultJWTFile("iliad"),
		SnapshotInterval:   defaultSnapshotInterval,
		SnapshotKeepRecent: defaultSnapshotKeepRecent,
		BackendType:        string(defaultDBBackend),
		MinRetainBlocks:    defaultMinRetainBlocks,
		PruningOption:      pruningtypes.PruningOptionDefault,
		EVMBuildDelay:      defaultEVMBuildDelay,
		EVMBuildOptimistic: false,
		APIEnable:          false,
		APIAddress:         "127.0.0.1:1317",
		EnableUnsafeCORS:   false,
		Tracer:             tracer.DefaultConfig(),
		RPCLaddr:           "tcp://127.0.0.1:26657",
		ExternalAddress:    "",
		Seeds:              "",
		SeedMode:           false,
	}
	LocalConfig = Config{
		HomeDir:            DefaultHomeDir(),
		Network:            "local",
		EngineEndpoint:     DefaultEngineEndpoint,
		EngineJWTFile:      DefaultJWTFile("local"),
		SnapshotInterval:   defaultSnapshotInterval,
		SnapshotKeepRecent: defaultSnapshotKeepRecent,
		BackendType:        string(defaultDBBackend),
		MinRetainBlocks:    defaultMinRetainBlocks,
		PruningOption:      pruningtypes.PruningOptionDefault,
		EVMBuildDelay:      defaultEVMBuildDelay,
		EVMBuildOptimistic: false,
		APIEnable:          false,
		APIAddress:         "127.0.0.1:1317",
		EnableUnsafeCORS:   false,
		Tracer:             tracer.DefaultConfig(),
		RPCLaddr:           "tcp://127.0.0.1:26657",
		ExternalAddress:    "",
		Seeds:              "",
		SeedMode:           false,
	}
)

// DefaultConfig returns the default story config.
func DefaultConfig() Config {
	return Config{
		HomeDir:            DefaultHomeDir(),
		Network:            "",                      // No default
		EngineEndpoint:     "http://localhost:8551", // No default
		EngineJWTFile:      "",                      // No default
		SnapshotInterval:   defaultSnapshotInterval,
		SnapshotKeepRecent: defaultSnapshotKeepRecent,
		BackendType:        string(defaultDBBackend),
		MinRetainBlocks:    defaultMinRetainBlocks,
		PruningOption:      defaultPruningOption,
		EVMBuildDelay:      defaultEVMBuildDelay,
		EVMBuildOptimistic: defaultEVMBuildOptimistic,
		APIEnable:          false,
		APIAddress:         "127.0.0.1:1317",
		EnableUnsafeCORS:   false,
		Tracer:             tracer.DefaultConfig(),
		RPCLaddr:           "tcp://127.0.0.1:26657",
		ExternalAddress:    "",
		Seeds:              "",
		SeedMode:           false,
	}
}

// DefaultJWTFile returns the default engine-api jwt file assumed to be used by the execution client.
func DefaultJWTFile(network string) string {
	return filepath.Join(baseDir(), "geth", network, "geth", "jwtsecret")
}

// DefaultHomeDir returns the default consensus client home directory.
func DefaultHomeDir() string {
	return filepath.Join(baseDir(), "story")
}

// baseDir generates the base directory path for the "Story" application.
func baseDir() string {
	home := homeDir()
	if home != "" {
		switch runtime.GOOS {
		case "darwin":
			return filepath.Join(home, "Library", "Story")
		case "windows":
			return filepath.Join(home, "AppData", "Roaming", "Story")
		default:
			return filepath.Join(home, ".story")
		}
	}

	return ""
}

// Config defines all story specific config.
type Config struct {
	HomeDir            string
	Network            netconf.ID
	EthKeyPassword     string
	EngineJWTFile      string
	EngineEndpoint     string
	SnapshotInterval   uint64 // See cosmossdk.io/store/snapshots/types/options.go
	SnapshotKeepRecent uint64 // See cosmossdk.io/store/snapshots/types/options.go
	BackendType        string // See cosmos-db/db.go
	MinRetainBlocks    uint64
	PruningOption      string // See cosmossdk.io/store/pruning/types/options.go
	EVMBuildDelay      time.Duration
	EVMBuildOptimistic bool
	APIEnable          bool
	APIAddress         string
	EnableUnsafeCORS   bool
	Tracer             tracer.Config
	RPCLaddr           string
	ExternalAddress    string
	Seeds              string
	SeedMode           bool
	RemoveBlock        bool // See cosmos-sdk/server/rollback.go
}

// ConfigFile returns the default path to the toml story config file.
func (c Config) ConfigFile() string {
	return filepath.Join(c.HomeDir, configDir, configFile)
}

func (c Config) DataDir() string {
	return filepath.Join(c.HomeDir, dataDir)
}

func (c Config) AppStateDir() string {
	return c.DataDir() // Maybe add a subdirectory for app state?
}

func (c Config) SnapshotDir() string {
	return filepath.Join(c.DataDir(), snapshotDataDir)
}

func (c Config) Verify() error {
	if c.EngineEndpoint == "" {
		return errors.New("flag --engine-endpoint is empty")
	} else if c.EngineJWTFile == "" {
		return errors.New("flag --engine-jwt-file is empty")
	} else if c.Network == "" {
		return errors.New("flag --network is empty")
	} else if err := c.Network.Verify(); err != nil {
		return err
	}

	return nil
}

//go:embed config.toml.tmpl
var tomlTemplate []byte

// WriteConfigTOML writes the toml story config to disk.
func WriteConfigTOML(cfg Config, logCfg log.Config) error {
	var buffer bytes.Buffer

	t, err := template.New("").Parse(string(tomlTemplate))
	if err != nil {
		return errors.Wrap(err, "parse template")
	}

	s := struct {
		Config
		Log     log.Config
		Version string
	}{
		Config:  cfg,
		Log:     logCfg,
		Version: buildinfo.Version(),
	}

	if err := t.Execute(&buffer, s); err != nil {
		return errors.Wrap(err, "execute template")
	}

	if err := cmtos.WriteFile(cfg.ConfigFile(), buffer.Bytes(), 0o644); err != nil {
		return errors.Wrap(err, "write config")
	}

	return nil
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}

	return ""
}
