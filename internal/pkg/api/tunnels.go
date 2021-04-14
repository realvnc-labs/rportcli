package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	TunnelsURL      = "/api/v1/clients/{client_id}/tunnels/{tunnel_id}"
	CreateTunnelURL = "/api/v1/clients/{client_id}/tunnels"
)

type TunnelResponse struct {
	Data *models.Tunnel
}

type TunnelCreatedResponse struct {
	Data *models.TunnelCreated
}

func (rp *Rport) CreateTunnel(
	ctx context.Context,
	clientID, local, remote, scheme, acl, checkPort string,
) (tunResp *TunnelCreatedResponse, err error) {
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

	tunResp = &TunnelCreatedResponse{}
	_, err = rp.CallBaseClient(req, tunResp)

	return tunResp, err
}

func (rp *Rport) DeleteTunnel(ctx context.Context, clientID, tunnelID string, force bool) (err error) {
	var req *http.Request
	u := strings.Replace(TunnelsURL, "{client_id}", clientID, 1)
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
	if force {
		q := req.URL.Query()
		q.Add("force", "1")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := rp.CallBaseClient(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNoContent {
		return
	}
	err = fmt.Errorf("unexpeted result received")

	return
}
