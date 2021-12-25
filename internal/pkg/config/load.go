package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
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
	DefaultServerURL = "http://localhost:3000"
)

func LoadParamsFromFileAndEnv(flags *pflag.FlagSet) (params *options.ParameterBag) {
	envValuesProvider := CreateEnvValuesProvider()
	fileValuesProvider, err := CreateFileValuesProvider()
	if err != nil {
		logrus.Warn(err)
		return options.New(envValuesProvider)
	}

	flagValuesProvider := CreateFlagValuesProvider(flags)

	valuesProvider := options.NewValuesProviderComposite(envValuesProvider, flagValuesProvider, fileValuesProvider)

	paramsToReturn := options.New(valuesProvider)

	return paramsToReturn
}

type FlagValuesProvider struct {
	flags *pflag.FlagSet
}

func CreateFlagValuesProvider(flags *pflag.FlagSet) options.ValuesProvider {
	return &FlagValuesProvider{flags: flags}
}

func (fvp *FlagValuesProvider) Dump(w io.Writer) (err error) {
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(fvp.ToKeyValues())
	return
}

func (fvp *FlagValuesProvider) ToKeyValues() map[string]interface{} {
	res := make(map[string]interface{})
	fvp.flags.VisitAll(func(flag *pflag.Flag) {
		res[flag.Name] = flag.Value.String()
	})

	return res
}

func (fvp *FlagValuesProvider) Read(name string) (val interface{}, found bool) {
	fl := fvp.flags.Lookup(name)
	if fl == nil || !fl.Changed {
		return nil, false
	}

	return fl.Value.String(), true
}

func CreateEnvValuesProvider() options.ValuesProvider {
	envsToRead := map[string]string{
		Password:            PasswordEnvVar,
		Login:               LoginEnvVar,
		ServerURL:           ServerURLEnvVar,
		PathForConfigEnvVar: PathForConfigEnvVar,
	}

	envMapValues := map[string]interface{}{}
	for paramName, envVarName := range envsToRead {
		envVarValue := env.ReadEnv(envVarName, "")
		if envVarValue != "" {
			envMapValues[paramName] = envVarValue
		}
	}

	return options.NewMapValuesProvider(envMapValues)
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

func DeleteConfig() (err error) {
	configLocation := getConfigLocation()

	if _, e := os.Stat(configLocation); e == nil {
		err = os.Remove(configLocation)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteConfig will write config values to file system
func WriteConfig(params *options.ParameterBag) (err error) {
	configLocation := getConfigLocation()

	configDir := filepath.Dir(configLocation)
	if _, e := os.Stat(configDir); os.IsNotExist(e) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return err
		}
	}

	configToWrite := map[string]interface{}{
		ServerURL: params.ReadString(ServerURL, ""),
		Token:     params.ReadString(Token, ""),
	}

	err = DeleteConfig()
	if err != nil {
		return err
	}

	fileToWrite, err := os.OpenFile(configLocation, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer io2.CloseResourceSecure("config file", fileToWrite)

	encoder := json.NewEncoder(fileToWrite)
	err = encoder.Encode(configToWrite)
	if err != nil {
		return err
	}

	logrus.Infof("created config at %s", configLocation)

	return nil
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

func DefineCommandInputs(c *cobra.Command, reqs []ParameterRequirement) {
	for _, req := range reqs {
		defaultStr := ""
		if req.Default != nil {
			defaultStr = fmt.Sprint(req.Default)
		}
		switch req.Type {
		case BoolRequirementType:
			boolValDefault := true
			if defaultStr == "" || defaultStr == "0" || defaultStr == "false" {
				boolValDefault = false
			}
			c.Flags().BoolP(req.Field, req.ShortName, boolValDefault, req.Description)
		case IntRequirementType:
			defaultInt, err := strconv.Atoi(defaultStr)
			if err == nil {
				c.Flags().IntP(req.Field, req.ShortName, defaultInt, req.Description)
			} else {
				c.Flags().IntP(req.Field, req.ShortName, 0, req.Description)
			}
		default:
			c.Flags().StringP(req.Field, req.ShortName, defaultStr, req.Description)
		}
	}
}

func LoadParamsFromFileAndEnvAndFlagsAndPrompt(
	c *cobra.Command,
	reqs []ParameterRequirement,
	promptReader PromptReader,
) (params *options.ParameterBag, err error) {
	envValuesProvider := CreateEnvValuesProvider()
	valueProviders := []options.ValuesProvider{
		envValuesProvider,
	}

	valuesProviderFromCommandAndPrompt, err := CollectParamsFromCommandAndPromptAndEnv(c, reqs, promptReader, envValuesProvider)
	if err != nil {
		return nil, err
	}
	valueProviders = append(valueProviders, valuesProviderFromCommandAndPrompt)

	jvp, err := CreateFileValuesProvider()
	if err != nil {
		logrus.Warn(err)
	} else {
		valueProviders = append(valueProviders, jvp)
	}

	mergedValuesProvider := options.NewValuesProviderComposite(valueProviders...)

	return options.New(mergedValuesProvider), nil
}

func CollectParamsFromCommandAndPromptAndEnv(
	c *cobra.Command,
	reqs []ParameterRequirement,
	promptReader PromptReader,
	altValuesProvider options.ValuesProvider,
) (vp options.ValuesProvider, err error) {
	paramsRaw := make(map[string]interface{}, len(reqs))
	for _, req := range reqs {
		envVal, isFound := altValuesProvider.Read(req.Field)
		if isFound {
			paramsRaw[req.Field] = envVal
			continue
		}

		switch req.Type {
		case BoolRequirementType:
			boolVal, e := c.Flags().GetBool(req.Field)
			if e != nil {
				return nil, e
			}
			paramsRaw[req.Field] = boolVal
		case IntRequirementType:
			intVal, e := c.Flags().GetInt(req.Field)
			if e != nil {
				return nil, e
			}
			paramsRaw[req.Field] = intVal
		default:
			strVal, e := c.Flags().GetString(req.Field)
			if e != nil {
				return nil, e
			}
			paramsRaw[req.Field] = strVal
		}
	}

	valuesProviderFromFlags := options.NewMapValuesProvider(paramsRaw)

	paramsFromFlags := options.New(valuesProviderFromFlags)

	missedRequirements := CheckRequirements(paramsFromFlags, reqs)
	if len(missedRequirements) == 0 {
		return valuesProviderFromFlags, nil
	}

	err = PromptRequiredValues(missedRequirements, paramsRaw, promptReader)
	if err != nil {
		return
	}

	return options.NewMapValuesProvider(paramsRaw), nil
}
