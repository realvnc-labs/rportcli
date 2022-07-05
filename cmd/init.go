package cmd

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"

	"github.com/spf13/cobra"
)

var isLogout bool

func init() {
	config.DefineCommandInputs(initCmd, getInitRequirements())
	initCmd.Flags().BoolVarP(&isLogout, "delete", "d", false, "Logout user")
	rootCmd.AddCommand(initCmd)

	// see help.go
	initCmd.SetUsageTemplate(usageTemplate + environmentVariables + serverAuthentication)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize your connection to the rportd API",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := buildContext(context.Background())
		defer cancel()

		if isLogout {
			return manageLogout(ctx, cmd)
		}

		return manageInit(ctx, cmd)
	},
	Long: `
The init command allows the user to authenticate with an rport server and cache the authentication token in a config.json file. The user will not need to re-authenticate until the token expires, at which point the init command must be run again. 

If 2fa is enabled then init will take the user through the 2fa flow and save the final token. The user will not need to complete 2fa until the token expires.`,
}

func manageLogout(ctx context.Context, cmd *cobra.Command) error {
	params, err := config.LoadParamsFromFileAndEnv(cmd.Flags())
	if err != nil {
		return err
	}

	rportAPI := buildRport(params)

	logoutController := controllers.NewLogoutController(rportAPI, config.DeleteConfig)

	return logoutController.Logout(ctx, params)
}

func manageInit(ctx context.Context, cmd *cobra.Command) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// when the RPORT_API_TOKEN env var is set, we shouldn't allow use of the init command
	hasAPIToken := config.HasAPIToken()
	if hasAPIToken {
		return errors.New("cannot init config when the RPORT_API_TOKEN is set. If you do not wish to use the API token, please unset RPORT_API_TOKEN and use RPORT_API_USER and RPORT_API_PASSWORD instead")
	}

	promptReader := &utils.PromptReader{
		Sc:              bufio.NewScanner(os.Stdin),
		SigChan:         sigs,
		PasswordScanner: utils.ReadPassword,
	}

	params, err := config.LoadParamsFromFileAndEnvAndFlagsAndPrompt(cmd, getInitRequirements(), promptReader)
	if err != nil {
		return err
	}

	initController := &controllers.InitController{
		ConfigWriter: config.WriteConfig,
		PromptReader: promptReader,
		TotPSecretRenderer: &output.TotPSecretRenderer{
			ColCountCalculator: utils.CalcTerminalColumnsCount,
			Writer:             os.Stdout,
			Format:             getOutputFormat(),
		},
		QrImageWriterProvider: output.GetQrImageFsWriter,
	}

	return initController.InitConfig(ctx, params)
}

func getInitRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		config.GetNoPromptFlagSpec(),
		{
			Field:       config.APIURL,
			Help:        "Enter Server Url",
			Validate:    config.RequiredValidate,
			Description: "Server address of rport to connect to",
			ShortName:   "s",
			Default:     "",
			Type:        config.StringRequirementType,
		},
		{
			Field:       config.APIUser,
			Help:        "Enter a valid login value",
			Validate:    config.RequiredValidate,
			Description: "Login to the rport server",
			ShortName:   "l",
		},
		{
			Field:       config.APIPassword,
			Help:        "Enter a valid password value",
			Validate:    config.RequiredValidate,
			Description: "Password to the rport server",
			ShortName:   "p",
			IsSecure:    true,
		},
	}
}
