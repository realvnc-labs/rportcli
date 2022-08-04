package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/mock"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

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

func (trm *TunnelRendererMock) RenderTunnel(t output.KvProvider) error {
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

type RDPWriterMock struct {
	FileInput      models.FileInput
	filePathToGive string
	errorToGive    error
}

func (rwm *RDPWriterMock) WriteRDPFile(fi models.FileInput) (filePath string, err error) {
	rwm.FileInput = fi
	return rwm.filePathToGive, rwm.errorToGive
}

type RDPExecutorMock struct {
	mock.Mock
}

func (rem *RDPExecutorMock) StartDefaultApp(filePath string) error {
	args := rem.Called(filePath)

	return args.Error(0)
}

func TestTunnelsController(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log145"
			pass = "pass144"
			return
		},
	}
	cl := api.New(srv.URL, apiAuth)

	buf := bytes.Buffer{}
	tController := &TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
	}

	err := tController.Tunnels(context.Background(), &options.ParameterBag{})
	require.NoError(t, err)
	if err != nil {
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, buf.Bytes(), "", "\t")
	require.NoError(t, err)
	t.Logf("Console response: %s", prettyJSON.String())
	expectedOutput := `[
		{
			"id": "1",
			"client_id": "123",
			"client_name": "Client 123",
			"lhost": "",
			"lport": "",
			"rhost": "",
			"rport": "",
			"lport_random": false,
			"scheme": "",
			"acl": "",
			"idle_timeout_minutes": 22
		}
	]`

	assert.JSONEq(t, expectedOutput, buf.String())
}

func TestTunnelDeleteByClientIDController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTU1OnBhc3MxNTU=", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/clients/cl1/tunnels/tun2?force=1", r.URL.String())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log155"
			pass = "pass155"
			return
		},
	}
	cl := api.New(srv.URL, apiAuth)
	buf := bytes.Buffer{}

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
	}

	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		config.ClientID:        "cl1",
		config.TunnelID:        "tun2",
		config.ClientNamesFlag: "",
		config.ForceDeletion:   "1",
	}))
	err := tController.Delete(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, `{"status":"Tunnel successfully deleted"}`, buf.String())
}

func TestInvalidInputForTunnelDelete(t *testing.T) {
	tController := TunnelController{}
	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		config.ClientID:        "",
		config.TunnelID:        "tunnel11",
		config.ClientNamesFlag: "",
	}))
	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, "no client id or name provided")
}

func TestTunnelCreateWithClientID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", r.Header.Get("Authorization"))
		assert.Equal(t, http.MethodPut, r.Method)

		assert.Equal(t, "/api/v1/clients/334/tunnels?acl=3.4.5.6&check_port=1&idle-timeout-minutes=7&local=lohost1%3A3300&remote=rhost2%3A3344&scheme=ssh", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
			ID:              "123",
			Lhost:           "lohost1",
			Lport:           "3300",
			Rhost:           "rhost2",
			Rport:           "3344",
			LportRandom:     true,
			Scheme:          utils.SSH,
			ACL:             "3.4.5.6",
			IdleTimeoutMins: 7,
		}})
		require.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
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
		config.ClientID:           "334",
		config.Local:              "lohost1:3300",
		config.Remote:             "rhost2:3344",
		config.Scheme:             utils.SSH,
		config.CheckPort:          "1",
		config.ServerURL:          "https://localhost.com:34",
		config.IdleTimeoutMinutes: "7",
	})
	err := tController.Create(context.Background(), params)

	require.NoError(t, err)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, buf.Bytes(), "", "\t")
	require.NoError(t, err)
	t.Logf("Console response: %s", prettyJSON.String())
	expectedOutput := `{
		"id": "123",
		"client_id": "334",
		"lhost": "lohost1",
		"lport": "3300",
		"rhost": "rhost2",
		"rport": "3344",
		"lport_random": true,
		"scheme": "ssh",
		"acl": "3.4.5.6",
		"usage": "ssh 127.0.0.1 -p 3300 <...more ssh options>",
		"idle_timeout_minutes": 7,
		"rport_server": "127.0.0.1"
	}`
	assert.JSONEq(t, expectedOutput, buf.String())
}

func TestTunnelCreateWithPortDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/1313/tunnels?acl=3.4.5.9&check_port=&local=lohost44%3A3302&remote=22&scheme=ssh", r.URL.String())
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:              "777",
				Lhost:           "lohost44",
				Rport:           "22",
				Scheme:          "ssh",
				ClientID:        "1313",
				IdleTimeoutMins: 5,
			}})
			require.NoError(t, e)
		}
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "logiin122", "passsii133", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.9",
		},
	}

	params := map[string]string{
		config.ClientID:  "1313",
		config.Local:     "lohost44:3302",
		config.Scheme:    utils.SSH,
		config.ServerURL: "http://some.com",
	}
	err := tController.Create(context.Background(), config.FromValues(params))
	require.NoError(t, err)

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, buf.Bytes(), "", "\t")
	require.NoError(t, err)
	t.Logf("Console response: %s", prettyJSON.String())

	expectedOutput := `{
		"id": "777",
		"client_id": "1313",
		"lhost": "lohost44",
		"lport": "",
		"rhost": "",
		"rport": "22",
		"lport_random": false,
		"scheme": "ssh",
		"acl": "",
		"usage": "ssh 127.0.0.1 -p  \u003c...more ssh options\u003e",
		"idle_timeout_minutes": 5,
		"rport_server": "127.0.0.1"
	}`
	assert.JSONEq(t, expectedOutput, buf.String())
}

func TestTunnelDeleteFailureWithActiveConnections(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
		assert.Equal(t, http.MethodDelete, r.Method)
		errs := models.ErrorResp{
			Errors: []models.Error{
				{
					Code:  "123",
					Title: "tunnel is still active: it has 1 active connection(s)",
				},
			},
		}
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(errs)
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "342314"
			pass = "gfgdgafd"
			return
		},
	}
	cl := api.New(srv.URL, apiAuth)
	buf := bytes.Buffer{}

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
	}
	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		config.ClientID: "cl1",
		config.TunnelID: "tun2",
	}))
	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, "tunnel is still active: it has 1 active connection(s), code: 123, details: , use -f to delete it anyway")
}
