package config

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			Field:       "slice",
			ShortName:   "s",
			Help:        "shows slice",
			Default:     "a slice",
			Description: "shows slice",
			IsSecure:    false,
			IsRequired:  false,
			Type:        StringSliceRequirementType,
		},
	}

	cmd := &cobra.Command{}
	DefineCommandInputs(cmd, reqs)

	_, err1 := cmd.Flags().GetBool("verbose")
	assert.NoError(t, err1)

	_, err2 := cmd.Flags().GetString("color")
	assert.NoError(t, err2)

	_, err3 := cmd.Flags().GetStringSlice("slice")
	assert.NoError(t, err3)
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

func TestFlagValuesProviderWithStringSlice(t *testing.T) {
	fl := &pflag.FlagSet{}
	fl.StringP("somekey", "s", "test-default", "")
	fl.StringSliceP("someslice", "l", []string{}, "")

	flagValuesProv := &FlagValuesProvider{
		flags: fl,
	}

	val, found := flagValuesProv.Read("somekey")
	assert.True(t, found)
	assert.Equal(t, "test-default", val)

	err := fl.Parse([]string{"--somekey", "someval"})
	require.NoError(t, err)

	val2, found2 := flagValuesProv.Read("somekey")
	assert.True(t, found2)
	assert.Equal(t, "someval", val2)

	err = fl.Parse([]string{"--someslice", "someval1", "--someslice", "someval2"})
	require.NoError(t, err)

	val3, found3, err := flagValuesProv.ReadFlag("someslice", StringSliceRequirementType)
	assert.NoError(t, err)
	assert.True(t, found3)
	slice := val3.([]string)
	assert.Equal(t, "someval1", slice[0])
	assert.Equal(t, "someval2", slice[1])
}
