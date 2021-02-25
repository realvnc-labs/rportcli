package api

import (
	"context"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

type TokenProviderMock struct {
	loginResp     LoginResponse
	tokenLifetime int
	errToGive     error
}

func (tpm *TokenProviderMock) GetToken(ctx context.Context, tokenLifetime int) (li LoginResponse, err error) {
	tpm.tokenLifetime = tokenLifetime
	return tpm.loginResp, tpm.errToGive
}

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
		tokenProv := &TokenProviderMock{
			loginResp: LoginResponse{
				Data: models.Token{
					Token: testCase.tokenToGive,
				},
			},
		}
		urlProvider := &WsCommandURLProvider{
			BaseURL:              testCase.inputURL,
			TokenProvider:        tokenProv,
			TokenValiditySeconds: testCase.tokenValidity,
		}

		actualWsURL, err := urlProvider.BuildWsURL(context.Background())
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(t, testCase.expectedURL, actualWsURL)
		assert.Equal(t, testCase.tokenValidity, tokenProv.tokenLifetime)
	}
}
