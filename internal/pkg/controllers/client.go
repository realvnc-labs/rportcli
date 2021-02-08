package controllers

import (
	"context"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type ClientRenderer interface {
	RenderClients(rw io.Writer, clients []*models.Client) error
}

type ClientController struct {
	Rport *api.Rport
	Cr ClientRenderer
}

func (cc *ClientController) Clients(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	return cc.Cr.RenderClients(rw, clResp.Data)
}
