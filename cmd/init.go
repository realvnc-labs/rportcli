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

func init() {
	config.DefineCommandInputs(initCmd, getInitRequirements())
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize your connection to the rportd API",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		promptReader := &utils.PromptReader{
			Sc:              bufio.NewScanner(os.Stdin),
			SigChan:         sigs,
			PasswordScanner: utils.ReadPassword,
		}
		params, err := config.CollectParams(cmd, getInitRequirements(), promptReader)
		if err != nil {
			return err
		}

		initController := &controllers.InitController{
			ConfigWriter: config.WriteConfig,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return initController.InitConfig(ctx, params)
	},
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
			Description: "GetToken to the rport server",
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
