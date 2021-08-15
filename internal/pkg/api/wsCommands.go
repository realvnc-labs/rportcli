package api

import (
	"context"
)

const (
	CommandsWSUri = "/api/v1/ws/commands"
)

type WsCommandURLProvider struct {
	*WsURLProvider
}

func (wup *WsCommandURLProvider) BuildWsURL(ctx context.Context) (wsURL string, err error) {
	return wup.buildWsFullURL(CommandsWSUri)
}
