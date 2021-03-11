package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

type TunnelRendererMock struct {
	Writer io.Writer
}

func (trm *TunnelRendererMock) RenderTunnels(tunnels []*models.Tunnel) error {
	jsonBytes, err := json.Marshal(tunnels)
	if err != nil {
		return err
	}

	_, err = trm.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func (trm *TunnelRendererMock) RenderDelete(s output.KvProvider) error {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = trm.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func (trm *TunnelRendererMock) RenderTunnel(t *models.Tunnel) error {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	_, err = trm.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

type IPProviderMock struct {
	IP string
}

func (ipm IPProviderMock) GetIP(ctx context.Context) (string, error) {
	return ipm.IP, nil
}

func TestTunnelsController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	apiAuth := &utils.BasicAuth{
		Login: "log1",
		Pass:  "pass1",
	}
	cl := api.New(srv.URL, apiAuth)

	buf := bytes.Buffer{}
	tController := TunnelController{Rport: cl, TunnelRenderer: &TunnelRendererMock{Writer: &buf}}

	err := tController.Tunnels(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"1","client":"123","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]`,
		buf.String(),
	)
}

func TestTunnelDeleteController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/clients/cl1/tunnels/tun2", r.URL.String())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	apiAuth := &utils.BasicAuth{
		Login: "log1",
		Pass:  "pass1",
	}
	cl := api.New(srv.URL, apiAuth)
	buf := bytes.Buffer{}
	tController := TunnelController{Rport: cl, TunnelRenderer: &TunnelRendererMock{Writer: &buf}}

	err := tController.Delete(context.Background(), "cl1", "tun2")
	assert.NoError(t, err)
	assert.Equal(t, `{"status":"OK"}`, buf.String())
}

func TestTunnelCreateController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodPut, r.Method)

		assert.Equal(t, "/api/v1/clients/334/tunnels?acl=3.4.5.6&check_port=1&local=lohost1%3A3300&remote=rhost2%3A3344&scheme=ssh", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(api.TunnelResponse{Data: &models.Tunnel{
			ID:          "123",
			Lhost:       "lohost1",
			Lport:       "3300",
			Rhost:       "rhost2",
			Rport:       "3344",
			LportRandom: true,
			Scheme:      "ssh",
			ACL:         "3.4.5.6",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.BasicAuth{
		Login: "log1",
		Pass:  "pass1",
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.6",
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:  "334",
		Local:     "lohost1:3300",
		Remote:    "rhost2:3344",
		Scheme:    "ssh",
		CheckPort: "1",
	})
	err := tController.Create(context.Background(), params)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":"123","client":"","lhost":"lohost1","lport":"3300","rhost":"rhost2","rport":"3344","lport_random":true,"scheme":"ssh","acl":"3.4.5.6"}`, buf.String())
}
