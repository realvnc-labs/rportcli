package cmd

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"syscall"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/spf13/cobra"
)

func init() {
	addClientsPaginationFlags(tunnelListCmd)
	addClientsSearchFlag(tunnelListCmd)
	tunnelListCmd.Flags().StringP(config.ClientNameFlag, "n", "", "Get tunnels of a client by name")
	tunnelListCmd.Flags().StringP(config.ClientID, "c", "", "Get tunnels of a client by client id")
	tunnelsCmd.AddCommand(tunnelListCmd)

	config.DefineCommandInputs(tunnelDeleteCmd, getDeleteTunnelRequirements())
	tunnelsCmd.AddCommand(tunnelDeleteCmd)

	config.DefineCommandInputs(tunnelCreateCmd, getCreateTunnelRequirements())
	tunnelsCmd.AddCommand(tunnelCreateCmd)

	rootCmd.AddCommand(tunnelsCmd)

	// see help.go
	tunnelsCmd.SetUsageTemplate(usageTemplate + serverAuthenticationRefer)
}

func getDeleteTunnelRequirements() []config.ParameterRequirement {
	return config.GetDeleteTunnelParamReqs()
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
		params, err := config.LoadParamsFromFileAndEnv(cmd.Flags())
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		tr := &output.TunnelRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		}

		tunnelController := &controllers.TunnelController{
			Rport:          rportAPI,
			TunnelRenderer: tr,
			IPProvider:     rportAPI,
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		return tunnelController.Tunnels(ctx, params)
	},
}

var tunnelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "terminates the specified tunnel of the specified client",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		params, err := readParams(cmd, getDeleteTunnelRequirements())
		if err != nil {
			return err
		}

		tunnelController, err := createTunnelController(params)
		if err != nil {
			return err
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		return tunnelController.Delete(ctx, params)
	},
}

var tunnelCreateCmd = &cobra.Command{
	Use:  "create",
	Long: config.CreateTunnelLong,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		params, err := readParams(cmd, getCreateTunnelRequirements())
		if err != nil {
			return err
		}

		tunnelController, err := createTunnelController(params)
		if err != nil {
			return err
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		return tunnelController.Create(ctx, params)
	},
}

func getCreateTunnelRequirements() []config.ParameterRequirement {
	return config.GetCreateTunnelParamReqs(IsRDPUserRequired)
}

// use nolint to quiet the lint error about error always returning nil
func createTunnelController(params *options.ParameterBag) (*controllers.TunnelController, error) { //nolint:unparam
	rportAPI := buildRport(params)

	tr := &output.TunnelRenderer{
		ColCountCalculator: utils.CalcTerminalColumnsCount,
		Writer:             os.Stdout,
		Format:             getOutputFormat(),
	}

	return &controllers.TunnelController{
		Rport:          rportAPI,
		TunnelRenderer: tr,
		IPProvider:     rportAPI,
	}, nil
}

func readParams(cmd *cobra.Command, reqs []config.ParameterRequirement) (*options.ParameterBag, error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	promptReader := &utils.PromptReader{
		Sc:              bufio.NewScanner(os.Stdin),
		SigChan:         sigs,
		PasswordScanner: utils.ReadPassword,
	}

	return config.LoadParamsFromFileAndEnvAndFlagsAndPrompt(cmd, reqs, promptReader, nil)
}
