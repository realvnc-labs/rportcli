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
	config.DefineCommandInputs(executeScript, getScriptRequirements())
	scriptCmd.AddCommand(executeScript)
	rootCmd.AddCommand(scriptCmd)
	addClientsSearchFlag(executeScript) // enable repeated '--search key=value' flag

	// see help.go
	scriptCmd.SetUsageTemplate(usageTemplate + serverAuthenticationRefer)
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "scripts management",
	Args:  cobra.ArbitraryArgs,
}

var executeScript = &cobra.Command{
	Use:   "execute",
	Short: "executes a remote script on rport client(s)",
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
		params, err := loadParams(cmd, getScriptRequirements(), promptReader, injected)
		if err != nil {
			return err
		}

		wsClient, err := newWsClient(ctx, params, makeWsScriptsURLProvider(params))
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		cmdExecutor := &controllers.ScriptsController{
			ExecutionHelper: newExecutionHelper(params, wsClient, rportAPI),
		}

		err = cmdExecutor.Start(ctx, params, promptReader, nil)

		return err
	},
}

func makeWsScriptsURLProvider(params *options.ParameterBag) (wsURLBuilder utils.WsURLBuilder) {
	baseRportURL := config.ReadAPIURL(params)
	urlProvider := &api.WsScriptsURLProvider{
		WsURLProvider: newWsURLProvider(params, baseRportURL),
	}

	return urlProvider.BuildWsURL
}

func getScriptRequirements() []config.ParameterRequirement {
	return config.GetScriptParamReqs()
}
