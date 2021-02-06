package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRemovePortFromUrl(t *testing.T) {
	testCases := []struct {
		Input          string
		ExpectedOutput string
	}{
		{
			Input:          "https://localhost:8182?aaa=:8182",
			ExpectedOutput: "https://localhost?aaa=:8182",
		},
		{
			Input:          "127.22.22.33:8182",
			ExpectedOutput: "127.22.22.33",
		},
		{
			Input:          "https://ya.ru:8182",
			ExpectedOutput: "https://ya.ru",
		},
		{
			Input:          "http://ya.ru?a=127&b=:33",
			ExpectedOutput: "http://ya.ru?a=127&b=:33",
		},
	}

	for _, testCase := range testCases {
		actualURL := RemovePortFromURL(testCase.Input)
		assert.Equal(t, testCase.ExpectedOutput, actualURL)
	}
}
