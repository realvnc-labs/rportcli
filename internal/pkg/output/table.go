package output

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type tableWidthColumnsCountMapping struct {
	minimalTableWidth int
	columnsCount      int
}

func buildTable(rw io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(rw)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(true)

	return table
}

func calcColumnsCount(widthMapping []tableWidthColumnsCountMapping) int {
	if len(widthMapping) == 0 {
		return 0
	}

	actualTerminalWidth, _, err := terminal.GetSize(0)
	if err != nil {
		logrus.Warnf("cannot determine terminal width: %v", err)
		return 0
	}

	logrus.Debugf("actual terminal width is %d", actualTerminalWidth)

	for _, widthInfo := range widthMapping {
		if actualTerminalWidth <= widthInfo.minimalTableWidth {
			logrus.Debugf("will show %d columns", widthInfo.columnsCount)
			return widthInfo.columnsCount
		}
	}

	return 0
}
