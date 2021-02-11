package output

import (
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type ClientRenderer struct{}

func (cr *ClientRenderer) RenderClients(rw io.Writer, clients []*models.Client) error {
	err := RenderHeader(rw, "Clients")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(clients))
	for _, cl := range clients {
		rowProviders = append(rowProviders, cl)
	}

	return RenderTable(rw, &models.Client{}, rowProviders)
}

func (cr *ClientRenderer) RenderClient(rw io.Writer, client *models.Client) error {
	if client == nil {
		return nil
	}
	err := RenderHeader(rw, fmt.Sprintf("Client [%s]\n", client.ID))
	if err != nil {
		return err
	}

	RenderKeyValues(rw, client)

	if len(client.Tunnels) == 0 {
		return nil
	}

	err = RenderHeader(rw, "\nTunnels")
	if err != nil {
		return err
	}

	rowProviders := make([]RowData, 0, len(client.Tunnels))
	for _, t := range client.Tunnels {
		rowProviders = append(rowProviders, t)
	}

	return RenderTable(rw, &models.Tunnel{}, rowProviders)
}
