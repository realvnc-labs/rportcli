package config

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"

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

	cfg := LoadParamsFromFileAndEnv()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "someValue", cfg.ReadString("somekey", ""))
	assert.Equal(t, 1, cfg.ReadInt("one", 0))
}

func TestLoadConfigFromEnvOrFile(t *testing.T) {
	rawJSON := []byte(`{"server":"https://10.10.10.11:3000"}`)
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

	cfg := LoadParamsFromFileAndEnv()

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

	params := LoadParamsFromFileAndEnv()
	assert.Equal(t, "", params.ReadString(Token, ""))
}

func TestWriteConfig(t *testing.T) {
	err := os.Setenv(PathForConfigEnvVar, "configToCheckAfter.json")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	params := &options.ParameterBag{
		BaseValuesProvider: options.NewMapValuesProvider(map[string]interface{}{
			ServerURL: "http://localhost:3000",
			Token:     "123",
		}),
	}
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

	err = WriteConfig(params)
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
	assert.Equal(t, `{"server":"http://localhost:3000","token":"123"}`+"\n", string(fileContents))
}


func TestCommandPopulation(t *testing.T) {
	reqs := []ParameterRequirement{
		{
			Field:       "color",
			ShortName:   "c",
			Help:        "shows color",
			Default:     "red",
			Description: "shows color me",
			IsSecure:    false,
			IsRequired:  false,
			Type:        StringRequirementType,
		},
		{
			Field:       "verbose",
			ShortName:   "v",
			Help:        "shows verbose output",
			Default:     "1",
			Description: "shows verbose output",
			IsSecure:    false,
			IsRequired:  false,
			Type:        BoolRequirementType,
		},
	}

	cmd := &cobra.Command{}
	DefineCommandInputs(cmd, reqs)

	_, err1 := cmd.Flags().GetBool("verbose")
	assert.NoError(t, err1)

	_, err2 := cmd.Flags().GetString("color")
	assert.NoError(t, err2)
}

func TestCollectParams(t *testing.T) {
	prm := &PromptReaderMock{
		Inputs: []string{},
		ReadOutputs: []string{
			"127.1.1.1",
		},
		PasswordReadOutputs: []string{
			"123",
		},
	}

	reqs := []ParameterRequirement{
		{
			Field:      "host",
			ShortName:  "h",
			IsSecure:   false,
			IsRequired: true,
			Type:       StringRequirementType,
			Validate:   RequiredValidate,
		},
		{
			Field:      "https",
			ShortName:  "s",
			Help:       "uses https protocol",
			Default:    "0",
			IsSecure:   false,
			IsRequired: false,
			Type:       BoolRequirementType,
		},
		{
			Field:      "password",
			ShortName:  "p",
			Help:       "provides password",
			IsSecure:   true,
			IsRequired: true,
			Type:       StringRequirementType,
			Validate:   RequiredValidate,
		},
	}

	cmd := &cobra.Command{}
	DefineCommandInputs(cmd, reqs)

	params, err := CollectParamsFromCommandAndPrompt(cmd, reqs, prm)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualHost, _ := params.Read("host")
	assert.Equal(t, "127.1.1.1", actualHost)

	actualHttpsParam, _ := params.Read("https")
	assert.Equal(t, "false", actualHttpsParam)

	actualPassword, _ := params.Read("password")
	assert.Equal(t, "123", actualPassword)
}

