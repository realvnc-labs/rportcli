package controllers

import (
	"context"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type ClientController struct {
	Rport *api.Rport
}

func (cc *ClientController) Clients(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	return output.RenderClients(rw, clResp.Data)
}
