package models

import (
	"strconv"
	"strings"

	"github.com/breathbath/go_utils/v2/pkg/testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

type Client struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Os             string   `json:"os"`
	OsArch         string   `json:"os_arch"`
	OsFamily       string   `json:"os_family"`
	OsKernel       string   `json:"os_kernel"`
	Hostname       string   `json:"hostname"`
	ConnState      string   `json:"connection_state"`
	DisconnectedAt string   `json:"disconnected_at"`
	ClientAuthID   string   `json:"client_auth_id"`
	Ipv4           []string `json:"ipv4"`
	Ipv6           []string `json:"ipv6"`
	Tags           []string `json:"tags"`
	Version        string   `json:"version"`
	Address        string   `json:"address"`
	Tunnels        []*Tunnel
}

func (c *Client) Headers() []string {
	return []string{
		"ID",
		"NAME",
		"TUNNELS",
		"REMOTE ADDRESS",
		"HOSTNAME",
		"OS_KERNEL",
		"S",
	}
}

func (c *Client) Row() []string {
	connState := ""
	if len(c.ConnState) > 0 {
		connState = c.ConnState[0:1]
	}
	return []string{
		c.ID,
		c.Name,
		strconv.Itoa(len(c.Tunnels)),
		utils.RemovePortFromURL(c.Address),
		c.Hostname,
		c.OsKernel,
		connState,
	}
}

func (c *Client) KeyValues() []testing.KeyValueStr {
	const sep = "\n"
	return []testing.KeyValueStr{
		{
			Key:   "ID",
			Value: c.ID,
		},
		{
			Key:   "Name",
			Value: c.Name,
		},
		{
			Key:   "Os",
			Value: c.Os,
		},
		{
			Key:   "OsArch",
			Value: c.OsArch,
		},
		{
			Key:   "OsFamily",
			Value: c.OsFamily,
		},
		{
			Key:   "OsKernel",
			Value: c.OsKernel,
		},
		{
			Key:   "Hostname",
			Value: c.Hostname,
		},
		{
			Key:   "Ipv4",
			Value: strings.Join(c.Ipv4, sep),
		},
		{
			Key:   "Ipv6",
			Value: strings.Join(c.Ipv6, sep),
		},
		{
			Key:   "Tags",
			Value: strings.Join(c.Tags, sep),
		},
		{
			Key:   "Version",
			Value: c.Version,
		},
		{
			Key:   "Address",
			Value: c.Address,
		},
		{
			Key:   "Connection State",
			Value: c.ConnState,
		},
		{
			Key:   "Disconnected At",
			Value: c.DisconnectedAt,
		},
		{
			Key:   "Client Auth ID",
			Value: c.ClientAuthID,
		},
	}
}
