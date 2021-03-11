package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	clientsCmd.AddCommand(clientsListCmd)
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
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rportAPI, err := buildRport()
		if err != nil {
			return err
		}
		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
		}

		return clientsController.Clients(context.Background())
	},
}

var clientCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "get all details about a specific client identified by its id",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("client id is not provided")
		}

		rportAPI, err := buildRport()
		if err != nil {
			return err
		}

		cr := &output.ClientRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		clientsController := &controllers.ClientController{
			Rport:          rportAPI,
			ClientRenderer: cr,
		}

		return clientsController.Client(context.Background(), args[0])
	},
}
