package config

import (
	"fmt"
	"io"
	"os"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/fs"
	"github.com/sirupsen/logrus"
)

type ValuesProvider struct {
	configFilePath     string
	envValuesProvider  options.ValuesProvider
	fileValuesProvider options.ValuesProvider
}

func NewValuesProvider(configFilePath string) *ValuesProvider {
	return &ValuesProvider{
		configFilePath:    configFilePath,
		envValuesProvider: &options.EnvValuesProvider{},
	}
}

func (cvp *ValuesProvider) Read(name string) (val interface{}, found bool) {
	var err error
	val, found = cvp.readFromEnv(name)
	if found {
		return val, found
	}

	if cvp.fileValuesProvider != nil {
		return cvp.fileValuesProvider.Read(name)
	}

	cvp.fileValuesProvider, err = cvp.createFileValuesProvider()
	if err != nil {
		logrus.Error(err)
	}
	if cvp.fileValuesProvider == nil {
		cvp.fileValuesProvider = options.NewMapValuesProvider(map[string]interface{}{})
	}

	return cvp.fileValuesProvider.Read(name)
}

func (cvp *ValuesProvider) readFromEnv(name string) (val interface{}, found bool) {
	envNameToRead := ""
	switch name {
	case Login:
		envNameToRead = LoginEnvVar
	case Password:
		envNameToRead = PasswordEnvVar
	case ServerURL:
		envNameToRead = ServerURLEnvVar
	}

	if envNameToRead == "" {
		return nil, false
	}

	return cvp.envValuesProvider.Read(envNameToRead)
}

func (cvp *ValuesProvider) Dump(w io.Writer) (err error) {
	return cvp.fileValuesProvider.Dump(w)
}

func (cvp *ValuesProvider) createFileValuesProvider() (options.ValuesProvider, error) {
	if !fs.FileExists(cvp.configFilePath) {
		logrus.Debugf("config file %s doesn't exist", cvp.configFilePath)
		return nil, nil
	}

	f, err := os.Open(cvp.configFilePath)
	if err != nil {
		err = fmt.Errorf("failed to open the file %s: %v", cvp.configFilePath, err)
		return nil, err
	}

	jvp, err := options.NewJSONValuesProvider(f)
	if err != nil {
		return nil, err
	}

	return jvp, nil
}
