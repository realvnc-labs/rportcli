package models

import (
	"fmt"
	"strconv"

	"github.com/breathbath/go_utils/v2/pkg/testing"
)

type OperationStatus struct {
	Status string `json:"status"`
}

func (os *OperationStatus) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "Status",
			Value: os.Status,
		},
	}
}

type Tunnel struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id" yaml:"client_id"`
	ClientName      string `json:"client_name" yaml:"client_name"`
	Lhost           string `json:"lhost" yaml:"local_host"`
	Lport           string `json:"lport" yaml:"local_port"`
	Rhost           string `json:"rhost" yaml:"remote_host"`
	Rport           string `json:"rport" yaml:"remote_port"`
	LportRandom     bool   `json:"lport_random" yaml:"local_port_random"`
	Scheme          string `json:"scheme" yaml:"scheme"`
	ACL             string `json:"acl" yaml:"acl"`
	IdleTimeoutMins int    `json:"idle_timeout_minutes" yaml:"idle_timeout_minutes"`
}

func (t *Tunnel) Headers() []string {
	return []string{
		"ID",
		"CLIENT_ID",
		"CLIENT_NAME",
		"LOCAL_HOST",
		"LOCAL_PORT",
		"REMOTE_HOST",
		"REMOTE_PORT",
		"LOCAL_PORT_RAND",
		"SCHEME",
		"ACL",
		"TIMEOUT",
	}
}

func (t *Tunnel) Row() []string {
	return []string{
		t.ID,
		t.ClientID,
		t.ClientName,
		t.Lhost,
		t.Lport,
		t.Rhost,
		t.Rport,
		fmt.Sprint(t.LportRandom),
		t.Scheme,
		t.ACL,
		strconv.Itoa(t.IdleTimeoutMins),
	}
}

func (t *Tunnel) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "ID",
			Value: t.ID,
		},
		{
			Key:   "CLIENT_ID",
			Value: t.ClientID,
		},
		{
			Key:   "CLIENT_NAME",
			Value: t.ClientName,
		},
		{
			Key:   "LOCAL_HOST",
			Value: t.Lhost,
		},
		{
			Key:   "LOCAL_PORT",
			Value: t.Lport,
		},
		{
			Key:   "REMOTE_HOST",
			Value: t.Rhost,
		},
		{
			Key:   "REMOTE_PORT",
			Value: t.Rport,
		},
		{
			Key:   "LOCAL_PORT RANDOM",
			Value: fmt.Sprintf("%v", t.LportRandom),
		},
		{
			Key:   "SCHEME",
			Value: t.Scheme,
		},
		{
			Key:   "ACL",
			Value: t.ACL,
		},
		{
			Key:   "TIMEOUT MINUTES",
			Value: strconv.Itoa(t.IdleTimeoutMins),
		},
	}
}

type TunnelCreated struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id" yaml:"client_id"`
	ClientName      string `json:"client_name,omitempty" yaml:"client_name,omitempty"`
	Lhost           string `json:"lhost" yaml:"local_host"`
	Lport           string `json:"lport" yaml:"local_port"`
	Rhost           string `json:"rhost" yaml:"remote_host"`
	Rport           string `json:"rport" yaml:"remote_port"`
	LportRandom     bool   `json:"lport_random" yaml:"local_port_random"`
	Scheme          string `json:"scheme" yaml:"scheme"`
	ACL             string `json:"acl" yaml:"acl"`
	Usage           string `json:"usage" yaml:"usage"`
	IdleTimeoutMins int    `json:"idle_timeout_minutes" yaml:"idle_timeout_minutes"`
	RportServer     string `json:"rport_server,omitempty" yaml:"rport_server,omitempty"`
}

func (tc *TunnelCreated) KeyValues() []testing.KeyValueStr {
	kvs := []testing.KeyValueStr{
		{
			Key:   "ID",
			Value: tc.ID,
		},
		{
			Key:   "CLIENT_ID",
			Value: tc.ClientID,
		},
		{
			Key:   "LOCAL_HOST",
			Value: tc.Lhost,
		},
		{
			Key:   "LOCAL_PORT",
			Value: tc.Lport,
		},
		{
			Key:   "REMOTE_HOST",
			Value: tc.Rhost,
		},
		{
			Key:   "REMOTE_PORT",
			Value: tc.Rport,
		},
		{
			Key:   "LOCAL_PORT RANDOM",
			Value: fmt.Sprintf("%v", tc.LportRandom),
		},
		{
			Key:   "SCHEME",
			Value: tc.Scheme,
		},
		{
			Key:   "IDLE TIMEOUT MINS",
			Value: fmt.Sprint(tc.IdleTimeoutMins),
		},
		{
			Key:   "ACL",
			Value: tc.ACL,
		},
		{
			Key:   "USAGE",
			Value: tc.Usage,
		},
	}

	return kvs
}
