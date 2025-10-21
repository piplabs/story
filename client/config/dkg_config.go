package config

import (
	"github.com/piplabs/story/lib/errors"
	"github.com/spf13/pflag"
)

type DKGConfig struct {
	// Enable enables or disables the DKG client
	Enable bool

	// TEEEndpoint is the URL/address of the TEE client
	TEEEndpoint string

	// EngineRPCEndpoint is the RPC endpoint of the execution layer
	EngineRPCEndpoint string
}

func DefaultDKGConfig() DKGConfig {
	return DKGConfig{
		Enable:            false,
		TEEEndpoint:       "127.0.0.1:50051",
		EngineRPCEndpoint: "http://127.0.0.1:8545",
	}
}

func BindDKGFlags(flags *pflag.FlagSet, cfg *DKGConfig) {
	flags.BoolVar(&cfg.Enable, "dkg-enable", cfg.Enable, "DKG client is enabled or not")
	flags.StringVar(&cfg.TEEEndpoint, "dkg-tee-endpoint", cfg.TEEEndpoint, "The endpoint of TEE client for DKG")
	flags.StringVar(&cfg.EngineRPCEndpoint, "dkg-engine-rpc-endpoint", cfg.EngineRPCEndpoint, "The RPC endpoint of execution layer")
}

func (c *DKGConfig) Validate() error {
	if !c.Enable {
		return nil
	}

	if c.TEEEndpoint == "" {
		return errors.New("tee endpoint should not be empty")
	}

	if c.EngineRPCEndpoint == "" {
		return errors.New("engine rpc endpoint should not be empty")
	}

	return nil
}
