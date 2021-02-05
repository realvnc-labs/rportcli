package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/spf13/cobra"
)

var (
	paramsFromArgumentsP map[string]*string
)

func init() {
	reqs := config.GetParameterRequirements()
	paramsFromArgumentsP = make(map[string]*string, len(reqs))
	for _, req := range reqs {
		paramVal := ""
		initCmd.Flags().StringVarP(&paramVal, req.Field, string(req.Field[0]), req.Default, req.Description)
		paramsFromArgumentsP[req.Field] = &paramVal
	}
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [arg=value ...]",
	Short: "Initialize the active profile of the config",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		paramsFromArguments := make(map[string]string, len(paramsFromArgumentsP))
		for k, valP := range paramsFromArgumentsP {
			paramsFromArguments[k] = *valP
		}
		config.Params = config.FromValues(paramsFromArguments)

		missedRequirements := config.CheckRequirements(config.Params, config.GetParameterRequirements())
		if len(missedRequirements) > 0 {
			reader := bufio.NewReader(os.Stdin)
			err := config.PromptRequiredValues(missedRequirements, paramsFromArguments, reader)
			if err != nil {
				return err
			}
			config.Params = config.FromValues(paramsFromArguments)
		}

		apiAuth := &api.BasicAuth{
			Login: config.Params.ReadString(config.Login, ""),
			Pass:  config.Params.ReadString(config.Password, ""),
		}
		cl := api.New(config.Params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)
		_, err := cl.Status(context.Background())
		if err != nil {
			return fmt.Errorf("config verification failed against the rport API: %v", err)
		}

		err = config.WriteConfig(config.Params)
		if err != nil {
			return err
		}

		return nil
	},
}
