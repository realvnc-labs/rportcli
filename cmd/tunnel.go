package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	tunnelsCmd.AddCommand(tunnelListCmd)
	tunnelsCmd.AddCommand(tunnelDeleteCmd)
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
		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
		}

		return tunnelController.Tunnels(context.Background(), os.Stdout)
	},
}

const minArgsCount = 2

var tunnelDeleteCmd = &cobra.Command{
	Use:   "delete <CLIENT_ID> <TUNNEL_ID>",
	Short: "Terminates the specified tunnel of the specified client",
	Args:  cobra.MinimumNArgs(minArgsCount),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < minArgsCount {
			return fmt.Errorf("either CLIENT_ID or TUNNEL_ID is not provided")
		}

		rportAPI, err := buildRport()
		if err != nil {
			return err
		}

		tr := &output.TunnelRenderer{}
		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
		}

		return tunnelController.Delete(context.Background(), os.Stdout, args[0], args[1])
	},
}
