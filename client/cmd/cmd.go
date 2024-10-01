// Package cmd provides the cli for running the story consensus client.
package cmd

import (
	"context"
	"fmt"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/app"
	storycfg "github.com/piplabs/story/client/config"
	"github.com/piplabs/story/lib/buildinfo"
	libcmd "github.com/piplabs/story/lib/cmd"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// New returns a new root cobra command that handles our command line tool.
func New() *cobra.Command {
	return libcmd.NewRootCmd(
		"story",
		"Story is a consensus client implementation for the Story L1 blockchain",
		newRunCmd("run", app.Run),
		newInitCmd(),
		buildinfo.NewVersionCmd(),
		newValidatorCmds(),
		newStatusCmd(),
		newRollbackCmd(app.CreateApp),
	)
}

// newRunCmd returns a new cobra command that runs the story consensus client.
func newRunCmd(name string, runFunc func(context.Context, app.Config) error) *cobra.Command {
	storyCfg := storycfg.DefaultConfig()
	logCfg := log.DefaultConfig()

	cmd := &cobra.Command{
		Use:   name,
		Short: "Runs the story consensus client",
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

			return runFunc(ctx, app.Config{
				Config: storyCfg,
				Comet:  cometCfg,
			})
		},
	}

	bindRunFlags(cmd, &storyCfg)
	log.BindFlags(cmd.Flags(), &logCfg)

	return cmd
}

// newRollbackCmd returns a new cobra command that rolls back one block of the story consensus client.
func newRollbackCmd(appCreateFunc func(context.Context, app.Config) *app.App) *cobra.Command {
	storyCfg := storycfg.DefaultConfig()
	logCfg := log.DefaultConfig()

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback Cosmos SDK and CometBFT state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when CometBFT has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also rolls back to height n - 1. No blocks are removed, so upon
restarting CometBFT the transactions in block n will be re-executed against the
application.
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

			app := appCreateFunc(ctx, app.Config{
				Config: storyCfg,
				Comet:  cometCfg,
			})
			height, hash, err := cmtcmd.RollbackState(&cometCfg, storyCfg.RemoveBlock)
			if err != nil {
				return errors.Wrap(err, "failed to rollback CometBFT state")
			}

			if err = app.CommitMultiStore().RollbackToVersion(height); err != nil {
				return errors.Wrap(err, "failed to rollback to version")
			}

			fmt.Printf("Rolled back state to height %d and hash %X", height, hash)

			return nil
		},
	}

	bindRunFlags(cmd, &storyCfg)
	bindRollbackFlags(cmd, &storyCfg)
	log.BindFlags(cmd.Flags(), &logCfg)

	return cmd
}
