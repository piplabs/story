package cmd

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/config"
	"github.com/piplabs/story/lib/errors"

	_ "embed"
)

const (
	contractAddressHex = "0xCCcCcC0000000000000000000000000000000001"
)

var (
	contractAddress common.Address
	contractABI     abi.ABI
)

//go:embed abi/IPTokenStaking.abi.json
var ipTokenStakingABI []byte

type baseConfig struct {
	RPC        string
	PrivateKey string
	Explorer   string
	ChainID    int64
}

type stakeConfig struct {
	baseConfig
	ValidatorPubKey string
	StakeAmount     string
}

type unstakeConfig struct {
	baseConfig
	RPC             string
	ValidatorPubKey string
	UnstakeAmount   string
}

type createValidatorConfig struct {
	baseConfig
	ValidatorKeyFile string
	StakeAmount      string
}

func init() {
	var err error
	contractAddress = common.HexToAddress(contractAddressHex)
	contractABI, err = abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found")
	}
}

func newValidatorCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Commands for validator operations",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newValidatorCreateCmd(),
		newValidatorKeyExportCmd(),
		newValidatorStakeCmd(),
		newValidatorUnstakeCmd(),
	)

	return cmd
}

func newValidatorCreateCmd() *cobra.Command {
	var cfg createValidatorConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new validator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			loadEnv()
			if err := validateCreateFlags(cfg); err != nil {
				fmt.Println("Debug: Entering cmd.Help()")
				_ = cmd.Help()

				return err
			}

			return createValidator(cmd.Context(), cfg)
		},
	}

	bindCreateValidatorConfig(cmd, &cfg)

	return cmd
}

func newValidatorStakeCmd() *cobra.Command {
	var cfg stakeConfig

	cmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake tokens on behalf of a delegator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			loadEnv()
			if err := validateStakeFlags(cfg); err != nil {
				fmt.Println("Debug: Entering cmd.Help()")
				_ = cmd.Help() // Print the help message

				return err
			}

			return stakeTokens(cmd.Context(), cfg)
		},
	}

	bindStakeConfig(cmd, &cfg)

	return cmd
}

func newValidatorUnstakeCmd() *cobra.Command {
	var cfg unstakeConfig

	cmd := &cobra.Command{
		Use:   "unstake",
		Short: "Unstake tokens on behalf of a delegator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			loadEnv()
			if err := validateUnstakeFlags(cfg); err != nil {
				fmt.Println("Debug: Entering cmd.Help()")
				_ = cmd.Help() // Print the help message

				return err
			}

			return unstakeTokens(cmd.Context(), cfg)
		},
	}

	bindUnstakeConfig(cmd, &cfg)

	return cmd
}

func newValidatorKeyExportCmd() *cobra.Command {
	var keyFilePath string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the EVM private key from the Tendermint key file",
		RunE: func(_ *cobra.Command, _ []string) error {
			loadEnv()
			return validatorKeyExport(keyFilePath)
		},
	}

	bindKeyConfig(cmd, &keyFilePath)

	return cmd
}

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

func validateCreateFlags(cfg createValidatorConfig) error {
	return validateFlags(map[string]string{
		"rpc":     cfg.RPC,
		"keyfile": cfg.ValidatorKeyFile,
		"stake":   cfg.StakeAmount,
	})
}

func validateStakeFlags(cfg stakeConfig) error {
	return validateFlags(map[string]string{
		"rpc":              cfg.RPC,
		"validator-pubkey": cfg.ValidatorPubKey,
		"stake":            cfg.StakeAmount,
	})
}

func validateUnstakeFlags(cfg unstakeConfig) error {
	return validateFlags(map[string]string{
		"rpc":              cfg.RPC,
		"validator-pubkey": cfg.ValidatorPubKey,
		"unstake":          cfg.UnstakeAmount,
	})
}

func bindBaseConfig(cmd *cobra.Command, cfg *baseConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://testnet.storyrpc.io", "RPC URL to connect to the testnet")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used for the transaction")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://testnet.storyscan.xyz", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1513, "Chain ID to use for the transaction (default 1513)")
}

func bindCreateValidatorConfig(cmd *cobra.Command, cfg *createValidatorConfig) {
	bindBaseConfig(cmd, &cfg.baseConfig)
	bindKeyConfig(cmd, &cfg.ValidatorKeyFile)
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount for the validator to self-delegate in wei")
}

func bindStakeConfig(cmd *cobra.Command, cfg *stakeConfig) {
	bindBaseConfig(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's 33 bytes compressed secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount to stake on behalf of the delegator in wei")
}

func bindUnstakeConfig(cmd *cobra.Command, cfg *unstakeConfig) {
	bindBaseConfig(cmd, &cfg.baseConfig)
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's 33 bytes compressed secp256k1 public key")
	cmd.Flags().StringVar(&cfg.UnstakeAmount, "unstake", "", "Amount to unstake on behalf of the delegator in wei")
}

func bindKeyConfig(cmd *cobra.Command, keyFilePath *string) {
	defaultKeyFilePath := filepath.Join(config.DefaultHomeDir(), "config", "priv_validator_key.json")
	cmd.Flags().StringVar(keyFilePath, "keyfile", defaultKeyFilePath, "Path to the Tendermint key file")
}

func createValidator(ctx context.Context, cfg createValidatorConfig) error {
	_, err := loadPrivateKey(&cfg.baseConfig)
	if err != nil {
		return err
	}

	keyFileBytes, err := os.ReadFile(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "invalid key file")
	}

	var keyFileData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyFileData); err != nil {
		return errors.Wrap(err, "failed to unmarshal priv_validator_key.json")
	}

	uncompressedPubKeyHex, err := decodeAndUncompressPubKey(keyFileData.PubKey.Value)
	if err != nil {
		return err
	}
	uncompressedPubKeyBytes, err := hex.DecodeString(uncompressedPubKeyHex)
	if err != nil {
		return errors.Wrap(err, "failed to decode uncompressed public key hex")
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "createValidatorOnBehalf", stakeAmount, uncompressedPubKeyBytes)
	if err != nil {
		return err
	}

	fmt.Println("Validator created successfully!")

	return nil
}

func stakeTokens(ctx context.Context, cfg stakeConfig) error {
	uncompressedPubKey, err := deriveUncompressedPublicKeyFromConfig(&cfg.baseConfig)
	if err != nil {
		return err
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "stake", stakeAmount, uncompressedPubKey, common.Hex2Bytes(cfg.ValidatorPubKey))
	if err != nil {
		return err
	}

	fmt.Println("Tokens staked successfully!")

	return nil
}

func unstakeTokens(ctx context.Context, cfg unstakeConfig) error {

	uncompressedPubKey, err := deriveUncompressedPublicKeyFromConfig(&cfg.baseConfig)
	if err != nil {
		return err
	}

	unstakeAmount, ok := new(big.Int).SetString(cfg.UnstakeAmount, 10)
	if !ok {
		return errors.New("invalid unstake amount", "amount", cfg.UnstakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unstake", big.NewInt(0), uncompressedPubKey, common.Hex2Bytes(cfg.ValidatorPubKey), unstakeAmount)
	if err != nil {
		return err
	}

	fmt.Println("Tokens unstaked successfully!")

	return nil
}

func prepareAndExecuteTransaction(ctx context.Context, cfg *baseConfig, methodName string, value *big.Int, args ...interface{}) error {
	data, err := contractABI.Pack(methodName, args...)
	if err != nil {
		return errors.Wrap(err, "failed to pack data")
	}

	return prepareAndSendTransaction(ctx, *cfg, contractAddress, value, data)
}

func deriveUncompressedPublicKeyFromConfig(cfg *baseConfig) ([]byte, error) {
	evmPrivKey, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, err
	}

	uncompressedPubKey, err := deriveUncompressedPublicKeyFromPrivateKey(evmPrivKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive uncompressed public key")
	}

	fmt.Printf("Uncompressed Delegator PubKey: %x\n", uncompressedPubKey)
	return uncompressedPubKey, nil
}

func loadPrivateKey(cfg *baseConfig) (*ecdsa.PrivateKey, error) {
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
		if cfg.PrivateKey == "" {
			return nil, errors.New("missing required flag", "private-key", "EVM private key")
		}
	}

	evmPrivKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "invalid EVM private key")
	}

	return evmPrivKey, nil
}
