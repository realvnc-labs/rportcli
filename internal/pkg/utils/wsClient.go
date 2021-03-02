package utils

import (
	"context"
	"io"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WsURLBuilder func(ctx context.Context) (url string, err error)

type Output struct {
	Payload []byte
	Error   error
}

type WsClient struct {
	WsURLBuilder WsURLBuilder
	Conn         *websocket.Conn
}

func NewWsClient(ctx context.Context, wsURLBuilder WsURLBuilder) (wsc *WsClient, err error) {
	wsURL, err := wsURLBuilder(ctx)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return &WsClient{
		WsURLBuilder: wsURLBuilder,
		Conn:         conn,
	}, nil
}

func (wc *WsClient) Close() error {
	if wc.Conn != nil {
		err := wc.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			logrus.Warnf("failed to write close message: %v", err)
		}
		return wc.Conn.Close()
	}

	return nil
}

func (wc *WsClient) Read() (msg []byte, err error) {
	_, msg, e := wc.Conn.ReadMessage()
	if e != nil {
		if _, ok := e.(*websocket.CloseError); ok {
			err = io.EOF
		}
		return
	}

	return msg, nil
}

func (wc *WsClient) Write(inputMsg []byte) (n int, err error) {
	err = wc.Conn.WriteMessage(websocket.TextMessage, inputMsg)
	if err == nil {
		logrus.Debugf("sent command message '%s' to the rport", string(inputMsg))
	}
	return len(inputMsg), err
}
