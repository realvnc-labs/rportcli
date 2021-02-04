package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PromptReaderMock struct {
	ReadCount   int
	ReadOutputs []string
	ErrToGive   error
}

func (prm *PromptReaderMock) ReadString(delim byte) (string, error) {
	prm.ReadCount++

	if len(prm.ReadOutputs) < prm.ReadCount {
		return "", prm.ErrToGive
	}

	return prm.ReadOutputs[prm.ReadCount-1], prm.ErrToGive
}

func TestPromptRequiredValues(t *testing.T) {
	readerMock := &PromptReaderMock{
		ReadCount: 0,
		ReadOutputs: []string{
			"server",
			"log1",
			"pass1",
			"la",
		},
	}

	requirements := []ParameterRequirement{
		{
			Field:    "one",
			Validate: RequiredValidate,
		},
		{
			Field:    "two",
			Validate: RequiredValidate,
		},
		{
			Field:    "three",
			Validate: RequiredValidate,
		},
		{
			Field:   "four",
			Default: "Four value",
		},
	}

	actualKV := map[string]string{}
	err := PromptRequiredValues(requirements, actualKV, readerMock)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		map[string]string{
			"one":   "server",
			"three": "pass1",
			"two":   "log1",
			"four":  "la",
		},
		actualKV,
	)
}
