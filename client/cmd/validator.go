package cmd

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/storyprotocol/iliad/client/config"
	"github.com/storyprotocol/iliad/lib/errors"

	_ "embed"
)

//go:embed abi/IPTokenStaking.abi.json
var ipTokenStakingABI []byte

type ValidatorKey struct {
	Address string  `json:"address"`
	PubKey  KeyInfo `json:"pub_key"`
	PrivKey KeyInfo `json:"priv_key"`
}

type KeyInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type stakeConfig struct {
	RPC             string
	ValidatorPubKey string
	StakeAmount     string
	PrivateKey      string
	Explorer        string
	ChainID         int64
}

type unstakeConfig struct {
	RPC              string
	ValidatorPubKey  string
	UnstakeAmount    string
	ExecutionAddress string
	PrivateKey       string
	Explorer         string
	ChainID          int64
}

type createValidatorConfig struct {
	RPC              string
	ValidatorKeyFile string
	StakeAmount      string
	PrivateKey       string
	Explorer         string
	ChainID          int64
}

func addKeyFileFlag(cmd *cobra.Command, keyFilePath *string) {
	defaultKeyFilePath := filepath.Join(config.DefaultHomeDir(), "config", "priv_validator_key.json")
	cmd.Flags().StringVar(keyFilePath, "keyfile", defaultKeyFilePath, "Path to the Tendermint key file")
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

	addKeyFileFlag(cmd, &keyFilePath)

	return cmd
}

func createValidator(ctx context.Context, cfg createValidatorConfig) error {
	// Read the priv_validator_key.json file
	keyFileBytes, err := os.ReadFile(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "invalid key file")
	}

	var keyFileData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyFileData); err != nil {
		return errors.Wrap(err, "failed to unmarshal priv_validator_key.json")
	}

	// Decode and uncompress the public key
	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(keyFileData.PubKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 key")
	}
	if len(compressedPubKeyBytes) != 33 {
		return errors.New("invalid compressed public key length", "length", len(compressedPubKeyBytes))
	}

	// Load the private key for funding the EVM transaction from the .env file if not provided as a flag
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
		if cfg.PrivateKey == "" {
			return errors.New("missing required flag", "private-key", "EVM private key")
		}
	}

	evmPrivKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "invalid EVM private key")
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum client")
	}

	// Get the balance of the account
	publicKey := evmPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Fetch balance
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to fetch balance")
	}

	fmt.Printf("Balance: %s wei\n", balance.String())

	// Convert the stake amount to a big.Int
	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	// Suggest gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to suggest gas price")
	}

	// Increase gas price by 20% to ensure faster confirmation
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	// Read the ABI file from embedded content
	contractABI, err := abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI")
	}

	contractAddress := common.HexToAddress("0xCCcCcC0000000000000000000000000000000001")

	// Get the nonce for the transaction
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get nonce")
	}

	fmt.Printf("Using Nonce: %d\n", nonce)

	// Prepare transaction data
	data, err := contractABI.Pack(
		"createValidatorOnBehalf",
		compressedPubKeyBytes,
	)
	if err != nil {
		return errors.Wrap(err, "failed to pack data")
	}

	fmt.Printf("Packed Data: %x\n", data)

	chainID := big.NewInt(cfg.ChainID)

	// Estimate gas limit
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &contractAddress,
		GasPrice: gasPrice,
		Value:    stakeAmount,
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to estimate gas")
	}

	// Define gas fee cap and gas tip cap dynamically
	gasTipCap := gasPrice
	gasFeeCap := new(big.Int).Mul(gasPrice, big.NewInt(2))

	gasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasFeeCap)
	totalTxCost := new(big.Int).Add(gasCost, stakeAmount)

	fmt.Printf("Stake Amount: %s wei\n", stakeAmount.String())
	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Price: %s wei\n", gasPrice.String())
	fmt.Printf("Gas Tip Cap: %s wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap: %s wei\n", gasFeeCap.String())
	fmt.Printf("Gas Cost: %s wei\n", gasCost.String())
	fmt.Printf("Total Transaction Cost: %s wei\n", totalTxCost.String())

	if balance.Cmp(totalTxCost) < 0 {
		return errors.New("insufficient funds for gas * price + value", "balance", balance.String(), "totalTxCost", totalTxCost.String())
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &contractAddress,
		Value:     stakeAmount,
		Data:      data,
	})

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), evmPrivKey)
	if err != nil {
		return errors.Wrap(err, "failed to sign transaction")
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)
	fmt.Printf("Explorer URL: %s/tx/%s\n", cfg.Explorer, txHash)

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return errors.Wrap(err, "failed to send transaction")
	}

	fmt.Println("Transaction sent, waiting for confirmation...")

	// Use bind.WaitMined to wait for the transaction receipt
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return errors.New("transaction failed", "status", receipt.Status)
	}

	fmt.Println("Transaction confirmed successfully!")
	fmt.Println("Validator created successfully!")

	return nil
}

func validatorKeyExport(keyFilePath string) error {
	// Read the key file
	keyFileBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read key file")
	}

	// Unmarshal the key file
	var keyData ValidatorKey
	if err := json.Unmarshal(keyFileBytes, &keyData); err != nil {
		return errors.Wrap(err, "failed to unmarshal key file")
	}

	// Decode the base64 encoded private key
	privKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PrivKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	// Convert to EVM private key
	privateKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "invalid private key")
	}

	// Derive EVM public key
	publicKeyInterface := privateKey.Public()
	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ecdsa.PublicKey")
	}
	evmPublicKey := crypto.PubkeyToAddress(*publicKey).Hex()

	// Print the EVM private key and the uncompressed public key
	evmPrivateKey := hex.EncodeToString(crypto.FromECDSA(privateKey))

	// Decode the base64 encoded compressed public key
	compressedPubKeyBytes, err := base64.StdEncoding.DecodeString(keyData.PubKey.Value)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64 public key")
	}
	if len(compressedPubKeyBytes) != 33 {
		return fmt.Errorf("invalid compressed public key length: %d", len(compressedPubKeyBytes))
	}

	curve := elliptic.P256()
	x, y := elliptic.UnmarshalCompressed(curve, compressedPubKeyBytes)
	if x == nil || y == nil {
		return errors.New("failed to unmarshal compressed public key")
	}

	// lint:ignore SA1019 ignoring deprecation warning for now
	uncompressedPubKeyBytes := elliptic.Marshal(curve, x, y)
	uncompressedPubKeyHex := hex.EncodeToString(uncompressedPubKeyBytes)
	compressedPubKeyHex := hex.EncodeToString(compressedPubKeyBytes)

	fmt.Println("------------------------------------------------------")
	fmt.Println("EVM Public Key")
	fmt.Println("------------------------------------------------------")
	fmt.Println(evmPublicKey)
	fmt.Println("------------------------------------------------------")
	fmt.Println("EVM Private Key:")
	fmt.Println("------------------------------------------------------")
	fmt.Println(evmPrivateKey)
	fmt.Println("------------------------------------------------------")
	fmt.Println("Compressed Public Key:")
	fmt.Println("------------------------------------------------------")
	fmt.Println(compressedPubKeyHex)
	fmt.Println("------------------------------------------------------")
	fmt.Println("Uncompressed Public Key:")
	fmt.Println("------------------------------------------------------")
	fmt.Println(uncompressedPubKeyHex)
	fmt.Println("------------------------------------------------------")

	return nil
}

func newValidatorCreateCmd() *cobra.Command {
	var cfg createValidatorConfig

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new validator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			loadEnv()
			if err := validateFlags(cfg); err != nil {
				fmt.Println("Debug: Entering cmd.Help()")
				_ = cmd.Help() // Print the help message

				return err
			}

			return createValidator(cmd.Context(), cfg)
		},
	}

	bindCreateValidatorConfig(cmd, &cfg)

	return cmd
}

func validateFlags(cfg createValidatorConfig) error {
	var missingFlags []string

	if cfg.RPC == "" {
		missingFlags = append(missingFlags, "rpc")
	}
	if cfg.ValidatorKeyFile == "" {
		missingFlags = append(missingFlags, "keyfile")
	}
	if cfg.StakeAmount == "" {
		missingFlags = append(missingFlags, "stake")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func bindCreateValidatorConfig(cmd *cobra.Command, cfg *createValidatorConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://rpc.partner.testnet.storyprotocol.net", "RPC URL to connect to the testnet")
	addKeyFileFlag(cmd, &cfg.ValidatorKeyFile)
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount for the validator to self-delegate in wei")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used to issue the validator creation transaction")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://explorer.testnet.storyprotocol.net", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1513, "Chain ID to use for the transaction (default 1513)")
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

func validateStakeFlags(cfg stakeConfig) error {
	var missingFlags []string

	if cfg.RPC == "" {
		missingFlags = append(missingFlags, "rpc")
	}
	if cfg.ValidatorPubKey == "" {
		missingFlags = append(missingFlags, "validator-pubkey")
	}
	if cfg.StakeAmount == "" {
		missingFlags = append(missingFlags, "stake")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func bindStakeConfig(cmd *cobra.Command, cfg *stakeConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://rpc.partner.testnet.storyprotocol.net", "RPC URL to connect to the testnet")
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's 33 bytes compressed secp256k1 public key")
	cmd.Flags().StringVar(&cfg.StakeAmount, "stake", "", "Amount to stake on behalf of the delegator in wei")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used to derive the delegator's public key")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://explorer.testnet.storyprotocol.net", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1513, "Chain ID to use for the transaction (default 1513)")
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

func validateUnstakeFlags(cfg unstakeConfig) error {
	var missingFlags []string

	if cfg.RPC == "" {
		missingFlags = append(missingFlags, "rpc")
	}
	if cfg.ValidatorPubKey == "" {
		missingFlags = append(missingFlags, "validator-pubkey")
	}
	if cfg.UnstakeAmount == "" {
		missingFlags = append(missingFlags, "unstake")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flag(s): %s", strings.Join(missingFlags, ", "))
	}

	return nil
}

func bindUnstakeConfig(cmd *cobra.Command, cfg *unstakeConfig) {
	cmd.Flags().StringVar(&cfg.RPC, "rpc", "https://rpc.partner.testnet.storyprotocol.net", "RPC URL to connect to the testnet")
	cmd.Flags().StringVar(&cfg.ValidatorPubKey, "validator-pubkey", "", "Validator's 33 bytes compressed secp256k1 public key")
	cmd.Flags().StringVar(&cfg.UnstakeAmount, "unstake", "", "Amount to unstake on behalf of the delegator in wei")
	cmd.Flags().StringVar(&cfg.PrivateKey, "private-key", "", "Private key used to derive the delegator's public key")
	cmd.Flags().StringVar(&cfg.Explorer, "explorer", "https://explorer.testnet.storyprotocol.net", "URL of the blockchain explorer")
	cmd.Flags().Int64Var(&cfg.ChainID, "chain-id", 1513, "Chain ID to use for the transaction (default 1513)")
}

func stakeTokens(ctx context.Context, cfg stakeConfig) error {
	// Decode the validator public key
	validatorPubKeyBytes := common.Hex2Bytes(cfg.ValidatorPubKey)
	if len(validatorPubKeyBytes) != 33 {
		return fmt.Errorf("invalid validator public key length: %d", len(validatorPubKeyBytes))
	}

	// Load the private key from the .env file if not provided as a flag
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
		if cfg.PrivateKey == "" {
			return errors.New("missing required flag", "private-key", "EVM private key")
		}
	}

	evmPrivKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return errors.New("invalid EVM private key", err)
	}

	// Derive the delegator's uncompressed public key
	pubKey := evmPrivKey.PublicKey
	// lint:ignore SA1019 ignoring deprecation warning for now
	uncompressedPubKey := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if len(uncompressedPubKey) != 65 {
		return fmt.Errorf("invalid uncompressed public key length: %d", len(uncompressedPubKey))
	}
	fmt.Printf("Uncompressed Delegator PubKey: %x\n", uncompressedPubKey)

	// Connect to the Ethereum client
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum client")
	}

	// Get the balance of the account
	publicKey := evmPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Fetch balance
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to fetch balance")
	}

	fmt.Printf("Balance: %s wei\n", balance.String())

	// Convert the stake amount to a big.Int
	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	// Suggest gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to suggest gas price")
	}

	// Increase gas price by 20% to ensure faster confirmation
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	// Read the ABI file from embedded content
	contractABI, err := abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI")
	}

	contractAddress := common.HexToAddress("0xCCcCcC0000000000000000000000000000000001")

	// Get the nonce for the transaction
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get nonce")
	}

	fmt.Printf("Using Nonce: %d\n", nonce)

	// Prepare transaction data
	data, err := contractABI.Pack(
		"stake",
		uncompressedPubKey,
		validatorPubKeyBytes,
	)
	if err != nil {
		return errors.Wrap(err, "failed to pack data")
	}

	fmt.Printf("Packed Data: %x\n", data)

	chainID := big.NewInt(cfg.ChainID)

	// Estimate gas limit
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &contractAddress,
		GasPrice: gasPrice,
		Value:    stakeAmount,
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to estimate gas")
	}

	// Define gas fee cap and gas tip cap dynamically
	gasTipCap := gasPrice
	gasFeeCap := new(big.Int).Mul(gasPrice, big.NewInt(2))

	gasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasFeeCap)
	totalTxCost := new(big.Int).Add(gasCost, stakeAmount)

	fmt.Printf("Stake Amount: %s wei\n", stakeAmount.String())
	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Price: %s wei\n", gasPrice.String())
	fmt.Printf("Gas Tip Cap: %s wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap: %s wei\n", gasFeeCap.String())
	fmt.Printf("Gas Cost: %s wei\n", gasCost.String())
	fmt.Printf("Total Transaction Cost: %s wei\n", totalTxCost.String())

	if balance.Cmp(totalTxCost) < 0 {
		return errors.New("insufficient funds for gas * price + value", "balance", balance.String(), "totalTxCost", totalTxCost.String())
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &contractAddress,
		Value:     stakeAmount,
		Data:      data,
	})

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), evmPrivKey)
	if err != nil {
		return errors.Wrap(err, "failed to sign transaction")
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)
	fmt.Printf("Explorer URL: %s/tx/%s\n", cfg.Explorer, txHash)

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return errors.Wrap(err, "failed to send transaction")
	}

	fmt.Println("Transaction sent, waiting for confirmation...")

	// Use bind.WaitMined to wait for the transaction receipt
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return errors.New("transaction failed", "status", receipt.Status)
	}

	fmt.Println("Transaction confirmed successfully!")
	fmt.Println("Tokens staked successfully!")

	return nil
}

func unstakeTokens(ctx context.Context, cfg unstakeConfig) error {
	// Decode the validator public key
	validatorPubKeyBytes := common.Hex2Bytes(cfg.ValidatorPubKey)
	if len(validatorPubKeyBytes) != 33 {
		return fmt.Errorf("invalid validator public key length: %d", len(validatorPubKeyBytes))
	}

	// Load the private key from the .env file if not provided as a flag
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
		if cfg.PrivateKey == "" {
			return errors.New("missing required flag", "private-key", "EVM private key")
		}
	}

	evmPrivKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "invalid EVM private key")
	}

	// Derive the delegator's uncompressed public key
	pubKey := evmPrivKey.PublicKey
	// lint:ignore SA1019 ignoring deprecation warning for now
	uncompressedPubKey := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if len(uncompressedPubKey) != 65 {
		return fmt.Errorf("invalid uncompressed public key length: %d", len(uncompressedPubKey))
	}
	fmt.Printf("Uncompressed Delegator PubKey: %x\n", uncompressedPubKey)

	// Convert the unstake amount to a big.Int
	unstakeAmount, ok := new(big.Int).SetString(cfg.UnstakeAmount, 10)
	if !ok {
		return errors.New("invalid unstake amount", "amount", cfg.UnstakeAmount)
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum client")
	}

	// Get the balance of the account
	publicKey := evmPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return errors.Wrap(err, "failed to fetch balance")
	}

	fmt.Printf("Balance: %s wei\n", balance.String())

	// Suggest gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to suggest gas price")
	}

	// Increase gas price by 20% to ensure faster confirmation
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	// Read the ABI file from embedded content
	contractABI, err := abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI")
	}

	contractAddress := common.HexToAddress("0xCCcCcC0000000000000000000000000000000001")

	// Get the nonce for the transaction
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get nonce")
	}

	fmt.Printf("Using Nonce: %d\n", nonce)

	// Prepare transaction data
	data, err := contractABI.Pack(
		"unstake",
		uncompressedPubKey,
		validatorPubKeyBytes,
		unstakeAmount,
	)
	if err != nil {
		return errors.Wrap(err, "failed to pack data")
	}

	fmt.Printf("Packed Data: %x\n", data)

	chainID := big.NewInt(cfg.ChainID)

	// Estimate gas limit
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &contractAddress,
		GasPrice: gasPrice,
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to estimate gas")
	}

	// Define gas fee cap and gas tip cap dynamically
	gasTipCap := gasPrice
	gasFeeCap := new(big.Int).Mul(gasPrice, big.NewInt(2))

	gasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasFeeCap)
	totalTxCost := gasCost // Only gas cost is considered

	fmt.Printf("Unstake Amount: %s wei\n", unstakeAmount.String())
	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Price: %s wei\n", gasPrice.String())
	fmt.Printf("Gas Tip Cap: %s wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap: %s wei\n", gasFeeCap.String())
	fmt.Printf("Gas Cost: %s wei\n", gasCost.String())
	fmt.Printf("Total Transaction Cost: %s wei\n", totalTxCost.String())

	if balance.Cmp(totalTxCost) < 0 {
		return errors.New("insufficient funds for gas * price + value", "balance", balance.String(), "totalTxCost", totalTxCost.String())
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &contractAddress,
		Value:     big.NewInt(0), // No value for unstake
		Data:      data,
	})

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), evmPrivKey)
	if err != nil {
		return errors.Wrap(err, "failed to sign transaction")
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction hash: %s\n", txHash)
	fmt.Printf("Explorer URL: %s/tx/%s\n", cfg.Explorer, txHash)

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return errors.Wrap(err, "failed to send transaction")
	}

	fmt.Println("Transaction sent, waiting for confirmation...")

	// Use bind.WaitMined to wait for the transaction receipt
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return errors.New("transaction failed", "status", receipt.Status)
	}

	fmt.Println("Transaction confirmed successfully!")
	fmt.Println("Tokens unstaked successfully!")

	return nil
}
