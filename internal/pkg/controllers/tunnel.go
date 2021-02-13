package controllers

import (
	"context"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type TunnelRenderer interface {
	RenderTunnels(rw io.Writer, tunnels []*models.Tunnel) error
}

type TunnelController struct {
	Rport          *api.Rport
	TunnelRenderer TunnelRenderer
}

func (cc *TunnelController) Tunnels(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	tunnels := make([]*models.Tunnel, 0)
	for _, cl := range clResp.Data {
		tunnels = append(tunnels, cl.Tunnels...)
	}

	return cc.TunnelRenderer.RenderTunnels(rw, tunnels)
}

func (cc *TunnelController) Delete(ctx context.Context, rw io.Writer, clientID, tunnelID string) error {
	err := cc.Rport.DeleteTunnel(ctx, clientID, tunnelID)
	if err != nil {
		return err
	}

	err = output.RenderHeader(rw, "OK")
	if err != nil {
		return err
	}

	return nil
}
