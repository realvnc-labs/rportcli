package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	io2 "github.com/breathbath/go_utils/v2/pkg/io"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/breathbath/go_utils/v2/pkg/fs"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/env"
	"github.com/sirupsen/logrus"
)

const (
	defaultPath      = ".config/rportcli/config.json"
	ServerURL        = "server"
	Login            = "login"
	Token            = "token"
	Password         = "password"
	DefaultServerURL = ""
	APIURL           = "api_url"
	APIUser          = "api_user"
	APIPassword      = "api_password"
	APIToken         = "api_token"
	NoPrompt         = "no-prompt"
	ReadYAML         = "read-yaml"
)

var (
	ErrAPIURLRequired = errors.New("please set the server API URL, either via RPORT_API_URL or the command line flags (if option available)")
)

func LoadParamsFromFileAndEnv(flags *pflag.FlagSet) (params *options.ParameterBag, err error) {
	var valuesProvider *options.ValuesProviderComposite

	envParams := readEnvVars()
	envValuesProvider := options.NewMapValuesProvider(envParams)

	flagValuesProvider := CreateFlagValuesProvider(flags)

	if HasAPIToken() {
		// ignore config file if using api token
		valuesProvider = options.NewValuesProviderComposite(envValuesProvider, flagValuesProvider)
	} else {
		fileValuesProvider, err := CreateFileValuesProvider()
		if err != nil {
			logrus.Warn(err)
			valuesProvider = options.NewValuesProviderComposite(envValuesProvider, flagValuesProvider)
		} else {
			valuesProvider = options.NewValuesProviderComposite(envValuesProvider, flagValuesProvider, fileValuesProvider)
		}
	}

	paramsToReturn := options.New(valuesProvider)

	if err := CheckIfMissingAPIURL(paramsToReturn); err != nil {
		return paramsToReturn, err
	}
	WarnIfLegacyConfig()

	return paramsToReturn, nil
}

func CreateFileValuesProvider() (options.ValuesProvider, error) {
	configFilePath := getConfigLocation()
	if !fs.FileExists(configFilePath) {
		return nil, fmt.Errorf("config file %s doesn't exist", configFilePath)
	}

	f, err := os.Open(configFilePath)
	if err != nil {
		err = fmt.Errorf("failed to open the file %s: %v", configFilePath, err)
		return nil, err
	}
	defer io2.CloseResourceSecure("config file", f)

	jvp, err := options.NewJSONValuesProvider(f)
	if err != nil {
		return nil, err
	}

	return jvp, nil
}

func LoadParamsFromFileAndEnvAndFlagsAndPrompt(
	c *cobra.Command,
	reqs []ParameterRequirement,
	promptReader PromptReader,
) (params *options.ParameterBag, err error) {
	flagsProvider := &FlagValuesProvider{
		flags: c.Flags(),
	}

	valuesProviderFromCommandAndPrompt, err := CollectParamsFromCommandAndPromptAndEnv(flagsProvider, reqs, promptReader)
	if err != nil {
		return nil, err
	}

	valueProviders := []options.ValuesProvider{
		valuesProviderFromCommandAndPrompt,
	}

	// ignore config file if has api token
	if !HasAPIToken() {
		jvp, err := CreateFileValuesProvider()
		if err != nil {
			logrus.Warn(err)
		} else {
			valueProviders = append(valueProviders, jvp)
		}
	}

	mergedValuesProvider := options.NewValuesProviderComposite(valueProviders...)
	paramsSoFar := options.New(mergedValuesProvider)

	if err := CheckIfMissingAPIURL(paramsSoFar); err != nil {
		return paramsSoFar, err
	}
	WarnIfLegacyConfig()

	return paramsSoFar, nil
}

func CollectParamsFromCommandAndPromptAndEnv(
	flagsProvider *FlagValuesProvider,
	reqs []ParameterRequirement,
	promptReader PromptReader,
) (vp options.ValuesProvider, err error) {
	// potentially revist the use of the env values provider. for now, the control of working
	// with the maps directly for merging etc seems best.
	envParams := readEnvVars()

	flagParams, err := readFlags(reqs, flagsProvider)
	if err != nil {
		return nil, err
	}

	rawParams := flagParams
	MergeMaps(rawParams, envParams)

	yParamFileList, hasYAMLParams := HasYAMLParams(flagParams)
	if hasYAMLParams {
		yParams, err := readYAML(yParamFileList, flagsProvider)
		if err != nil {
			return nil, err
		}

		if yParams != nil {
			MergeMaps(rawParams, yParams)
		}
	}

	vp = options.NewMapValuesProvider(rawParams)
	paramsSoFar := options.New(vp)

	if err := CheckRequiredParams(paramsSoFar, reqs); err != nil {
		return nil, err
	}

	if err := CheckTargetingParams(paramsSoFar); err != nil {
		return nil, err
	}

	// if the no-prompt cli flag is set, then do not prompt for missing values
	noPrompt := paramsSoFar.ReadBool(NoPrompt, false)
	if noPrompt || promptReader == nil {
		return vp, nil
	}

	missedRequirements := CheckRequirements(paramsSoFar, reqs)
	if len(missedRequirements) == 0 {
		return vp, nil
	}

	err = PromptRequiredValues(missedRequirements, rawParams, promptReader)
	if err != nil {
		return nil, err
	}

	return options.NewMapValuesProvider(rawParams), nil
}

func MergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

func readEnvVars() (rawParams map[string]interface{}) {
	envsToRead := map[string]string{
		Password:            PasswordEnvVar,
		Login:               LoginEnvVar,
		ServerURL:           ServerURLEnvVar,
		PathForConfigEnvVar: PathForConfigEnvVar,
		APIURL:              APIURLEnvVar,
		APIUser:             APIUserEnvVar,
		APIPassword:         APIPasswordEnvVar,
		APIToken:            APITokenEnvVar,
	}

	rawParams = map[string]interface{}{}
	for paramName, envVarName := range envsToRead {
		envVarValue := env.ReadEnv(envVarName, "")
		if envVarValue != "" {
			rawParams[paramName] = envVarValue
			continue
		}

		if envVarName == APIUserEnvVar {
			envVarValue = env.ReadEnv(LoginEnvVar, "")
		}
		if envVarName == APIPassword {
			envVarValue = env.ReadEnv(PasswordEnvVar, "")
		}
		if envVarName == APIURL {
			envVarValue = env.ReadEnv(ServerURLEnvVar, "")
		}
		if envVarValue != "" {
			rawParams[envVarName] = envVarValue
		}
	}

	return rawParams
}

func readFlags(reqs []ParameterRequirement,
	flagsProvider *FlagValuesProvider) (rawParams map[string]interface{}, err error) {
	rawParams = make(map[string]interface{}, len(reqs))
	for _, req := range reqs {
		flagVal, isFound, err := flagsProvider.ReadFlag(req.Field, req.Type)
		if err != nil {
			return nil, err
		}
		if isFound {
			rawParams[req.Field] = flagVal
		}
	}
	return rawParams, nil
}

func readYAML(yParamFileList []string,
	flagsChecker UsedFlagsChecker) (rawParams map[string]interface{}, err error) {
	yParams, err := ReadYAMLExecuteParams(yParamFileList, flagsChecker)
	if err != nil {
		return nil, err
	}
	// TEST: --read-yaml set but no params found
	if yParams == nil {
		return nil, nil
	}

	return yParams, nil
}

func getConfigLocation() (configPath string) {
	configPathFromEnv := env.ReadEnv(PathForConfigEnvVar, "")
	if configPathFromEnv != "" {
		configPath = configPathFromEnv
		return
	}

	usr, err := user.Current()
	if err != nil {
		logrus.Warnf("failed to read current user data: %v", err)
		configPath = "config.yaml"
		return
	}

	pathParts := []string{usr.HomeDir}
	pathParts = append(pathParts, strings.Split(defaultPath, "/")...)
	configPath = filepath.Join(pathParts...)
	return
}
