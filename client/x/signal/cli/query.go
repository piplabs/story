package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/piplabs/story/client/x/signal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the CLI query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGetUpgrade())
	return cmd
}

func CmdGetUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Query for the upgrade information if an upgrade is pending",
		Args:    cobra.NoArgs,
		Example: "upgrade",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GetUpgrade(cmd.Context(), &types.QueryGetUpgradeRequest{})
			if err != nil {
				return err
			}

			if resp.Upgrade != nil {
				return clientCtx.PrintString(fmt.Sprintf("An upgrade is pending to app version %d at height %d.\n", resp.Upgrade.AppVersion, resp.Upgrade.UpgradeHeight))
			}
			return clientCtx.PrintString("No upgrade is pending.\n")
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
