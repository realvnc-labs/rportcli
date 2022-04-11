package api

import (
	"context"
	"fmt"
	"net/http"
	url2 "net/url"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/breathbath/go_utils/v2/pkg/url"
)

const (
	ClientsURL          = "/api/v1/clients"
	ClientURL           = "/api/v1/clients/%s"
	ClientsLimitDefault = 50
	ClientsLimitMax     = 500
)

type ClientsResponse struct {
	Data []*models.Client
}

func (rp *Rport) Clients(ctx context.Context, pagination Pagination, filters Filters) (cr *ClientsResponse, err error) {
	var req *http.Request
	u, err := url2.Parse(url.JoinURL(rp.BaseURL, ClientsURL))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("fields[clients]", "id,name,timezone,tunnels,address,hostname,os_kernel,connection_state,disconnected_at")
	pagination.Apply(q)
	filters.Apply(q)
	u.RawQuery = q.Encode()

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return
	}

	cr = &ClientsResponse{}
	_, err = rp.CallBaseClient(req, cr)

	return
}

type ClientResponse struct {
	Data *models.Client
}

func (rp *Rport) Client(ctx context.Context, id string) (*models.Client, error) {
	var req *http.Request
	u, err := url2.Parse(url.JoinURL(rp.BaseURL, fmt.Sprintf(ClientURL, id)))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	u.RawQuery = q.Encode()

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	cr := &ClientResponse{}
	_, err = rp.CallBaseClient(req, cr)
	if err != nil {
		return nil, err
	}

	return cr.Data, nil
}
