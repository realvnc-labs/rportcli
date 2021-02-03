package cmd

import (
	"github.com/cloudradar-monitoring/rportcli/config"
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

		missedRequirements := config.GetNotMatchedRequirements(config.Params)
		if len(missedRequirements) > 0 {
			err := config.PromptRequiredValues(missedRequirements, paramsFromArguments)
			if err != nil {
				return err
			}
			config.Params = config.FromValues(paramsFromArguments)
		}

		return nil
	},
}
