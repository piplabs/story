package config

import (
	"fmt"
	"time"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/netconf"
)

// DKGConfig represents the configuration for the DKG service.
type DKGConfig struct {
	// Enable determines whether the DKG service is enabled
	Enable bool `mapstructure:"enable"`

	// TEEEndpoint is the URL/address of the TEE client
	TEEEndpoint string `mapstructure:"tee_endpoint"`

	// RPCEndpoint is the Cosmos SDK consensus client RPC endpoint
	CosmosRPCEndpoint string `mapstructure:"cosmos_rpc_endpoint"`

	// Cosmos chain ID
	CosmosChainID string `mapstructure:"cosmos_chain_id"`

	// Event polling interval for listening to blockchain events
	EventPollingInterval time.Duration `mapstructure:"event_polling_interval"`

	// Maximum number of retries for TEE operations
	MaxTEERetries int `mapstructure:"max_tee_retries"`

	// TEE operation timeout
	TEETimeout time.Duration `mapstructure:"tee_timeout"`

	// Local storage path for DKG state
	DataDir string `mapstructure:"data_dir"`

	// Log level for DKG service
	LogLevel string `mapstructure:"log_level"`
}

// DefaultDKGConfig returns a default DKG configuration.
func DefaultDKGConfig() *DKGConfig {
	return &DKGConfig{
		Enable:               false,
		TEEEndpoint:          "",
		CosmosRPCEndpoint:    "http://localhost:26657",
		CosmosChainID:        netconf.TestChainID,
		EventPollingInterval: 5 * time.Second, // Poll every 5 seconds
		MaxTEERetries:        3,
		TEETimeout:           30 * time.Second,
		DataDir:              "data/dkg",
		LogLevel:             "info",
	}
}

// Validate validates the DKG configuration.
func (c *DKGConfig) Validate() error {
	if !c.Enable {
		return nil // Skip validation if DKG is disabled
	}

	if c.TEEEndpoint == "" {
		return errors.New("tee endpoint is required")
	}

	if c.CosmosRPCEndpoint == "" {
		return errors.New("rpc endpoint is required when dkg is enabled")
	}

	if c.CosmosChainID == "" {
		return errors.New("cosmos chain id is required")
	}

	if c.EventPollingInterval <= 0 {
		return errors.New("event polling interval must be positive")
	}

	if c.MaxTEERetries < 1 {
		return errors.New("max tee retries must be at least 1")
	}

	if c.TEETimeout <= 0 {
		return errors.New("tee timeout must be positive")
	}

	if c.DataDir == "" {
		return errors.New("data directory is required")
	}

	return nil
}

// String returns a string representation of the configuration.
func (c *DKGConfig) String() string {
	return fmt.Sprintf(
		"DKGConfig{Enable:%v, TEEEndpoint:%s, CosmosRPCEndpoint:%s, EventPollingInterval:%v, MaxTEERetries:%d, "+
			"TEETimeout:%v, DataDir:%s, LogLevel:%s}",
		c.Enable, c.TEEEndpoint, c.CosmosRPCEndpoint, c.EventPollingInterval, c.MaxTEERetries,
		c.TEETimeout, c.DataDir, c.LogLevel,
	)
}
