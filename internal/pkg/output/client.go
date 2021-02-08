package output

import (
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type ClientRenderer struct{}

func (cr *ClientRenderer) RenderClients(rw io.Writer, clients []*models.Client) error {
	cl := &models.Client{}
	table := buildTable(rw)

	colsCount := calcColumnsCount([]tableWidthColumnsCountMapping{
		{
			minimalTableWidth: 70,
			columnsCount:      2,
		},
		{
			minimalTableWidth: 80,
			columnsCount:      3,
		},
		{
			minimalTableWidth: 100,
			columnsCount:      4,
		},
		{
			minimalTableWidth: 120,
			columnsCount:      5,
		},
	})
	table.SetHeader(cl.Headers(colsCount))

	for _, clnt := range clients {
		table.Append(clnt.Row(colsCount))
	}

	table.Render()

	return nil
}

func (cr *ClientRenderer) RenderClient(rw io.Writer, client *models.Client) error {
	if client == nil {
		return nil
	}

	caption := fmt.Sprintf(`Client [%s]
`, client.ID)

	_, err := rw.Write([]byte(caption))
	if err != nil {
		return err
	}

	tableClient := buildTable(rw)
	tableClient.SetHeader([]string{"KEY", "VALUE"})

	for _, kv := range client.ToKv("\n") {
		tableClient.Append([]string{kv.Key + ":", kv.Value})
	}
	tableClient.Render()
	if len(client.Tunnels) == 0 {
		return nil
	}

	tableTunnels := buildTable(rw)
	tunnelForHeaders := &models.Tunnel{}
	tableTunnels.SetHeader(tunnelForHeaders.Headers(0))

	caption2 := "\nTunnels\n"
	_, err = rw.Write([]byte(caption2))
	if err != nil {
		return err
	}

	for _, tunl := range client.Tunnels {
		tableTunnels.Append(tunl.Row(0))
	}
	tableTunnels.Render()

	return nil
}
