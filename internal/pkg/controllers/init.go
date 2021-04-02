package controllers

import (
	"context"
	"fmt"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/env"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

type ConfigWriter func(params *options.ParameterBag) (err error)

type InitController struct {
	ConfigWriter ConfigWriter
	PromptReader config.PromptReader
}

func (ic *InitController) InitConfig(ctx context.Context, params *options.ParameterBag) error {
	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = params.ReadString(config.Login, "")
			pass = params.ReadString(config.Password, "")
			return
		},
	}

	cl := api.New(params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)
	loginResp, err := cl.GetToken(ctx, env.ReadEnvInt(config.SessionValiditySecondsEnvVar, api.DefaultTokenValiditySeconds))
	if err != nil {
		return fmt.Errorf("config verification failed against the rport: %v", err)
	}
	if loginResp.Data.Token == "" {
		return fmt.Errorf("empty token received from rport")
	}

	valuesProvider := options.NewMapValuesProvider(map[string]interface{}{
		config.ServerURL: params.ReadString(config.ServerURL, ""),
		config.Token:     loginResp.Data.Token,
	})

	err = ic.ConfigWriter(options.New(valuesProvider))
	if err != nil {
		return err
	}

	return nil
}
