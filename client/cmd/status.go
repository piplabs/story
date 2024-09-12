// Package cmd provides commands for interacting with the Story consensus client.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	storycfg "github.com/piplabs/story/client/config"
)

// StatusConfig is the config for the status command.
type StatusConfig struct {
	HomeDir  string
	RPCLaddr string
}

// NewStatusCmd returns a cobra command that fetches the status of the local chain.
func newStatusCmd() *cobra.Command {
	cfg := StatusConfig{
		HomeDir:  storycfg.DefaultHomeDir(),
		RPCLaddr: "http://localhost:26657", // Default CometBFT RPC address
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Fetch status of the Story chain",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			cometCfg, err := parseCometConfig(ctx, cfg.HomeDir)
			if err != nil {
				return err
			}
			cfg.RPCLaddr = strings.Replace(cometCfg.RPC.ListenAddress, "tcp://", "http://", 1)

			return CheckStatus(cmd.Context(), cfg)
		},
	}

	bindStatusFlags(cmd.Flags(), &cfg)

	return cmd
}

// checkStatus queries the status endpoint of the Story chain and fetches the latest block height.
func CheckStatus(ctx context.Context, cfg StatusConfig) error {
	rpcClient, err := http.New(cfg.RPCLaddr, "")
	if err != nil {
		return errors.Wrap(err, "failed to create RPC client")
	}

	status, err := rpcClient.Status(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to query cometBFT status endpoint")
	}

	blockHeight := strconv.FormatInt(status.SyncInfo.LatestBlockHeight, 10)

	responseJSON := map[string]any{
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
