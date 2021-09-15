package cmd

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"syscall"

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
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize your connection to the rportd API",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if isLogout {
			return manageLogout(ctx, cmd)
		}

		return manageInit(ctx, cmd)
	},
}

func manageLogout(ctx context.Context, cmd *cobra.Command) error {
	params := config.LoadParamsFromFileAndEnv(cmd.Flags())

	rportAPI := buildRport(params)

	logoutController := controllers.NewLogoutController(rportAPI, config.DeleteConfig)

	return logoutController.Logout(ctx, params)
}

func manageInit(ctx context.Context, cmd *cobra.Command) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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
	}

	return initController.InitConfig(ctx, params)
}

func getInitRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       config.ServerURL,
			Help:        "Enter Server Url",
			Validate:    config.RequiredValidate,
			Description: "Server address of rport to connect to",
			ShortName:   "s",
		},
		{
			Field:       config.Login,
			Help:        "Enter a valid login value",
			Validate:    config.RequiredValidate,
			Description: "Login to the rport server",
			ShortName:   "l",
		},
		{
			Field:       config.Password,
			Help:        "Enter a valid password value",
			Validate:    config.RequiredValidate,
			Description: "Password to the rport server",
			ShortName:   "p",
			IsSecure:    true,
		},
	}
}
