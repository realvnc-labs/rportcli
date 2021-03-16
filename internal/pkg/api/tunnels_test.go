package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateTunnel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodPut, r.Method)

		assert.Equal(t, "/api/v1/clients/334/tunnels?acl=127.0.0.1&check_port=1&local=lohost1%3A3300&remote=rhost2%3A3344&scheme=ssh", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(TunnelResponse{Data: &models.Tunnel{
			ID:          "123",
			Lhost:       "lohost1",
			Lport:       "3300",
			Rhost:       "rhost2",
			Rport:       "3344",
			LportRandom: true,
			Scheme:      "ssh",
			ACL:         "127.0.0.1",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	}

	cl := New(srv.URL, apiAuth)

	clientsResp, err := cl.CreateTunnel(
		context.Background(),
		"334",
		"lohost1:3300",
		"rhost2:3344",
		"ssh",
		"127.0.0.1",
		"1",
	)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualTunnel := clientsResp.Data
	assert.Equal(t, "123", actualTunnel.ID)
	assert.Equal(t, "lohost1", actualTunnel.Lhost)
	assert.Equal(t, "3300", actualTunnel.Lport)
	assert.Equal(t, "rhost2", actualTunnel.Rhost)
	assert.Equal(t, "3344", actualTunnel.Rport)
	assert.Equal(t, "ssh", actualTunnel.Scheme)
	assert.Equal(t, "127.0.0.1", actualTunnel.ACL)
}

func TestDeleteTunnel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodDelete, r.Method)

		assert.Equal(t, "/api/v1/clients/123/tunnels/345", r.URL.String())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()


	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	}
	cl := New(srv.URL, apiAuth)

	err := cl.DeleteTunnel(
		context.Background(),
		"123",
		"345",
	)
	assert.NoError(t, err)
	if err != nil {
		return
	}
}
