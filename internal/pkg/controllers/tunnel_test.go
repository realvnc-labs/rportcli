package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

type TunnelRendererMock struct{}

func (trm *TunnelRendererMock) RenderTunnels(rw io.Writer, tunnels []*models.Tunnel) error {
	jsonBytes, err := json.Marshal(tunnels)
	if err != nil {
		return err
	}

	_, err = rw.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func TestTunnelsController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	tController := TunnelController{Rport: cl, TunnelRenderer: &TunnelRendererMock{}}

	buf := bytes.Buffer{}
	err := tController.Tunnels(context.Background(), &buf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"1","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]`,
		buf.String(),
	)
}

func TestTunnelDeleteController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/clients/cl1/tunnels/tun2", r.URL.String())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	tController := TunnelController{Rport: cl, TunnelRenderer: &TunnelRendererMock{}}

	buf := bytes.Buffer{}
	err := tController.Delete(context.Background(), &buf, "cl1", "tun2")
	assert.NoError(t, err)
	assert.Equal(t, "OK\n", buf.String())
}
