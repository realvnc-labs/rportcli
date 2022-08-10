package cmd

import (
	"bufio"
	"os"
	"strings"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	config.DefineCommandInputs(executeCmd, getCommandRequirements())
	commandCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(commandCmd)
	addClientsSearchFlag(executeCmd) // enable repeated '--search key=value' flag

	// see help.go
	commandCmd.SetUsageTemplate(usageTemplate + serverAuthenticationRefer)
}

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "command management",
	Args:  cobra.ArbitraryArgs,
}

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "executes a remote command on an rport client(s)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel, sigs := makeRunContext()
		defer cancel()

		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}

		var injected map[string]string
		if len(searchFlags) > 0 {
			injected = map[string]string{"combined-search": strings.Join(searchFlags, "&")}
		}
		params, err := loadParams(cmd, getCommandRequirements(), promptReader, injected)
		if err != nil {
			return err
		}

		wsClient, err := newWsClient(ctx, params, makeWsCommandURLProvider(params))
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		cmdExecutor := &controllers.CommandsController{
			ExecutionHelper: newExecutionHelper(params, wsClient, rportAPI),
		}

		err = cmdExecutor.Start(ctx, params, promptReader, nil)

		return err
	},
}

func makeWsCommandURLProvider(params *options.ParameterBag) (wsURLBuilder utils.WsURLBuilder) {
	baseRportURL := config.ReadAPIURL(params)
	urlProvider := &api.WsCommandURLProvider{
		WsURLProvider: newWsURLProvider(params, baseRportURL),
	}

	return urlProvider.BuildWsURL
}

func getCommandRequirements() []config.ParameterRequirement {
	return config.GetCommandParamReqs()
}
