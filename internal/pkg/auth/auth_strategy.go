package auth

import (
	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

func GetUsernameAndPassword(params *options.ParameterBag) (login, pass string, err error) {
	login = config.ReadAPIUser(params)
	pass = config.ReadAPIPassword(params)
	apiToken := params.ReadString(config.APIToken, "")
	if pass != "" && apiToken != "" {
		return "", "", utils.ErrAPIPasswordAndAPITokenAreBothSet
	}
	if apiToken != "" {
		return login, apiToken, err
	}
	return login, pass, err
}

func GetToken(params *options.ParameterBag) (token string, err error) {
	return params.ReadString(config.Token, ""), nil
}

func GetAuthStrategy(params *options.ParameterBag) (auth *utils.FallbackAuth) {
	auth = &utils.FallbackAuth{
		PrimaryAuth: &utils.StorageBasicAuth{
			AuthProvider: func() (login, pass string, err error) {
				return GetUsernameAndPassword(params)
			},
		},
		FallbackAuth: &utils.BearerAuth{
			TokenProvider: func() (string, error) {
				return GetToken(params)
			},
		},
	}

	return auth
}