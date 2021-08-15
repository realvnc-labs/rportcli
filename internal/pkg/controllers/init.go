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
	login := params.ReadString(config.Login, "")
	serverURL := params.ReadString(config.ServerURL, config.DefaultServerURL)

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (l, p string, err error) {
			p = params.ReadString(config.Password, "")
			return login, p, nil
		},
	}

	tokenValidity := env.ReadEnvInt(config.SessionValiditySecondsEnvVar, api.DefaultTokenValiditySeconds)

	cl := api.New(serverURL, apiAuth)
	loginResp, err := cl.GetToken(ctx, tokenValidity)
	if err != nil {
		return fmt.Errorf("config verification failed against the rport: %v", err)
	}

	if loginResp.Data.TwoFA.SentTo != "" {
		twoFACl := api.New(serverURL, nil)
		loginResp, err = ic.process2FA(ctx, twoFACl, loginResp.Data.TwoFA.SentTo, login, tokenValidity)
		if err != nil {
			return fmt.Errorf("2 factor login to rport failed: %v", err)
		}
	}
	if loginResp.Data.Token == "" {
		return fmt.Errorf("no auth token received from rport")
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

func (ic *InitController) process2FA(
	ctx context.Context,
	cl *api.Rport,
	twoFASentTo, login string,
	tokenLifetime int,
) (li api.LoginResponse, err error) {
	req := config.ParameterRequirement{
		Field: "code",
		Help: fmt.Sprintf(
			"2 factor auth is enabled, please provide code that was sent to %s",
			twoFASentTo,
		),
		Validate:   config.RequiredValidate,
		IsRequired: true,
		Type:       config.StringRequirementType,
	}
	resultMap := map[string]interface{}{}
	err = config.PromptRequiredValues([]config.ParameterRequirement{req}, resultMap, ic.PromptReader)
	if err != nil {
		return li, err
	}

	li, err = cl.GetTokenBy2FA(ctx, resultMap["code"].(string), login, tokenLifetime)

	return li, err
}
