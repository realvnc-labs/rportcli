package controllers

import (
	"context"
	"fmt"
	options "github.com/breathbath/go_utils/utils/config"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

func GetInitRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       config.ServerURL,
			Help:        "Enter Server Url",
			Validate:    config.RequiredValidate,
			Description: "Server address of rport to connect to",
			ShortName:   "s",
		},
		{
			Field:       config.Login,
			Help:        "Enter a valid login value",
			Validate:    config.RequiredValidate,
			Description: "GetToken to the rport server",
			ShortName:   "l",
		},
		{
			Field:       config.Password,
			Help:        "Enter a valid password value",
			Validate:    config.RequiredValidate,
			Description: "Password to the rport server",
			ShortName:   "p",
			IsSecure:    true,
		},
	}
}

type ConfigWriter func(params *options.ParameterBag) (err error)

type InitController struct {
	ConfigWriter ConfigWriter
	PromptReader config.PromptReader
}

func (ic *InitController) InitConfig(ctx context.Context, parametersFromArguments map[string]*string) error {
	paramsFromArguments := make(map[string]string, len(parametersFromArguments))
	for k, valP := range parametersFromArguments {
		paramsFromArguments[k] = *valP
	}
	params := config.FromValues(paramsFromArguments)

	missedRequirements := config.CheckRequirements(params, GetInitRequirements())
	if len(missedRequirements) > 0 {
		err := config.PromptRequiredValues(missedRequirements, paramsFromArguments, ic.PromptReader)
		if err != nil {
			return err
		}
		params = config.FromValues(paramsFromArguments)
	}

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = params.ReadString(config.Login, "")
			pass = params.ReadString(config.Password, "")
			return
		},
	}

	cl := api.New(params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)
	loginResp, err := cl.GetToken(ctx, env.ReadEnvInt(config.SessionValiditySecondsEnvVar, int(api.DefaultTokenValiditySeconds)))
	if err != nil {
		return fmt.Errorf("config verification failed against the rport: %v", err)
	}
	if loginResp.Data.Token == "" {
		return fmt.Errorf("empty token received from rport")
	}

	paramsToSave := config.FromValues(map[string]string{
		config.ServerURL: params.ReadString(config.ServerURL, ""),
		config.Token:     loginResp.Data.Token,
	})

	err = ic.ConfigWriter(paramsToSave)
	if err != nil {
		return err
	}

	return nil
}
