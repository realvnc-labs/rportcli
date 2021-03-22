package api

import (
	"context"
	"strings"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	CommandsWSUri = "/api/v1/ws/commands"
	HTTPSPrefix   = "https"
	HTTPPrefix    = "http"
	WssPrefix     = "wss"
	WsPrefix      = "ws"
)

type WsCommandURLProvider struct {
	BaseURL              string
	TokenProvider        func() (token string, err error)
	TokenValiditySeconds int
}

func (wup *WsCommandURLProvider) BuildWsURL(ctx context.Context) (wsURL string, err error) {
	token, err := wup.TokenProvider()
	if err != nil {
		return "", err
	}

	wsURL = wup.buildWsURL(token, wup.BaseURL)

	return
}

func (wup *WsCommandURLProvider) buildWsURL(token, baseURL string) string {
	baseURL = wup.replaceHTTPWithWsProtocolPrefix(baseURL)
	return url.JoinURL(baseURL, CommandsWSUri) + "?access_token=" + token
}

func (wup *WsCommandURLProvider) replaceHTTPWithWsProtocolPrefix(u string) string {
	if strings.HasPrefix(u, HTTPSPrefix) {
		return strings.Replace(u, HTTPSPrefix, WssPrefix, 1)
	}
	if strings.HasPrefix(u, HTTPPrefix) {
		return strings.Replace(u, HTTPPrefix, WsPrefix, 1)
	}

	return u
}
