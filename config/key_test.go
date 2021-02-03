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
