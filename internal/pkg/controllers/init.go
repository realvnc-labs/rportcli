package controllers

import (
	"context"
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

func GetInitRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       config.ServerURL,
			Help:        "Enter Server Url",
			Default:     config.DefaultServerURL,
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

type InitController struct{}

func (cc *InitController) InitConfig(ctx context.Context, parametersFromArguments map[string]*string) error {
	paramsFromArguments := make(map[string]string, len(parametersFromArguments))
	for k, valP := range parametersFromArguments {
		paramsFromArguments[k] = *valP
	}
	config.Params = config.FromValues(paramsFromArguments)

	missedRequirements := config.CheckRequirements(config.Params, GetInitRequirements())
	if len(missedRequirements) > 0 {
		err := config.PromptRequiredValues(missedRequirements, paramsFromArguments, &utils.PromptReader{})
		if err != nil {
			return err
		}
		config.Params = config.FromValues(paramsFromArguments)
	}

	apiAuth := &utils.BasicAuth{
		Login: config.Params.ReadString(config.Login, ""),
		Pass:  config.Params.ReadString(config.Password, ""),
	}

	cl := api.New(config.Params.ReadString(config.ServerURL, config.DefaultServerURL), apiAuth)
	_, err := cl.Status(context.Background())
	if err != nil {
		return fmt.Errorf("config verification failed against the rport API: %v", err)
	}

	err = config.WriteConfig(config.Params)
	if err != nil {
		return err
	}

	_, err = cl.Status(ctx)
	if err != nil {
		return fmt.Errorf("config verification failed against the rport API: %v", err)
	}

	return nil
}
