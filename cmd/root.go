package cmd

import (
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/applog"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	Verbose = false
	rootCmd = &cobra.Command{
		Use:           "rportcli",
		Short:         "Rport cli",
		RunE:          initCmd.RunE,
		Version:       version(),
		SilenceErrors: true,
		SilenceUsage:  true,
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

func buildRport() (*api.Rport, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	apiAuth := &utils.BasicAuth{
		Login: cfg.ReadString(config.Login, ""),
		Pass:  cfg.ReadString(config.Password, ""),
	}
	rportAPI := api.New(config.Params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)

	return rportAPI, nil
}
