package api

import (
	"context"
	"net/http"

	"github.com/breathbath/go_utils/utils/url"
)

const (
	ClientsUrl  = "/api/v1/clients"
)

type Tunnel struct {
	ID          string `json:"id"`
	Lhost       string `json:"lhost"`
	Lport       string `json:"lport"`
	Rhost       string `json:"rhost"`
	Rport       string `json:"rport"`
	LportRandom bool   `json:"lport_random"`
	Scheme      string `json:"scheme"`
	ACL         string `json:"acl"`
}

type Client struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Os       string   `json:"os"`
	OsArch   string   `json:"os_arch"`
	OsFamily string   `json:"os_family"`
	OsKernel string   `json:"os_kernel"`
	Hostname string   `json:"hostname"`
	Ipv4     []string `json:"ipv4"`
	Ipv6     []string `json:"ipv6"`
	Tags     []string `json:"tags"`
	Version  string   `json:"version"`
	Address  string   `json:"address"`
	Tunnels  []Tunnel
}

type ClientsResponse struct {
	Data []Client
}

func (rp *Rport) Clients(ctx context.Context) (cr ClientsResponse, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.JoinURL(rp.BaseURL, ClientsUrl),
		nil,
	)
	if err != nil {
		return
	}

	cl := &BaseClient{}
	cl.WithAuth(rp.Auth)

	err = cl.Call(req, &cr)
	return
}
