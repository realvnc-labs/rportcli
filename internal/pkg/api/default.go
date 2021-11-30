package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

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
