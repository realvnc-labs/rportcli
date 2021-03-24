package output

import (
	"io"
	"regexp"
	"strings"

	"github.com/breathbath/go_utils/v2/pkg/testing"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

var columnsCountToTerminalWidthMap = []tableWidthColumnsCountMapping{
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
}

type CalcTerminalColumnsCount func() int

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

func calcColumnsCount(widthMapping []tableWidthColumnsCountMapping, calc CalcTerminalColumnsCount) int {
	if len(widthMapping) == 0 {
		return 0
	}

	actualTerminalWidth := calc()

	logrus.Debugf("actual terminal width is %d", actualTerminalWidth)

	for _, widthInfo := range widthMapping {
		if actualTerminalWidth <= widthInfo.minimalTableWidth {
			logrus.Debugf("will show %d columns", widthInfo.columnsCount)
			return widthInfo.columnsCount
		}
	}

	return 0
}

type RowData interface {
	Row() []string
}

type ColumnsData interface {
	Headers() []string
}

type KvProvider interface {
	KeyValues() []testing.KeyValueStr
}

func RenderTable(rw io.Writer, col ColumnsData, rowProviders []RowData, calc CalcTerminalColumnsCount) error {
	table := buildTable(rw)

	colsCount := calcColumnsCount(columnsCountToTerminalWidthMap, calc)
	allHeaders := col.Headers()

	if colsCount > len(allHeaders) || colsCount == 0 {
		colsCount = len(allHeaders)
	}
	table.SetHeader(allHeaders[0:colsCount])

	for _, rowProvider := range rowProviders {
		row := rowProvider.Row()
		if colsCount > len(row) || colsCount == 0 {
			colsCount = len(row)
		}
		table.Append(row[0:colsCount])
	}

	table.Render()

	return nil
}

func RenderHeader(rw io.Writer, header string) error {
	_, err := rw.Write([]byte(header + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func RenderKeyValues(rw io.Writer, kvP KvProvider) {
	tableClient := buildTable(rw)
	tableClient.SetHeader([]string{"KEY", "VALUE"})
	tableClient.SetAutoWrapText(false)

	for _, kv := range kvP.KeyValues() {
		tableClient.Append([]string{kv.Key + ":", kv.Value})
	}
	tableClient.Render()
}

func RemoveEmptySpaces(input string) string {
	r := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(r.ReplaceAllString(input, " "))
}
