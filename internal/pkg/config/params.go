package config

import (
	"errors"
	"fmt"
	"reflect"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/env"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	ErrMultipleTargetingOptions = fmt.Errorf("multiple client targeting options. "+
		"Please only specify one of --%s, --%s, --%s, --%s, %s, or list of --%s",
		ClientIDs, GroupIDs, ClientNameFlag, ClientNamesFlag, ClientCombinedSearchFlag, ClientSearchFlag,
	)
	ErrInvalidSchemeForHTTPProxy = errors.New("--http-proxy can only be used with the http or https schemes")
)

func ReadClientNames(params *options.ParameterBag) (names string) {
	names = params.ReadString(ClientNamesFlag, "")
	if names == "" {
		names = params.ReadString(ClientNameFlag, "")
	}
	return names
}

func ExecLogRequested(params *options.ParameterBag) (requested bool, logFilename string) {
	logFilename = params.ReadString(WriteExecLog, "")
	return logFilename != "", logFilename
}

func SourceExecLog(params *options.ParameterBag) (requested bool, logFilename string) {
	logFilename = params.ReadString(ReadExecLog, "")
	return logFilename != "", logFilename
}

func ReadNoPrompt(params *options.ParameterBag) (noPrompt bool) {
	// currently just reuse the NoPrompt flag
	noPrompt = params.ReadBool(NoPrompt, false)
	return noPrompt
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

func HasHTTPProxy(params *options.ParameterBag) (hasHTTPProxy bool) {
	return params.ReadBool(UseHTTPProxy, false)
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

func CheckTargetingParams(params *options.ParameterBag) (err error) {
	paramList := []string{ClientIDs, GroupIDs, ClientNameFlag, ClientNamesFlag, ClientCombinedSearchFlag}

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
		if !req.IsRequired {
			continue
		}
		errorHint := "Specify either via the command line or include in a yaml file"
		val, found := params.Read(req.Field, nil)
		switch {
		case req.IsEnabled != nil:
			// Process all parameters with a custom IsEnabled() function that signals either the parameter is
			// literally present or an equivalent replacement is used.
			if req.IsEnabled(params) && (!found || val == "") {
				return fmt.Errorf("required option '--%s' or equivalent is missing. %s", req.Field, errorHint)
			}
		case !found || val == "":
			return fmt.Errorf("required option '--%s' is missing. %s", req.Field, errorHint)
		default:
			return nil
		}
	}
	return nil
}

func checkCorrectSchemeForHTTPProxy(params *options.ParameterBag) (isCorrect bool) {
	scheme := params.ReadString(Scheme, "")
	if scheme == utils.HTTP || scheme == utils.HTTPS {
		return true
	}
	return false
}
