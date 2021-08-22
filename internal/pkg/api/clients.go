package api

import (
	"context"
	"net/http"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	ClientsURL = "/api/v1/clients"
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

	cr = &ClientsResponse{}
	_, err = rp.CallBaseClient(req, cr)

	return
}

func (rp *Rport) GetClients(ctx context.Context) (cls []*models.Client, err error) {
	var cr *ClientsResponse
	cr, err = rp.Clients(ctx)

	if err != nil {
		return
	}

	return cr.Data, nil
}
