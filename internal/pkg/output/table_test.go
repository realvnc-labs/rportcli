package output

import (
	"bytes"
	"testing"

	testing2 "github.com/breathbath/go_utils/v2/pkg/testing"
	"github.com/stretchr/testify/assert"
)

type KVProviderStub struct {
}

func (kv KVProviderStub) KeyValues() []testing2.KeyValueStr {
	return []testing2.KeyValueStr{
		{
			Key:   "one",
			Value: "1",
		},
		{
			Key:   "two",
			Value: "2",
		},
	}
}

type ColumnsDataStub struct{}

func (cds ColumnsDataStub) Headers() []string {
	return []string{"col1", "col2", "col3"}
}

type RowDataStub struct{}

func (rds RowDataStub) Row() []string {
	return []string{"val1", "val2", "val3"}
}

func TestRenderHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	err := RenderHeader(buf, "some header")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "some header\n", buf.String())
}

func TestRenderKeyValues(t *testing.T) {
	buf := &bytes.Buffer{}
	RenderKeyValues(buf, KVProviderStub{})
	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(t, "KEY VALUE one: 1 two: 2", actualRenderResult)
}

func TestRenderTable(t *testing.T) {
	buf := &bytes.Buffer{}

	terminalWidthCalcFucn := func() int {
		return 150
	}
	err := RenderTable(buf, ColumnsDataStub{}, []RowData{RowDataStub{}}, terminalWidthCalcFucn)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(t, "COL1 COL2 COL3 val1 val2 val3", actualRenderResult)
}
