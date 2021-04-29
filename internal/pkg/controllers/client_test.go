package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

type ClientRendererMock struct {
	Writer io.Writer
}

var clientStub = &models.Client{
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
	ConnState: "connected",
}

func (crm *ClientRendererMock) RenderClients(clients []*models.Client) error {
	jsonBytes, err := json.Marshal(clients)
	if err != nil {
		return err
	}

	_, err = crm.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func (crm *ClientRendererMock) RenderClient(client *models.Client) error {
	jsonBytes, err := json.Marshal(client)
	if err != nil {
		return err
	}

	_, err = crm.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func TestClientsController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}
	clController := ClientController{
		Rport:          cl,
		ClientRenderer: &ClientRendererMock{Writer: &buf},
	}

	err := clController.Clients(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","Tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]}]`,
		buf.String(),
	)
}

func TestClientFoundByIDController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}

	clController := ClientController{
		Rport:          cl,
		ClientRenderer: &ClientRendererMock{Writer: &buf},
	}

	err := clController.Client(context.Background(), &options.ParameterBag{}, "123", "")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","Tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]}`,
		buf.String(),
	)
}

func TestClientFoundByNameController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}

	clSearch := &ClientSearchMock{
		searchTermGiven: "",
		clientsToGive:   []models.Client{*clientStub},
		errorToGive:     nil,
	}
	clController := ClientController{
		Rport:          cl,
		ClientRenderer: &ClientRendererMock{Writer: &buf},
		ClientSearch:   clSearch,
	}

	err := clController.Client(context.Background(), &options.ParameterBag{}, "", "Client 123")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","Tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]}`,
		buf.String(),
	)
}

func TestClientNotFoundController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}

	clController := ClientController{
		Rport:          cl,
		ClientRenderer: &ClientRendererMock{Writer: &buf},
		ClientSearch: &ClientSearchMock{
			errorToGive: errors.New("client not found by the provided id '434' or name ''"),
		},
	}

	err := clController.Client(context.Background(), &options.ParameterBag{}, "434", "")
	assert.EqualError(t, err, `client not found by the provided id '434' or name ''`)

	err = clController.Client(context.Background(), &options.ParameterBag{}, "", "some unknown name")
	assert.EqualError(t, err, `client not found by the provided id '434' or name ''`)
}

func TestInvalidInputForClients(t *testing.T) {
	clController := ClientController{}

	err := clController.Client(context.Background(), &options.ParameterBag{}, "", "")
	assert.EqualError(t, err, "no client id nor name provided")
}

func startClientsServer() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		clientsStub := []*models.Client{clientStub}
		e := jsonEnc.Encode(api.ClientsResponse{Data: clientsStub})
		if e != nil {
			rw.WriteHeader(500)
		}
	}))

	return srv
}
