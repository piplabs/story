// Package cmd provides the cli for running the story consensus client.
package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/piplabs/story/client/app"
	storycfg "github.com/piplabs/story/client/config"
	"github.com/piplabs/story/lib/buildinfo"
	libcmd "github.com/piplabs/story/lib/cmd"
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
		newKeyCmds(),
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
