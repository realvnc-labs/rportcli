package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/exec"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/recorder"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTunnelCreateWithSSH(t *testing.T) {
	CmdRecorder := recorder.NewCmdRecorder()
	defer CmdRecorder.Stop()
	randomPort := "234567"
	isTunnelDeleted := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		jsonEnc := json.NewEncoder(rw)
		if r.Method == http.MethodPut {
			assert.Equal(t, "/api/v1/clients/1314/tunnels?acl=3.4.5.10&check_port=&idle-timeout-minutes=5&local=lohost77%3A3303&remote=2000&scheme=ssh", r.URL.String())
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:              "777",
				Lhost:           "lohost77",
				ClientID:        "1314",
				Lport:           randomPort,
				Scheme:          utils.SSH,
				IdleTimeoutMins: 5,
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

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.10",
		},
	}
	sshParams := "-l root -i somefile"
	params := config.FromValues(map[string]string{
		config.ClientID:           "1314",
		config.Local:              "lohost77:3303",
		config.Scheme:             utils.SSH,
		config.Remote:             "2000",
		config.ServerURL:          "http://rport-url.com",
		config.LaunchSSH:          sshParams,
		config.IdleTimeoutMinutes: "5",
	})
	err := tController.Create(context.Background(), params)
	require.NoError(t, err)
	t.Logf("Got response: %s", buf.String())

	assert.Contains(t, buf.String(), "Tunnel successfully deleted")
	assert.Contains(t, buf.String(), "ssh 127.0.0.1 -p "+randomPort)
	assert.True(t, isTunnelDeleted)
	cmdExecuted := CmdRecorder.GetRecords()
	assert.Len(t, cmdExecuted, 1)
	t.Logf("Command recorded: %s", cmdExecuted[0])
	cmdExpected := fmt.Sprintf("ssh 127.0.0.1 %s -p %s", sshParams, randomPort)
	assert.Equal(t, cmdExpected, cmdExecuted[0])
}

func TestTunnelCreateUriHandler(t *testing.T) {
	testCases := []struct {
		scheme          string
		expectedHandler string
		expectedUsage   string
	}{
		{
			scheme:          "vnc",
			expectedHandler: "vnc",
			expectedUsage:   "Connect a vnc viewer to server address '",
		},
		{
			scheme:          "realvnc",
			expectedHandler: "com.realvnc.vncviewer.connect",
			expectedUsage:   "Connect VNCViewer to VNCServer address '",
		},
		{
			scheme:          "http",
			expectedHandler: "http",
			expectedUsage:   "Open the following address with a browser 'http://",
		},
		{
			scheme:          "https",
			expectedHandler: "https",
			expectedUsage:   "Open the following address with a browser 'https://",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scheme, func(t *testing.T) {
			CmdRecorder := recorder.NewCmdRecorder()
			defer CmdRecorder.Stop()
			randomPort := "168759"
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				jsonEnc := json.NewEncoder(rw)
				if r.Method == http.MethodPut {
					expectedURLCalled := fmt.Sprintf(
						"/api/v1/clients/1314/tunnels?acl=3.4.5.10&check_port=&idle-timeout-minutes=5&local=lohost77%%3A3303&remote=%d&scheme=%s",
						utils.GetPortByScheme(tc.scheme),
						tc.scheme,
					)
					assert.Equal(t, expectedURLCalled, r.URL.String())
					e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
						ID:              "777",
						Lhost:           "lohost77",
						ClientID:        "1314",
						Lport:           randomPort,
						Scheme:          tc.scheme,
						IdleTimeoutMins: 5,
					}})
					assert.NoError(t, e)
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
			tController := TunnelController{
				Rport:          cl,
				TunnelRenderer: &TunnelRendererMock{Writer: &buf},
				IPProvider: IPProviderMock{
					IP: "3.4.5.10",
				},
			}
			params := config.FromValues(map[string]string{
				config.ClientID:           "1314",
				config.Local:              "lohost77:3303",
				config.Scheme:             tc.scheme,
				config.ServerURL:          "http://rport-url.com",
				config.LaunchURIHandler:   "true",
				config.IdleTimeoutMinutes: "5",
			})
			err := tController.Create(context.Background(), params)
			require.NoError(t, err)
			t.Logf("Got response: %s", buf.String())

			// Check the right usage is displayed
			assert.Contains(
				t,
				buf.String(),
				fmt.Sprintf("%s127.0.0.1:%s", tc.expectedUsage, randomPort),
				"wrong usage returned",
			)
			// Check the right command is executed
			cmdExpected := fmt.Sprintf("%s %s://127.0.0.1:%s", exec.OpenCmd, tc.expectedHandler, randomPort)
			cmdExecuted := CmdRecorder.GetRecords()
			assert.Len(t, cmdExecuted, 1)
			t.Logf("Command recorded: %s", cmdExecuted[0])
			assert.Equal(t, cmdExpected, cmdExecuted[0], "wrong command executed")
		})
	}
}

func TestTunnelCreateWithRDP(t *testing.T) {
	CmdRecorder := recorder.NewCmdRecorder()
	defer CmdRecorder.Stop()
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.TunnelCreatedResponse{Data: &models.TunnelCreated{
				ID:              "777",
				Lhost:           "lohost77",
				ClientID:        "1314",
				Lport:           "3344",
				Scheme:          utils.RDP,
				IdleTimeoutMins: 5,
				RportServer:     "http://rport-url123.com",
			}})
			assert.NoError(t, e)
		}
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer srv.Close()

	apiAuth := &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) { return "dfasf", "34123", nil },
	}

	buf := bytes.Buffer{}

	cl := api.New(srv.URL, apiAuth)

	tController := TunnelController{
		Rport:          cl,
		TunnelRenderer: &TunnelRendererMock{Writer: &buf},
		IPProvider: IPProviderMock{
			IP: "3.4.5.166",
		},
	}
	params := config.FromValues(map[string]string{
		config.ClientID:           "1314",
		config.Local:              "lohost88:3304",
		config.Scheme:             utils.RDP,
		config.ServerURL:          "http://rport-url123.com",
		config.LaunchRDP:          "1",
		config.RDPUser:            "Administrator",
		config.RDPWidth:           "1090",
		config.RDPHeight:          "990",
		config.IdleTimeoutMinutes: "5",
	})
	err := tController.Create(context.Background(), params)
	assert.NoError(t, err)

	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, buf.Bytes(), "", "\t"); err == nil {
		t.Logf("Console response: %s", prettyJSON.String())
	} else {
		t.Logf("Console response: %s", buf.String())
	}

	assert.Contains(t, buf.String(), "Connect remote desktop to remote pc '127.0.0.1:3344'")
	cmdExecuted := CmdRecorder.GetRecords()
	assert.Len(t, cmdExecuted, 1)
	t.Logf("Command recorded: %s", cmdExecuted[0])
	assert.Regexp(t, exec.OpenCmd+".*client-id-1314.rdp$", cmdExecuted[0])

	// Read the RDP file and check the content
	rdpFileName := strings.Split(cmdExecuted[0], " ")[1]
	assert.FileExists(t, rdpFileName)
	rdpFile, err := os.ReadFile(rdpFileName)
	require.NoError(t, err)
	assert.Contains(t, string(rdpFile), "username:s:Administrator")
	assert.Contains(t, string(rdpFile), "desktopwidth:i:1090")
	assert.Contains(t, string(rdpFile), "desktopheight:i:990")
	assert.Contains(t, string(rdpFile), "full address:s:127.0.0.1:3344")
	err = os.Remove(rdpFileName)
	require.NoError(t, err)
}
