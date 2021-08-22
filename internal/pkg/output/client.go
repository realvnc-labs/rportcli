package output

import (
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type ClientRenderer struct {
	ColCountCalculator CalcTerminalColumnsCount
	Writer             io.Writer
	Format             string
}

func (cr *ClientRenderer) RenderClients(clients []*models.Client) error {
	return RenderByFormat(
		cr.Format,
		cr.Writer,
		clients,
		func() error {
			return cr.renderClientsToHumanFormat(clients)
		},
	)
}

func (cr *ClientRenderer) renderClientsToHumanFormat(clients []*models.Client) error {
	err := RenderHeader(cr.Writer, "GetClients")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(clients))
	for _, cl := range clients {
		rowProviders = append(rowProviders, cl)
	}

	return RenderTable(cr.Writer, &models.Client{}, rowProviders, cr.ColCountCalculator)
}

func (cr *ClientRenderer) RenderClient(client *models.Client, renderDetails bool) error {
	return RenderByFormat(
		cr.Format,
		cr.Writer,
		client,
		func() error {
			return cr.renderClientToHumanFormat(client, renderDetails)
		},
	)
}

func (cr *ClientRenderer) renderClientToHumanFormat(client *models.Client, renderDetails bool) error {
	if client == nil {
		return nil
	}
	err := RenderHeader(cr.Writer, fmt.Sprintf("Client [%s]\n", client.ID))
	if err != nil {
		return err
	}

	RenderKeyValues(cr.Writer, client)

	if !renderDetails {
		return nil
	}

	err = cr.renderTunnels(client.Tunnels)
	if err != nil {
		return err
	}

	err = cr.renderUpdatesStatus(client.UpdatesStatus)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ClientRenderer) renderUpdatesStatus(updateStatus *models.UpdatesStatus) (err error) {
	if updateStatus == nil {
		return nil
	}

	err = RenderHeader(cr.Writer, "Update status:")
	if err != nil {
		return err
	}
	RenderKeyValues(cr.Writer, updateStatus)

	err = cr.renderUpdateStatusSummaries(updateStatus.UpdateSummaries)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ClientRenderer) renderUpdateStatusSummaries(summaries []models.UpdateSummary) (err error) {
	if len(summaries) == 0 {
		return nil
	}

	err = RenderHeader(cr.Writer, "\nUpdate summaries")
	if err != nil {
		return err
	}

	rows := make([]RowData, 0, len(summaries))
	for _, s := range summaries {
		rows = append(rows, &s)
	}

	return RenderTable(cr.Writer, &models.UpdateSummary{}, rows, cr.ColCountCalculator)
}

func (cr *ClientRenderer) renderTunnels(tunnels []*models.Tunnel) (err error) {
	if len(tunnels) == 0 {
		return nil
	}

	err = RenderHeader(cr.Writer, "\nTunnels")
	if err != nil {
		return err
	}

	rows := make([]RowData, 0, len(tunnels))
	for _, t := range tunnels {
		rows = append(rows, t)
	}

	return RenderTable(cr.Writer, &models.Tunnel{}, rows, cr.ColCountCalculator)
}
