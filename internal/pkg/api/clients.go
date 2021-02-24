package api

import (
	"context"
	"net/http"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/utils/url"
)

const (
	ClientsURL = "/api/v1/clients"
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

type ClientsResponse struct {
	Data []*models.Client
}

func (rp *Rport) Clients(ctx context.Context) (cr *ClientsResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, ClientsURL),
		nil,
	)
	if err != nil {
		return
	}

	cr = &ClientsResponse{}
	_, err = rp.CallBaseClient(req, cr)

	return
}
