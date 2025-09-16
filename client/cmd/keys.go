package cmd

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/tempfile"
	"github.com/cometbft/cometbft/privval"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/app"
)

type keyConfig struct {
	ValidatorKeyFile      string
	PrivateKeyFile        string
	PubKeyHex             string
	PubKeyBase64          string
	PubKeyHexUncompressed string
	EncPrivKeyFile        string
}

func newKeyCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Commands for key management",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newKeyConvertCmd(),
		newKeyGenPrivKeyJSONCmd(),
		newKeyEncryptCmd(),
		newKeyShowEncryptedCmd(),
	)

	return cmd
}

func newKeyConvertCmd() *cobra.Command {
	var cfg keyConfig

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert between various key formats",
		Args:  cobra.NoArgs,
		RunE: runValidatorCommand(
			validateKeyConvertFlags,
			func(ctx context.Context) error { return convertKey(ctx, cfg) },
		),
	}

	bindKeyConvertFlags(cmd, &cfg)

	return cmd
}

func newKeyGenPrivKeyJSONCmd() *cobra.Command {
	var cfg genPrivKeyJSONConfig

	cmd := &cobra.Command{
		Use:   "gen-priv-key-json",
		Short: "Generate a priv_validator_key.json file from EVM private key",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeBaseConfig(&cfg.baseConfig)
		},
		RunE: runValidatorCommand(
			func(_ *cobra.Command) error {
				return validateGenPrivKeyJSONFlags(&cfg)
			},
			func(ctx context.Context) error { return genValidatorPrivKeyJSON(ctx, cfg) },
		),
	}

	bindKeyGenPrivKeyJSONFlags(cmd, &cfg)

	return cmd
}

func newKeyEncryptCmd() *cobra.Command {
	var cfg baseConfig

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the private key stored in .env",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				return validateEncryptFlags(cmd, &cfg)
			},
			func(_ context.Context) error { return encryptPrivKey(cfg) },
		),
	}

	bindValidatorBaseFlags(cmd, &cfg)

	return cmd
}

func newKeyShowEncryptedCmd() *cobra.Command {
	var cfg showEncryptedConfig

	cmd := &cobra.Command{
		Use:   "show-encrypted",
		Short: "Show the encrypted private key after decryption",
		Args:  cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		RunE: runValidatorCommand(
			func(cmd *cobra.Command) error {
				return validateShowEncryptedFlags(cmd, &cfg)
			},
			func(_ context.Context) error { return showEncryptedKey(cfg) },
		),
	}

	bindKeyShowEncryptedFlags(cmd, &cfg)

	return cmd
}

func convertKey(_ context.Context, cfg keyConfig) error {
	var compressedPubKeyBytes []byte
	var err error

	switch {
	case cfg.ValidatorKeyFile != "":
		compressedPubKeyBytes, err = validatorKeyFileToCmpPubKey(cfg.ValidatorKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to load validator private key")
		}
	case cfg.PrivateKeyFile != "":
		compressedPubKeyBytes, err = privKeyFileToCmpPubKey(cfg.PrivateKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to load private key file")
		}
	case cfg.PubKeyHex != "":
		pubKeyHex := strings.TrimPrefix(cfg.PubKeyHex, "0x")
		compressedPubKeyBytes, err = hex.DecodeString(pubKeyHex)
		if err != nil {
			return errors.Wrap(err, "failed to decode hex public key")
		}
	case cfg.PubKeyBase64 != "":
		compressedPubKeyBytes, err = base64.StdEncoding.DecodeString(cfg.PubKeyBase64)
		if err != nil {
			return errors.Wrap(err, "failed to decode base64 public key")
		}
	case cfg.PubKeyHexUncompressed != "":
		pubKeyHex := strings.TrimPrefix(cfg.PubKeyHexUncompressed, "0x")
		uncompressedPubKeyBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			return errors.Wrap(err, "failed to decode hex public key")
		}
		compressedPubKeyBytes, err = uncmpPubKeyToCmpPubKey(uncompressedPubKeyBytes)
		if err != nil {
			return errors.Wrap(err, "failed to convert uncompressed pub key")
		}
	case cfg.EncPrivKeyFile != "":
		password, err := app.InputPassword(
			app.PasswordPromptText,
			"",
			false,
			app.ValidatePasswordInput,
		)
		if err != nil {
			return errors.Wrap(err, "error occurred while input password")
		}

		pv, err := app.LoadEncryptedPrivKey(password, cfg.EncPrivKeyFile)
		if err != nil {
			return errors.Wrap(err, "failed to load encrypted private key")
		}
		compressedPubKeyBytes = pv.PubKey.Bytes()
	default:
		return errors.New("no valid key input provided")
	}

	return printKeyFormats(compressedPubKeyBytes)
}

func genValidatorPrivKeyJSON(_ context.Context, cfg genPrivKeyJSONConfig) error {
	privKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	privKey := k1.PrivKey(privKeyBytes)
	newPV := &privval.FilePVKey{
		Address: privKey.PubKey().Address(),
		PubKey:  privKey.PubKey(),
		PrivKey: privKey,
	}

	jsonBytes, err := cmtjson.MarshalIndent(newPV, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal pv data")
	}

	if err := tempfile.WriteFileAtomic(cfg.ValidatorKeyFile, jsonBytes, 0600); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}

func encryptPrivKey(cfg baseConfig) error {
	password, err := app.InputPassword(
		app.NewKeyPasswordPromptText,
		app.ConfirmPasswordPromptText,
		true, /* Should confirm password */
		app.ValidatePasswordInput,
	)
	if err != nil {
		return errors.Wrap(err, "error occurred while input password")
	}

	privKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode private key")
	}

	pk := k1.PrivKey(privKeyBytes)
	pv := privval.FilePVKey{
		PrivKey: pk,
		PubKey:  pk.PubKey(),
		Address: pk.PubKey().Address(),
	}

	if err := app.EncryptAndStoreKey(pv, password, cfg.EncPrivKeyFile); err != nil {
		return errors.Wrap(err, "failed to encrypt and store the key")
	}

	return nil
}

func showEncryptedKey(cfg showEncryptedConfig) error {
	password, err := app.InputPassword(
		app.PasswordPromptText,
		"",
		false,
		app.ValidatePasswordInput,
	)
	if err != nil {
		return errors.Wrap(err, "error occurred while input password")
	}

	pv, err := app.LoadEncryptedPrivKey(password, cfg.EncPrivKeyFile)
	if err != nil {
		return errors.Wrap(err, "failed to load encrypted private key")
	}

	cmpPubKeyBytes, err := privKeyToCmpPubKey(pv.PrivKey.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to get compressed public key from private key")
	}

	if err := printKeyFormats(cmpPubKeyBytes); err != nil {
		return errors.Wrap(err, "failed to print key formats")
	}

	if cfg.ShowPrivate {
		fmt.Println("Private Key (hex):", hex.EncodeToString(pv.PrivKey.Bytes()))
	}

	return nil
}
