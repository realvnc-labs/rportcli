package models

import (
	"fmt"

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
	ID          string `json:"id"`
	ClientID    string `json:"client_id" yaml:"client_id"`
	ClientName  string `json:"client_name" yaml:"client_name"`
	Lhost       string `json:"local_host" yaml:"local_host"`
	Lport       string `json:"local_port" yaml:"local_port"`
	Rhost       string `json:"remote_host" yaml:"remote_host"`
	Rport       string `json:"remote_port" yaml:"remote_port"`
	LportRandom bool   `json:"local_port_random" yaml:"local_port_random"`
	Scheme      string `json:"scheme" yaml:"scheme"`
	ACL         string `json:"acl" yaml:"acl"`
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
	}
}
