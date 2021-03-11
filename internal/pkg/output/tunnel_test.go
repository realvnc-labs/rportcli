package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTunnels(t *testing.T) {
	tunnels := []*models.Tunnel{
		{
			ID:          "id22",
			Lhost:       "lhost",
			Lport:       "123",
			Rhost:       "rhost",
			Rport:       "124",
			LportRandom: false,
			Scheme:      "ssh",
			ACL:         "0.0.0.0",
		},
	}

	buf := &bytes.Buffer{}
	tr := &TunnelRenderer{
		ColCountCalculator: func() int {
			return 150
		},
		Writer: buf,
		Format: FormatHuman,
	}

	err := tr.RenderTunnels(tunnels)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Tunnels ID CLIENT LHOST LPORT RHOST RPORT LPORTRAND SCHEME ACL id22 lhost 123 rhost 124 false ssh 0.0.0.0 ",
		actualRenderResult,
	)
}
func TestRenderTunnel(t *testing.T) {
	tunnel := &models.Tunnel{
		ID:          "id22",
		Lhost:       "lhost",
		Lport:       "123",
		Rhost:       "rhost",
		Rport:       "124",
		LportRandom: false,
		Scheme:      "ssh",
		ACL:         "0.0.0.0",
	}

	buf := &bytes.Buffer{}
	tr := &TunnelRenderer{
		ColCountCalculator: func() int {
			return 150
		},
		Writer: buf,
		Format: FormatHuman,
	}

	err := tr.RenderTunnel(tunnel)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Tunnel KEY VALUE ID: id22 CLIENT: LHOST: lhost LPORT: 123 RHOST: rhost RPORT: 124 LPORT RANDOM: false SCHEME: ssh ACL: 0.0.0.0 ",
		actualRenderResult,
	)
}
