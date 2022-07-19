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
		assert.Equal(t, "Basic bG9nMWRmYWRmOnBhc3MxZGZhc2ZhZGY=", r.Header.Get("Authorization"))
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
			Scheme:      utils.SSH,
			ACL:         "127.0.0.1",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (l, p string, err error) {
			l = "log1dfadf"
			p = "pass1dfasfadf"
			return
		},
	}

	cl := New(srv.URL, apiAuth)

	clientsResp, err := cl.CreateTunnel(
		context.Background(),
		"334",
		"lohost1:3300",
		"rhost2:3344",
		utils.SSH,
		"127.0.0.1",
		"1",
		0,
		false,
		false,
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
	assert.Equal(t, utils.SSH, actualTunnel.Scheme)
	assert.Equal(t, "127.0.0.1", actualTunnel.ACL)
}

func TestCreateTunnelWithHTTPProxy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMWRmYWRmOnBhc3MxZGZhc2ZhZGY=", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodPut, r.Method)

		assert.Equal(t, "/api/v1/clients/1234567890/tunnels?acl=78.34.189.65&check_port=1&http_proxy=true&local=0.0.0.0%3A20793&remote=127.0.0.1%3A80&scheme=http", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(TunnelResponse{Data: &models.Tunnel{
			ID:          "10",
			Lhost:       "0.0.0.0",
			Lport:       "20793",
			Rhost:       "127.0.0.1",
			Rport:       "80",
			LportRandom: true,
			Scheme:      utils.HTTP,
			ACL:         "78.34.189.65",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (l, p string, err error) {
			l = "log1dfadf"
			p = "pass1dfasfadf"
			return
		},
	}

	cl := New(srv.URL, apiAuth)

	clientsResp, err := cl.CreateTunnel(
		context.Background(),
		"1234567890",
		"0.0.0.0:20793",
		"127.0.0.1:80",
		utils.HTTP,
		"78.34.189.65",
		"1",
		0,
		false,
		true,
	)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualTunnel := clientsResp.Data
	assert.Equal(t, "10", actualTunnel.ID)
	assert.Equal(t, "0.0.0.0", actualTunnel.Lhost)
	assert.Equal(t, "20793", actualTunnel.Lport)
	assert.Equal(t, "127.0.0.1", actualTunnel.Rhost)
	assert.Equal(t, "80", actualTunnel.Rport)
	assert.Equal(t, utils.HTTP, actualTunnel.Scheme)
	assert.Equal(t, "78.34.189.65", actualTunnel.ACL)
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
		false,
	)
	assert.NoError(t, err)
	if err != nil {
		return
	}
}
