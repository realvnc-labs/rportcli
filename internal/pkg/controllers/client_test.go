package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ClientRendererMock struct {}

func (crm *ClientRendererMock) RenderClients(rw io.Writer, clients []*models.Client) error {
	jsonBytes, err := json.Marshal(clients)
	if err != nil {
		return err
	}

	_, err = rw.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func TestClientsController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		clientsStub := []*models.Client{
			{
				ID:       "123",
				Name:     "Client 123",
				Os:       "Windows XP",
				OsArch:   "386",
				OsFamily: "Windows",
				OsKernel: "windows",
				Hostname: "localhost",
				Tags:     []string{"one"},
				Address:  "12.2.2.3:80",
				Tunnels: []*models.Tunnel{
					{
						ID: "1",
					},
				},
			},
		}
		e := jsonEnc.Encode(api.ClientsResponse{Data: clientsStub})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	clController := ClientController{Rport: cl, Cr: &ClientRendererMock{}}

	buf := bytes.Buffer{}
	err := clController.Clients(context.Background(), &buf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","Tunnels":[{"id":"1","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]}]`,
		buf.String(),
	)
}
