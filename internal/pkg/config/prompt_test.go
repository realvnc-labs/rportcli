package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PromptReaderMock struct {
	ReadCount           int
	PasswordReadCount   int
	ReadOutputs         []string
	PasswordReadOutputs []string
	ErrToGive           error
	Inputs              []string
}

func (prm *PromptReaderMock) ReadString() (string, error) {
	prm.ReadCount++

	if len(prm.ReadOutputs) < prm.ReadCount {
		return "", prm.ErrToGive
	}

	return prm.ReadOutputs[prm.ReadCount-1], prm.ErrToGive
}

func (prm *PromptReaderMock) ReadPassword() (string, error) {
	prm.PasswordReadCount++

	if len(prm.PasswordReadOutputs) < prm.PasswordReadCount {
		return "", prm.ErrToGive
	}

	return prm.PasswordReadOutputs[prm.PasswordReadCount-1], prm.ErrToGive
}

func (prm *PromptReaderMock) Output(text string) {
	prm.Inputs = append(prm.Inputs, text)
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

	actualKV := map[string]interface{}{}
	err := PromptRequiredValues(requirements, actualKV, readerMock)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		map[string]interface{}{
			"one":   "server",
			"three": "pass1",
			"two":   "log1",
			"four":  "la",
		},
		actualKV,
	)
}

func TestPromptPassword(t *testing.T) {
	readerMock := &PromptReaderMock{
		PasswordReadCount:   0,
		PasswordReadOutputs: []string{"123"},
	}

	requirements := []ParameterRequirement{
		{
			Field:    "password",
			Validate: RequiredValidate,
			IsSecure: true,
		},
	}

	actualKV := map[string]interface{}{}
	err := PromptRequiredValues(requirements, actualKV, readerMock)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		map[string]interface{}{
			"password": "123",
		},
		actualKV,
	)
}
