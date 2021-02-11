package models

import "fmt"

type Tunnel struct {
	ID          string `json:"id"`
	Lhost       string `json:"lhost"`
	Lport       string `json:"lport"`
	Rhost       string `json:"rhost"`
	Rport       string `json:"rport"`
	LportRandom bool   `json:"lport_random"`
	Scheme      string `json:"scheme"`
	ACL         string `json:"acl"`
}

func (t *Tunnel) Headers() []string {
	return []string{
		"ID",
		"LHOST",
		"LPORT",
		"RHOST",
		"RPORT",
		"LPORTRAND",
		"SCHEME",
		"ACL",
	}
}

func (t *Tunnel) Row() []string {
	return []string{
		t.ID,
		t.Lhost,
		t.Lport,
		t.Rhost,
		t.Rport,
		fmt.Sprint(t.LportRandom),
		t.Scheme,
		t.ACL,
	}
}
