package controllers

import (
	"context"

	"github.com/sirupsen/logrus"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
)

type LogoutAPI interface {
	Logout(ctx context.Context) (err error)
}

type LogoutController struct {
	api           LogoutAPI
	configDeleter func() (err error)
}

func NewLogoutController(api LogoutAPI, configDeleter func() (err error)) *LogoutController {
	return &LogoutController{
		api:           api,
		configDeleter: configDeleter,
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

	err = lc.configDeleter()
	if err != nil {
		return err
	}
	logrus.Info("Logout success")

	return nil
}
