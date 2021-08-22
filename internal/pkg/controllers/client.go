package controllers

import (
	"context"
	"fmt"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

const (
	ClientNameFlag = "name"
)

type ClientRenderer interface {
	RenderClients(clients []*models.Client) error
	RenderClient(client *models.Client, renderDetails bool) error
}

type ClientController struct {
	Rport          *api.Rport
	ClientSearch   ClientSearch
	ClientRenderer ClientRenderer
}

func (cc *ClientController) Clients(ctx context.Context) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	return cc.ClientRenderer.RenderClients(clResp.Data)
}

func (cc *ClientController) Client(ctx context.Context, params *options.ParameterBag, id, name string) error {
	if id == "" && name == "" {
		return fmt.Errorf("no client id nor name provided")
	}

	renderDetails := params.ReadBool("all", false)

	if id != "" {
		clResp, err := cc.Rport.Clients(ctx)
		if err != nil {
			return err
		}
		for _, cl := range clResp.Data {
			if cl.ID == id {
				return cc.ClientRenderer.RenderClient(cl, renderDetails)
			}
		}
	} else {
		cl, err := cc.ClientSearch.FindOne(ctx, name, params)
		if err != nil {
			return err
		}
		return cc.ClientRenderer.RenderClient(cl, renderDetails)
	}

	return fmt.Errorf("client not found by the provided id '%s' or name '%s'", id, name)
}
