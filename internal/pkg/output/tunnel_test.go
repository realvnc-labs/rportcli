package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTunnels(t *testing.T) {
	testCases := []struct{
		Format string
		ExpectedOutput string
		ColCountToGive int
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Tunnels
ID   CLIENT LHOST LPORT RHOST RPORT LPORTRAND SCHEME ACL     
id22        lhost 123   rhost 124   false     ssh    0.0.0.0 
`,
			ColCountToGive: 150,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `[{"id":"id22","client":"","lhost":"lhost","lport":"123","rhost":"rhost","rport":"124","lport_random":false,"scheme":"ssh","acl":"0.0.0.0"}]
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `[
  {
    "id": "id22",
    "client": "",
    "lhost": "lhost",
    "lport": "123",
    "rhost": "rhost",
    "rport": "124",
    "lport_random": false,
    "scheme": "ssh",
    "acl": "0.0.0.0"
  }
]
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `- id: id22
  client: ""
  lhost: lhost
  lport: "123"
  rhost: rhost
  rport: "124"
  lportrandom: false
  scheme: ssh
  acl: 0.0.0.0
`,
			ColCountToGive: 10,
		},
	}

	tunnels := []*models.Tunnel{
		{
			ID:          "id22",
			Lhost:       "lhost",
			Lport:       "123",
			Rhost:       "rhost",
			Rport:       "124",
			LportRandom: false,
			Scheme:      "ssh",
			ACL:         "0.0.0.0",
		},
	}
	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		tr := &TunnelRenderer{
			ColCountCalculator: func() int {
				return testCase.ColCountToGive
			},
			Writer: buf,
			Format: testCase.Format,
		}

		err := tr.RenderTunnels(tunnels)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(
			t,
			testCase.ExpectedOutput,
			buf.String(),
		)
	}
}
func TestRenderTunnel(t *testing.T) {
	testCases := []struct{
		Format string
		ExpectedOutput string
		ColCountToGive int
	}{
		{
			Format:         FormatHuman,
			ExpectedOutput: `Tunnel
KEY           VALUE   
ID:           id22    
CLIENT:               
LHOST:        lhost   
LPORT:        123     
RHOST:        rhost   
RPORT:        124     
LPORT RANDOM: false   
SCHEME:       ssh     
ACL:          0.0.0.0 
`,
			ColCountToGive: 150,
		},
		{
			Format:         FormatJSON,
			ExpectedOutput: `{"id":"id22","client":"","lhost":"lhost","lport":"123","rhost":"rhost","rport":"124","lport_random":false,"scheme":"ssh","acl":"0.0.0.0"}
`,
			ColCountToGive: 10,
		},
		{
			Format:         FormatJSONPretty,
			ExpectedOutput: `{
  "id": "id22",
  "client": "",
  "lhost": "lhost",
  "lport": "123",
  "rhost": "rhost",
  "rport": "124",
  "lport_random": false,
  "scheme": "ssh",
  "acl": "0.0.0.0"
}
`,
			ColCountToGive: 10,
		},
		{
			Format:         FormatYAML,
			ExpectedOutput: `id: id22
client: ""
lhost: lhost
lport: "123"
rhost: rhost
rport: "124"
lportrandom: false
scheme: ssh
acl: 0.0.0.0
`,
			ColCountToGive: 10,
		},
	}
	tunnel := &models.Tunnel{
		ID:          "id22",
		Lhost:       "lhost",
		Lport:       "123",
		Rhost:       "rhost",
		Rport:       "124",
		LportRandom: false,
		Scheme:      "ssh",
		ACL:         "0.0.0.0",
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		tr := &TunnelRenderer{
			ColCountCalculator: func() int {
				return testCase.ColCountToGive
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
	}
}
