package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/utils/url"
)

const (
	ClientsURL       = "/api/v1/clients"
	ClientTunnelsURL = "/api/v1/clients/{client_id}/tunnels/{tunnel_id}"
	CreateTunnelURL  = "/api/v1/clients/{client_id}/tunnels"
)

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

	cl := &BaseClient{}
	cl.WithAuth(rp.Auth)

	cr = &ClientsResponse{}

	resp, err := cl.Call(req, cr)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			logrus.Error(closeErr)
		}
	}()

	return
}

type TunnelResponse struct {
	Data *models.Tunnel
}

func (rp *Rport) CreateTunnel(
	ctx context.Context,
	clientID, local, remote, scheme, acl, checkPort string,
) (tunResp *TunnelResponse, err error) {
	var req *http.Request
	u := strings.Replace(CreateTunnelURL, "{client_id}", clientID, 1)
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url.JoinURL(rp.BaseURL, u),
		nil,
	)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("local", local)
	q.Add("remote", remote)
	q.Add("scheme", scheme)
	q.Add("acl", acl)
	q.Add("check_port", checkPort)
	req.URL.RawQuery = q.Encode()

	cl := &BaseClient{}
	cl.WithAuth(rp.Auth)

	tunResp = &TunnelResponse{}

	resp, err := cl.Call(req, tunResp)
	if err != nil {
		return nil, err
	}
	defer closeRespBody(resp)

	return tunResp, nil
}

func (rp *Rport) DeleteTunnel(ctx context.Context, clientID, tunnelID string) (err error) {
	var req *http.Request
	u := strings.Replace(ClientTunnelsURL, "{client_id}", clientID, 1)
	u = strings.Replace(u, "{tunnel_id}", tunnelID, 1)
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url.JoinURL(rp.BaseURL, u),
		nil,
	)
	if err != nil {
		return
	}

	cl := &BaseClient{}
	cl.WithAuth(rp.Auth)

	resp, err := cl.Call(req, nil)

	if err != nil {
		return err
	}
	defer closeRespBody(resp)

	if resp.StatusCode == http.StatusNoContent {
		return
	}
	err = fmt.Errorf("unexpeted response code received: %d, expected code is %d", resp.StatusCode, http.StatusNoContent)

	return
}
