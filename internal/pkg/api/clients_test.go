package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var clientsStub = []*models.Client{
	{
		ID:       "123",
		Name:     "Client 123",
		Os:       "Windows XP",
		OsArch:   "386",
		OsFamily: "Windows",
		OsKernel: "windows",
		Hostname: "localhost",
		Ipv4:     []string{"127.0.0.1"},
		Ipv6:     nil,
		Tags:     []string{"one"},
		Version:  "1",
		Address:  "12.2.2.3",
		Tunnels: []*models.Tunnel{
			{
				ID:          "1",
				Lhost:       "localhost",
				Lport:       "80",
				Rhost:       "rhost",
				Rport:       "81",
				LportRandom: false,
				Scheme:      "https",
				ACL:         "acl123",
			},
		},
	},
	{
		ID:       "124",
		Name:     "Client 124",
		Os:       "Linux Ubuntu",
		OsArch:   "x64",
		OsFamily: "Linux",
		OsKernel: "ubuntu",
		Hostname: "localhost",
		Ipv4:     []string{"127.0.0.1", "127.0.0.2"},
		Ipv6:     nil,
		Tags:     []string{"one", "two"},
		Version:  "2",
		Address:  "12.2.2.4",
		Tunnels: []*models.Tunnel{
			{
				ID:          "1",
				Lhost:       "localhost",
				Lport:       "80",
				Rhost:       "rhost",
				Rport:       "81",
				LportRandom: false,
				Scheme:      "https",
				ACL:         "acl123",
			},
			{
				ID:          "2",
				Lhost:       "localhost",
				Lport:       "66",
				Rhost:       "somehost",
				Rport:       "67",
				LportRandom: true,
				Scheme:      "http",
				ACL:         "acl124",
			},
		},
	},
}

func TestClientsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", authHeader)

		assert.Equal(t, ClientsURL, r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(ClientsResponse{Data: clientsStub})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.BasicAuth{
		Login: "log1",
		Pass:  "pass1",
	})

	clientsResp, err := cl.Clients(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, clientsStub, clientsResp.Data)
}
