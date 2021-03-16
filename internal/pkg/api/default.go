package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/utils/url"
)

const (
	LoginURL                    = "/api/v1/login"
	MeURL                       = "/api/v1/me"
	StatusURL                   = "/api/v1/status"
	DefaultTokenValiditySeconds = 90 * 24 * time.Hour //90 days is max value
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
