package cmd

import (
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/applog"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	Verbose      = false
	OutputFormat = output.FormatHuman
	IsJSONPretty = false
	rootCmd      = &cobra.Command{
		Use:           "rportcli",
		Short:         "Rport cli",
		RunE:          initCmd.RunE,
		Version:       version(),
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if OutputFormat != "" && OutputFormat != output.FormatHuman && OutputFormat != output.FormatYAML && OutputFormat != output.FormatJSON {
				return fmt.Errorf(
					"unknown format '%s', supported formats are %s, %s, %s",
					OutputFormat,
					output.FormatJSON,
					output.FormatJSONPretty,
					output.FormatYAML,
				)
			}
			return nil
		},
	}
)

func getOutputFormat() string {
	if OutputFormat == output.FormatJSON {
		if IsJSONPretty {
			return output.FormatJSONPretty
		}
		return output.FormatJSON
	}

	return OutputFormat
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&IsJSONPretty,
		"json-pretty",
		"j",
		false,
		"in combination with json format this flag will pretty print the json data",
	)
	rootCmd.PersistentFlags().StringVarP(
		&OutputFormat,
		"output",
		"o",
		output.FormatHuman,
		fmt.Sprintf("Output format: %s, %s or %s", output.FormatJSON, output.FormatYAML, output.FormatHuman),
	)
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

func buildRport() *api.Rport {
	auth := &utils.FallbackAuth{
		PrimaryAuth: &utils.StorageBasicAuth{
			AuthProvider: func() (login, pass string, err error) {
				login = config.Params.ReadString(config.Login, "")
				pass = config.Params.ReadString(config.Password, "")
				return
			},
		},
		FallbackAuth: &utils.BearerAuth{
			TokenProvider: func() (string, error) {
				return config.Params.ReadString(config.Token, ""), nil
			},
		},
	}

	rportAPI := api.New(config.Params.ReadString(config.ServerURL, config.DefaultServerURL), auth)

	return rportAPI
}
