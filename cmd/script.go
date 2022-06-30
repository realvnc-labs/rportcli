package cmd

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/breathbath/go_utils/v2/pkg/env"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/auth"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

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
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}

		params, err := config.LoadParamsFromFileAndEnvAndFlagsAndPrompt(cmd, getScriptRequirements(), promptReader)
		if err != nil {
			return err
		}

		baseRportURL := config.ReadApiURL(params)
		tokenValidity := env.ReadEnvInt(config.SessionValiditySecondsEnvVar, api.DefaultTokenValiditySeconds)
		wsURLBuilder := &api.WsScriptsURLProvider{
			WsURLProvider: &api.WsURLProvider{
				BaseURL: baseRportURL,
				TokenProvider: func() (token string, err error) {
					return auth.GetToken(params)
				},
				TokenValiditySeconds: tokenValidity,
			},
		}

		ctx, cancel := buildContext(context.Background())
		defer cancel()

		wsClient, err := utils.NewWsClient(ctx, wsURLBuilder.BuildWsURL)
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		isFullJobOutput := params.ReadBool(controllers.IsFullOutput, false)
		cmdExecutor := &controllers.ScriptsController{
			ExecutionHelper: &controllers.ExecutionHelper{
				ReadWriter: wsClient,
				JobRenderer: &output.JobRenderer{
					Writer:       os.Stdout,
					Format:       getOutputFormat(),
					IsFullOutput: isFullJobOutput,
				},
				Rport: rportAPI,
			},
		}

		err = cmdExecutor.Start(ctx, params)

		return err
	},
}

func getScriptRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
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
