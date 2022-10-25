package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/oauth"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	LoginURL                    = "/api/v1/login"
	LogoutURL                   = "/api/v1/logout"
	TwoFaURL                    = "/api/v1/verify-2fa"
	MeURL                       = "/api/v1/me"
	MeIPURL                     = "/api/v1/me/ip"
	StatusURL                   = "/api/v1/status"
	DefaultTokenValiditySeconds = 20 * 24 * 60 * 60 // 30 days is max value
)

type LoginResponse struct {
	Data models.Token `json:"data"`
}

func (rp *Rport) GetToken(ctx context.Context, tokenLifetime int) (li LoginResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, LoginURL),
		nil,
	)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("token-lifetime", strconv.Itoa(tokenLifetime))
	req.URL.RawQuery = q.Encode()

	_, err = rp.CallBaseClient(req, &li)

	return
}

const loginMsg = "To sign in, use a web browser to open the page %s and enter the code %s to authenticate.\n"

func (rp *Rport) GetTokenViaOAuth(ctx context.Context, tokenLifetime int) (token string, err error) {
	providerInfo, _, err := getAuthProviderInfo(ctx, rp.BaseURL)
	if err != nil {
		return "", err
	}

	authSettings, _, err := getAuthSettings(ctx, rp.BaseURL, providerInfo)
	if err != nil {
		return "", err
	}

	abortOnCtrlC()

	loginInfo := authSettings.LoginInfo
	authInfo := loginInfo.DeviceAuthInfo

	fmt.Println("Provider:      ", authSettings.AuthProvider)
	fmt.Println("Authorize URL: ", authInfo.VerificationURI)
	fmt.Println("User Code:     ", authInfo.UserCode)

	if authInfo.Message != "" {
		fmt.Println(authInfo.Message)
	} else {
		fmt.Printf(loginMsg, authInfo.VerificationURI, authInfo.UserCode)
	}
	fmt.Print("\nWaiting for OAuth provider response ... ")

	token, err = pollLogin(ctx, rp.BaseURL, loginInfo, oauth.MaxOAuthRetries, tokenLifetime)
	if err != nil {
		return "", err
	}
	fmt.Println("OK")

	return token, nil
}

func getAuthProviderInfo(ctx context.Context, baseURL string) (providerInfo *oauth.AuthProviderInfo, statusCode int, err error) {
	providerInfo, statusCode, err = oauth.GetAuthProviderInfo(ctx, baseURL)
	if statusCode != http.StatusOK || err != nil {
		if err != nil {
			return nil, statusCode, fmt.Errorf("unable to get auth provider info: %d, %w", statusCode, err)
		}
		return nil, statusCode, fmt.Errorf("unable to get auth provider info: %d", statusCode)
	}
	return providerInfo, statusCode, nil
}

func getAuthSettings(
	ctx context.Context,
	baseURL string,
	providerInfo *oauth.AuthProviderInfo) (authSettings *oauth.DeviceAuthSettings, statusCode int, err error) {
	authSettings, statusCode, err = oauth.GetDeviceAuthSettings(ctx, baseURL, providerInfo)

	if statusCode != http.StatusOK || err != nil {
		if err != nil {
			return nil, statusCode, fmt.Errorf("unable to get auth login info: %d, %w", statusCode, err)
		}
		return nil, statusCode, fmt.Errorf("unable to get auth login info: %d", statusCode)
	}

	return authSettings, statusCode, nil
}

func pollLogin(ctx context.Context, baseURL string, loginInfo oauth.DeviceLoginInfo, retries int, tokenLifetime int) (
	token string, err error) {
	// allow an extra second so not racing with the provider
	interval := loginInfo.DeviceAuthInfo.Interval + 1
	// don't poll faster than the MinIntervalTime and if the interval time is longer
	// than the min then just leave the interval and use that.
	if interval < oauth.MinIntervalTime {
		interval = oauth.MinIntervalTime
	}

	for attempt := 0; attempt < retries; attempt++ {
		loginResponse, statusCode, err := oauth.GetDeviceLogin(
			ctx,
			baseURL,
			loginInfo.LoginURI,
			loginInfo.DeviceAuthInfo.DeviceCode,
			tokenLifetime)
		if err != nil {
			return "", fmt.Errorf("unable to login: %d, %w", statusCode, err)
		}

		// if all ok and no errors then we should have a token
		if statusCode == http.StatusOK && loginResponse.ErrorCode == "" {
			return loginResponse.Token, nil
		}

		if loginResponse.ErrorCode != "" {
			if strings.Contains(loginResponse.ErrorCode, "slow") {
				// back off if asked to slow down
				interval *= 2
			} else if !strings.Contains(loginResponse.ErrorCode, "pending") {
				// if not pending and not slow down, then we have an error, otherwise fall through to sleep and retry
				return "", fmt.Errorf("NOT OK\nunable to login: %s\n%s\n%s",
					loginResponse.ErrorCode,
					loginResponse.ErrorMessage,
					loginResponse.ErrorURI)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return "", fmt.Errorf("max login attempts (%d) exceeded", retries)
}

func abortOnCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// no need to be graceful
		os.Exit(1)
	}()
}

type TwoFaLogin struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (rp *Rport) GetTokenBy2FA(ctx context.Context, twoFACode, login string, tokenLifetime int) (li LoginResponse, err error) {
	loginModel := TwoFaLogin{
		Username: login,
		Token:    twoFACode,
	}

	buf := &bytes.Buffer{}
	loginModelRaw := json.NewEncoder(buf)
	err = loginModelRaw.Encode(loginModel)
	if err != nil {
		return li, err
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url.JoinURL(rp.BaseURL, TwoFaURL),
		buf,
	)
	if err != nil {
		return li, err
	}

	q := req.URL.Query()
	q.Add("token-lifetime", strconv.Itoa(tokenLifetime))
	req.URL.RawQuery = q.Encode()

	_, err = rp.CallBaseClient(req, &li)

	return li, err
}

type MetaPart struct {
	Meta struct{} `json:"meta"`
}

type UserResponse struct {
	Data models.Me `json:"data"`
}

type IPResponse struct {
	Data models.IP `json:"data"`
}

func (rp *Rport) Me(ctx context.Context) (user UserResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, MeURL),
		nil,
	)
	if err != nil {
		return
	}

	_, err = rp.CallBaseClient(req, &user)

	return
}

func (rp *Rport) MeIP(ctx context.Context) (ipResp IPResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, MeIPURL),
		nil,
	)
	if err != nil {
		return
	}

	_, err = rp.CallBaseClient(req, &ipResp)

	return
}

func (rp *Rport) GetIP(ctx context.Context) (string, error) {
	resp, err := rp.MeIP(ctx)
	if err != nil {
		return "", err
	}

	return resp.Data.IP, nil
}

type StatusResponse struct {
	MetaPart
	Data models.Status `json:"data"`
}

func (rp *Rport) Status(ctx context.Context) (st StatusResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, StatusURL),
		nil,
	)
	if err != nil {
		return
	}

	_, err = rp.CallBaseClient(req, &st)

	return
}

func (rp *Rport) Logout(ctx context.Context) (err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url.JoinURL(rp.BaseURL, LogoutURL),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := rp.CallBaseClient(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response code %d, %d is expected", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
