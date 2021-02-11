package cmd

import (
	"context"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	tunnelsCmd.AddCommand(tunnelListCmd)
	rootCmd.AddCommand(tunnelsCmd)
}

var tunnelsCmd = &cobra.Command{
	Use:   "tunnel [command]",
	Short: "Tunnel API",
	Args:  cobra.ArbitraryArgs,
}

var tunnelListCmd = &cobra.Command{
	Use:   "list",
	Short: "Tunnel List API",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rportAPI, err := buildRport()
		if err != nil {
			return err
		}

		tr := &output.TunnelRenderer{}
		clientsController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
		}

		return clientsController.Tunnels(context.Background(), os.Stdout)
	},
}
