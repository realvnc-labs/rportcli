package output

import (
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type TunnelRenderer struct{}

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

	return RenderTable(rw, &models.Tunnel{}, rowProviders)
}
