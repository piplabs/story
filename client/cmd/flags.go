package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cosmossdk.io/math"

	cmtos "github.com/cometbft/cometbft/libs/os"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/piplabs/story/client/config"
	apisvr "github.com/piplabs/story/client/server"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tracer"
)

func bindRunFlags(cmd *cobra.Command, cfg *config.Config) {
	flags := cmd.Flags()

	libcmd.BindHomeFlag(flags, &cfg.HomeDir)
	tracer.BindFlags(flags, &cfg.Tracer)
	netconf.BindFlag(flags, &cfg.Network)
	apisvr.BindFlags(flags, &cfg.API)
	flags.StringVar(&cfg.EngineEndpoint, "engine-endpoint", cfg.EngineEndpoint, "An EVM execution client Engine API http endpoint")
	flags.StringVar(&cfg.EngineJWTFile, "engine-jwt-file", cfg.EngineJWTFile, "The path to the Engine API JWT file")
	flags.Uint64Var(&cfg.SnapshotInterval, "state-sync.snapshot-interval", cfg.SnapshotInterval, "State sync snapshot interval")
	flags.Uint64Var(&cfg.SnapshotKeepRecent, "state-sync.snapshot-keep-recent", cfg.SnapshotKeepRecent, "State sync snapshot to keep")
	flags.Uint64Var(&cfg.MinRetainBlocks, "min-retain-blocks", cfg.MinRetainBlocks, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	flags.StringVar(&cfg.BackendType, "app-db-backend", cfg.BackendType, "The type of database for application and snapshots databases")
	flags.StringVar(&cfg.PruningOption, "pruning", cfg.PruningOption, "Pruning strategy (default|nothing|everything)")
	flags.DurationVar(&cfg.EVMBuildDelay, "evm-build-delay", cfg.EVMBuildDelay, "Minimum delay between triggering and fetching a EVM payload build")
	flags.BoolVar(&cfg.EVMBuildOptimistic, "evm-build-optimistic", cfg.EVMBuildOptimistic, "Enables optimistic building of EVM payloads on previous block finalize")
	flags.BoolVar(&cfg.WithComet, "with-comet", true, "Run abci app embedded in-process with CometBFT")
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
	flags.BoolVar(&cfg.EncryptPrivKey, "encrypt-priv-key", false, "Encrypt the validator's private key")
}

func bindValidatorBaseFlags(cmd *cobra.Command, cfg *baseConfig) {
	libcmd.BindHomeFlag(cmd.Flags(), &cfg.HomeDir)
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://mainnet.storyrpc.io", "RPC URL to connect to the network")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://storyscan.xyz", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1514, "Chain ID to use for the transaction")
	cmd.Flags().StringVar(&cfg.StoryAPI, "story-api", "", "URL of Story API server for some validations")
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
	cmd.Flags().StringVar(&cfg.WithdrawalAddress, "withdrawal-address", "", "Address to receive stake withdrawals")
}

func bindSetRewardsAddressFlags(cmd *cobra.Command, cfg *rewardsConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.RewardsAddress, "rewards-address", "", "Address to receive rewards")
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
	cmd.Flags().StringVar(&cfg.DelegatorAddress, "delegator-address", "", "Delegator's EVM address")
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
	cmd.Flags().StringVar(&cfg.DelegatorAddress, "delegator-address", "", "Delegator's EVM address")
	cmd.Flags().StringVar(&cfg.StakeAmount, "unstake", "", "Amount to unstake in wei")
	cmd.Flags().Uint32Var(&cfg.DelegationID, "delegation-id", 0, "The delegation ID (0 for flexible staking)")
}

func bindValidatorRedelegateFlags(cmd *cobra.Command, cfg *redelegateConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorSrcPubKey, "validator-src-pubkey", "", "Src validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.ValidatorDstPubKey, "validator-dst-pubkey", "", "Dst validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "redelegate", "", "Amount to redelegate in wei")
	cmd.Flags().Uint32Var(&cfg.DelegationID, "delegation-id", 0, "The delegation ID (0 for flexible staking)")
}

func bindValidatorRedelegateOnBehalfFlags(cmd *cobra.Command, cfg *redelegateConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.DelegatorAddress, "delegator-address", "", "Delegator's EVM address")
	cmd.Flags().StringVar(&cfg.ValidatorSrcPubKey, "validator-src-pubkey", "", "Src validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.ValidatorDstPubKey, "validator-dst-pubkey", "", "Dst validator's hex-encoded compressed 33-byte secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "redelegate", "", "Amount to redelegate in wei")
	cmd.Flags().Uint32Var(&cfg.DelegationID, "delegation-id", 0, "The delegation ID (0 for flexible staking)")
}

func bindValidatorKeyExportFlags(cmd *cobra.Command, cfg *exportKeyConfig) {
	bindValidatorKeyFlags(cmd, &cfg.ValidatorKeyFile)
	defaultEVMKeyFilePath := filepath.Join(config.DefaultHomeDir(), "config", "private_key.txt")
	cmd.Flags().BoolVar(&cfg.ExportEVMKey, "export-evm-key", false, "Export the EVM private key")
	cmd.Flags().StringVar(&cfg.EvmKeyFile, "evm-key-path", defaultEVMKeyFilePath, "Path to save the exported EVM private key")
}

func bindKeyGenPrivKeyJSONFlags(cmd *cobra.Command, cfg *genPrivKeyJSONConfig) {
	bindValidatorKeyFlags(cmd, &cfg.ValidatorKeyFile)
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
}

func bindKeyShowEncryptedFlags(cmd *cobra.Command, cfg *showEncryptedConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().BoolVar(&cfg.ShowPrivate, "show-private", false, "Show private key")
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

func bindValidatorUpdateCommissionFlags(cmd *cobra.Command, cfg *updateCommissionConfig) {
	bindValidatorBaseFlags(cmd, &cfg.baseConfig)
	cmd.Flags().Uint32Var(&cfg.CommissionRate, "commission-rate", 0, "Commission rate to update (e.g. 1000 for 10%)")
}

var ErrValidatorNotFound = errors.New("the validator doesn't exist")

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

	if err := validateCommissionRate(ctx, cfg); err != nil {
		return err
	}

	validatorPubKey, err := validatorKeyFileToCmpPubKey(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to extract compressed pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("the validator already exist")
	}

	return nil
}

func validateOperatorFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{})
}

func validateWithdrawalFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{
		"withdrawal-address",
	})
}

func validateRewardsFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{
		"rewards-address",
	})
}

func validateValidatorStakeFlags(ctx context.Context, cmd *cobra.Command, cfg *stakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "stake"}); err != nil {
		return errors.Wrap(err, "failed to validate stake flags")
	}

	if err := validateMinStakeAmount(ctx, cfg); err != nil {
		return err
	}

	validatorPubKey, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateValidatorStakeOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *stakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "delegator-address", "stake"}); err != nil {
		return errors.Wrap(err, "failed to validate stake-on-behalf flags")
	}

	if err := validateMinStakeAmount(ctx, cfg); err != nil {
		return err
	}

	validatorPubKey, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateValidatorUnstakeFlags(ctx context.Context, cmd *cobra.Command, cfg *unstakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "unstake"}); err != nil {
		return errors.Wrap(err, "failed to validate unstake flags")
	}

	if err := validateMinUnstakeAmount(ctx, cfg); err != nil {
		return err
	}

	validatorPubKey, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateValidatorUnstakeOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *unstakeConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey", "delegator-address", "unstake"}); err != nil {
		return errors.Wrap(err, "failed to validate unstake-on-behalf flags")
	}

	if err := validateMinUnstakeAmount(ctx, cfg); err != nil {
		return err
	}

	validatorPubKey, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateValidatorRedelegateFlags(ctx context.Context, cmd *cobra.Command, cfg *redelegateConfig) error {
	if err := validateFlags(cmd, []string{"validator-src-pubkey", "validator-dst-pubkey", "redelegate"}); err != nil {
		return errors.Wrap(err, "failed to validate redelegate flags")
	}

	if err := validateMinRedelegateAmount(ctx, cfg); err != nil {
		return err
	}

	validatorSrcPubKey, err := hex.DecodeString(cfg.ValidatorSrcPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	validatorDstPubKey, err := hex.DecodeString(cfg.ValidatorDstPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	existSrc, err := isValidatorFound(ctx, cfg.StoryAPI, validatorSrcPubKey)
	if err != nil {
		return err
	}

	if !existSrc {
		return errors.New("the src validator doesn't exist")
	}

	existDst, err := isValidatorFound(ctx, cfg.StoryAPI, validatorDstPubKey)
	if err != nil {
		return err
	}

	if !existDst {
		return errors.New("the dst validator doesn't exist")
	}

	return nil
}

func validateValidatorRedelegateOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *redelegateConfig) error {
	if err := validateFlags(cmd, []string{"delegator-address", "validator-src-pubkey", "validator-dst-pubkey", "redelegate"}); err != nil {
		return errors.Wrap(err, "failed to validate redelegate-on-behalf flags")
	}

	if err := validateMinRedelegateAmount(ctx, cfg); err != nil {
		return err
	}

	validatorSrcPubKey, err := hex.DecodeString(cfg.ValidatorSrcPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	validatorDstPubKey, err := hex.DecodeString(cfg.ValidatorDstPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	existSrc, err := isValidatorFound(ctx, cfg.StoryAPI, validatorSrcPubKey)
	if err != nil {
		return err
	}

	if !existSrc {
		return errors.New("the src validator doesn't exist")
	}

	existDst, err := isValidatorFound(ctx, cfg.StoryAPI, validatorDstPubKey)
	if err != nil {
		return err
	}

	if !existDst {
		return errors.New("the dst validator doesn't exist")
	}

	return nil
}

func validateKeyConvertFlags(cmd *cobra.Command) error {
	return validateFlags(cmd, []string{})
}

func validateGenPrivKeyJSONFlags(cfg *genPrivKeyJSONConfig) error {
	// if there is an existing priv_validator_key.json file, do not overwrite it.
	if _, err := os.Stat(cfg.ValidatorKeyFile); err == nil {
		return errors.New("priv_validator_key.json file already exists")
	}

	return nil
}

func validateEncryptFlags(cfg *baseConfig) error {
	if cmtos.FileExists(cfg.EncPrivKeyFile()) {
		return errors.New("already encrypted private key exists")
	}

	loadEnv()
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		return errors.New("no private key is provided")
	}

	if _, err := crypto.HexToECDSA(pk); err != nil {
		return errors.New("invalid secp256k1 private key")
	}

	cfg.PrivateKey = pk

	return nil
}

func validateShowEncryptedFlags(cfg *showEncryptedConfig) error {
	if !cmtos.FileExists(cfg.EncPrivKeyFile()) {
		return errors.New("no encrypted private key file")
	}

	return nil
}

func validateValidatorUnjailFlags(ctx context.Context, cmd *cobra.Command, cfg *unjailConfig) error {
	if err := validateFlags(cmd, []string{}); err != nil {
		return err
	}

	privKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	validatorPubKey, err := privKeyToCmpPubKey(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to get compressed pub key from private key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateValidatorUnjailOnBehalfFlags(ctx context.Context, cmd *cobra.Command, cfg *unjailConfig) error {
	if err := validateFlags(cmd, []string{"validator-pubkey"}); err != nil {
		return err
	}

	validatorPubKey, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded validator public key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return nil
}

func validateUpdateValidatorCommissionFlags(ctx context.Context, cmd *cobra.Command, cfg *updateCommissionConfig) error {
	if err := validateFlags(cmd, []string{"commission-rate"}); err != nil {
		return err
	}

	privKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	validatorPubKey, err := privKeyToCmpPubKey(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to get compressed pub key from private key")
	}

	if cfg.StoryAPI == "" {
		fmt.Println("No staking API is provided. Skip validator existence check and validation of new commission rate.")
		return nil
	}

	exist, err := isValidatorFound(ctx, cfg.StoryAPI, validatorPubKey)
	if err != nil {
		return err
	}

	if !exist {
		return ErrValidatorNotFound
	}

	return validateNewCommissionRate(ctx, cfg, validatorPubKey)
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

func validateMinRedelegateAmount(ctx context.Context, cfg *redelegateConfig) error {
	redelegateAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return fmt.Errorf("invalid redelegate amount: %s", cfg.StakeAmount)
	}

	minRedelegateAmount, err := getUint256(ctx, &cfg.baseConfig, "minStakeAmount")
	if err != nil {
		return errors.Wrap(err, "failed to retrieve minimum redelegate amount")
	}

	if redelegateAmount.Cmp(minRedelegateAmount) < 0 {
		return fmt.Errorf("redelegate amount is less than the minimum required: %s", minRedelegateAmount.String())
	}

	return nil
}

func validateCommissionRate(ctx context.Context, cfg *createValidatorConfig) error {
	commissionRate := new(big.Int).SetUint64(uint64(cfg.CommissionRate))
	maxCommissionChangeRate := new(big.Int).SetUint64(uint64(cfg.MaxCommissionChangeRate))

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

	if maxCommissionChangeRate.Cmp(maxCommissionRate) > 0 {
		return fmt.Errorf("commission change rate cannot be more than the max rate: maxCommissionChangeRate: %s, maxCommissionRate: %s", maxCommissionChangeRate.String(), maxCommissionRate.String())
	}

	return nil
}

func validateNewCommissionRate(ctx context.Context, cfg *updateCommissionConfig, pubKey []byte) error {
	validator, err := getValidatorByEVMAddr(ctx, cfg.StoryAPI, pubKey)
	if err != nil {
		return err
	}

	if validator.OperatorAddress == "" {
		return errors.New("validator not found")
	}

	prevCommission := validator.Commission
	newCommission := math.LegacyNewDecWithPrec(int64(cfg.CommissionRate), 4)

	if time.Since(prevCommission.UpdateTime).Hours() < 24 {
		return stypes.ErrCommissionUpdateTime
	}

	if newCommission.GT(prevCommission.MaxRate) {
		return stypes.ErrCommissionGTMaxRate
	}

	if newCommission.Sub(prevCommission.Rate).GT(prevCommission.MaxChangeRate) {
		return stypes.ErrCommissionGTMaxChangeRate
	}

	return nil
}

type Validator struct {
	OperatorAddress string            `json:"operator_address"`
	Commission      stypes.Commission `json:"commission"`
}

type ValidatorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Msg   struct {
		Validator Validator `json:"validator"`
	} `json:"msg"`
}

// isValidatorFound checks whether a validator with the given public key exists.
func isValidatorFound(ctx context.Context, endpoint string, pubKey []byte) (bool, error) {
	validator, err := getValidatorByEVMAddr(ctx, endpoint, pubKey)
	if err != nil {
		return false, err
	}

	return validator.OperatorAddress != "", nil
}

// getValidatorByEVMAddr gets a validator with the given public key.
func getValidatorByEVMAddr(ctx context.Context, endpoint string, pubKey []byte) (Validator, error) {
	valEVMAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKey)
	if err != nil {
		return Validator{}, errors.Wrap(err, "failed to convert pub key to evm address")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/staking/validators/%s", endpoint, valEVMAddr), nil)
	if err != nil {
		return Validator{}, errors.Wrap(err, "failed to create request for getting validator")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Validator{}, errors.Wrap(err, "failed to get validator")
	}
	if resp.StatusCode != http.StatusOK {
		return Validator{}, errors.New("failed to get validator")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Validator{}, errors.Wrap(err, "failed to read response body")
	}

	var response ValidatorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return Validator{}, errors.Wrap(err, "failed to unmarshal response")
	}

	return response.Msg.Validator, nil
}
