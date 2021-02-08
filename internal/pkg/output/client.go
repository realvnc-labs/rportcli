package output

import (
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
	table.SetHeader(cl.HeadersShort(colsCount))

	for _, clnt := range clients {
		table.Append(clnt.RowShort(colsCount))
	}

	table.Render()

	return nil
}
