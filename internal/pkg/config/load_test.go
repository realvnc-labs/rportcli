package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/pflag"

	"github.com/stretchr/testify/mock"

	"github.com/spf13/cobra"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestOSFileLocationFromEnv(t *testing.T) {
	SetEnvVar(t, PathForConfigEnvVar, "lala")
	defer ResetEnvVar(t, PathForConfigEnvVar)

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

	WriteTestConfigFile(t, "config.json", rawJSON)
	defer RemoveTestConfigFile(t, "config.json")

	// required, otherwise an error will be returned by LoadParamsFromFileAndEnv
	SetEnvVar(t, APIURLEnvVar, "http://localhost:3000")
	defer ResetEnvVar(t, APIURLEnvVar)

	SetEnvVar(t, PathForConfigEnvVar, "config.json")
	defer ResetEnvVar(t, PathForConfigEnvVar)

	cfg, err := LoadParamsFromFileAndEnv(&pflag.FlagSet{})
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

	WriteTestConfigFile(t, filePath, rawJSON)
	defer RemoveTestConfigFile(t, filePath)

	envs := map[string]string{
		PathForConfigEnvVar: filePath,
		PasswordEnvVar:      "somepass",
		LoginEnvVar:         "log1",
	}

	for k, v := range envs {
		SetEnvVar(t, k, v)
	}

	defer func() {
		for k := range envs {
			ResetEnvVar(t, k)
		}
	}()

	cfg, err := LoadParamsFromFileAndEnv(&pflag.FlagSet{})
	assert.NoError(t, err)

	assert.Equal(t, "somepass", cfg.ReadString(Password, ""))
	assert.Equal(t, "log1", cfg.ReadString(Login, ""))

	assert.Equal(t, "https://10.10.10.11:3000", cfg.ReadString(ServerURL, ""))
}

func TestLoadEnvPreferredOverFile(t *testing.T) {
	rawJSON := []byte(`{"server":"https://10.10.10.11:3000"}`)
	filePath := "config123.json"

	WriteTestConfigFile(t, filePath, rawJSON)
	defer RemoveTestConfigFile(t, filePath)

	envs := map[string]string{
		PathForConfigEnvVar: filePath,
		ServerURLEnvVar:     "https://10.10.10.11:4000",
	}

	for k, v := range envs {
		SetEnvVar(t, k, v)
	}

	defer func() {
		for k := range envs {
			ResetEnvVar(t, k)
		}
	}()

	cfg, err := LoadParamsFromFileAndEnv(&pflag.FlagSet{})
	assert.NoError(t, err)

	assert.Equal(t, "https://10.10.10.11:4000", cfg.ReadString(ServerURL, ""))
}

func TestNoErrorWhenMissingConfigFile(t *testing.T) {
	// required, otherwise an error will be returned by LoadParamsFromFileAndEnv
	SetEnvVar(t, APIURLEnvVar, "http://localhost:3000")
	defer ResetEnvVar(t, APIURLEnvVar)

	SetEnvVar(t, PathForConfigEnvVar, "configNotExisting.json")
	defer ResetEnvVar(t, PathForConfigEnvVar)

	params, err := LoadParamsFromFileAndEnv(&pflag.FlagSet{})
	assert.NoError(t, err)

	assert.Equal(t, "", params.ReadString(Token, ""))
}

func TestLoadConfigErrorWhenNoAPIURL(t *testing.T) {
	rawJSON := []byte(`{"token":"1234"}`)
	filePath := "config1234.json"

	WriteTestConfigFile(t, filePath, rawJSON)
	defer RemoveTestConfigFile(t, filePath)

	SetEnvVar(t, PathForConfigEnvVar, filePath)
	defer ResetEnvVar(t, PathForConfigEnvVar)

	_, err := LoadParamsFromFileAndEnv(&pflag.FlagSet{})
	assert.ErrorIs(t, err, ErrAPIURLRequired)
}

func TestWriteConfig(t *testing.T) {
	filePath := "configToCheckAfter.json"

	SetEnvVar(t, PathForConfigEnvVar, filePath)
	defer ResetEnvVar(t, PathForConfigEnvVar)

	params := &options.ParameterBag{
		BaseValuesProvider: options.NewMapValuesProvider(map[string]interface{}{
			ServerURL: "http://localhost:3000",
			Token:     "123",
		}),
	}

	defer RemoveTestConfigFile(t, filePath)

	err := WriteConfig(params)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.FileExists(t, filePath)
	fileContents, err := ioutil.ReadFile(filePath)
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

type ValueProviderMock struct {
	mock.Mock
}

func (vpm *ValueProviderMock) Read(name string) (val interface{}, found bool) {
	args := vpm.Called(name)

	return args.Get(0), args.Bool(1)
}

func (vpm *ValueProviderMock) Dump(w io.Writer) (err error) {
	args := vpm.Called(w)

	return args.Error(0)
}

func (vpm *ValueProviderMock) ToKeyValues() map[string]interface{} {
	args := vpm.Called()

	return args.Get(0).(map[string]interface{})
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
		{
			Field:      "token",
			ShortName:  "o",
			Help:       "token given",
			IsSecure:   false,
			IsRequired: false,
		},
	}

	cmd := &cobra.Command{}
	DefineCommandInputs(cmd, reqs)

	vp := new(ValueProviderMock)
	vp.On("Read", "token").Return("tokVal", true)
	vp.On("Read", mock.Anything).Return("", false)

	params, err := CollectParamsFromCommandAndPromptAndEnv(cmd, reqs, prm, vp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualHost, _ := params.Read("host")
	assert.Equal(t, "127.1.1.1", actualHost)

	actualHTTPSParam, _ := params.Read("https")
	assert.Equal(t, false, actualHTTPSParam)

	actualPassword, _ := params.Read("password")
	assert.Equal(t, "123", actualPassword)

	actualToken, _ := params.Read("token")
	assert.Equal(t, "tokVal", actualToken)
}

func TestFlagValuesProvider(t *testing.T) {
	fl := &pflag.FlagSet{}
	fl.StringP("somekey", "s", "test-default", "")

	flagValuesProv := CreateFlagValuesProvider(fl)

	val, found := flagValuesProv.Read("somekey")
	assert.True(t, found)
	assert.Equal(t, "test-default", val)

	err := fl.Parse([]string{"--somekey", "someval"})
	require.NoError(t, err)

	val2, found2 := flagValuesProv.Read("somekey")
	assert.True(t, found2)
	assert.Equal(t, "someval", val2)

	actualKeyValues := flagValuesProv.ToKeyValues()
	assert.Equal(t, map[string]interface{}{"somekey": "someval"}, actualKeyValues)

	buf := &bytes.Buffer{}
	err = flagValuesProv.Dump(buf)
	require.NoError(t, err)
	assert.Equal(t, `{"somekey":"someval"}`+"\n", buf.String())
}

func SetEnvVar(t *testing.T, envVar, value string) {
	err := os.Setenv(envVar, value)
	assert.NoError(t, err)
}

func ResetEnvVar(t *testing.T, envVar string) {
	err := os.Unsetenv(envVar)
	assert.NoError(t, err)
}

func WriteTestConfigFile(t *testing.T, filePath string, rawJSON []byte) {
	err := ioutil.WriteFile(filePath, rawJSON, 0600)
	assert.NoError(t, err)
}

func RemoveTestConfigFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}
