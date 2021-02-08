package models

import (
	"strconv"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

type Client struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Os       string   `json:"os"`
	OsArch   string   `json:"os_arch"`
	OsFamily string   `json:"os_family"`
	OsKernel string   `json:"os_kernel"`
	Hostname string   `json:"hostname"`
	Ipv4     []string `json:"ipv4"`
	Ipv6     []string `json:"ipv6"`
	Tags     []string `json:"tags"`
	Version  string   `json:"version"`
	Address  string   `json:"address"`
	Tunnels  []*Tunnel
}

func (c *Client) HeadersShort(count int) []string {
	allHeaders := []string{
		"ID",
		"NAME",
		"NUM_TUNNELS",
		"REMOTE ADDRESS",
		"HOSTNAME",
		"OS_KERNEL",
	}

	if count > len(allHeaders) || count == 0 {
		count = len(allHeaders)
	}

	return allHeaders[0 : count]
}

func (c *Client) RowShort(count int) []string {
	allRowItems := []string{
		c.ID,
		c.Name,
		strconv.Itoa(len(c.Tunnels)),
		utils.RemovePortFromURL(c.Address),
		c.Hostname,
		c.OsKernel,
	}

	if count > len(allRowItems) || count == 0 {
		count = len(allRowItems)
	}

	return allRowItems[0 : count]
}
