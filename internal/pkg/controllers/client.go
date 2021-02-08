package controllers

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type ClientRenderer interface {
	RenderClients(rw io.Writer, clients []*models.Client) error
	RenderClient(rw io.Writer, client *models.Client) error
}

type ClientController struct {
	Rport          *api.Rport
	ClientRenderer ClientRenderer
}

func (cc *ClientController) Clients(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	return cc.ClientRenderer.RenderClients(rw, clResp.Data)
}

func (cc *ClientController) Client(ctx context.Context, id string, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	for _, cl := range clResp.Data {
		if cl.ID == id {
			return cc.ClientRenderer.RenderClient(rw, cl)
		}
	}

	fmt.Println("client not found")

	return nil
}
