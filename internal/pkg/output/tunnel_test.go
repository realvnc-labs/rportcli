package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTunnels(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
		ColCountToGive int
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Tunnels
ID   CLIENT ID CLIENT NAME LOCAL HOST LOCAL PORT REMOTE HOST REMOTE PORT LOCAL PORT RAND SCHEME ACL     TIMEOUT 
id22                       lhost      123        rhost       124         false           ssh    0.0.0.0 33      
`,
			ColCountToGive: 150,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `[{"id":"id22","client_id":"","client_name":"","lhost":"lhost","lport":"123","rhost":"rhost","rport":"124","lport_random":false,"scheme":"ssh","acl":"0.0.0.0","idle_timeout_minutes":33}]
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `[
  {
    "id": "id22",
    "client_id": "",
    "client_name": "",
    "lhost": "lhost",
    "lport": "123",
    "rhost": "rhost",
    "rport": "124",
    "lport_random": false,
    "scheme": "ssh",
    "acl": "0.0.0.0",
    "idle_timeout_minutes": 33
  }
]
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `- id: id22
  client_id: ""
  client_name: ""
  local_host: lhost
  local_port: "123"
  remote_host: rhost
  remote_port: "124"
  local_port_random: false
  scheme: ssh
  acl: 0.0.0.0
  idle_timeout_minutes: 33
`,
			ColCountToGive: 10,
		},
	}

	tunnels := []*models.Tunnel{
		{
			ID:              "id22",
			Lhost:           "lhost",
			Lport:           "123",
			Rhost:           "rhost",
			Rport:           "124",
			LportRandom:     false,
			Scheme:          utils.SSH,
			ACL:             "0.0.0.0",
			IdleTimeoutMins: 33,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(testCase.Format, func(t *testing.T) {
			buf := &bytes.Buffer{}
			colCountToGive := tc.ColCountToGive
			tr := &TunnelRenderer{
				ColCountCalculator: func() int {
					return colCountToGive
				},
				Writer: buf,
				Format: tc.Format,
			}

			err := tr.RenderTunnels(tunnels)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			assert.Equal(
				t,
				tc.ExpectedOutput,
				buf.String(),
			)
		})
	}
}
func TestRenderTunnel(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
		ColCountToGive int
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Tunnel
KEY                VALUE                           
ID:                id22                            
CLIENT_ID:                                         
LOCAL_HOST:        lhost                           
LOCAL_PORT:        123                             
REMOTE_HOST:       rhost                           
REMOTE_PORT:       124                             
LOCAL_PORT RANDOM: false                           
SCHEME:            ssh                             
IDLE TIMEOUT MINS: 7                               
ACL:               0.0.0.0                         
USAGE:             ssh -p 123 123.22.22.33 -l root 
`,
			ColCountToGive: 150,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"id":"id22","client_id":"","client_name":"","lhost":"lhost","lport":"123","rhost":"rhost","rport":"124","lport_random":false,"scheme":"ssh","acl":"0.0.0.0","usage":"ssh -p 123 123.22.22.33 -l root","idle_timeout_minutes":7}
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "id": "id22",
  "client_id": "",
  "client_name": "",
  "lhost": "lhost",
  "lport": "123",
  "rhost": "rhost",
  "rport": "124",
  "lport_random": false,
  "scheme": "ssh",
  "acl": "0.0.0.0",
  "usage": "ssh -p 123 123.22.22.33 -l root",
  "idle_timeout_minutes": 7
}
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `id: id22
client_id: ""
client_name: ""
local_host: lhost
local_port: "123"
remote_host: rhost
remote_port: "124"
local_port_random: false
scheme: ssh
acl: 0.0.0.0
usage: ssh -p 123 123.22.22.33 -l root
idle_timeout_minutes: 7
`,
			ColCountToGive: 10,
		},
	}
	tunnel := &models.TunnelCreated{
		ID:              "id22",
		Lhost:           "lhost",
		Lport:           "123",
		Rhost:           "rhost",
		Rport:           "124",
		LportRandom:     false,
		Scheme:          utils.SSH,
		ACL:             "0.0.0.0",
		Usage:           "ssh -p 123 123.22.22.33 -l root",
		IdleTimeoutMins: 7,
	}

	for _, testCase := range testCases {
		t.Run(testCase.Format, func(t *testing.T) {
			buf := &bytes.Buffer{}
			colCountToGive := testCase.ColCountToGive
			tr := &TunnelRenderer{
				ColCountCalculator: func() int {
					return colCountToGive
				},
				Writer: buf,
				Format: testCase.Format,
			}

			err := tr.RenderTunnel(tunnel)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			actualResult := buf.String()

			assert.Equal(t, testCase.ExpectedOutput, actualResult)
		})
	}
}
