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
	err := os.Setenv("CONFIG_PATH", "lala")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv("CONFIG_PATH")
		if e != nil {
			logrus.Error(e)
		}
	}()

	assert.Equal(t, "lala", getConfigLocation())
}

func TestOSFileLocationFromHome(t *testing.T) {
	assert.Contains(t, getConfigLocation(), ".config/rportcli/config.json")
}

func TestReadConfig(t *testing.T) {
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

	err = os.Setenv("CONFIG_PATH", "config.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer func() {
		e := os.Unsetenv("CONFIG_PATH")
		if e != nil {
			logrus.Error(e)
		}
	}()

	cfg, err := GetConfig()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "someValue", cfg.ReadString("somekey", ""))
	assert.Equal(t, 1, cfg.ReadInt("one", 0))
}

func TestReadConfigError(t *testing.T) {
	_, err := GetConfig()
	assert.Error(t, err)
	if err == nil {
		return
	}
	assert.Contains(t, err.Error(), "config.json doesn't exist")
}
