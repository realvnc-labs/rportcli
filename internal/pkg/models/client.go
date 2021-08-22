package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/breathbath/go_utils/v2/pkg/testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

type UpdatesStatus struct {
	Refreshed                time.Time       `json:"refreshed"`
	UpdatesAvailable         int             `json:"updates_available"`
	SecurityUpdatesAvailable int             `json:"security_updates_available"`
	UpdateSummaries          []UpdateSummary `json:"update_summaries"`
	RebootPending            bool            `json:"reboot_pending"`
	Error                    string          `json:"error,omitempty"`
	Hint                     string          `json:"hint,omitempty"`
}

func (us *UpdatesStatus) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "Refreshed",
			Value: us.Refreshed.Format(time.RFC3339),
		},
		{
			Key:   "UpdatesAvailable",
			Value: strconv.Itoa(us.UpdatesAvailable),
		},
		{
			Key:   "SecurityUpdatesAvailable",
			Value: strconv.Itoa(us.SecurityUpdatesAvailable),
		},
		{
			Key:   "RebootPending",
			Value: fmt.Sprintf("%v", us.SecurityUpdatesAvailable),
		},
		{
			Key:   "Error",
			Value: us.Error,
		},
		{
			Key:   "Hint",
			Value: us.Hint,
		},
	}
}

type UpdateSummary struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	RebootRequired   bool   `json:"reboot_required"`
	IsSecurityUpdate bool   `json:"is_security_update"`
}

func (us *UpdateSummary) Headers() []string {
	return []string{
		"TITLE",
		"DESCRIPTION",
		"REBOOT REQUIRED",
		"IS SECURITY UPDATE",
	}
}

func (us *UpdateSummary) Row() []string {
	return []string{
		us.Title,
		us.Description,
		fmt.Sprintf("%v", us.RebootRequired),
		fmt.Sprintf("%v", us.IsSecurityUpdate),
	}
}

type Client struct {
	ID                     string         `json:"id"`
	Name                   string         `json:"name"`
	Os                     string         `json:"os"`
	OsArch                 string         `json:"os_arch"`
	OsFamily               string         `json:"os_family"`
	OsKernel               string         `json:"os_kernel"`
	Hostname               string         `json:"hostname"`
	ConnState              string         `json:"connection_state"`
	DisconnectedAt         string         `json:"disconnected_at"`
	ClientAuthID           string         `json:"client_auth_id"`
	Ipv4                   []string       `json:"ipv4"`
	Ipv6                   []string       `json:"ipv6"`
	Tags                   []string       `json:"tags"`
	Version                string         `json:"version"`
	Address                string         `json:"address"`
	Tunnels                []*Tunnel      `json:"tunnels"`
	OSFullName             string         `json:"os_full_name"`
	OSVersion              string         `json:"os_version"`
	OSVirtualizationSystem string         `json:"os_virtualization_system"`
	OSVirtualizationRole   string         `json:"os_virtualization_role"`
	CPUFamily              string         `json:"cpu_family"`
	CPUModel               string         `json:"cpu_model"`
	CPUModelName           string         `json:"cpu_model_name"`
	CPUVendor              string         `json:"cpu_vendor"`
	NumCPUs                int            `json:"num_cpus"`
	MemoryTotal            uint64         `json:"mem_total"`
	Timezone               string         `json:"timezone"`
	AllowedUserGroups      []string       `json:"allowed_user_groups"`
	UpdatesStatus          *UpdatesStatus `json:"updates_status"`
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
		connState = strings.ToUpper(c.ConnState[0:1])
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
		{
			Key:   "OSFullName",
			Value: c.OSFullName,
		},
		{
			Key:   "OSVersion",
			Value: c.OSVersion,
		},
		{
			Key:   "OSVirtualizationSystem",
			Value: c.OSVirtualizationSystem,
		},
		{
			Key:   "OSVirtualizationRole",
			Value: c.OSVirtualizationRole,
		},
		{
			Key:   "CPUFamily",
			Value: c.CPUFamily,
		},
		{
			Key:   "CPUModel",
			Value: c.CPUModel,
		},
		{
			Key:   "CPUModelName",
			Value: c.CPUModelName,
		},
		{
			Key:   "CPUVendor",
			Value: c.CPUVendor,
		},
		{
			Key:   "NumCPUs",
			Value: strconv.Itoa(c.NumCPUs),
		},
		{
			Key:   "MemoryTotal",
			Value: humanize.Bytes(c.MemoryTotal),
		},
		{
			Key:   "Timezone",
			Value: c.Timezone,
		},
		{
			Key:   "AllowedUserGroups",
			Value: strings.Join(c.AllowedUserGroups, sep),
		},
	}
}
