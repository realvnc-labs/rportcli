package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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

func LoadParamsFromFileAndEnv() (params *options.ParameterBag, err error) {
	envMapValues := map[string]interface{}{
		Password:            env.ReadEnv(PasswordEnvVar, ""),
		Login:               env.ReadEnv(LoginEnvVar, ""),
		ServerURL:           env.ReadEnv(ServerURLEnvVar, ""),
		PathForConfigEnvVar: env.ReadEnv(PathForConfigEnvVar, ""),
		CacheValidityEnvVar: env.ReadEnv(CacheValidityEnvVar, ""),
		CacheFolderEnvVar:   env.ReadEnv(CacheFolderEnvVar, ""),
	}

	envValuesProvider := options.NewMapValuesProvider(envMapValues)

	configFilePath := getConfigLocation()
	if !fs.FileExists(configFilePath) {
		logrus.Debugf("config file %s doesn't exist", configFilePath)
		return options.New(envValuesProvider), nil
	}

	f, err := os.Open(configFilePath)
	if err != nil {
		err = fmt.Errorf("failed to open the file %s: %v", configFilePath, err)
		return nil, err
	}
	jvp, err := options.NewJSONValuesProvider(f)
	if err != nil {
		return nil, err
	}

	valuesProvider := options.NewValuesProviderComposite(jvp, envValuesProvider)

	return options.New(valuesProvider), nil
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

	fileToWrite, err := os.OpenFile(configLocation, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(fileToWrite)
	err = encoder.Encode(configToWrite)

	logrus.Infof("created config at %s", configLocation)

	return
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
