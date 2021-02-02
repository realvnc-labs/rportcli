package cmd

import (
	"github.com/cloudradar-monitoring/rportcli/applog"
	"github.com/spf13/cobra"
)

var (
	Verbose = false
	rootCmd = &cobra.Command{
		Use:     "rportcli",
		Short:   "Rport cli",
		RunE:    exeCmd.RunE,
		Version: version(),
	}
)

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func initLog() {
	applog.Init(Verbose)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
