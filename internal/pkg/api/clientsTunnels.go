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

	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("unexpeted response code received: %d, expected code is %d", resp.StatusCode, http.StatusNoContent)
	}

	return
}
