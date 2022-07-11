package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	io2 "github.com/breathbath/go_utils/v2/pkg/io"
	"github.com/sirupsen/logrus"
)

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
