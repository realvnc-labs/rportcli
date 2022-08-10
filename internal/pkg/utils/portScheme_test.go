package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPortByScheme(t *testing.T) {
	for _, i := range PortSchemesMap {
		assert.Equal(t, i.Port, GetPortByScheme(i.Scheme))
	}
}

func TestGetHandlerByScheme(t *testing.T) {
	for _, i := range PortSchemesMap {
		assert.Equal(t, i.HandlerScheme, GetHandlerByScheme(i.Scheme))
	}
}

func TestGetUsageByScheme(t *testing.T) {
	assert.Equal(
		t,
		"ssh example.com -p 22 <...more ssh options>",
		GetUsageByScheme("ssh", "example.com", "22"),
	)
	assert.Equal(
		t,
		"Connect remote desktop to remote pc 'example.com:3389'",
		GetUsageByScheme("rdp", "example.com", "3389"),
	)
	assert.Equal(
		t,
		"Connect a vnc viewer to server address 'example.com:5900'",
		GetUsageByScheme("vnc", "example.com", "5900"),
	)
	assert.Equal(
		t,
		"Connect VNCViewer to VNCServer address 'example.com:5900'",
		GetUsageByScheme("realvnc", "example.com", "5900"),
	)
}
