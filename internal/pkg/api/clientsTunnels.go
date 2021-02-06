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
	err = cl.Call(req, cr)
	return
}
