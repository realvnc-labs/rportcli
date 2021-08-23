package controllers

import (
	"context"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
)

type LogoutAPI interface {
	Logout(ctx context.Context) (err error)
}

type LogoutController struct {
	api          LogoutAPI
	configWriter ConfigWriter
}

func NewLogoutController(api LogoutAPI, configWriter ConfigWriter) *LogoutController {
	return &LogoutController{
		api:          api,
		configWriter: configWriter,
	}
}

func (lc *LogoutController) Logout(ctx context.Context, params *options.ParameterBag) error {
	var err error

	err = lc.api.Logout(ctx)
	if err != nil {
		return err
	}

	serverURL := params.ReadString(config.ServerURL, "")
	if serverURL == "" {
		return nil
	}

	valuesProvider := options.NewMapValuesProvider(map[string]interface{}{
		config.ServerURL: params.ReadString(config.ServerURL, ""),
		config.Token:     "",
	})

	err = lc.configWriter(options.New(valuesProvider))
	if err != nil {
		return err
	}

	return nil
}
