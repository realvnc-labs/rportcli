package api

import (
	"context"
)

const (
	ScriptsWSUri = "/api/v1/ws/scripts"
)

type WsScriptsURLProvider struct {
	*WsURLProvider
}

func (wup *WsScriptsURLProvider) BuildWsURL(ctx context.Context) (wsURL string, err error) {
	return wup.buildWsFullURL(ScriptsWSUri)
}
