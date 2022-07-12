package config

import (
	"errors"
	"fmt"
	"reflect"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/env"
	"github.com/sirupsen/logrus"
)

func ReadClientNames(params *options.ParameterBag) (names string) {
	names = params.ReadString(ClientNamesFlag, "")
	if names == "" {
		names = params.ReadString(ClientNameFlag, "")
	}
	return names
}

func ReadAPIUser(params *options.ParameterBag) (user string) {
	user = params.ReadString(APIUser, "")
	if user == "" {
		user = params.ReadString(Login, "")
	}
	return user
}

func ReadAPIPassword(params *options.ParameterBag) (pass string) {
	pass = params.ReadString(APIPassword, "")
	if pass == "" {
		pass = params.ReadString(Password, "")
	}
	return pass
}

func ReadAPIURL(params *options.ParameterBag) (url string) {
	url = ReadAPIURLWithDefault(params, DefaultServerURL)
	return url
}

func ReadAPIURLWithDefault(params *options.ParameterBag, defaultURL string) (url string) {
	url = params.ReadString(APIURL, "")
	if url == "" {
		url = params.ReadString(ServerURL, defaultURL)
	}
	return url
}

func HasAPIToken() (hasAPIToken bool) {
	apiToken := env.ReadEnv(APITokenEnvVar, "")
	return apiToken != ""
}

func CheckIfMissingAPIURL(params *options.ParameterBag) (err error) {
	APIURL := ReadAPIURL(params)
	if APIURL == "" {
		return ErrAPIURLRequired
	}
	return nil
}

func WarnIfLegacyConfig() {
	login := env.ReadEnv(LoginEnvVar, "")
	pass := env.ReadEnv(PasswordEnvVar, "")
	serverURL := env.ReadEnv(ServerURLEnvVar, "")
	if login != "" || pass != "" || (serverURL != "") {
		logrus.Warn("use of RPORT_USER, RPORT_PASSWORD and RPORT_SERVER_URL will be removed in a future release. " +
			"Please use RPORT_API_USER, RPORT_API_PASSWORD and RPORT_API_URL instead")
	}
}

func HasYAMLParams(flagParams map[string]interface{}) (yFileList []string, hasYAMLParams bool) {
	if flagParams[ReadYAML] == nil {
		return nil, false
	}
	readYAMLParam := flagParams[ReadYAML]
	// just being extra cautious
	if reflect.TypeOf(readYAMLParam) == reflect.TypeOf([]string{}) {
		yFiles := flagParams[ReadYAML].([]string)
		return yFiles, true
	}

	return nil, false
}

var (
	ErrMultipleTargetingOptions = errors.New("multiple targeting options. Please only specify one of --cids, --gids, --name or --names")
)

func CheckTargetingParams(params *options.ParameterBag) (err error) {
	paramList := []string{ClientIDs, GroupIDs, ClientNameFlag, ClientNamesFlag}

	count := 0
	for _, param := range paramList {
		if params.ReadString(param, "") != "" {
			count++
		}
	}

	if count > 1 {
		return ErrMultipleTargetingOptions
	}

	return nil
}

func CheckRequiredParams(params *options.ParameterBag, reqs []ParameterRequirement) (err error) {
	for _, req := range reqs {
		// TODO: this needs refactoring to DRY out
		if req.IsRequired {
			if req.IsEnabled == nil {
				val, found := params.Read(req.Field, nil)
				if !found || val == "" {
					return fmt.Errorf("required option (--%s or equivalent) is missing. "+
						"It must be specified either via the command line or included in a yaml params file", req.Field)
				}
			} else {
				val, found := params.Read(req.Field, nil)
				if req.IsEnabled(params) && (!found || val == "") {
					return fmt.Errorf("required option (--%s or equivalent) is missing. "+
						"It must be specified either via the command line or included in a yaml params file", req.Field)
				}
			}
		}
	}
	return nil
}
