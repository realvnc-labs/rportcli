package config

import (
	"encoding/json"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/spf13/pflag"

	"github.com/stretchr/testify/require"

	"github.com/spf13/cobra"

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

	SetEnvSet(t, envs)
	defer ResetEnvSet(t, envs)

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

	SetEnvSet(t, envs)
	defer ResetEnvSet(t, envs)

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
			IsRequired: false, // covered by separate test
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
			IsRequired: false, // covered by separate test
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

	fl := cmd.Flags()
	fp := &FlagValuesProvider{
		flags: fl,
	}

	err := SetCLIFlagString(t, fl, "token", "tokVal")
	require.NoError(t, err)

	params, err := CollectParamsFromCommandAndPromptAndEnv(fp, reqs, prm)
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

func TestErrOnMultipleTargetingOptions(t *testing.T) {
	cases := []struct {
		Name            string
		CIDValue        string
		GIDValue        string
		ClientNameValue string
		ShouldErr       bool
		ErrValue        error
	}{
		{
			Name:            "CIDs Only",
			CIDValue:        "1234",
			GIDValue:        "",
			ClientNameValue: "",
			ShouldErr:       false,
			ErrValue:        nil,
		},
		{
			Name:            "GIDs Only",
			CIDValue:        "",
			GIDValue:        "1234",
			ClientNameValue: "",
			ShouldErr:       false,
			ErrValue:        nil,
		},
		{
			Name:            "Name Only",
			CIDValue:        "",
			GIDValue:        "",
			ClientNameValue: "Name",
			ShouldErr:       false,
			ErrValue:        nil,
		},
		{
			Name:            "Error on CIDs and GIDs",
			CIDValue:        "1234",
			GIDValue:        "4567",
			ClientNameValue: "",
			ShouldErr:       true,
			ErrValue:        ErrMultipleTargetingOptions,
		},
		{
			Name:            "Error on CIDs and Name",
			CIDValue:        "1234",
			GIDValue:        "",
			ClientNameValue: "4567",
			ShouldErr:       true,
			ErrValue:        ErrMultipleTargetingOptions,
		},
		{
			Name:            "Error on GIDs and Name",
			CIDValue:        "",
			GIDValue:        "1234",
			ClientNameValue: "4567",
			ShouldErr:       true,
			ErrValue:        ErrMultipleTargetingOptions,
		},
		{
			Name:            "Error on All Set",
			CIDValue:        "1234",
			GIDValue:        "4567",
			ClientNameValue: "4321",
			ShouldErr:       true,
			ErrValue:        ErrMultipleTargetingOptions,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			prm := &PromptReaderMock{
				Inputs: []string{},
				ReadOutputs: []string{
					"127.1.1.1",
				},
				PasswordReadOutputs: []string{
					"123",
				},
			}

			cmd := &cobra.Command{}
			reqs := getMultipleTargetFlagSpecs()
			DefineCommandInputs(cmd, reqs)

			fl := cmd.Flags()
			fp := &FlagValuesProvider{
				flags: fl,
			}

			if tc.CIDValue != "" {
				err := SetCLIFlagString(t, fl, ClientIDs, tc.CIDValue)
				require.NoError(t, err)
			}

			if tc.GIDValue != "" {
				err := SetCLIFlagString(t, fl, GroupIDs, tc.GIDValue)
				require.NoError(t, err)
			}

			if tc.ClientNameValue != "" {
				err := SetCLIFlagString(t, fl, ClientNameFlag, tc.ClientNameValue)
				require.NoError(t, err)
			}

			_, err := CollectParamsFromCommandAndPromptAndEnv(fp, reqs, prm)
			if tc.ShouldErr {
				assert.ErrorIs(t, err, ErrMultipleTargetingOptions)
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}
		})
	}
}

func TestCheckRequiredParams(t *testing.T) {
	cases := []struct {
		Name          string
		IncludeParams []string
		ShouldErr     bool
		ErrText       string
	}{
		{
			Name:          "Missing Mandatory Params",
			IncludeParams: []string{},
			ShouldErr:     true,
			ErrText:       "missing",
		},
		{
			Name:          "Include Some Mandatory Params",
			IncludeParams: []string{"--cids"},
			ShouldErr:     true,
			ErrText:       "missing",
		},
		{
			Name:          "Include All Mandatory Params",
			IncludeParams: []string{"--cids", "--command"},
			ShouldErr:     false,
			ErrText:       "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			prm := &PromptReaderMock{
				Inputs: []string{},
				ReadOutputs: []string{
					"127.1.1.1",
				},
				PasswordReadOutputs: []string{
					"123",
				},
			}

			cmd := &cobra.Command{}
			reqs := GetCommandFlagSpecs()
			DefineCommandInputs(cmd, reqs)

			fl := cmd.Flags()
			fp := &FlagValuesProvider{
				flags: fl,
			}

			if len(tc.IncludeParams) > 0 {
				for _, p := range tc.IncludeParams {
					err := fl.Parse([]string{p, "1234"})
					require.NoError(t, err)
				}
			}

			_, err := CollectParamsFromCommandAndPromptAndEnv(fp, reqs, prm)
			if tc.ShouldErr {
				assert.NotNil(t, err)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ErrText)
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				return
			}
		})
	}
}

func TestCheckMultipleYAMLFiles(t *testing.T) {
	cmd := &cobra.Command{}
	reqs := GetCommandFlagSpecs()
	DefineCommandInputs(cmd, reqs)

	fl := cmd.Flags()
	fp := &FlagValuesProvider{
		flags: fl,
	}

	err := SetCLIFlagString(t, fl, ReadYAML, "../../../testdata/test1-ok.yaml")
	require.NoError(t, err)

	err = SetCLIFlagString(t, fl, ReadYAML, "../../../testdata/test3-ok.yaml")
	require.NoError(t, err)

	vp, err := CollectParamsFromCommandAndPromptAndEnv(fp, reqs, nil)
	assert.NoError(t, err)

	params := options.New(vp)

	assert.False(t, params.ReadBool(IsFullOutput, false))
	assert.Equal(t, params.ReadString(Command, ""), "pwd")

	cids, found := params.Read(ClientIDs, []string{})
	assert.True(t, found)
	assert.Equal(t, cids, "differentserver")
}

func TestShouldPreferCLIToYAML(t *testing.T) {
	cmd := &cobra.Command{}
	reqs := GetCommandFlagSpecs()
	DefineCommandInputs(cmd, reqs)

	fl := cmd.Flags()
	fp := &FlagValuesProvider{
		flags: fl,
	}

	err := SetCLIFlagString(t, fl, ReadYAML, "../../../testdata/test1-ok.yaml")
	require.NoError(t, err)

	err = SetCLIFlagString(t, fl, ClientIDs, "anotherserver")
	require.NoError(t, err)

	vp, err := CollectParamsFromCommandAndPromptAndEnv(fp, reqs, nil)
	assert.NoError(t, err)

	params := options.New(vp)

	assert.True(t, params.ReadBool(IsFullOutput, false))
	assert.Equal(t, params.ReadString("command", ""), "ls")

	cids, found := params.Read("cids", []string{})
	assert.True(t, found)
	assert.Equal(t, cids, "anotherserver")
}

func getMultipleTargetFlagSpecs() (flagSpecs []ParameterRequirement) {
	return []ParameterRequirement{
		GetNoPromptFlagSpec(),
		GetReadYAMLFlagSpec(),
		{
			Field: ClientIDs,
			Help:  "Enter comma separated client IDs",
			Description: "[required] Comma separated client ids for which the command should be executed. " +
				"Alternatively use -n to execute a command by client name(s), or use --search flag.",
			ShortName: "d",
		},
		{
			Field:       ClientNameFlag,
			Description: "Comma separated client names for which the command should be executed",
			ShortName:   "n",
		},
		{
			Field:       Command,
			Help:        "Enter command",
			Description: "[required] Command which should be executed on the clients",
			ShortName:   "c",
		},
		{
			Field:       GroupIDs,
			Help:        "Enter comma separated group IDs",
			Description: "Comma separated client group IDs",
			ShortName:   "g",
		},
	}
}
