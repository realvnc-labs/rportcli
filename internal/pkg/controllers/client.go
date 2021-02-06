package controllers

import (
	"context"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/olekukonko/tablewriter"
)

type ClientController struct {
	Rport *api.Rport
}

func (cc *ClientController) Clients(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	cl := models.Client{}
	table := tablewriter.NewWriter(rw)
	table.SetHeader(cl.HeadersShort(0))

	for _, clnt := range clResp.Data {
		table.Append(clnt.RowShort(0))
	}

	table.Render()

	return nil
}
