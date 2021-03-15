package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	options "github.com/breathbath/go_utils/utils/config"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/sirupsen/logrus"
)

const (
	defaultPath         = ".config/rportcli/config.json"
	ServerURL           = "server_url"
	Login               = "login"
	Password            = "password"
	DefaultServerURL    = "http://localhost:3000"
	PathForConfigEnvVar = "CONFIG_PATH"
	LoginEnvVar         = "RPORT_USER"
	PasswordEnvVar      = "RPORT_PASSWORD"
	ServerURLEnvVar     = "RPORT_SERVER_URL"
)

func LoadConfig() (params *options.ParameterBag, err error) {
	configLocation := getConfigLocation()

	return &options.ParameterBag{
		BaseValuesProvider: NewValuesProvider(configLocation),
	}, nil
}

func AuthConfigProvider() (login, pass string, err error) {
	login, pass = Params.ReadString(Login, ""), Params.ReadString(Password, "")
	return
}

func GetDefaultConfig() (params *options.ParameterBag) {
	vp := options.NewMapValuesProvider(map[string]interface{}{
		ServerURL: DefaultServerURL,
		Password:  "",
		Login:     "",
	})

	return &options.ParameterBag{BaseValuesProvider: vp}
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

	fileToWrite, err := os.OpenFile(configLocation, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	err = params.BaseValuesProvider.Dump(fileToWrite)
	if err != nil {
		return err
	}

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
