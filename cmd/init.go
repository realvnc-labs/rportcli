package cmd

import (
	"context"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"

	"github.com/spf13/cobra"
)

var (
	paramsFromArgumentsP map[string]*string
)

func init() {
	reqs := controllers.GetInitRequirements()
	paramsFromArgumentsP = make(map[string]*string, len(reqs))
	for _, req := range reqs {
		paramVal := ""
		initCmd.Flags().StringVarP(&paramVal, req.Field, req.ShortName, req.Default, req.Description)
		paramsFromArgumentsP[req.Field] = &paramVal
	}
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize your connection to the rportd API",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		initController := &controllers.InitController{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return initController.InitConfig(ctx, paramsFromArgumentsP)
	},
}
