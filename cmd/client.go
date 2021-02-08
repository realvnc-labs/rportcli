package cmd

import (
	"context"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	clientsCmd.AddCommand(clientsListCmd)
	rootCmd.AddCommand(clientsCmd)
}

var clientsCmd = &cobra.Command{
	Use:   "client [command]",
	Short: "Client API",
	Args:  cobra.ArbitraryArgs,
}

var clientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Client List API",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		paramsFromArguments := make(map[string]string, len(paramsFromArgumentsP))
		for k, valP := range paramsFromArgumentsP {
			paramsFromArguments[k] = *valP
		}
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		apiAuth := &api.BasicAuth{
			Login: cfg.ReadString(config.Login, ""),
			Pass:  cfg.ReadString(config.Password, ""),
		}
		rportAPI := api.New(config.Params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)
		cr := &output.ClientRenderer{}
		clientsController := &controllers.ClientController{
			Rport: rportAPI,
			Cr:    cr,
		}

		return clientsController.Clients(context.Background(), os.Stdout)
	},
}
