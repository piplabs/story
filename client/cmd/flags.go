package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/piplabs/story/client/config"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tracer"

	// Used for ABI embedding of the staking contract.
	_ "embed"
)

func bindRunFlags(cmd *cobra.Command, cfg *config.Config) {
	flags := cmd.Flags()

	libcmd.BindHomeFlag(flags, &cfg.HomeDir)
	tracer.BindFlags(flags, &cfg.Tracer)
	netconf.BindFlag(flags, &cfg.Network)
	flags.StringVar(&cfg.EngineEndpoint, "engine-endpoint", cfg.EngineEndpoint, "An EVM execution client Engine API http endpoint")
	flags.StringVar(&cfg.EngineJWTFile, "engine-jwt-file", cfg.EngineJWTFile, "The path to the Engine API JWT file")
	flags.Uint64Var(&cfg.SnapshotInterval, "snapshot-interval", cfg.SnapshotInterval, "State sync snapshot interval")
	flags.Uint64Var(&cfg.SnapshotKeepRecent, "snapshot-keep-recent", cfg.SnapshotKeepRecent, "State sync snapshot to keep")
	flags.Uint64Var(&cfg.MinRetainBlocks, "min-retain-blocks", cfg.MinRetainBlocks, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	flags.StringVar(&cfg.BackendType, "app-db-backend", cfg.BackendType, "The type of database for application and snapshots databases")
	flags.StringVar(&cfg.PruningOption, "pruning", cfg.PruningOption, "Pruning strategy (default|nothing|everything)")
	flags.DurationVar(&cfg.EVMBuildDelay, "evm-build-delay", cfg.EVMBuildDelay, "Minimum delay between triggering and fetching a EVM payload build")
	flags.BoolVar(&cfg.EVMBuildOptimistic, "evm-build-optimistic", cfg.EVMBuildOptimistic, "Enables optimistic building of EVM payloads on previous block finalize")
	flags.BoolVar(&cfg.APIEnable, "api-enable", cfg.APIEnable, "Define if the API server should be enabled")
	flags.StringVar(&cfg.APIAddress, "api-address", cfg.APIAddress, "The API server address to listen on")
	flags.BoolVar(&cfg.EnableUnsafeCORS, "enabled-unsafe-cors", cfg.EnableUnsafeCORS, "Enable unsafe CORS for API server")
}

func bindInitFlags(flags *pflag.FlagSet, cfg *InitConfig) {
	libcmd.BindHomeFlag(flags, &cfg.HomeDir)
	netconf.BindFlag(flags, &cfg.Network)
	flags.BoolVar(&cfg.TrustedSync, "trusted-sync", cfg.TrustedSync, "Initialize trusted state-sync height and hash by querying the Story RPC")
	flags.BoolVar(&cfg.Force, "force", cfg.Force, "Force initialization (overwrite existing files)")
	flags.BoolVar(&cfg.Clean, "clean", cfg.Clean, "Delete home directory before initialization")
	flags.StringVar(&cfg.RPCLaddr, "rpc-laddr", "", "Override the RPC listening address")
	flags.StringVar(&cfg.ExternalAddress, "external-address", "", "Override the P2P external address")
	flags.StringVar(&cfg.Seeds, "seeds", "", "Override the P2P seeds (comma-separated)")
	flags.BoolVar(&cfg.SeedMode, "seed-mode", false, "Enable seed mode")
	flags.StringVar(&cfg.Moniker, "moniker", cfg.Moniker, "Custom moniker name for this node")
}

func bindValidatorBaseFlags(cmd *cobra.Command, cfg *baseConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://testnet.storyrpc.io", "RPC URL to connect to the testnet")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used for the transaction")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://testnet.storyscan.xyz", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1513, "Chain ID to use for the transaction (default 1513)")
}

func bindValidatorCreateFlags(cmd *cobra.Command, cfg *createValidatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	bindValidatorKeyFlags(cmd, &cfg.ValidatorKeyFile)
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount for the validator to self-delegate in wei")
}

func bindAddOperatorFlags(cmd *cobra.Command, cfg *operatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.Operator, "operator", "", "Adds an operator to your delegator")
}

func bindRemoveOperatorFlags(cmd *cobra.Command, cfg *operatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.Operator, "operator", "", "Removes an operator from your delegator")
}

func bindSetWithdrawalAddressFlags(cmd *cobra.Command, cfg *withdrawalConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.WithdrawalAddress, "withdrawal-address", "", "Address to receive staking and reward withdrawals")
}

func bindValidatorStakeFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount to stake on behalf of the delegator in wei")
}

func bindValidatorStakeOnBehalfFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.DelegatorPubKey, "delegator-pubkey", "", "Delegator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount to stake on behalf of the delegator in wei")
}

func bindValidatorUnstakeFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "unstake", "", "Amount to unstake on behalf of the delegator in wei")
}

func bindValidatorUnstakeOnBehalfFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.DelegatorPubKey, "delegator-pubkey", "", "Delegator's base64-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "unstake", "", "Amount to unstake on behalf of the delegator in wei")
}

func bindValidatorKeyExportFlags(cmd *cobra.Command, cfg *exportKeyConfig) {
	bindValidatorKeyFlags(cmd, &cfg.ValidatorKeyFile)
	defaultEVMKeyFilePath := filepath.Join(config.DefaultHomeDir(), "config", "private_key.txt")
	cmd.Flags().BoolVar(&cfg.ExportEVMKey, "export-evm-key", false, "Export the EVM private key")
	cmd.Flags().StringVar(&cfg.EvmKeyFile, "evm-key-path", defaultEVMKeyFilePath, "Path to save the exported EVM private key")
}

func bindValidatorKeyFlags(cmd *cobra.Command, keyFilePath *string) {
	defaultKeyFilePath := filepath.Join(config.DefaultHomeDir(), "config", "priv_validator_key.json")
	cmd.Flags().StringVar(keyFilePath, "keyfile", defaultKeyFilePath, "Path to the Tendermint key file")
}

func bindStatusFlags(flags *pflag.FlagSet, cfg *StatusConfig) {
	libcmd.BindHomeFlag(flags, &cfg.HomeDir)
}

func bindRollbackFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().BoolVar(&cfg.RemoveBlock, "hard", false, "remove last block as well as state")
}

// Flag Validation

func validateFlags(flags map[string]string) error {
	var missingFlags []string

	for flag, value := range flags {
		if value == "" {
			missingFlags = append(missingFlags, flag)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func validateValidatorCreateFlags(cfg createValidatorConfig) error {
	return validateFlags(map[string]string{
		"rpc":     cfg.RPC,
		"keyfile": cfg.ValidatorKeyFile,
		"stake":   cfg.StakeAmount,
	})
}

func validateOperatorFlags(cfg operatorConfig) error {
	return validateFlags(map[string]string{
		"rpc":      cfg.RPC,
		"operator": cfg.Operator,
	})
}

func validateWithdrawalFlags(cfg withdrawalConfig) error {
	return validateFlags(map[string]string{
		"rpc":                cfg.RPC,
		"withdrawal-address": cfg.WithdrawalAddress,
	})
}

func validateValidatorStakeFlags(cfg stakeConfig) error {
	return validateFlags(map[string]string{
		"rpc":              cfg.RPC,
		"validator-pubkey": cfg.ValidatorPubKey,
		"stake":            cfg.StakeAmount,
	})
}

func validateValidatorStakeOnBehalfFlags(cfg stakeConfig) error {
	return validateFlags(map[string]string{
		"rpc":              cfg.RPC,
		"validator-pubkey": cfg.ValidatorPubKey,
		"delegator-pubkey": cfg.DelegatorPubKey,
		"stake":            cfg.StakeAmount,
	})
}

func validateValidatorUnstakeOnBehalfFlags(cfg stakeConfig) error {
	return validateFlags(map[string]string{
		"rpc":              cfg.RPC,
		"validator-pubkey": cfg.ValidatorPubKey,
		"delegator-pubkey": cfg.DelegatorPubKey,
		"unstake":          cfg.StakeAmount,
	})
}
