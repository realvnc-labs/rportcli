package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredValidate(t *testing.T) {
	testCases := []struct {
		inputValue    interface{}
		fieldName     string
		expectedError string
	}{
		{
			inputValue:    nil,
			fieldName:     "someVal",
			expectedError: "value 'someVal' is required and cannot be empty",
		},
		{
			inputValue:    "",
			fieldName:     "someVal",
			expectedError: "value 'someVal' is required and cannot be empty",
		},
		{
			inputValue: "a",
			fieldName:  "someVal",
		},
	}

	for _, testCase := range testCases {
		actualErr := RequiredValidate(testCase.fieldName, testCase.inputValue)
		if testCase.expectedError == "" {
			assert.NoError(t, actualErr)
			continue
		}

		assert.EqualError(t, actualErr, testCase.expectedError)
	}
}

func TestCheckRequirementsAllMatched(t *testing.T) {
	requirementsToCheck := []ParameterRequirement{
		{
			Field:    "one",
			Validate: RequiredValidate,
		},
		{
			Field:    "two",
			Validate: RequiredValidate,
		},
		{
			Field:   "three",
			Default: "3",
		},
	}

	params := FromValues(map[string]string{
		"one": "1",
		"two": "2",
	})

	missedRequirements := CheckRequirements(params, requirementsToCheck)
	assert.Len(t, missedRequirements, 0)
}

func TestFromValues(t *testing.T) {
	config := FromValues(map[string]string{"one": "1", "two": "2"})
	assert.Equal(t, "1", config.ReadString("one", ""))
	assert.Equal(t, "2", config.ReadString("two", ""))
}

func TestCheckRequirementsMissed(t *testing.T) {
	requirementsToCheck := []ParameterRequirement{
		{
			Field:    "one",
			Validate: RequiredValidate,
		},
	}

	params := FromValues(map[string]string{})

	missedRequirements := CheckRequirements(params, requirementsToCheck)
	assert.Len(t, missedRequirements, 1)
	assert.Equal(t, "one", missedRequirements[0].Field)
}
