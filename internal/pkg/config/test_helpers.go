package config

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

func SetEnvVar(t *testing.T, envVar, value string) {
	err := os.Setenv(envVar, value)
	assert.NoError(t, err)
}

func ResetEnvVar(t *testing.T, envVar string) {
	err := os.Unsetenv(envVar)
	assert.NoError(t, err)
}

func SetEnvSet(t *testing.T, envs map[string]string) {
	for k, v := range envs {
		SetEnvVar(t, k, v)
	}
}

func ResetEnvSet(t *testing.T, envs map[string]string) {
	for k := range envs {
		ResetEnvVar(t, k)
	}
}

func SetCLIFlagString(t *testing.T, fl *pflag.FlagSet, flagName, flagVal string) (err error) {
	err = fl.Parse([]string{fmt.Sprintf("--%s", flagName), flagVal})
	require.NoError(t, err)
	return err
}

func WriteTestConfigFile(t *testing.T, filePath string, rawJSON []byte) {
	err := os.WriteFile(filePath, rawJSON, 0600)
	assert.NoError(t, err)
}

func RemoveTestConfigFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}
