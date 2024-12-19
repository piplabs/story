package cmd

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/piplabs/story/client/config"
	apisvr "github.com/piplabs/story/client/server"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/errors"
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
	apisvr.BindFlags(flags, &cfg.API)
	flags.StringVar(&cfg.EngineEndpoint, "engine-endpoint", cfg.EngineEndpoint, "An EVM execution client Engine API http endpoint")
	flags.StringVar(&cfg.EngineJWTFile, "engine-jwt-file", cfg.EngineJWTFile, "The path to the Engine API JWT file")
	flags.Uint64Var(&cfg.SnapshotInterval, "snapshot-interval", cfg.SnapshotInterval, "State sync snapshot interval")
	flags.Uint64Var(&cfg.SnapshotKeepRecent, "snapshot-keep-recent", cfg.SnapshotKeepRecent, "State sync snapshot to keep")
	flags.Uint64Var(&cfg.MinRetainBlocks, "min-retain-blocks", cfg.MinRetainBlocks, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	flags.StringVar(&cfg.BackendType, "app-db-backend", cfg.BackendType, "The type of database for application and snapshots databases")
	flags.StringVar(&cfg.PruningOption, "pruning", cfg.PruningOption, "Pruning strategy (default|nothing|everything)")
	flags.DurationVar(&cfg.EVMBuildDelay, "evm-build-delay", cfg.EVMBuildDelay, "Minimum delay between triggering and fetching a EVM payload build")
	flags.BoolVar(&cfg.EVMBuildOptimistic, "evm-build-optimistic", cfg.EVMBuildOptimistic, "Enables optimistic building of EVM payloads on previous block finalize")
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
	flags.StringVar(&cfg.PersistentPeers, "persistent-peers", "", "Override the persistent peers (comma-separated)")
	flags.StringVar(&cfg.Moniker, "moniker", "", "Declare a custom moniker for your node")
}

func bindValidatorBaseFlags(cmd *cobra.Command, cfg *baseConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://storyrpc.io", "RPC URL to connect to the network")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used for the transaction")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://storyscan.xyz", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1415, "Chain ID to use for the transaction")
}

func bindValidatorCreateFlags(cmd *cobra.Command, cfg *createValidatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	bindValidatorKeyFlags(cmd, &cfg.ValidatorKeyFile)
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "1024000000000000000000", "Amount for the validator to self-delegate in wei")
	cmd.Flags().Uint32Var(&cfg.CommissionRate, "commission-rate", 1000, "The validator commission rate in bips (e.g. 1000 for 10%)")
	cmd.Flags().Uint32Var(&cfg.MaxCommissionRate, "max-commission-rate", 5000, "The maximum validator commission rate in bips, e.g. 5000 for 50%")
	cmd.Flags().Uint32Var(&cfg.MaxCommissionChangeRate, "max-commission-change-rate", 1000, "The maximum validator commission change rate in bips, e.g. 100 for 1%")
	cmd.Flags().BoolVar(&cfg.Unlocked, "unlocked", true, "Whether to support unlocked token staking (true for unlocked staking, false for locked staking)")
	cmd.Flags().StringVar(&cfg.Moniker, "moniker", "", "Custom moniker name for this node")
}

func bindSetOperatorFlags(cmd *cobra.Command, cfg *operatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.Operator, "operator", "", "Sets an operator to your delegator")
}

func bindUnsetOperatorFlags(cmd *cobra.Command, cfg *operatorConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
}

func bindSetWithdrawalAddressFlags(cmd *cobra.Command, cfg *withdrawalConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.WithdrawalAddress, "withdrawal-address", "", "Address to receive staking and reward withdrawals")
}

func bindValidatorStakeFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount for the validator to self-delegate in wei")
	cmd.Flags().Var(&cfg.StakePeriod, "staking-period", `Staking period (options: "flexible", "short", "medium", "long")`)
}

func bindValidatorStakeOnBehalfFlags(cmd *cobra.Command, cfg *stakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.DelegatorPubKey, "delegator-pubkey", "", "Delegator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount for the validator to self-delegate in wei")
	cmd.Flags().Var(&cfg.StakePeriod, "staking-period", `Staking period (options: "flexible", "short", "medium", "long")`)
}

func bindValidatorUnstakeFlags(cmd *cobra.Command, cfg *unstakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "unstake", "", "Amount to unstake in wei")
	cmd.Flags().Uint32Var(&cfg.DelegationID, "delegation-id", 0, "The delegation ID (0 for flexible staking)")
}

func bindValidatorUnstakeOnBehalfFlags(cmd *cobra.Command, cfg *unstakeConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.DelegatorPubKey, "delegator-pubkey", "", "Delegator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "unstake", "", "Amount to unstake in wei")
	cmd.Flags().Uint32Var(&cfg.DelegationID, "delegation-id", 0, "The delegation ID (0 for flexible staking)")
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

func bindKeyConvertFlags(cmd *cobra.Command, cfg *keyConfig) {
	cmd.Flags().StringVar(&cfg.ValidatorKeyFile, "validator-key-file", "", "Path to the validator key file")
	cmd.Flags().StringVar(&cfg.PrivateKeyFile, "private-key-file", "", "Path to the EVM private key env file")
	cmd.Flags().StringVar(&cfg.PubKeyHex, "pubkey-hex", "", "Public key in hex format")
	cmd.Flags().StringVar(&cfg.PubKeyBase64, "pubkey-base64", "", "Public key in base64 format")
	cmd.Flags().StringVar(&cfg.PubKeyHexUncompressed, "pubkey-hex-uncompressed", "", "Uncompressed public key in hex format")
}

func bindRollbackFlags(cmd *cobra.Command, cfg *config.RollbackConfig) {
	cmd.Flags().Uint64VarP(&cfg.RollbackHeights, "number", "n", 1, "number of blocks to rollback")
}

func bindValidatorUnjailFlags(cmd *cobra.Command, cfg *unjailConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
}

func bindValidatorUnjailOnBehalfFlags(cmd *cobra.Command, cfg *unjailConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's hex-encoded compressed 33-byte secp256k1 public key")
}

// Flag Validation

func validateFlags(cmd *cobra.Command, flags []string) error {
	var missingFlags []string

	for _, flag := range flags {
		if !cmd.Flags().Changed(flag) {
			missingFlags = append(missingFlags, flag)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func validateValidatorCreateFlags(ctx context.Context, cmd *cobra.Command, cfg *createValidatorConfig) error {
	if err := validateFlags(cmd, []string{"moniker"}); err != nil {
		return errors.Wrap(err, "failed to validate create flags")
	}

	if err := validateMinStakeAmount(ctx, &cfg.stakeConfig); err != nil {
		return err
	}

	return validateCommissionRate(ctx, cfg)
}

func validateOperatorFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{})
}

func validateWithdrawalFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{
		"withdrawal-address",
	})
}

func validateValidatorStakeFlags(ctx context.Context, cmd *cobra.Command, cfg *stakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "stake"}); err != nil {
		return errors.Wrap(err, "failed to validate stake flags")
	}

	return validateMinStakeAmount(ctx, cfg)
}

func validateValidatorStakeOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *stakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "delegator-pubkey", "stake"}); err != nil {
		return errors.Wrap(err, "failed to validate stake-on-behalf flags")
	}

	return validateMinStakeAmount(ctx, cfg)
}

func validateValidatorUnstakeFlags(ctx context.Context, cmd *cobra.Command, cfg *unstakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "unstake"}); err != nil {
		return errors.Wrap(err, "failed to validate unstake flags")
	}

	return validateMinUnstakeAmount(ctx, cfg)
}

func validateValidatorUnstakeOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *unstakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "delegator-pubkey", "unstake"}); err != nil {
		return errors.Wrap(err, "failed to validate unstake-on-behalf flags")
	}

	return validateMinUnstakeAmount(ctx, cfg)
}

func validateKeyConvertFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{})
}

func validateValidatorUnjailFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{})
}

func validateValidatorUnjailOnBehalfFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{"validator-pubkey"})
}

func validateMinStakeAmount(ctx context.Context, cfg *stakeConfig) error {
	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return fmt.Errorf("invalid stake amount: %s", cfg.StakeAmount)
	}

	minStakeAmount, err := getUint256(ctx, &cfg.baseConfig, "minStakeAmount")
	if err != nil {
		return errors.Wrap(err, "failed to retrieve minimum stake amount")
	}

	if stakeAmount.Cmp(minStakeAmount) < 0 {
		return fmt.Errorf("stake amount is less than the minimum required: %s", minStakeAmount.String())
	}

	return nil
}

func validateMinUnstakeAmount(ctx context.Context, cfg *unstakeConfig) error {
	unstakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return fmt.Errorf("invalid unstake amount: %s", cfg.StakeAmount)
	}

	minUnstakeAmount, err := getUint256(ctx, &cfg.baseConfig, "minUnstakeAmount")
	if err != nil {
		return errors.Wrap(err, "failed to retrieve minimum unstake amount")
	}

	if unstakeAmount.Cmp(minUnstakeAmount) < 0 {
		return fmt.Errorf("unstake amount is less than the minimum required: %s", minUnstakeAmount.String())
	}

	return nil
}

func validateCommissionRate(ctx context.Context, cfg *createValidatorConfig) error {
	commissionRate := new(big.Int).SetUint64(uint64(cfg.CommissionRate))

	minCommissionRate, err := getUint256(ctx, &cfg.baseConfig, "minCommissionRate")
	if err != nil {
		return errors.Wrap(err, "failed to retrieve minimum commission rate")
	}

	if commissionRate.Cmp(minCommissionRate) < 0 {
		return fmt.Errorf("commission rate is less than the minimum required: %s", minCommissionRate.String())
	}

	maxCommissionRate := new(big.Int).SetUint64(uint64(cfg.MaxCommissionRate))

	if commissionRate.Cmp(maxCommissionRate) > 0 {
		return fmt.Errorf("commission rate exceeds the maximum allowed: %s", maxCommissionRate.String())
	}

	return nil
}
