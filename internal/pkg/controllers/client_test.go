package controllers

import (
	"bytes"
	"context"
	"encoding/json"
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
	Writer             io.Writer
	renderDetailsGiven bool
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
			ID:              "1",
			IdleTimeoutMins: 22,
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

func (crm *ClientRendererMock) RenderClient(client *models.Client, renderDetails bool) error {
	crm.renderDetailsGiven = renderDetails

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

	err := clController.Clients(context.Background(), options.New(nil))
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","idle_timeout_minutes":22}],"os_full_name":"","os_version":"","os_virtualization_system":"","os_virtualization_role":"","cpu_family":"","cpu_model":"","cpu_model_name":"","cpu_vendor":"","num_cpus":0,"mem_total":0,"timezone":"","allowed_user_groups":null,"updates_status":null}]`,
		buf.String(),
	)
}

func TestClientFoundByIDController(t *testing.T) {
	srv := startClientServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}

	renderMock := &ClientRendererMock{Writer: &buf}
	clController := ClientController{
		Rport:          cl,
		ClientRenderer: renderMock,
	}

	paramsProv := options.NewMapValuesProvider(map[string]interface{}{"all": true})
	params := options.New(paramsProv)

	err := clController.Client(context.Background(), params, "123", "")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","idle_timeout_minutes":22}],"os_full_name":"","os_version":"","os_virtualization_system":"","os_virtualization_role":"","cpu_family":"","cpu_model":"","cpu_model_name":"","cpu_vendor":"","num_cpus":0,"mem_total":0,"timezone":"","allowed_user_groups":null,"updates_status":null}`,
		buf.String(),
	)
	assert.True(t, renderMock.renderDetailsGiven)
}

func TestClientFoundByNameController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	cl := api.New(srv.URL, nil)
	buf := bytes.Buffer{}

	clController := ClientController{
		Rport:          cl,
		ClientRenderer: &ClientRendererMock{Writer: &buf},
	}

	err := clController.Client(context.Background(), &options.ParameterBag{}, "", "Client 123")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`{"id":"123","name":"Client 123","os":"Windows XP","os_arch":"386","os_family":"Windows","os_kernel":"windows","hostname":"localhost","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":["one"],"version":"","address":"12.2.2.3:80","tunnels":[{"id":"1","client_id":"","client_name":"","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","idle_timeout_minutes":22}],"os_full_name":"","os_version":"","os_virtualization_system":"","os_virtualization_role":"","cpu_family":"","cpu_model":"","cpu_model_name":"","cpu_vendor":"","num_cpus":0,"mem_total":0,"timezone":"","allowed_user_groups":null,"updates_status":null}`,
		buf.String(),
	)
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

func startClientServer() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(api.ClientResponse{Data: clientStub})
		if e != nil {
			rw.WriteHeader(500)
		}
	}))

	return srv
}
