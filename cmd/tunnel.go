package cmd

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/cache"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/client"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	tunnelsCmd.AddCommand(tunnelListCmd)

	config.DefineCommandInputs(tunnelDeleteCmd, getDeleteTunnelRequirements())
	tunnelsCmd.AddCommand(tunnelDeleteCmd)

	config.DefineCommandInputs(tunnelCreateCmd, getCreateTunnelRequirements())
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
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		clientSearch := &client.Search{
			DataProvider: rportAPI,
			Cache:        &cache.ClientsCache{},
		}

		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider:     rportAPI,
			ClientSearch:   clientSearch,
		}

		return tunnelController.Tunnels(context.Background())
	},
}

var tunnelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "terminates the specified tunnel of the specified client",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}
		params, err := config.CollectParams(cmd, getDeleteTunnelRequirements(), promptReader)
		if err != nil {
			return err
		}

		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		clientSearch := &client.Search{
			DataProvider: rportAPI,
			Cache:        &cache.ClientsCache{},
		}

		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider:     rportAPI,
			ClientSearch:   clientSearch,
		}

		clientID := params.ReadString(controllers.ClientID, "")
		tunnelID := params.ReadString(controllers.TunnelID, "")
		clientName := params.ReadString(controllers.ClientNameFlag, "")
		return tunnelController.Delete(context.Background(), clientID, clientName, tunnelID)
	},
}

var tunnelCreateCmd = &cobra.Command{
	Use: "create",
	Long: `creates a new tunnel, e.g.
rportcli tunnel create -l 0.0.0.0:22 -r 3394 -d bc0b705d-b5fb-4df5-84e3-82dba437bbef -s ssh --acl 10.1.2.3
this example opens port 3394 on the rport server and forwards to port 22 of the client bc0b705d-b5fb-4df5-84e3-82dba437bbef
with ssh url scheme and an IP address 10:1:2:3 allowed to access the tunnel
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}
		params, err := config.CollectParams(cmd, getCreateTunnelRequirements(), promptReader)
		if err != nil {
			return err
		}

		rportAPI := buildRport()

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		clientSearch := &client.Search{
			DataProvider: rportAPI,
			Cache:        &cache.ClientsCache{},
		}

		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider:     rportAPI,
			ClientSearch:   clientSearch,
		}

		return tunnelController.Create(context.Background(), params)
	},
}

func getCreateTunnelRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       controllers.ClientID,
			Description: "[conditionally required] client id, if not provided, client name should be given",
			Validate:    config.RequiredValidate,
			ShortName:   "d",
			IsRequired:  true,
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(controllers.ClientNameFlag, "") == ""
			},
			Help: "Enter a client ID",
		},
		{
			Field:       controllers.ClientNameFlag,
			Description: `client name, if no client id is provided`,
			ShortName:   "n",
		},
		{
			Field: controllers.Local,
			Description: `refers to the ports of the rport server address to use for a new tunnel, e.g. '3390' or '0.0.0.0:3390'. 
If local is not specified, a random server port will be assigned automatically`,
			ShortName: "l",
		},
		{
			Field: controllers.Remote,
			Description: "[required] the ports are defined from the servers' perspective. " +
				"'Remote' refers to the ports and interfaces of the client., e.g. '3389'",
			ShortName:  "r",
			IsRequired: true,
			Validate:   config.RequiredValidate,
			Help:       "Enter a remote port value",
		},
		{
			Field:       controllers.Scheme,
			Description: "URI scheme to be used. For example, 'ssh', 'rdp', etc.",
			ShortName:   "s",
		},
		{
			Field:       controllers.ACL,
			Description: "ACL, IP addresses who is allowed to use the tunnel. For example, '142.78.90.8,201.98.123.0/24,'",
			Default:     controllers.DefaultACL,
			ShortName:   "a",
		},
		{
			Field:       controllers.CheckPort,
			Description: "A flag whether to check availability of a public port. By default check is disabled.",
			ShortName:   "p",
			Type:        config.BoolRequirementType,
			Default:     "0",
		},
	}
}

func getDeleteTunnelRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       controllers.ClientID,
			Description: "[conditionally required] client id, if not provided, client name should be given",
			Validate:    config.RequiredValidate,
			ShortName:   "i",
			IsRequired:  true,
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(controllers.ClientNameFlag, "") == ""
			},
			Help: "Enter a client id",
		},
		{
			Field:       controllers.ClientNameFlag,
			Description: `client name, if no client id is provided`,
			ShortName:   "n",
		},
		{
			Field:       controllers.TunnelID,
			Description: "[required]  tunnel id to delete",
			ShortName:   "t",
			IsRequired:  true,
			Validate:    config.RequiredValidate,
			Help:        "Enter a tunnel id",
		},
	}
}
