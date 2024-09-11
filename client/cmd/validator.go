package cmd

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/lib/errors"

	_ "embed"
)

const (
	contractAddressHex = "0xCCcCcC0000000000000000000000000000000001"
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
	DelegatorPubKey string
	ValidatorPubKey string
	StakeAmount     string
}

type operatorConfig struct {
	baseConfig
	Operator string
}

type withdrawalConfig struct {
	baseConfig
	WithdrawalAddress string
}

type createValidatorConfig struct {
	baseConfig
	ValidatorKeyFile string
	StakeAmount      string
}

type exportKeyConfig struct {
	ValidatorKeyFile string
	EvmKeyFile       string
	ExportEVMKey     bool
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
		newValidatorStakeOnBehalfCmd(),
		newValidatorUnstakeCmd(),
		newValidatorUnstakeOnBehalfCmd(),
		newValidatorAddOperatorCmd(),
		newValidatorRemoveOperatorCmd(),
		newValidatorSetWithdrawalAddressCmd(),
		newValidatorStatusCmd(),
	)

	return cmd
}

func newValidatorCreateCmd() *cobra.Command {
	var cfg createValidatorConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new validator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateValidatorCreateFlags(cfg) },
			func(ctx context.Context) error { return createValidator(ctx, cfg) },
		),
	}

	bindValidatorCreateFlags(cmd, &cfg)

	return cmd
}

func newValidatorAddOperatorCmd() *cobra.Command {
	var cfg operatorConfig

	cmd := &cobra.Command{
		Use:   "add-operator",
		Short: "Add a new operator to your delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateOperatorFlags(cfg) },
			func(ctx context.Context) error { return addOperator(ctx, cfg) },
		),
	}

	bindAddOperatorFlags(cmd, &cfg)

	return cmd
}

func newValidatorRemoveOperatorCmd() *cobra.Command {
	var cfg operatorConfig

	cmd := &cobra.Command{
		Use:   "remove-operator",
		Short: "Removes an existing operator from your delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateOperatorFlags(cfg) },
			func(ctx context.Context) error { return removeOperator(ctx, cfg) },
		),
	}

	bindRemoveOperatorFlags(cmd, &cfg)

	return cmd
}

func newValidatorSetWithdrawalAddressCmd() *cobra.Command {
	var cfg withdrawalConfig

	cmd := &cobra.Command{
		Use:   "set-withdrawal-address",
		Short: "Updates the withdrawal address that receives stake and reward withdrawals",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateWithdrawalFlags(cfg) },
			func(ctx context.Context) error { return setWithdrawalAddress(ctx, cfg) },
		),
	}

	bindSetWithdrawalAddressFlags(cmd, &cfg)

	return cmd
}

func newValidatorStakeCmd() *cobra.Command {
	var cfg stakeConfig

	cmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake tokens as the delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error {
				return validateValidatorStakeFlags(cfg)
			},
			func(ctx context.Context) error { return stake(ctx, cfg) },
		),
	}

	bindValidatorStakeFlags(cmd, &cfg)

	return cmd
}

func newValidatorStakeOnBehalfCmd() *cobra.Command {
	var cfg stakeConfig

	cmd := &cobra.Command{
		Use:   "stake-on-behalf",
		Short: "Stake tokens on behalf of a delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error {
				return validateValidatorStakeOnBehalfFlags(cfg)
			},
			func(ctx context.Context) error { return stakeOnBehalf(ctx, cfg) },
		),
	}

	bindValidatorStakeOnBehalfFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnstakeCmd() *cobra.Command {
	var cfg stakeConfig

	cmd := &cobra.Command{
		Use:   "unstake",
		Short: "Unstake tokens as the delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateValidatorStakeFlags(cfg) },
			func(ctx context.Context) error { return unstake(ctx, cfg) },
		),
	}

	bindValidatorUnstakeFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnstakeOnBehalfCmd() *cobra.Command {
	var cfg stakeConfig

	cmd := &cobra.Command{
		Use:   "unstake-on-behalf",
		Short: "Unstake tokens on behalf of a delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return loadAndValidatePrivateKey(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func() error { return validateValidatorUnstakeOnBehalfFlags(cfg) },
			func(ctx context.Context) error { return unstakeOnBehalf(ctx, cfg) },
		),
	}

	bindValidatorUnstakeOnBehalfFlags(cmd, &cfg)

	return cmd
}

func newValidatorKeyExportCmd() *cobra.Command {
	var cfg exportKeyConfig

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the EVM private key from the Tendermint key file",
		RunE: runValidatorCommand(
			func() error { return nil },
			func(ctx context.Context) error { return exportKey(ctx, cfg) },
		),
	}

	bindValidatorKeyExportFlags(cmd, &cfg)

	return cmd
}

func newValidatorStatusCmd() *cobra.Command {
	var cfg baseConfig

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Fetch status of the Story chain",
		Args:  cobra.NoArgs,
		RunE: runValidatorCommand(
			func() error { return validateValidatorStatusFlags(cfg) },
			func(ctx context.Context) error { return checkStatus(ctx, cfg) },
		),
	}

	bindValidatorBaseFlags(cmd, &cfg)

	return cmd
}

func runValidatorCommand(
	validate func() error,
	execute func(ctx context.Context) error,
) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		if err := validate(); err != nil {
			_ = cmd.Help()
			return err
		}

		return execute(cmd.Context())
	}
}

func exportKey(_ context.Context, cfg exportKeyConfig) error {
	keyFileBytes, err := os.ReadFile(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to read key file")
	}

	var keyData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyData); err != nil {
		return errors.Wrap(err, "failed to unmarshal key file")
	}

	privKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PrivKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	privateKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "invalid private key")
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ecdsa.PublicKey")
	}
	evmPublicKey := crypto.PubkeyToAddress(*publicKey).Hex()

	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PubKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 pub key")
	}
	compressedPubKeyHex := hex.EncodeToString(compressedPubKeyBytes)

	uncompressedPubKeyHex, err := uncompressPubKey(keyData.PubKey.Value)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println("EVM Public Key:", evmPublicKey)
	fmt.Println("Compressed Public Key (base64):", keyData.PubKey.Value)
	fmt.Println("Compressed Public Key (hex):", compressedPubKeyHex)
	fmt.Println("Uncompressed Public Key:", uncompressedPubKeyHex)
	fmt.Println("------------------------------------------------------")

	if cfg.ExportEVMKey {
		evmPrivateKey := hex.EncodeToString(crypto.FromECDSA(privateKey))
		keyContent := "PRIVATE_KEY=" + evmPrivateKey
		if err := os.WriteFile(cfg.EvmKeyFile, []byte(keyContent), 0600); err != nil {
			return errors.Wrap(err, "failed to export private key")
		}

		fmt.Printf("EVM Private Key saved to: %s\n", cfg.EvmKeyFile)
		fmt.Println("WARNING: The EVM private key is highly sensitive. Store this file in a secure location.")
	}

	return nil
}

func createValidator(ctx context.Context, cfg createValidatorConfig) error {
	keyFileBytes, err := os.ReadFile(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "invalid key file")
	}

	var keyFileData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyFileData); err != nil {
		return errors.Wrap(err, "failed to unmarshal priv_validator_key.json")
	}

	uncompressedPubKeyHex, err := uncompressPubKey(keyFileData.PubKey.Value)
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

func setWithdrawalAddress(ctx context.Context, cfg withdrawalConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	withdrawalAddress := common.HexToAddress(cfg.WithdrawalAddress)

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "setWithdrawalAddress", big.NewInt(0), uncompressedPubKey, withdrawalAddress)
	if err != nil {
		return err
	}

	fmt.Println("Withdrawal address successfully set!")

	return nil
}

func addOperator(ctx context.Context, cfg operatorConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	operatorAddress := common.HexToAddress(cfg.Operator)

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "addOperator", big.NewInt(0), uncompressedPubKey, operatorAddress)
	if err != nil {
		return err
	}

	fmt.Println("Operator added successfully!")

	return nil
}

func removeOperator(ctx context.Context, cfg operatorConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	operatorAddress := common.HexToAddress(cfg.Operator)

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "removeOperator", big.NewInt(0), uncompressedPubKey, operatorAddress)
	if err != nil {
		return err
	}

	fmt.Println("Operator removed successfully!")

	return nil
}

func stake(ctx context.Context, cfg stakeConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	validatorPubKeyBytes, err := base64.StdEncoding.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 pub key")
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "stakeOnBehalf", stakeAmount, uncompressedPubKey, validatorPubKeyBytes)
	if err != nil {
		return err
	}

	fmt.Println("Tokens staked successfully!")

	return nil
}

func stakeOnBehalf(ctx context.Context, cfg stakeConfig) error {
	uncompressedDelegatorPubKeyHex, err := uncompressPubKey(cfg.DelegatorPubKey)
	if err != nil {
		return err
	}
	uncompressedDelegatorPubKeyBytes, err := hex.DecodeString(uncompressedDelegatorPubKeyHex)
	if err != nil {
		return errors.Wrap(err, "failed to decode uncompressed delegator public key")
	}

	validatorPubKeyBytes, err := base64.StdEncoding.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode validator public key")
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "stakeOnBehalf", stakeAmount, uncompressedDelegatorPubKeyBytes, validatorPubKeyBytes)
	if err != nil {
		return err
	}

	fmt.Println("Tokens staked on behalf of delegator successfully!")

	return nil
}

func unstake(ctx context.Context, cfg stakeConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	validatorPubKeyBytes, err := base64.StdEncoding.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 pub key")
	}

	unstakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid unstake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unstake", big.NewInt(0), uncompressedPubKey, validatorPubKeyBytes, unstakeAmount)
	if err != nil {
		return err
	}

	fmt.Println("Tokens unstaked successfully!")

	return nil
}

func unstakeOnBehalf(ctx context.Context, cfg stakeConfig) error {
	delegatorPubKeyBytes, err := base64.StdEncoding.DecodeString(cfg.DelegatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 delegator pub key")
	}

	validatorPubKeyBytes, err := base64.StdEncoding.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 validator pub key")
	}

	unstakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid unstake amount", "amount", cfg.StakeAmount)
	}

	err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unstakeOnBehalf", big.NewInt(0), delegatorPubKeyBytes, validatorPubKeyBytes, unstakeAmount)
	if err != nil {
		return err
	}

	fmt.Println("Tokens unstaked on behalf of delegator successfully!")

	return nil
}

func checkStatus(_ context.Context, cfg baseConfig) error {
	resp, err := http.Get(fmt.Sprintf("%s/status", cfg.RPC))
	if err != nil {
		return errors.Wrap(err, "failed to query cometBFT status endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected response code: " + strconv.Itoa(resp.StatusCode))
	}

	type statusResponse struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	var statusResp statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return errors.Wrap(err, "failed to decode JSON response")
	}

	blockHeight, err := strconv.ParseInt(statusResp.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid block height")
	}

	responseJSON := map[string]interface{}{
		"sync_info": map[string]any{
			"latest_block_height": blockHeight,
		},
	}

	output, err := json.Marshal(responseJSON)
	if err != nil {
		return errors.Wrap(err, "failed to marshal output JSON")
	}

	fmt.Println(string(output))

	return nil
}

func prepareAndExecuteTransaction(ctx context.Context, cfg *baseConfig, methodName string, value *big.Int, args ...any) error {
	contractAddress := common.HexToAddress(contractAddressHex)
	contractABI, err := abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI")
	}
	data, err := contractABI.Pack(methodName, args...)
	if err != nil {
		return errors.Wrap(err, "failed to pack data")
	}

	return prepareAndSendTransaction(ctx, *cfg, contractAddress, value, data)
}

func loadAndValidatePrivateKey(cfg *baseConfig) error {
	if cfg.PrivateKey == "" {
		loadEnv()
		cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
		if cfg.PrivateKey == "" {
			return errors.New("missing required flag", "private-key", "EVM private key")
		}
	}

	_, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "invalid EVM private key")
	}

	return nil
}
