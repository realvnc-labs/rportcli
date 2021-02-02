package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exeCmd)
}

var exeCmd = &cobra.Command{
	Use:   "exec {{PATH_TO_SCRIPT}}",
	Short: "",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		logrus.Infof("will execute script")

		return nil
	},
}
