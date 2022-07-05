package cmd

import (
	"strconv"

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

		params, err := loadParams(cmd, sigs, getScriptRequirements())
		if err != nil {
			return err
		}

		wsClient, err := newWsClient(ctx, makeWsScriptsURLProvider(params))
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		cmdExecutor := &controllers.ScriptsController{
			ExecutionHelper: newExecutionHelper(params, wsClient, rportAPI),
		}

		err = cmdExecutor.Start(ctx, params)

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
	return []config.ParameterRequirement{
		config.GetNoPromptFlagSpec(),
		config.GetReadYAMLFlagSpec(),
		{
			Field:    controllers.ClientIDs,
			Help:     "Enter comma separated client IDs",
			Validate: config.RequiredValidate,
			Description: "[required] Comma separated client ids on which the script should be executed. " +
				"Alternatively use -n to execute a script by client name(s), or use --search flag.",
			ShortName: "d",
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(controllers.ClientNameFlag, "") == "" && providedParams.ReadString(controllers.SearchFlag, "") == ""
			},
			IsRequired: true,
		},
		{
			Field:       controllers.ClientNameFlag,
			Description: "Comma separated client names on which the script should be executed",
			ShortName:   "n",
		},
		{
			Field:       controllers.SearchFlag,
			Description: "Search clients on all fields, supports wildcards (*).",
		},
		{
			Field:       controllers.Script,
			Help:        "Enter script path",
			Validate:    config.RequiredValidate,
			Description: "[required] Path to the script file",
			ShortName:   "s",
			IsRequired:  true,
		},
		{
			Field:       controllers.Timeout,
			Help:        "Enter timeout in seconds",
			Description: "timeout in seconds that was used to observe the script execution",
			Default:     strconv.Itoa(controllers.DefaultCmdTimeoutSeconds),
			ShortName:   "t",
		},
		{
			Field:       controllers.GroupIDs,
			Help:        "Enter comma separated group IDs",
			Description: "Comma separated client group IDs",
			ShortName:   "g",
		},
		{
			Field:       controllers.ExecConcurrently,
			Help:        "execute the script concurrently on multiple clients",
			Description: "execute the script concurrently on multiple clients",
			ShortName:   "r",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.IsFullOutput,
			Help:        "output detailed information of a script execution",
			Description: "output detailed information of a script execution",
			ShortName:   "f",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.IsSudo,
			Help:        "execute script as sudo",
			Description: "execute script as sudo",
			ShortName:   "u",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.Interpreter,
			Help:        "enter interpreter/shell name for the script execution",
			Description: "interpreter/shell name for the script execution",
			ShortName:   "i",
			Type:        config.StringRequirementType,
		},
		{
			Field:       controllers.Cwd,
			Help:        "enter current working directory",
			Description: "current working directory",
			ShortName:   "w",
			Type:        config.StringRequirementType,
			Default:     "",
		},
	}
}
