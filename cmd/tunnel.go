package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

var (
	tunnelCreateRequirementsP map[string]*string
)

func init() {
	tunnelsCmd.AddCommand(tunnelListCmd)
	tunnelsCmd.AddCommand(tunnelDeleteCmd)

	tunnelCreateRequirements := controllers.GetCreateTunnelRequirements()
	tunnelCreateRequirementsP = make(map[string]*string, len(tunnelCreateRequirements))
	for _, req := range tunnelCreateRequirements {
		paramVal := ""
		tunnelCreateCmd.Flags().StringVarP(&paramVal, req.Field, req.ShortName, req.Default, req.Description)
		if req.IsRequired {
			err := tunnelCreateCmd.MarkFlagRequired(req.Field)
			if err != nil {
				logrus.Error(err)
			}
		}
		tunnelCreateRequirementsP[req.Field] = &paramVal
	}
	tunnelsCmd.AddCommand(tunnelCreateCmd)

	rootCmd.AddCommand(tunnelsCmd)
}

var tunnelsCmd = &cobra.Command{
	Use:   "tunnel [command]",
	Short: "manage tunnels of connected clients",
	Args:  cobra.ArbitraryArgs,
}

var tunnelListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all active tunnels created with rport",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider: utils.APIIPProvider{
				URL: utils.IPCheckerURL,
			},
		}

		return tunnelController.Tunnels(context.Background())
	},
}

const minArgsCount = 2

var tunnelDeleteCmd = &cobra.Command{
	Use:   "delete <CLIENT_ID> <TUNNEL_ID>",
	Short: "terminates the specified tunnel of the specified client",
	Args:  cobra.MinimumNArgs(minArgsCount),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < minArgsCount {
			return fmt.Errorf("either CLIENT_ID or TUNNEL_ID is not provided")
		}

		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider: utils.APIIPProvider{
				URL: utils.IPCheckerURL,
			},
		}

		return tunnelController.Delete(context.Background(), args[0], args[1])
	},
}

var tunnelCreateCmd = &cobra.Command{
	Use: "create",
	Long: `creates a new tunnel, e.g.
rportcli tunnel create -l 0.0.0.0:22 -r 3394 -d bc0b705d-b5fb-4df5-84e3-82dba437bbef -s ssh --acl 10.1.2.3
this example opens port 3394 on the rport server and forwards to port 22 of the client bc0b705d-b5fb-4df5-84e3-82dba437bbef
with ssh url scheme and an IP address 10:1:2:3 allowed to access the tunnel
`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tunnelCreateRequirements := make(map[string]string, len(tunnelCreateRequirementsP))
		for k, valP := range tunnelCreateRequirementsP {
			tunnelCreateRequirements[k] = *valP
		}
		params := config.FromValues(tunnelCreateRequirements)

		err := config.CheckRequirementsError(params, controllers.GetCreateTunnelRequirements())
		if err != nil {
			return err
		}

		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}
		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider: utils.APIIPProvider{
				URL: utils.IPCheckerURL,
			},
		}

		return tunnelController.Create(context.Background(), params)
	},
}
