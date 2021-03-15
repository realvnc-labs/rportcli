package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOSFileLocationFromEnv(t *testing.T) {
	err := os.Setenv(PathForConfigEnvVar, "lala")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv(PathForConfigEnvVar)
		if e != nil {
			logrus.Error(e)
		}
	}()

	assert.Equal(t, "lala", getConfigLocation())
}

func TestOSFileLocationFromHome(t *testing.T) {
	assert.Contains(t, getConfigLocation(), ".config/rportcli/config.json")
}

func TestLoadConfigFromFile(t *testing.T) {
	config := map[string]interface{}{
		"somekey": "someValue",
		"one":     1,
	}
	rawJSON, err := json.Marshal(config)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("config.json", rawJSON, 0600)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		e := os.Remove("config.json")
		if e != nil {
			logrus.Error(e)
		}
	}()

	err = os.Setenv(PathForConfigEnvVar, "config.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv(PathForConfigEnvVar)
		if e != nil {
			logrus.Error(e)
		}
	}()

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "someValue", cfg.ReadString("somekey", ""))
	assert.Equal(t, 1, cfg.ReadInt("one", 0))
}

func TestLoadConfigFromEnvOrFile(t *testing.T) {
	rawJSON := []byte(`{"server_url":"https://10.10.10.11:3000"}`)
	filePath := "config123.json"

	err := ioutil.WriteFile(filePath, rawJSON, 0600)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		e := os.Remove(filePath)
		if e != nil {
			logrus.Error(e)
		}
	}()

	envs := map[string]string{
		PathForConfigEnvVar: filePath,
		PasswordEnvVar:      "somepass",
		LoginEnvVar:         "log1",
	}

	for k, v := range envs {
		err = os.Setenv(k, v)
		assert.NoError(t, err)
		if err != nil {
			return
		}
	}

	defer func() {
		for k := range envs {
			e := os.Unsetenv(k)
			if e != nil {
				logrus.Error(e)
			}
		}
	}()

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "somepass", cfg.ReadString(Password, ""))
	assert.Equal(t, "log1", cfg.ReadString(Login, ""))
	assert.Equal(t, "https://10.10.10.11:3000", cfg.ReadString(ServerURL, ""))
}

func TestLoadConfigFromFileError(t *testing.T) {
	err := os.Setenv(PathForConfigEnvVar, "configNotExisting.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv(PathForConfigEnvVar)
		if e != nil {
			logrus.Error(e)
		}
	}()

	_, err = LoadConfig()
	assert.NoError(t, err)
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	assert.Equal(t, DefaultServerURL, config.ReadString(ServerURL, ""))
	assert.Equal(t, "", config.ReadString(Login, ""))
	assert.Equal(t, "", config.ReadString(Password, ""))
}

func TestWriteConfig(t *testing.T) {
	config := GetDefaultConfig()

	err := os.Setenv(PathForConfigEnvVar, "configToCheckAfter.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv(PathForConfigEnvVar)
		if e != nil {
			logrus.Error(e)
		}
	}()

	err = WriteConfig(config)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Remove("configToCheckAfter.json")
		if e != nil {
			logrus.Error(e)
		}
	}()
	assert.FileExists(t, "configToCheckAfter.json")
	fileContents, err := ioutil.ReadFile("configToCheckAfter.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, `{"login":"","password":"","server_url":"http://localhost:3000"}`+"\n", string(fileContents))
}
