package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/lib/errors"

	_ "embed"
)

type StakingPeriod int

const (
	FLEXIBLE StakingPeriod = iota
	SHORT
	MEDIUM
	LONG
)

func (sp *StakingPeriod) String() string {
	switch *sp {
	case FLEXIBLE:
		return "flexible"
	case SHORT:
		return "short"
	case MEDIUM:
		return "medium"
	case LONG:
		return "long"
	default:
		return "unknown"
	}
}

func (sp *StakingPeriod) Set(v string) error {
	switch strings.ToLower(v) {
	case "flexible":
		*sp = FLEXIBLE
	case "short":
		*sp = SHORT
	case "medium":
		*sp = MEDIUM
	case "long":
		*sp = LONG
	default:
		return errors.New(`staking period must be one of "flexible", "short", "medium", or "long"`)
	}

	return nil
}

func (*StakingPeriod) Type() string {
	return "stakingPeriod"
}

//go:embed abi/IPTokenStaking.abi.json
var ipTokenStakingABI []byte

type baseConfig struct {
	RPC          string
	PrivateKey   string
	Explorer     string
	ChainID      int64
	ABI          *abi.ABI
	ContractAddr common.Address
}

type createValidatorConfig struct {
	stakeConfig
	ValidatorKeyFile        string
	Moniker                 string
	CommissionRate          uint32
	MaxCommissionRate       uint32
	MaxCommissionChangeRate uint32
	Unlocked                bool
}

type stakeConfig struct {
	baseConfig
	DelegatorPubKey string
	ValidatorPubKey string
	StakeAmount     string
	StakePeriod     StakingPeriod
}

type unstakeConfig struct {
	stakeConfig
	DelegationID uint32
}

type unjailConfig struct {
	baseConfig
	ValidatorPubKey string
}

type operatorConfig struct {
	baseConfig
	Operator string
}

type withdrawalConfig struct {
	baseConfig
	WithdrawalAddress string
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
		newValidatorSetOperatorCmd(),
		newValidatorUnsetOperatorCmd(),
		newValidatorSetWithdrawalAddressCmd(),
		newValidatorUnjailCmd(),
		newValidatorUnjailOnBehalfCmd(),
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
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				ctx := cmd.Context()
				return validateValidatorCreateFlags(ctx, cmd, &cfg)
			},
			func(ctx context.Context) error { return createValidator(ctx, cfg) },
		),
	}

	bindValidatorCreateFlags(cmd, &cfg)

	return cmd
}

func newValidatorSetOperatorCmd() *cobra.Command {
	var cfg operatorConfig

	cmd := &cobra.Command{
		Use:   "set-operator",
		Short: "set an operator to your delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			validateOperatorFlags,
			func(ctx context.Context) error { return setOperator(ctx, cfg) },
		),
	}

	bindSetOperatorFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnsetOperatorCmd() *cobra.Command {
	var cfg operatorConfig

	cmd := &cobra.Command{
		Use:   "unset-operator",
		Short: "Unsets an existing operator from your delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			validateOperatorFlags,
			func(ctx context.Context) error { return unsetOperator(ctx, cfg) },
		),
	}

	bindUnsetOperatorFlags(cmd, &cfg)

	return cmd
}

func newValidatorSetWithdrawalAddressCmd() *cobra.Command {
	var cfg withdrawalConfig

	cmd := &cobra.Command{
		Use:   "set-withdrawal-address",
		Short: "Updates the withdrawal address that receives stake and reward withdrawals",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			validateWithdrawalFlags,
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
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				ctx := cmd.Context()
				return validateValidatorStakeFlags(ctx, cmd, &cfg)
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
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				ctx := cmd.Context()
				return validateValidatorStakeOnBehalfFlags(ctx, cmd, &cfg)
			},
			func(ctx context.Context) error { return stakeOnBehalf(ctx, cfg) },
		),
	}

	bindValidatorStakeOnBehalfFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnstakeCmd() *cobra.Command {
	var cfg unstakeConfig

	cmd := &cobra.Command{
		Use:   "unstake",
		Short: "Unstake tokens as the delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				ctx := cmd.Context()
				return validateValidatorUnstakeFlags(ctx, cmd, &cfg)
			},
			func(ctx context.Context) error { return unstake(ctx, cfg) },
		),
	}

	bindValidatorUnstakeFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnstakeOnBehalfCmd() *cobra.Command {
	var cfg unstakeConfig

	cmd := &cobra.Command{
		Use:   "unstake-on-behalf",
		Short: "Unstake tokens on behalf of a delegator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				ctx := cmd.Context()
				return validateValidatorUnstakeOnBehalfFlags(ctx, cmd, &cfg)
			},
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
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		RunE: runValidatorCommand(
			func(_ *cobra.Command) error { return nil },
			func(ctx context.Context) error { return exportKey(ctx, cfg) },
		),
	}

	bindValidatorKeyExportFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnjailCmd() *cobra.Command {
	var cfg unjailConfig

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Unjail the validator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			validateValidatorUnjailFlags,
			func(ctx context.Context) error { return unjail(ctx, cfg) },
		),
	}

	bindValidatorUnjailFlags(cmd, &cfg)

	return cmd
}

func newValidatorUnjailOnBehalfCmd() *cobra.Command {
	var cfg unjailConfig

	cmd := &cobra.Command{
		Use:   "unjail-on-behalf",
		Short: "Unjail the validator on behalf of a validator",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			validateValidatorUnjailOnBehalfFlags,
			func(ctx context.Context) error { return unjailOnBehalf(ctx, cfg) },
		),
	}

	bindValidatorUnjailOnBehalfFlags(cmd, &cfg)

	return cmd
}

func runValidatorCommand(
	validate func(cmd *cobra.Command) error,
	execute func(ctx context.Context) error,
) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		if err := validate(cmd); err != nil {
			_ = cmd.Help()
			return err
		}

		return execute(cmd.Context())
	}
}

func exportKey(_ context.Context, cfg exportKeyConfig) error {
	privKeyBytes, err := loadValidatorFile(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to load validator key file")
	}

	compressedPubKeyBytes, err := privKeyToCmpPubKey(privKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to decode compressed pub key")
	}

	if err := printKeyFormats(compressedPubKeyBytes); err != nil {
		return err
	}

	if cfg.ExportEVMKey {
		privateKey, err := crypto.ToECDSA(privKeyBytes)
		if err != nil {
			return errors.Wrap(err, "invalid private key")
		}
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
	compressedPubKeyBytes, err := validatorKeyFileToCmpPubKey(cfg.ValidatorKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to extract compressed pub key")
	}

	uncompressedPubKeyBytes, err := cmpPubKeyToUncmpPubKey(compressedPubKeyBytes)
	if err != nil {
		return err
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	_, err = prepareAndExecuteTransaction(
		ctx,
		&cfg.baseConfig,
		"createValidator",
		stakeAmount,
		uncompressedPubKeyBytes,
		cfg.Moniker,
		cfg.CommissionRate,
		cfg.MaxCommissionRate,
		cfg.MaxCommissionChangeRate,
		cfg.Unlocked,
		[]byte{},
	)
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

	fee, err := getUint256(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}
	fmt.Printf("Fee for withdrawing: %s wei\n", fee.String())

	_, err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "setWithdrawalAddress", fee, uncompressedPubKey, withdrawalAddress)
	if err != nil {
		return err
	}

	fmt.Println("Withdrawal address successfully set!")

	return nil
}

func setOperator(ctx context.Context, cfg operatorConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	operatorAddress := common.HexToAddress(cfg.Operator)

	fee, err := getUint256(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	fmt.Printf("Fee for setting operator: %s wei\n", fee.String())

	_, err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "setOperator", fee, uncompressedPubKey, operatorAddress)
	if err != nil {
		return err
	}

	fmt.Println("Operator set successfully!")

	return nil
}

func unsetOperator(ctx context.Context, cfg operatorConfig) error {
	uncompressedPubKey, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	result, err := prepareAndReadContract(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	var unsetOperatorFee *big.Int
	err = cfg.ABI.UnpackIntoInterface(&unsetOperatorFee, "fee", result)
	if err != nil {
		return errors.Wrap(err, "failed to unpack unsetOperatorFee")
	}

	_, err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unsetOperator", unsetOperatorFee, uncompressedPubKey)
	if err != nil {
		return err
	}

	fmt.Println("Operator unset successfully!")

	return nil
}

func stake(ctx context.Context, cfg stakeConfig) error {
	uncompressedDelegatorPubKeyBytes, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	validatorPubKeyBytes, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded pub key")
	}
	uncompressedValidatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(validatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress validator public key")
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	receipt, err := prepareAndExecuteTransaction(
		ctx,
		&cfg.baseConfig,
		"stakeOnBehalf",
		stakeAmount,
		uncompressedDelegatorPubKeyBytes,
		uncompressedValidatorPubKeyBytes,
		uint8(cfg.StakePeriod),
		[]byte{},
	)
	if err != nil {
		return err
	}

	fmt.Println("Tokens staked successfully! Extracting delegation ID...")

	delegationID, err := extractDelegationIDFromStake(&cfg, receipt)
	if err != nil {
		return err
	}

	fmt.Printf("Delegation ID: %s\n", delegationID.String())

	return nil
}

func stakeOnBehalf(ctx context.Context, cfg stakeConfig) error {
	delegatorPubKeyBytes, err := hex.DecodeString(cfg.DelegatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded delegator public key")
	}

	uncompressedDelegatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(delegatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress delegator public key")
	}

	validatorPubKeyBytes, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoed validator public key")
	}
	uncompressedValidatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(validatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress validator public key")
	}

	stakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	receipt, err := prepareAndExecuteTransaction(
		ctx,
		&cfg.baseConfig,
		"stakeOnBehalf",
		stakeAmount,
		uncompressedDelegatorPubKeyBytes,
		uncompressedValidatorPubKeyBytes,
		uint8(cfg.StakePeriod),
		[]byte{},
	)
	if err != nil {
		return err
	}

	fmt.Println("Tokens staked on behalf of delegator successfully! Extracting delegation ID...")

	delegationID, err := extractDelegationIDFromStake(&cfg, receipt)
	if err != nil {
		return err
	}

	fmt.Printf("Delegation ID: %s\n", delegationID.String())

	return nil
}

func unstake(ctx context.Context, cfg unstakeConfig) error {
	uncompressedDelegatorPubKeyBytes, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return err
	}

	validatorPubKeyBytes, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded validator pub key")
	}
	uncompressedValidatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(validatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress validator public key")
	}

	unstakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	delegationID := new(big.Int).SetUint64(uint64(cfg.DelegationID))

	result, err := prepareAndReadContract(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	var unstakeFee *big.Int
	err = cfg.ABI.UnpackIntoInterface(&unstakeFee, "fee", result)
	if err != nil {
		return errors.Wrap(err, "failed to unpack unstakeFee")
	}

	_, err = prepareAndExecuteTransaction(
		ctx,
		&cfg.baseConfig,
		"unstake",
		unstakeFee,
		uncompressedDelegatorPubKeyBytes,
		uncompressedValidatorPubKeyBytes,
		delegationID,
		unstakeAmount,
		[]byte{},
	)
	if err != nil {
		return err
	}

	fmt.Println("Tokens unstaked successfully!")

	return nil
}

func unstakeOnBehalf(ctx context.Context, cfg unstakeConfig) error {
	delegatorPubKeyBytes, err := hex.DecodeString(cfg.DelegatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded delegator pub key")
	}

	uncompressedDelegatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(delegatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress delegator public key")
	}

	validatorPubKeyBytes, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded validator pub key")
	}
	uncompressedValidatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(validatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress validator public key")
	}

	unstakeAmount, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		return errors.New("invalid stake amount", "amount", cfg.StakeAmount)
	}

	delegationID := new(big.Int).SetUint64(uint64(cfg.DelegationID))

	result, err := prepareAndReadContract(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	var unstakeOnBehalfFee *big.Int
	err = cfg.ABI.UnpackIntoInterface(&unstakeOnBehalfFee, "fee", result)
	if err != nil {
		return errors.Wrap(err, "failed to unpack unstakeOnBehalfFee")
	}

	_, err = prepareAndExecuteTransaction(
		ctx,
		&cfg.baseConfig,
		"unstakeOnBehalf",
		unstakeOnBehalfFee,
		uncompressedDelegatorPubKeyBytes,
		uncompressedValidatorPubKeyBytes,
		delegationID,
		unstakeAmount,
		[]byte{},
	)
	if err != nil {
		return err
	}

	fmt.Println("Tokens unstaked on behalf of delegator successfully!")

	return nil
}

func unjail(ctx context.Context, cfg unjailConfig) error {
	uncompressedValidatorPubKeyBytes, err := uncompressPrivateKey(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to get uncompressed pub key from private key")
	}

	result, err := prepareAndReadContract(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	var unjailFee *big.Int
	err = cfg.ABI.UnpackIntoInterface(&unjailFee, "fee", result)
	if err != nil {
		return errors.Wrap(err, "failed to unpack unjailFee")
	}

	fmt.Printf("Unjail fee: %s\n", unjailFee.String())

	_, err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unjail", unjailFee, uncompressedValidatorPubKeyBytes, []byte{})
	if err != nil {
		return err
	}

	fmt.Println("Validator successfully unjailed!")

	return nil
}

func unjailOnBehalf(ctx context.Context, cfg unjailConfig) error {
	validatorPubKeyBytes, err := hex.DecodeString(cfg.ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex-encoded validator public key")
	}

	uncompressedValidatorPubKeyBytes, err := cmpPubKeyToUncmpPubKey(validatorPubKeyBytes)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress validator public key")
	}

	result, err := prepareAndReadContract(ctx, &cfg.baseConfig, "fee")
	if err != nil {
		return err
	}

	var unjailFee *big.Int
	err = cfg.ABI.UnpackIntoInterface(&unjailFee, "fee", result)
	if err != nil {
		return errors.Wrap(err, "failed to unpack unjailFee")
	}

	fmt.Printf("Unjail fee: %s\n", unjailFee.String())

	_, err = prepareAndExecuteTransaction(ctx, &cfg.baseConfig, "unjailOnBehalf", unjailFee, uncompressedValidatorPubKeyBytes, []byte{})
	if err != nil {
		return err
	}

	fmt.Println("Validator successfully unjailed on behalf of validator!")

	return nil
}

func initializeBaseConfig(cfg *baseConfig) error {
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

	contractABI, err := abi.JSON(strings.NewReader(string(ipTokenStakingABI)))
	if err != nil {
		return errors.Wrap(err, "failed to parse contract ABI")
	}

	cfg.ABI = &contractABI
	cfg.ContractAddr = common.HexToAddress(predeploys.IPTokenStaking)

	return nil
}

func extractDelegationIDFromStake(cfg *stakeConfig, receipt *types.Receipt) (*big.Int, error) {
	event := cfg.ABI.Events["Deposit"]
	eventSignature := event.ID
	for _, vLog := range receipt.Logs {
		if vLog.Topics[0] == eventSignature {
			eventData := struct {
				DelegatorUncmpPubkey []byte
				ValidatorUncmpPubkey []byte
				StakeAmount          *big.Int
				StakingPeriod        *big.Int
				DelegationId         *big.Int //nolint:revive,stylecheck // Definition comes from ABI
				OperatorAddress      common.Address
				Data                 []byte
			}{}

			err := cfg.ABI.UnpackIntoInterface(&eventData, "Deposit", vLog.Data)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unpack deposit event")
			}

			return eventData.DelegationId, nil
		}
	}

	return nil, errors.New("deposit event not found in transaction logs")
}
