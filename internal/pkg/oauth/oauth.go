package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetAuthProviderInfoURL = "/api/v1/auth/provider"
	MaxOAuthRetries        = 20
)

// AuthProviderInfo is used to tell the client where to get info for authorizing
type AuthProviderInfo struct {
	AuthProvider      string `json:"auth_provider"`
	SettingsURI       string `json:"settings_uri"`
	DeviceSettingsURI string `json:"device_settings_uri"`
}

// AuthProviderInfoResponse is the provider info response from the server
type AuthProviderInfoResponse struct {
	Data AuthProviderInfo `json:"data"`
}

// DeviceAuthInfo contains info for allowing authorization of a device/cli
type DeviceAuthInfo struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

// DeviceLoginInfo contains the auth info and which rport uri should be used for login
type DeviceLoginInfo struct {
	LoginURI string `json:"login_uri"`

	DeviceAuthInfo *DeviceAuthInfo `json:"auth_info"`
}

// DeviceAuthSettings contains the login info and which oauth provider is being used
type DeviceAuthSettings struct {
	AuthProvider string          `json:"auth_provider"`
	LoginInfo    DeviceLoginInfo `json:"details"`
}

// DeviceAuthSettingsResponse is the device auth setting/info response from the server
type DeviceAuthSettingsResponse struct {
	Data DeviceAuthSettings `json:"data"`
}

// DeviceLoginDetails contains the login details after a device login attempt
type DeviceLoginDetails struct {
	Token        string `json:"token"`
	StatusCode   int    `json:"status_code"`
	ErrorCode    string `json:"error"`
	ErrorMessage string `json:"error_description"`
	ErrorURI     string `json:"error_uri"`
}

// DeviceLoginDetailsResponse is the login details response from the server
type DeviceLoginDetailsResponse struct {
	Data DeviceLoginDetails `json:"data"`
}

// GetAuthProviderInfo gets the provider info from the rportd server
func GetAuthProviderInfo(
	ctx context.Context,
	baseURL string) (providerInfo *AuthProviderInfo, statusCode int, err error) {
	providerInfoRes, err := doGet(ctx, baseURL+GetAuthProviderInfoURL)
	if err != nil {
		return nil, 0, err
	}
	defer providerInfoRes.Body.Close()

	if providerInfoRes.StatusCode != http.StatusOK {
		return nil, providerInfoRes.StatusCode, fmt.Errorf("unable to get auth provider info: %d", providerInfoRes.StatusCode)
	}

	var providerInfoResp AuthProviderInfoResponse
	if err := json.NewDecoder(providerInfoRes.Body).Decode(&providerInfoResp); err != nil {
		return nil, http.StatusOK, err
	}

	return &providerInfoResp.Data, http.StatusOK, nil
}

// GetDeviceAuthSettings gets the details to be used for an authorization from the server.
func GetDeviceAuthSettings(
	ctx context.Context,
	baseURL string,
	providerInfo *AuthProviderInfo,
) (authSettings *DeviceAuthSettings, statusCode int, err error) {
	loginInfoRes, err := doGet(ctx, baseURL+providerInfo.DeviceSettingsURI)
	if err != nil {
		return nil, 0, err
	}
	defer loginInfoRes.Body.Close()

	if loginInfoRes.StatusCode != http.StatusOK {
		return nil, loginInfoRes.StatusCode, fmt.Errorf("unable to get auth login info: %d", loginInfoRes.StatusCode)
	}

	var loginInfoResponse DeviceAuthSettingsResponse
	if err := json.NewDecoder(loginInfoRes.Body).Decode(&loginInfoResponse); err != nil {
		return nil, 0, err
	}

	return &loginInfoResponse.Data, http.StatusOK, nil
}

// GetDeviceLogin will atempt to authorize with the rportd server. The login attempt may
// fail if the user hasn't completed authorization. The login can be retried.
func GetDeviceLogin(ctx context.Context, baseURL string, loginURI string, deviceCode string, tokenLifetime int) (
	loginResponse *DeviceLoginDetails, statusCode int, err error) {
	loginReq, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+loginURI, nil)
	if err != nil {
		return nil, 0, err
	}

	params := url.Values{
		"device_code":    {deviceCode},
		"token-lifetime": {strconv.Itoa(tokenLifetime)},
	}
	loginReq.URL.RawQuery = params.Encode()

	loginReq.Header.Set("accept", "application/json")

	loginRes, err := http.DefaultClient.Do(loginReq)
	if err != nil {
		return nil, 0, err
	}
	defer loginRes.Body.Close()

	var loginDetailsResponse DeviceLoginDetailsResponse
	if err := json.NewDecoder(loginRes.Body).Decode(&loginDetailsResponse); err != nil {
		return nil, loginRes.StatusCode, err
	}

	return &loginDetailsResponse.Data, loginRes.StatusCode, nil
}

// doGet is a simple GET helper function
func doGet(ctx context.Context, targetURL string) (res *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	return http.DefaultClient.Do(req)
}
