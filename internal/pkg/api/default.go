package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/breathbath/go_utils/utils/url"
)

const (
	LoginURL  = "/api/v1/login"
	MeURL     = "/api/v1/me"
	StatusURL = "/api/v1/status"
)

type LoginInfo struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (rp *Rport) Login(ctx context.Context, login, pass string, tokenLifetime int) (li LoginInfo, err error) {
	ba := &BasicAuth{
		Login: login,
		Pass:  pass,
	}

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

	cl := &Client{}
	cl.WithAuth(ba)

	err = cl.Call(req, &li)
	return
}

type UserInfo struct {
	Meta struct {
	} `json:"meta"`
	Data struct {
		User   string   `json:"user"`
		Groups []string `json:"groups"`
	} `json:"data"`
}

func (rp *Rport) Me() (user UserInfo, err error) {
	var req *http.Request
	req, err = http.NewRequest(
		http.MethodGet,
		url.JoinURL(rp.BaseURL, MeURL),
		nil,
	)
	if err != nil {
		return
	}

	cl := &Client{}
	cl.WithAuth(rp.Auth)

	err = cl.Call(req, &user)
	return
}

type Status struct {
	Meta struct {
	} `json:"meta"`
	Data struct {
		SessionsCount int    `json:"sessions_count"`
		Version       string `json:"version"`
		Fingerprint   string `json:"fingerprint"`
		ConnectURL    string `json:"connect_url"`
	} `json:"data"`
}

func (rp *Rport) Status(ctx context.Context) (st Status, err error) {
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

	cl := &Client{}
	cl.WithAuth(rp.Auth)

	err = cl.Call(req, &st)
	return
}
