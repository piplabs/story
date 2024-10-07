package cmd

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type keyConfig struct {
	ValidatorKeyFile string
	PrivateKeyFile   string
	PubKeyHex        string
	PubKeyBase64     string
}

func newKeyCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Commands for key management",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newKeyConvertCmd(),
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
			func() error { return validateKeyConvertFlags(cfg) },
			func(ctx context.Context) error { return convertKey(ctx, cfg) },
		),
	}

	bindKeyConvertFlags(cmd, &cfg)

	return cmd
}

func convertKey(ctx context.Context, cfg keyConfig) error {
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
	default:
		return fmt.Errorf("no valid key input provided")
	}

	if err != nil {
		return err
	}

	return printKeyFormats(compressedPubKeyBytes)
}
