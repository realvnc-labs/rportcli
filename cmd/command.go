package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/breathbath/go_utils/utils/env"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	isInteractive                  bool
	commandExecutionFromArgumentsP map[string]*string
)

func init() {
	commandsCmd.Flags().BoolVarP(&isInteractive, "interactive", "i", false, "opens interactive session for command execution")

	reqs := controllers.GetCommandRequirements()
	commandExecutionFromArgumentsP = make(map[string]*string, len(reqs))
	for _, req := range reqs {
		paramVal := ""
		commandsCmd.Flags().StringVarP(&paramVal, req.Field, req.ShortName, req.Default, req.Description)
		commandExecutionFromArgumentsP[req.Field] = &paramVal
	}

	rootCmd.AddCommand(commandsCmd)
}

var commandsCmd = &cobra.Command{
	Use:   "command",
	Short: "executes remote command on rport client",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isInteractive {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			basicAuth := &utils.StorageBasicAuth{
				AuthProvider: config.AuthConfigProvider,
			}
			baseRportURL := cfg.ReadString(config.ServerURL, config.DefaultServerURL)

			rportCl := api.New(baseRportURL, basicAuth)

			wsURLBuilder := &api.WsCommandURLProvider{
				TokenProvider:        rportCl,
				BaseURL:              baseRportURL,
				TokenValiditySeconds: env.ReadEnvInt("SESSION_VALIDITY_SECONDS", api.DefaultTokenValiditySeconds),
			}
			wsClient, err := utils.NewWsClient(ctx, wsURLBuilder.BuildWsURL)
			if err != nil {
				return err
			}

			cmdExecutor := &controllers.InteractiveCommandsController{
				ReadWriter:   wsClient,
				PromptReader: &utils.PromptReader{},
				Spinner:      utils.NewSpinner(),
				JobRenderer:  &output.JobRenderer{},
				OutputWriter: os.Stdin,
			}

			err = cmdExecutor.Start(ctx, commandExecutionFromArgumentsP)

			return err
		}
		// todo implement
		return errors.New("non interactive command execution is not implemented yet")
	},
}
