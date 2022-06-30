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
	config.DefineCommandInputs(executeCmd, getCommandRequirements())
	commandCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(commandCmd)
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
		ctx, cancel := buildContext(context.Background())
		defer cancel()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}

		params, err := config.LoadParamsFromFileAndEnvAndFlagsAndPrompt(cmd, getCommandRequirements(), promptReader)
		if err != nil {
			return err
		}

		tokenValidity := env.ReadEnvInt(config.SessionValiditySecondsEnvVar, api.DefaultTokenValiditySeconds)

		baseRportURL := config.ReadApiURL(params)
		wsURLBuilder := &api.WsCommandURLProvider{
			WsURLProvider: &api.WsURLProvider{
				TokenProvider: func() (token string, err error) {
					return auth.GetToken(params)
				},
				BaseURL:              baseRportURL,
				TokenValiditySeconds: tokenValidity,
			},
		}
		wsClient, err := utils.NewWsClient(ctx, wsURLBuilder.BuildWsURL)
		if err != nil {
			return err
		}

		rportAPI := buildRport(params)

		isFullJobOutput := params.ReadBool(controllers.IsFullOutput, false)
		cmdExecutor := &controllers.CommandsController{
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

func getCommandRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:    controllers.ClientIDs,
			Help:     "Enter comma separated client IDs",
			Validate: config.RequiredValidate,
			Description: "[required] Comma separated client ids for which the command should be executed. " +
				"Alternatively use -n to execute a command by client name(s), or use --search flag.",
			ShortName: "d",
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(controllers.ClientNameFlag, "") == "" && providedParams.ReadString(controllers.SearchFlag, "") == ""
			},
			IsRequired: true,
		},
		{
			Field:       controllers.ClientNameFlag,
			Description: "Comma separated client names for which the command should be executed",
			ShortName:   "n",
		},
		{
			Field:       controllers.SearchFlag,
			Description: "Search clients on all fields, supports wildcards (*).",
		},
		{
			Field:       controllers.Command,
			Help:        "Enter command",
			Validate:    config.RequiredValidate,
			Description: "[required] Command which should be executed on the clients",
			ShortName:   "c",
			IsRequired:  true,
		},
		{
			Field:       controllers.Timeout,
			Help:        "Enter timeout in seconds",
			Description: "timeout in seconds that was used to observe the command execution",
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
			Help:        "execute the command concurrently on multiple clients",
			Description: "execute the command concurrently on multiple clients",
			ShortName:   "r",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.IsFullOutput,
			Help:        "output detailed information of a job execution",
			Description: "output detailed information of a job execution",
			ShortName:   "f",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.IsSudo,
			Help:        "should execute command as sudo",
			Description: "execute command as sudo",
			ShortName:   "u",
			Type:        config.BoolRequirementType,
			Default:     false,
		},
		{
			Field:       controllers.Interpreter,
			Help:        "enter interpreter/shell name for the command execution",
			Description: "interpreter/shell name for the command execution",
			ShortName:   "i",
			Type:        config.StringRequirementType,
		},
		{
			Field:       controllers.AbortOnError,
			Description: "if true and command fails on one client, it's not executed on others",
			Help:        "should abort command if it fails on any client",
			ShortName:   "a",
			Type:        config.BoolRequirementType,
			Default:     false,
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
