package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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

	params, err := CollectParams(cmd, reqs, prm)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "127.1.1.1", params.ReadString("host", ""))
	assert.False(t, params.ReadBool("https", true))
	assert.Equal(t, "123", params.ReadString("password", ""))
}
