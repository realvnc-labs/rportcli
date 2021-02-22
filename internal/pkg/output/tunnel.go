package output

import (
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type TunnelRenderer struct {
	ColCountCalculator CalcTerminalColumnsCount
}

func (cr *TunnelRenderer) RenderTunnels(rw io.Writer, tunnels []*models.Tunnel) error {
	if len(tunnels) == 0 {
		return nil
	}

	err := RenderHeader(rw, "Tunnels")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(tunnels))
	for _, t := range tunnels {
		rowProviders = append(rowProviders, t)
	}

	return RenderTable(rw, &models.Tunnel{}, rowProviders, cr.ColCountCalculator)
}

func (cr *TunnelRenderer) RenderTunnel(rw io.Writer, t *models.Tunnel) error {
	if t == nil {
		return nil
	}

	err := RenderHeader(rw, "Created Tunnel")
	if err != nil {
		return err
	}

	RenderKeyValues(rw, t)

	return nil
}
