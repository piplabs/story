package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/app"
	storycfg "github.com/piplabs/story/client/config"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/log"
)

type rollbackConfig struct {
	RemoveBlock     bool // See cosmos-sdk/server/rollback.go
	RollbackEVM     bool
	RollbackHeights int64
}

// newRollbackCmd returns a new cobra command that rolls back one block of the story consensus client.
func newRollbackCmd(appCreateFunc func(context.Context, app.Config) *app.App) *cobra.Command {
	rollbackCfg := rollbackConfig{
		RemoveBlock:     false,
		RollbackEVM:     false,
		RollbackHeights: 1,
	}
	storyCfg := storycfg.DefaultConfig()
	logCfg := log.DefaultConfig()

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback Cosmos SDK and CometBFT state by X height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when CometBFT has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - X.
The application also rolls back to height n - X. If --hard=true, the block
itself will also be deleted and re-downloaded from the p2p network. Note that
different blocks from n - X to n cannot be re-built/re-proposed since that would result in validator slashing.
If --hard=false, No blocks are removed, so upon restarting CometBFT the transactions in blocks [n - X + 1, n]
will be re-executed against the application.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := log.Init(cmd.Context(), logCfg)
			if err != nil {
				return err
			}
			if err := libcmd.LogFlags(ctx, cmd.Flags()); err != nil {
				return err
			}

			cometCfg, err := parseCometConfig(ctx, storyCfg.HomeDir)
			if err != nil {
				return err
			}

			appCfg := app.Config{
				Config: storyCfg,
				Comet:  cometCfg,
			}
			a := appCreateFunc(ctx, appCfg)
			lastHeight, lastHash, err := app.RollbackCometAndAppState(ctx, a, appCfg, cometCfg, rollbackCfg.RollbackHeights, rollbackCfg.RemoveBlock, rollbackCfg.RollbackEVM)
			if err != nil {
				return err
			}

			fmt.Printf("Rolled back state to height %d and hash %X", lastHeight, lastHash)

			return nil
		},
	}

	bindRunFlags(cmd, &storyCfg)
	bindRollbackFlags(cmd, &rollbackCfg)
	log.BindFlags(cmd.Flags(), &logCfg)

	return cmd
}
