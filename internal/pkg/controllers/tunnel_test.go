package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/rdp"

	options "github.com/breathbath/go_utils/v2/pkg/config"

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
	isSSHExecuted := false
	tController := &TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		SSHFunc: func(sshParams []string) error {
			isSSHExecuted = true
			return nil
		},
		ClientSearch:   &ClientSearchMock{},
	}
	assert.False(t, isSSHExecuted)

	err := tController.Tunnels(context.Background(), &options.ParameterBag{})
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(
		t,
		`[{"id":"1","client_id":"123","client_name":"Client 123","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]`,
		buf.String(),
	)
}

func TestTunnelsControllerByClient(t *testing.T) {
	srv := startClientsServer()
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1455"
			pass = "pass1446"
			return
		},
	}
	cl := api.New(srv.URL, apiAuth)
	buf := bytes.Buffer{}
	searchMock := &ClientSearchMock{
		clientsToGive: []*models.Client{
			{
				ID:   "cl2",
				Name: "client 354351",
				Tunnels: []*models.Tunnel{
					{
						ID: "23",
					},
				},
			},
		},
	}
	tController := &TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		SSHFunc: func(sshParams []string) error {
			return nil
		},
		ClientSearch:   searchMock,
	}

	paramProv := options.NewMapValuesProvider(map[string]interface{}{
		ClientNameFlag: "client 354351",
	})
	err := tController.Tunnels(context.Background(), options.New(paramProv))
	require.NoError(t, err)

	assert.Equal(
		t,
		`[{"id":"23","client_id":"cl2","client_name":"client 354351","lhost":"","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":""}]`,
		buf.String(),
	)

	assert.Equal(t, "client 354351", searchMock.searchTermGiven)

	paramProv2 := options.NewMapValuesProvider(map[string]interface{}{
		ClientID: "cl2",
	})
	err = tController.Tunnels(context.Background(), options.New(paramProv2))
	require.NoError(t, err)
	assert.Equal(t, "cl2", searchMock.searchTermGiven)
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

	isSSHExecuted := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		ClientSearch:   &ClientSearchMock{},
		SSHFunc: func(sshParams []string) error {
			isSSHExecuted = true
			return nil
		},
	}
	assert.False(t, isSSHExecuted)

	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID:       "cl1",
		TunnelID:       "tun2",
		ClientNameFlag: "",
		ForceDeletion:  "1",
	}))
	err := tController.Delete(context.Background(), params)
	assert.NoError(t, err)
	assert.Equal(t, `{"status":"Tunnel successfully deleted"}`, buf.String())
}

func TestTunnelDeleteByClientNameController(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/clients/cl2/tunnels/tun4", r.URL.String())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log24124"
			pass = "pass341324"
			return
		},
	}
	cl := api.New(srv.URL, apiAuth)
	buf := bytes.Buffer{}
	searchMock := &ClientSearchMock{
		clientsToGive: []*models.Client{
			{
				ID:   "cl2",
				Name: "some client",
			},
		},
	}
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		ClientSearch:   searchMock,
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}

	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID:       "",
		TunnelID:       "tun4",
		ClientNameFlag: "some client",
	}))

	err := tController.Delete(context.Background(), params)
	assert.NoError(t, err)
	assert.Equal(t, `{"status":"Tunnel successfully deleted"}`, buf.String())
}

func TestTunnelDeleteByAmbiguousClientName(t *testing.T) {
	searchMock := &ClientSearchMock{
		clientsToGive: []*models.Client{
			{
				ID:   "cl1",
				Name: "some client 1",
			},
			{
				ID:   "cl2",
				Name: "some client 2",
			},
		},
	}
	tController := TunnelController{
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}

	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID:       "",
		TunnelID:       "tun3",
		ClientNameFlag: "some client",
	}))

	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, `client identified by 'some client' is ambiguous, use a more precise name or use the client id`)
}

func TestTunnelDeleteNotFoundClientName(t *testing.T) {
	searchMock := &ClientSearchMock{
		clientsToGive: []*models.Client{},
	}
	tController := TunnelController{
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}

	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID:       "",
		TunnelID:       "tun5",
		ClientNameFlag: "some client",
	}))

	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, `unknown client 'some client'`)
}

func TestInvalidInputForTunnelDelete(t *testing.T) {
	tController := TunnelController{}
	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID:       "",
		TunnelID:       "tunnel11",
		ClientNameFlag: "",
	}))
	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, "no client id nor name provided")
}

func TestTunnelCreateWithClientID(t *testing.T) {
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
			Scheme:      utils.SSH,
			ACL:         "3.4.5.6",
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

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)
	isSSHExecuted := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.6",
		},
		SSHFunc: func(sshParams []string) error {
			isSSHExecuted = true
			return nil
		},
	}
	assert.False(t, isSSHExecuted)

	params := config.FromValues(map[string]string{
		ClientID:         "334",
		Local:            "lohost1:3300",
		Remote:           "rhost2:3344",
		Scheme:           utils.SSH,
		CheckPort:        "1",
		config.ServerURL: "https://localhost.com:34",
	})
	err := tController.Create(context.Background(), params)
	assert.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		`{"id":"123","client_id":"334","client_name":"","lhost":"lohost1","lport":"3300","rhost":"rhost2","rport":"3344","lport_random":true,"scheme":"ssh","acl":"3.4.5.6","usage":"ssh -p 3300 localhost.com -l ${USER}","rport_server":"%s"}`,
		srv.URL,
	)
	assert.Equal(t, expectedOutput, buf.String())
}

func TestTunnelCreateWithClientName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/clients/444/tunnels?acl=3.4.5.7&check_port=1&local=lohost2%3A3301&remote=rhost4%3A3345&scheme=ssh", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(api.TunnelResponse{Data: &models.Tunnel{
			ID:          "444",
			Lhost:       "lohost2",
			Lport:       "3301",
			Rhost:       "rhost4",
			Rport:       "3345",
			LportRandom: true,
			Scheme:      utils.SSH,
			ACL:         "3.4.5.7",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "someloggg", "somepaaas", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	searchMock := &ClientSearchMock{
		clientsToGive: []*models.Client{
			{
				ID:   "444",
				Name: "some client 444",
			},
		},
	}

	isSSHExecuted := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.7",
		},
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			isSSHExecuted = true
			return nil
		},
	}
	assert.False(t, isSSHExecuted)

	params := config.FromValues(map[string]string{
		ClientID:         "",
		ClientNameFlag:   "some client 444",
		Local:            "lohost2:3301",
		Remote:           "rhost4:3345",
		Scheme:           utils.SSH,
		CheckPort:        "1",
		config.ServerURL: "http://11.11.11.11:33",
	})
	err := tController.Create(context.Background(), params)
	assert.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		`{"id":"444","client_id":"some client 444","client_name":"","lhost":"lohost2","lport":"3301","rhost":"rhost4","rport":"3345","lport_random":true,"scheme":"ssh","acl":"3.4.5.7","usage":"ssh -p 3301 11.11.11.11 -l ${USER}","rport_server":"%s"}`,
		srv.URL,
	)
	assert.Equal(t, expectedOutput, buf.String())
}

func TestInvalidInputForTunnelCreate(t *testing.T) {
	tController := TunnelController{}
	params := config.FromValues(map[string]string{
		ClientID:       "",
		ClientNameFlag: "",
		Local:          "lohost1:3300",
		Remote:         "rhost2:3344",
		Scheme:         utils.SSH,
		CheckPort:      "1",
	})
	err := tController.Create(context.Background(), params)
	assert.EqualError(t, err, "no client id nor name provided")
}

func TestTunnelCreateNotFoundClientName(t *testing.T) {
	searchMock := &ClientSearchMock{
		errorToGive: errors.New("unknown client 'some client'"),
	}
	tController := TunnelController{
		ClientSearch: searchMock,
	}

	params := config.FromValues(map[string]string{
		ClientNameFlag: "some client",
	})
	err := tController.Create(context.Background(), params)
	assert.EqualError(t, err, `unknown client 'some client'`)
}

func TestTunnelCreateWithSchemeDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/32312/tunnels?acl=3.4.5.8&check_port=&local=lohost33%3A3301&remote=rhost5%3A22&scheme=ssh", r.URL.String())
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.TunnelResponse{Data: &models.Tunnel{
				ID:       "444",
				Lhost:    "lohost33",
				ClientID: "32312",
			}})
			assert.NoError(t, e)
		}
		if r.Method == http.MethodDelete {
			rw.WriteHeader(http.StatusNoContent)
		}
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "logiin1", "passsii1", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	searchMock := &ClientSearchMock{clientsToGive: []*models.Client{}}

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.8",
		},
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}

	params := map[string]string{
		ClientID:         "32312",
		Local:            "lohost33:3301",
		Remote:           "rhost5:22",
		config.ServerURL: "http://ya.ru",
	}
	err := tController.Create(context.Background(), config.FromValues(params))
	assert.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		`{"id":"444","client_id":"32312","client_name":"","lhost":"lohost33","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","usage":"ssh ya.ru -l ${USER}","rport_server":"%s"}`,
		srv.URL,
	)

	assert.Equal(
		t,
		expectedOutput,
		buf.String(),
	)
}

func TestTunnelCreateWithPortDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/1313/tunnels?acl=3.4.5.9&check_port=&local=lohost44%3A3302&remote=22&scheme=ssh", r.URL.String())
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:       "777",
				Lhost:    "lohost44",
				ClientID: "1313",
			}})
			assert.NoError(t, e)
		}
		if r.Method == http.MethodDelete {
			assert.Equal(t, "/api/v1/clients/1313/tunnels/777", r.URL.String())
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "logiin122", "passsii133", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	searchMock := &ClientSearchMock{clientsToGive: []*models.Client{}}

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.9",
		},
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}

	params := map[string]string{
		ClientID:         "1313",
		Local:            "lohost44:3302",
		Scheme:           utils.SSH,
		config.ServerURL: "http://some.com",
	}
	err := tController.Create(context.Background(), config.FromValues(params))
	assert.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		`{"id":"777","client_id":"1313","client_name":"","lhost":"lohost44","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","usage":"ssh some.com -l ${USER}","rport_server":"%s"}`,
		srv.URL,
	)

	assert.Equal(
		t,
		expectedOutput,
		buf.String(),
	)
	buf = bytes.Buffer{}

	delete(params, Scheme)
	params[LaunchSSH] = "-l root"
	err = tController.Create(context.Background(), config.FromValues(params))
	assert.NoError(t, err)

	expectedOutput2 := fmt.Sprintf(
		`{"id":"777","client_id":"1313","client_name":"","lhost":"lohost44","lport":"","rhost":"","rport":"","lport_random":false,"scheme":"","acl":"","usage":"ssh some.com -l ${USER}","rport_server":"%s"}{"status":"Tunnel successfully deleted"}`,
		srv.URL,
	)
	assert.Equal(
		t,
		expectedOutput2,
		buf.String(),
	)
}

func TestTunnelCreateWithSSH(t *testing.T) {
	isTunnelDeleted := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/1314/tunnels?acl=3.4.5.10&check_port=&local=lohost77%3A3303&remote=22&scheme=ssh", r.URL.String())
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:       "777",
				Lhost:    "lohost77",
				ClientID: "1314",
				Lport:    "22",
				Scheme:   utils.SSH,
			}})
			assert.NoError(t, e)
			return
		}
		if r.Method == http.MethodDelete {
			isTunnelDeleted = true
			assert.Equal(t, "/api/v1/clients/1314/tunnels/777", r.URL.String())
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		rw.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "364872364", "3463284", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	searchMock := &ClientSearchMock{clientsToGive: []*models.Client{}}

	isSSHCalled := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.10",
		},
		ClientSearch: searchMock,
		SSHFunc: func(sshParams []string) error {
			isSSHCalled = true
			assert.Equal(t, []string{"rport-url.com", "-p", "22", "-l", "root", "-i", "somefile"}, sshParams)
			return nil
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:         "1314",
		Local:            "lohost77:3303",
		Scheme:           utils.SSH,
		config.ServerURL: "http://rport-url.com",
		LaunchSSH:        "-l root -i somefile",
	})
	err := tController.Create(context.Background(), params)
	assert.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		`{"id":"777","client_id":"1314","client_name":"","lhost":"lohost77","lport":"22","rhost":"","rport":"","lport_random":false,"scheme":"ssh","acl":"","usage":"ssh -p 22 rport-url.com -l ${USER}","rport_server":"%s"}{"status":"Tunnel successfully deleted"}`,
		srv.URL,
	)

	assert.Equal(
		t,
		expectedOutput,
		buf.String(),
	)

	assert.True(t, isSSHCalled)
	assert.True(t, isTunnelDeleted)
}

func TestTunnelCreateWithSSHFailure(t *testing.T) {
	isTunnelDeleted := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/1316/tunnels?acl=3.4.5.16&check_port=&local=lohost776%3A3306&remote=22&scheme=ssh", r.URL.String())
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:       "6666",
				Lhost:    "lohost66",
				ClientID: "1316",
			}})
			assert.NoError(t, e)
			return
		}
		if r.Method == http.MethodDelete {
			isTunnelDeleted = true
			assert.Equal(t, "/api/v1/clients/1316/tunnels/6666", r.URL.String())
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		rw.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "sdfafj", "34234", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.16",
		},
		ClientSearch: &ClientSearchMock{clientsToGive: []*models.Client{}},
		SSHFunc: func(sshParams []string) error {
			return errors.New("ssh failure")
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:         "1316",
		Local:            "lohost776:3306",
		config.ServerURL: "http://rport-url2.com",
		LaunchSSH:        "-l root",
	})
	err := tController.Create(context.Background(), params)
	assert.EqualError(t, err, "ssh failure")
	assert.True(t, isTunnelDeleted)
}

func TestTunnelCreateWithRDP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:       "777",
				Lhost:    "lohost77",
				ClientID: "1314",
				Lport:    "3344",
				Scheme:   utils.SSH,
			}})
			assert.NoError(t, e)
		}
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "dfasf", "34123", nil },
	}

	renderBuf := bytes.Buffer{}
	cmdOutput := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	const filePathGiven = "/tmp/somefile.rdp"
	fileWriter := &RDPWriterMock{
		filePathToGive: filePathGiven,
		errorToGive:    nil,
	}
	isRDPCalled := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &renderBuf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.166",
		},
		ClientSearch: &ClientSearchMock{clientsToGive: []*models.Client{}},
		RDPWriter:    fileWriter,
		RDPExecutor: &rdp.Executor{
			CommandProvider: func(filePath string) (cmd string, args []string) {
				isRDPCalled = true
				assert.Equal(t, filePathGiven, filePath)
				return "echo", []string{"rdp executed"}
			},
			StdOut: &cmdOutput,
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:         "1314",
		Local:            "lohost88:3304",
		Scheme:           utils.RDP,
		config.ServerURL: "http://rport-url123.com",
		LaunchRDP:        "1",
		RDPUser:          "Administrator",
		RDPWidth:         "1090",
		RDPHeight:        "990",
	})
	err := tController.Create(context.Background(), params)
	expectedFileInput := models.FileInput{
		Address:      "rport-url123.com:3344",
		ScreenHeight: 990,
		ScreenWidth:  1090,
		UserName:     "Administrator",
	}
	assert.Equal(t, expectedFileInput.Address, fileWriter.FileInput.Address)
	assert.Equal(t, expectedFileInput.ScreenHeight, fileWriter.FileInput.ScreenHeight)
	assert.Equal(t, expectedFileInput.ScreenWidth, fileWriter.FileInput.ScreenWidth)
	assert.Equal(t, expectedFileInput.UserName, fileWriter.FileInput.UserName)
	assert.NoError(t, err)
	assert.Equal(
		t,
		`rdp executed
`,
		cmdOutput.String(),
	)

	assert.True(t, isRDPCalled)
}

func TestTunnelCreateWithRDPIncompatibleFlags(t *testing.T) {
	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "dfasf", "34123", nil },
	}

	renderBuf := bytes.Buffer{}
	cmdOutput := bytes.Buffer{}

	cl := api.New("localhost", apiAuth)

	isRDPCalled := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &renderBuf},
		IPProvider:     IPProviderMock{},
		ClientSearch:   &ClientSearchMock{clientsToGive: []*models.Client{}},
		RDPWriter:      nil,
		RDPExecutor: &rdp.Executor{
			CommandProvider: func(filePath string) (cmd string, args []string) {
				isRDPCalled = true
				return
			},
			StdOut: &cmdOutput,
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:         "1319",
		Local:            "lohost88:3305",
		Scheme:           utils.RDP,
		config.ServerURL: "http://rport-url123.com",
		LaunchSSH:        "-l root",
		ACL:              "0.0.0.0",
	})
	err := tController.Create(context.Background(), params)
	assert.EqualError(t, err, fmt.Sprintf("scheme rdp is not compatible with the %s option", LaunchSSH))
	assert.False(t, isRDPCalled)
}

func TestTunnelCreateWithSSHIncompatibleFlags(t *testing.T) {
	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "2123", "34124", nil },
	}

	renderBuf := bytes.Buffer{}

	cl := api.New("localhost", apiAuth)

	isSSHCalled := false
	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &renderBuf},
		IPProvider:     IPProviderMock{},
		ClientSearch:   &ClientSearchMock{clientsToGive: []*models.Client{}},
		SSHFunc: func(sshParams []string) error {
			isSSHCalled = true
			return nil
		},
	}

	params := config.FromValues(map[string]string{
		ClientID:         "1320",
		Local:            "lohost88:3309",
		Scheme:           utils.SSH,
		config.ServerURL: "http://rport-url125.com",
		LaunchRDP:        "1",
	})
	err := tController.Create(context.Background(), params)
	assert.EqualError(t, err, fmt.Sprintf("scheme ssh is not compatible with the %s option", LaunchRDP))
	assert.False(t, isSSHCalled)
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
		ClientSearch:   &ClientSearchMock{},
		SSHFunc: func(sshParams []string) error {
			return nil
		},
	}
	params := options.New(options.NewMapValuesProvider(map[string]interface{}{
		ClientID: "cl1",
		TunnelID: "tun2",
	}))
	err := tController.Delete(context.Background(), params)
	assert.EqualError(t, err, "tunnel is still active: it has 1 active connection(s), code: 123, details: , use -f to delete it anyway")
}
