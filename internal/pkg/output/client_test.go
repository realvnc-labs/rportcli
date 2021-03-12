package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderClients(t *testing.T) {
	clients := []*models.Client{
		{
			ID:   "123",
			Name: "SomeName",
		},
		{
			ID:   "124",
			Name: "SomeOtherName",
		},
	}

	buf := &bytes.Buffer{}
	cr := &ClientRenderer{
		ColCountCalculator: func() int {
			return 150
		},
		Writer: buf,
	}

	err := cr.RenderClients(clients)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Clients ID NAME NUM TUNNELS REMOTE ADDRESS HOSTNAME OS KERNEL 123 SomeName 0 124 SomeOtherName 0",
		actualRenderResult,
	)
}

func TestRenderClient(t *testing.T) {
	client := &models.Client{
		ID:   "123",
		Name: "SomeName",
	}

	buf := &bytes.Buffer{}
	cr := &ClientRenderer{
		ColCountCalculator: func() int {
			return 150
		},
		Writer: buf,
	}

	err := cr.RenderClient(client)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	expectedResult := `Client [123] KEY VALUE ID: 123 Name: SomeName Os: OsArch: OsFamily: OsKernel: Hostname: Ipv4: Ipv6: Tags: Version: Address:`
	assert.Equal(t, expectedResult, actualRenderResult)
}
