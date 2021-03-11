package controllers

import (
	"context"
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type ClientRenderer interface {
	RenderClients(clients []*models.Client) error
	RenderClient(client *models.Client) error
}

type ClientController struct {
	Rport          *api.Rport
	ClientRenderer ClientRenderer
}

func (cc *ClientController) Clients(ctx context.Context) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	return cc.ClientRenderer.RenderClients(clResp.Data)
}

func (cc *ClientController) Client(ctx context.Context, id string) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	for _, cl := range clResp.Data {
		if cl.ID == id {
			return cc.ClientRenderer.RenderClient(cl)
		}
	}

	fmt.Println("client not found")

	return nil
}
