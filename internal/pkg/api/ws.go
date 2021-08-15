package api

import (
	"fmt"
	"strings"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	HTTPSPrefix = "https"
	HTTPPrefix  = "http"
	WssPrefix   = "wss"
	WsPrefix    = "ws"
)

type WsURLProvider struct {
	BaseURL              string
	TokenProvider        func() (token string, err error)
	TokenValiditySeconds int
}

func (wup *WsURLProvider) buildWsFullURL(uriPath string) (wsURL string, err error) {
	token, err := wup.TokenProvider()
	if err != nil {
		return "", err
	}

	if token == "" {
		err = fmt.Errorf("no auth data stored to use this command, please call init to login")
		return
	}

	wsURL = wup.buildWsURL(token, wup.BaseURL, uriPath)

	return
}

func (wup *WsURLProvider) buildWsURL(token, baseURL, uriPath string) string {
	baseURL = wup.replaceHTTPWithWsProtocolPrefix(baseURL)
	return url.JoinURL(baseURL, uriPath) + "?access_token=" + token
}

func (wup *WsURLProvider) replaceHTTPWithWsProtocolPrefix(u string) string {
	if strings.HasPrefix(u, HTTPSPrefix) {
		return strings.Replace(u, HTTPSPrefix, WssPrefix, 1)
	}
	if strings.HasPrefix(u, HTTPPrefix) {
		return strings.Replace(u, HTTPPrefix, WsPrefix, 1)
	}

	return u
}
