package auth_test

import (
	"context"
	"net/http"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/auth"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type ParameterBagInput map[string]interface{}

func getTestRequest(t *testing.T) (req *http.Request) {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), "post", "/", nil)
	assert.NoError(t, err)
	return req
}

func TestBasicAuthStrategy(t *testing.T) {
	cases := []struct {
		name           string
		params         ParameterBagInput
		expectedResult string
	}{
		{
			name: "WithRegularCredentials",
			params: ParameterBagInput{
				config.APIUser:     "admin",
				config.APIPassword: "foobaz",
			},
			expectedResult: "Basic YWRtaW46Zm9vYmF6",
		},
		{
			name: "WithApiToken",
			params: ParameterBagInput{
				config.APIUser:  "admin",
				config.APIToken: "12345678901234567890",
			},
			expectedResult: "Basic YWRtaW46MTIzNDU2Nzg5MDEyMzQ1Njc4OTA=",
		},
		{
			name: "WithLegacyCredentials",
			params: ParameterBagInput{
				config.Login:    "admin",
				config.Password: "foobaz",
			},
			expectedResult: "Basic YWRtaW46Zm9vYmF6",
		},
		{
			name: "PreferApiCredentials",
			params: ParameterBagInput{
				config.APIUser:     "admin",
				config.APIPassword: "foobaz",
				config.Login:       "admin1",
				config.Password:    "foobaz1",
			},
			expectedResult: "Basic YWRtaW46Zm9vYmF6",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := getTestRequest(t)
			params := &options.ParameterBag{
				BaseValuesProvider: options.NewMapValuesProvider(tc.params),
			}

			authStrategy := auth.GetAuthStrategy(params)
			primaryAuth := authStrategy.PrimaryAuth
			assert.NotNil(t, primaryAuth)

			err := primaryAuth.AuthRequest(req)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedResult, req.Header.Get("Authorization"))
		})
	}
}

func TestBasicAuthStrategyErrors(t *testing.T) {
	cases := []struct {
		name        string
		params      ParameterBagInput
		expectedErr error
	}{
		{
			name: "WhenApiTokenAndPassword",
			params: ParameterBagInput{
				config.APIUser:     "admin",
				config.APIPassword: "foobaz",
				config.APIToken:    "1234",
			},
			expectedErr: utils.ErrAPIPasswordAndAPITokenAreBothSet,
		},
		{
			name: "WhenApiTokenAndLegacyPassword",
			params: ParameterBagInput{
				config.APIUser:  "admin",
				config.Password: "foobaz",
				config.APIToken: "1234",
			},
			expectedErr: utils.ErrAPIPasswordAndAPITokenAreBothSet,
		},
		{
			name: "WhenApiTokenAndLegacyUser",
			params: ParameterBagInput{
				config.Login:       "admin",
				config.APIPassword: "foobaz",
				config.APIToken:    "1234",
			},
			expectedErr: utils.ErrAPIPasswordAndAPITokenAreBothSet,
		},
		{
			name: "WhenApiTokenAndLegacyUserAndLegacyPassword",
			params: ParameterBagInput{
				config.Login:    "admin",
				config.Password: "foobaz",
				config.APIToken: "1234",
			},
			expectedErr: utils.ErrAPIPasswordAndAPITokenAreBothSet,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := getTestRequest(t)
			params := &options.ParameterBag{
				BaseValuesProvider: options.NewMapValuesProvider(tc.params),
			}

			authStategy := auth.GetAuthStrategy(params)
			primaryAuth := authStategy.PrimaryAuth
			assert.NotNil(t, primaryAuth)

			err := primaryAuth.AuthRequest(req)
			if assert.Error(t, err) {
				assert.Equal(t, tc.expectedErr, err)
			}
		})
	}
}
