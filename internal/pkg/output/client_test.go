package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderClients(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `GetClients
ID  NAME          TUNNELS REMOTE ADDRESS HOSTNAME OS KERNEL S 
123 SomeName      0                                         C 
124 SomeOtherName 0                                         C 
`,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `[{"id":"123","name":"SomeName","os":"","os_arch":"","os_family":"","os_kernel":"","hostname":"","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":null,"version":"","address":"","Tunnels":null},{"id":"124","name":"SomeOtherName","os":"","os_arch":"","os_family":"","os_kernel":"","hostname":"","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":null,"version":"","address":"","Tunnels":null}]
`,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `[
  {
    "id": "123",
    "name": "SomeName",
    "os": "",
    "os_arch": "",
    "os_family": "",
    "os_kernel": "",
    "hostname": "",
    "connection_state": "connected",
    "disconnected_at": "",
    "client_auth_id": "",
    "ipv4": null,
    "ipv6": null,
    "tags": null,
    "version": "",
    "address": "",
    "Tunnels": null
  },
  {
    "id": "124",
    "name": "SomeOtherName",
    "os": "",
    "os_arch": "",
    "os_family": "",
    "os_kernel": "",
    "hostname": "",
    "connection_state": "connected",
    "disconnected_at": "",
    "client_auth_id": "",
    "ipv4": null,
    "ipv6": null,
    "tags": null,
    "version": "",
    "address": "",
    "Tunnels": null
  }
]
`,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `- id: "123"
  name: SomeName
  os: ""
  osarch: ""
  osfamily: ""
  oskernel: ""
  hostname: ""
  connstate: connected
  disconnectedat: ""
  clientauthid: ""
  ipv4: []
  ipv6: []
  tags: []
  version: ""
  address: ""
  tunnels: []
- id: "124"
  name: SomeOtherName
  os: ""
  osarch: ""
  osfamily: ""
  oskernel: ""
  hostname: ""
  connstate: connected
  disconnectedat: ""
  clientauthid: ""
  ipv4: []
  ipv6: []
  tags: []
  version: ""
  address: ""
  tunnels: []
`,
		},
	}

	clients := []*models.Client{
		{
			ID:        "123",
			Name:      "SomeName",
			ConnState: "connected",
		},
		{
			ID:        "124",
			Name:      "SomeOtherName",
			ConnState: "connected",
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		cr := &ClientRenderer{
			ColCountCalculator: func() int {
				return 150
			},
			Writer: buf,
			Format: testCase.Format,
		}

		err := cr.RenderClients(clients)
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

func TestRenderClient(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Client [123]

KEY               VALUE     
ID:               123       
Name:             SomeName  
Os:                         
OsArch:                     
OsFamily:                   
OsKernel:                   
Hostname:                   
Ipv4:                       
Ipv6:                       
Tags:                       
Version:                    
Address:                    
Connection State: connected 
Disconnected At:            
Client Auth ID:             
`,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"id":"123","name":"SomeName","os":"","os_arch":"","os_family":"","os_kernel":"","hostname":"","connection_state":"connected","disconnected_at":"","client_auth_id":"","ipv4":null,"ipv6":null,"tags":null,"version":"","address":"","Tunnels":null}
`,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "id": "123",
  "name": "SomeName",
  "os": "",
  "os_arch": "",
  "os_family": "",
  "os_kernel": "",
  "hostname": "",
  "connection_state": "connected",
  "disconnected_at": "",
  "client_auth_id": "",
  "ipv4": null,
  "ipv6": null,
  "tags": null,
  "version": "",
  "address": "",
  "Tunnels": null
}
`,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `id: "123"
name: SomeName
os: ""
osarch: ""
osfamily: ""
oskernel: ""
hostname: ""
connstate: connected
disconnectedat: ""
clientauthid: ""
ipv4: []
ipv6: []
tags: []
version: ""
address: ""
tunnels: []
`,
		},
	}
	client := &models.Client{
		ID:        "123",
		Name:      "SomeName",
		ConnState: "connected",
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		cr := &ClientRenderer{
			ColCountCalculator: func() int {
				return 150
			},
			Writer: buf,
			Format: testCase.Format,
		}

		err := cr.RenderClient(client)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(t, testCase.ExpectedOutput, buf.String())
	}
}
