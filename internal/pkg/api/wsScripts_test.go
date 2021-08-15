package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildScriptsWsURL(t *testing.T) {
	testCases := []struct {
		inputURL      string
		expectedURL   string
		tokenValidity int
		tokenToGive   string
	}{
		{
			inputURL:      "http://scripts.url",
			expectedURL:   "ws://scripts.url/api/v1/ws/scripts?access_token=tok123",
			tokenValidity: 100,
			tokenToGive:   "tok123",
		},
		{
			inputURL:      "https://scripts.url",
			expectedURL:   "wss://scripts.url/api/v1/ws/scripts?access_token=tok1234",
			tokenValidity: 0,
			tokenToGive:   "tok1234",
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		urlProvider := &WsScriptsURLProvider{
			WsURLProvider: &WsURLProvider{
				BaseURL: testCase.inputURL,
				TokenProvider: func() (token string, err error) {
					token = tc.tokenToGive
					return
				},
				TokenValiditySeconds: testCase.tokenValidity,
			},
		}

		actualWsURL, err := urlProvider.BuildWsURL(context.Background())
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(t, testCase.expectedURL, actualWsURL)
	}
}
