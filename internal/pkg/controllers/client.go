package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

const (
	ClientNameFlag = "name"
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

func (cc *ClientController) Client(ctx context.Context, id, name string) error {
	if id == "" && name == "" {
		return fmt.Errorf("no client id nor name provided")
	}

	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	for _, cl := range clResp.Data {
		if id != "" && cl.ID == id {
			return cc.ClientRenderer.RenderClient(cl)
		}
		if name != "" && strings.HasPrefix(strings.ToLower(cl.Name), strings.ToLower(name)) {
			return cc.ClientRenderer.RenderClient(cl)
		}
	}

	return fmt.Errorf("client not found by the provided id '%s' or name '%s'", id, name)
}
