package cmd

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			baseRportURL := config.Params.ReadString(config.ServerURL, config.DefaultServerURL)
			wsURLBuilder := &api.WsCommandURLProvider{
				TokenProvider: func() (token string, err error) {
					token = config.Params.ReadString(config.Token, "")
					return
				},
				BaseURL: baseRportURL,
			}
			wsClient, err := utils.NewWsClient(ctx, wsURLBuilder.BuildWsURL)
			if err != nil {
				return err
			}

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			var spinner controllers.Spinner = &utils.NullSpinner{}
			if Verbose {
				spinner = utils.NewSpinner()
			}
			cmdExecutor := &controllers.InteractiveCommandsController{
				ReadWriter: wsClient,
				PromptReader: &utils.PromptReader{
					Sc:              bufio.NewScanner(os.Stdin),
					SigChan:         sigs,
					PasswordScanner: utils.ReadPassword,
				},
				Spinner: spinner,
				JobRenderer: &output.JobRenderer{
					Writer: os.Stdin,
					Format: getOutputFormat(),
				},
			}

			err = cmdExecutor.Start(ctx, commandExecutionFromArgumentsP)

			return err
		}
		// todo implement
		return errors.New("non interactive command execution is not implemented yet")
	},
}
