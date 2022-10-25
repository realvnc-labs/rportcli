package config

import (
	"os"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/stretchr/testify/assert"
)

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
	fileContents, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	if err != nil {
		return
	}
	assert.Equal(t, `{"server":"http://localhost:3000","token":"123"}`+"\n", string(fileContents))
}
