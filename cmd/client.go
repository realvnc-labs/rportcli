package cmd

import (
	"context"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	addClientsPaginationFlags(clientsListCmd)
	addClientsSearchFlag(clientsListCmd)
	clientsCmd.AddCommand(clientsListCmd)
	clientCmd.Flags().StringP(config.ClientNameFlag, "n", "", "[deprecated] Get client by name")
	clientCmd.Flags().StringP(config.ClientNamesFlag, "", "", "Get client by name(s)")
	clientCmd.Flags().BoolP("all", "a", false, "Show client info with additional details")
	clientsCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(clientsCmd)

	// see help.go
	clientsCmd.SetUsageTemplate(usageTemplate + serverAuthenticationRefer)
}

var clientsCmd = &cobra.Command{
	Use:   "client [command]",
	Short: "manage rport clients",
	Args:  cobra.ArbitraryArgs,
}

var clientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all connected and disconnected rport clients",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		params, err := config.LoadParamsFromFileAndEnv(cmd.Flags())
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)
		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		return clientsController.Clients(ctx, params)
	},
}

var clientCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "get all details about a specific client identified by its id or flags like names",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var clientName string
		var clientID string
		if len(args) == 0 {
			cn, err := cmd.Flags().GetString(config.ClientNamesFlag)
			if err != nil {
				return err
			}
			if cn == "" {
				cn, err = cmd.Flags().GetString(config.ClientNameFlag)
				if err != nil {
					return err
				}
			}
			clientName = cn
		} else {
			clientID = args[0]
		}

		params, err := config.LoadParamsFromFileAndEnv(cmd.Flags())
		if err != nil {
			return err
		}
		rportAPI := buildRport(params)

		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		return clientsController.Client(ctx, params, clientID, clientName)
	},
}

func addClientsPaginationFlags(cmd *cobra.Command) {
	// TODO: why isn't this getting picked up
	cmd.Flags().IntP(api.PaginationLimit, "", api.ClientsLimitDefault, "Number of clients to fetch")
	cmd.Flags().IntP(api.PaginationOffset, "", 0, "Offset for clients fetch")
}

func addClientsSearchFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("search", "", "", "Search clients on all fields, supports wildcards (*).")
}
