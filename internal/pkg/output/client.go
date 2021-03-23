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

func (cr *ClientRenderer) RenderClient(client *models.Client) error {
	return RenderByFormat(
		cr.Format,
		cr.Writer,
		client,
		func() error {
			return cr.renderClientToHumanFormat(client)
		},
	)
}

func (cr *ClientRenderer) renderClientToHumanFormat(client *models.Client) error {
	if client == nil {
		return nil
	}
	err := RenderHeader(cr.Writer, fmt.Sprintf("Client [%s]\n", client.ID))
	if err != nil {
		return err
	}

	RenderKeyValues(cr.Writer, client)

	if len(client.Tunnels) == 0 {
		return nil
	}

	err = RenderHeader(cr.Writer, "\nTunnels")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(client.Tunnels))
	for _, t := range client.Tunnels {
		rowProviders = append(rowProviders, t)
	}

	return RenderTable(cr.Writer, &models.Tunnel{}, rowProviders, cr.ColCountCalculator)
}
