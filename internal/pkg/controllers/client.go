package controllers

import (
	"context"
	"fmt"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type ClientRenderer interface {
	RenderClients(clients []*models.Client) error
	RenderClient(client *models.Client, renderDetails bool) error
}

type ClientController struct {
	Rport          *api.Rport
	ClientRenderer ClientRenderer
}

func (cc *ClientController) Clients(ctx context.Context, params *options.ParameterBag) error {
	clResp, err := cc.Rport.Clients(
		ctx,
		api.NewPaginationFromParams(params),
		api.NewFilters("*", params.ReadString(config.ClientSearchFlag, "")),
	)
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
		client, err := cc.Rport.Client(ctx, id)
		if err != nil {
			return err
		}
		return cc.ClientRenderer.RenderClient(client, renderDetails)
	}

	clients, err := cc.Rport.Clients(ctx, api.NewPaginationWithLimit(2), api.NewFilters("name", name))
	if err != nil {
		return err
	}
	if len(clients.Data) < 1 {
		return fmt.Errorf("unknown client with name %q", name)
	}
	if len(clients.Data) > 1 {
		return fmt.Errorf("client with name %q is ambiguous, use a more precise name or use the client id", name)
	}

	return cc.ClientRenderer.RenderClient(clients.Data[0], renderDetails)
}
