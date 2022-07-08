package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(handle))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsCl, err := NewWsClient(ctx, func(ctx context.Context) (url string, err error) {
		u := strings.Replace(srv.URL, "http:", "ws:", 1)
		return u, nil
	}, nil)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer wsCl.Close()

	_, err = wsCl.Write([]byte("some"))
	assert.NoError(t, err)
	if err != nil {
		return
	}

	msg, err := wsCl.Read()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "some", string(msg))
}

func handle(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logrus.Error(err)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			logrus.Error(err)
			break
		}
	}
}
