package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	LoginURL                    = "/api/v1/login"
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

type MetaPart struct {
	Meta struct{} `json:"meta"`
}

type UserResponse struct {
	MetaPart
	Data models.User `json:"data"`
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
