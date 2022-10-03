package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/oauth"
)

func TestShouldGetAuthProviderInfo(t *testing.T) {
	expectedInfoResponse := oauth.AuthProviderInfoResponse{
		Data: oauth.AuthProviderInfo{
			AuthProvider:      "google",
			SettingsURI:       "/ext/oauth/settings",
			DeviceSettingsURI: "/ext/oauth/settings/device",
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, oauth.GetAuthProviderInfoURL, r.URL.String())
		writeExpectedResponse(t, w, http.StatusOK, expectedInfoResponse)
	}))
	defer s.Close()

	ctx := context.Background()
	cl := api.New(s.URL, nil)

	authInfo, statusCode, err := oauth.GetAuthProviderInfo(ctx, cl.BaseURL)
	require.NoError(t, err)
	require.NotNil(t, authInfo)
	assert.Equal(t, http.StatusOK, statusCode)

	assert.Equal(t, expectedInfoResponse.Data.AuthProvider, authInfo.AuthProvider)
	assert.Equal(t, expectedInfoResponse.Data.SettingsURI, authInfo.SettingsURI)
	assert.Equal(t, expectedInfoResponse.Data.DeviceSettingsURI, authInfo.DeviceSettingsURI)
}

func TestShouldErrorWhenGetAuthProviderInfoFails(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeExpectedResponse(t, w, http.StatusInternalServerError, nil)
	}))
	defer s.Close()

	ctx := context.Background()
	cl := api.New(s.URL, nil)

	_, statusCode, err := oauth.GetAuthProviderInfo(ctx, cl.BaseURL)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Error(t, err)
}

func TestShouldGetAuthSettingsInfo(t *testing.T) {
	providerInfo := &oauth.AuthProviderInfo{
		AuthProvider:      "google",
		SettingsURI:       "/api/v1/ext/oauth/settings",
		DeviceSettingsURI: "/api/v1/ext/oauth/settings/device",
	}
	expectedLoginInfoResponse := oauth.DeviceAuthSettingsResponse{
		Data: oauth.DeviceAuthSettings{
			LoginInfo: oauth.DeviceLoginInfo{
				LoginURI: "/api/v1/ext/oauth/device/login",
				DeviceAuthInfo: &oauth.DeviceAuthInfo{
					UserCode:        "1234",
					DeviceCode:      "1234",
					VerificationURI: "1234",
					ExpiresIn:       333,
					Interval:        4,
					Message:         "1234",
				},
			},
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, providerInfo.DeviceSettingsURI, r.URL.String())
		writeExpectedResponse(t, w, http.StatusOK, expectedLoginInfoResponse)
	}))
	defer s.Close()

	ctx := context.Background()
	cl := api.New(s.URL, nil)

	authSettings, statusCode, err := oauth.GetDeviceAuthSettings(ctx, cl.BaseURL, providerInfo)
	require.NoError(t, err)

	loginInfo := authSettings.LoginInfo
	require.NotNil(t, loginInfo)
	assert.Equal(t, http.StatusOK, statusCode)

	assert.Equal(t, expectedLoginInfoResponse.Data.LoginInfo.LoginURI, loginInfo.LoginURI)
	assert.Equal(t, expectedLoginInfoResponse.Data.LoginInfo.DeviceAuthInfo.UserCode, loginInfo.DeviceAuthInfo.UserCode)
}

func TestShouldErrorWhenGetAuthSettingsFails(t *testing.T) {
	providerInfo := &oauth.AuthProviderInfo{
		AuthProvider:      "google",
		SettingsURI:       "/ext/oauth/settings",
		DeviceSettingsURI: "/ext/oauth/settings/device",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeExpectedResponse(t, w, http.StatusInternalServerError, nil)
	}))
	defer s.Close()

	ctx := context.Background()
	cl := api.New(s.URL, nil)

	_, statusCode, err := oauth.GetDeviceAuthSettings(ctx, cl.BaseURL, providerInfo)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Error(t, err)
}

func TestShouldLogin(t *testing.T) {
	cases := []struct {
		name          string
		loginResponse *oauth.DeviceLoginDetailsResponse
		statusCode    int
	}{
		{
			name: "happy path",
			loginResponse: &oauth.DeviceLoginDetailsResponse{
				Data: oauth.DeviceLoginDetails{
					Token: "1234",
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "pending error with OK status",
			loginResponse: &oauth.DeviceLoginDetailsResponse{
				Data: oauth.DeviceLoginDetails{
					Token:        "",
					ErrorCode:    "pending error",
					ErrorMessage: "there was a pending error",
					ErrorURI:     "http://error-uri.com",
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "bad device code with not OK status",
			loginResponse: &oauth.DeviceLoginDetailsResponse{
				Data: oauth.DeviceLoginDetails{
					Token:        "",
					ErrorCode:    "bad device code error",
					ErrorMessage: "there was a bad device code error",
					ErrorURI:     "http://error-uri.com",
				},
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/oauth/login/device", r.URL.EscapedPath())
				assert.Equal(t, r.URL.Query().Get("device_code"), "123456")
				writeExpectedResponse(t, w, tc.statusCode, tc.loginResponse)
			}))
			defer s.Close()

			ctx := context.Background()
			cl := api.New(s.URL, nil)

			loginResponse, statusCode, err := oauth.GetDeviceLogin(ctx, cl.BaseURL, "/oauth/login/device", "123456")

			if statusCode == http.StatusOK {
				require.NoError(t, err)
				require.NotNil(t, loginResponse)
			}

			if loginResponse != nil {
				if loginResponse.ErrorCode == "" {
					assert.Equal(t, http.StatusOK, statusCode)
					assert.Equal(t, tc.loginResponse.Data.Token, loginResponse.Token)
				} else {
					assert.Equal(t, tc.statusCode, statusCode)
					assert.Equal(t, tc.loginResponse.Data.ErrorCode, loginResponse.ErrorCode)
					assert.Equal(t, tc.loginResponse.Data.ErrorMessage, loginResponse.ErrorMessage)
					assert.Equal(t, tc.loginResponse.Data.ErrorURI, loginResponse.ErrorURI)
				}
			}
		})
	}
}

func writeExpectedResponse(t *testing.T, w http.ResponseWriter, expectedStatusCode int, expectedResponse interface{}) {
	t.Helper()

	if expectedResponse == nil {
		w.WriteHeader(expectedStatusCode)
		return
	}

	jsonBytes, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	w.WriteHeader(expectedStatusCode)
	n, err := w.Write(jsonBytes)

	require.NotZero(t, n)
	require.NoError(t, err)
}
