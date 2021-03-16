package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildWsURL(t *testing.T) {
	testCases := []struct {
		inputURL      string
		expectedURL   string
		tokenValidity int
		tokenToGive   string
	}{
		{
			inputURL:      "http://some.url",
			expectedURL:   "ws://some.url/api/v1/ws/commands?access_token=sometoken",
			tokenValidity: 100,
			tokenToGive:   "sometoken",
		},
		{
			inputURL:      "https://some.url",
			expectedURL:   "wss://some.url/api/v1/ws/commands?access_token=someothertoken",
			tokenValidity: 0,
			tokenToGive:   "someothertoken",
		},
	}
	for _, testCase := range testCases {
		urlProvider := &WsCommandURLProvider{
			BaseURL:              testCase.inputURL,
			TokenProvider: func() (token string, err error) {
				token = testCase.tokenToGive
				return
			},
			TokenValiditySeconds: testCase.tokenValidity,
		}

		actualWsURL, err := urlProvider.BuildWsURL(context.Background())
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(t, testCase.expectedURL, actualWsURL)
	}
}
