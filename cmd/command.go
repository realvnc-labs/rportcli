package cmd

import (
	"context"
	"errors"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	isInteractive bool
)

func init() {
	commandsCmd.Flags().BoolVarP(&isInteractive, "interactive", "i", false, "opens interactive session for command execution")

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

			wsURLBuilder := &api.WsCommandURLProvider{
				AuthProvider: config.AuthConfigProvider,
				BaseURL:      cfg.ReadString(config.ServerURL, config.DefaultServerURL),
			}
			wsClient, err := utils.NewWsClient(ctx, wsURLBuilder.BuildWsURL)
			if err != nil {
				return err
			}

			cmdExecutor := &api.InteractiveCommandExecutor{
				ReadWriter:      wsClient,
				UserInputReader: &utils.PromptReader{},
			}

			err = cmdExecutor.Start(ctx)

			return err
		}
		// todo implement
		return errors.New("non interactive command execution is not implemented yet")
	},
}
