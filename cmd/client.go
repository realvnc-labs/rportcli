package cmd

import (
	"context"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/cache"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/client"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	clientsCmd.AddCommand(clientsListCmd)
	clientCmd.Flags().StringP(controllers.ClientNameFlag, "n", "", "Get client by name")
	clientsCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(clientsCmd)
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
		params := config.LoadParamsFromFileAndEnv()

		rportAPI := buildRport(params)
		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		clientSearch := &client.Search{
			DataProvider: rportAPI,
			Cache:        &cache.ClientsCache{},
		}

		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
			ClientSearch:   clientSearch,
		}

		return clientsController.Clients(context.Background())
	},
}

var clientCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "get all details about a specific client identified by its id or flags like name",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var clientName string
		var clientID string
		if len(args) == 0 {
			cn, err := cmd.Flags().GetString(controllers.ClientNameFlag)
			if err != nil {
				return err
			}
			clientName = cn
		} else {
			clientID = args[0]
		}

		params := config.LoadParamsFromFileAndEnv()
		rportAPI := buildRport(params)

		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		clientSearch := &client.Search{
			DataProvider: rportAPI,
			Cache:        &cache.ClientsCache{},
		}
		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
			ClientSearch:   clientSearch,
		}

		return clientsController.Client(context.Background(), params, clientID, clientName)
	},
}
