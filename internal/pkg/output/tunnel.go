package output

import (
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type TunnelRenderer struct {
	ColCountCalculator CalcTerminalColumnsCount
	Writer             io.Writer
	Format             string
}

func (tr *TunnelRenderer) RenderTunnels(tunnels []*models.Tunnel) error {
	return RenderByFormat(
		tr.Format,
		tr.Writer,
		tunnels,
		func() error {
			return tr.renderTunnelsInHumanFormat(tunnels)
		},
	)
}

func (tr *TunnelRenderer) renderTunnelsInHumanFormat(tunnels []*models.Tunnel) error {
	if len(tunnels) == 0 {
		return nil
	}

	err := RenderHeader(tr.Writer, "Tunnels")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(tunnels))
	for _, t := range tunnels {
		rowProviders = append(rowProviders, t)
	}

	return RenderTable(tr.Writer, &models.Tunnel{}, rowProviders, tr.ColCountCalculator)
}

func (tr *TunnelRenderer) RenderTunnel(t *models.Tunnel) error {
	return RenderByFormat(
		tr.Format,
		tr.Writer,
		t,
		func() error {
			return tr.renderTunnelInHumanFormat(t)
		},
	)
}

func (tr *TunnelRenderer) renderTunnelInHumanFormat(t *models.Tunnel) error {
	if t == nil {
		return nil
	}

	err := RenderHeader(tr.Writer, "Tunnel")
	if err != nil {
		return err
	}

	RenderKeyValues(tr.Writer, t)

	return nil
}

func (tr *TunnelRenderer) RenderDelete(os KvProvider) error {
	return RenderByFormat(
		tr.Format,
		tr.Writer,
		os,
		func() error {
			RenderKeyValues(tr.Writer, os)
			return nil
		},
	)
}
