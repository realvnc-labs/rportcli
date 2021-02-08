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

func (t *Tunnel) Headers(count int) []string {
	allHeaders := []string{
		"ID",
		"LHOST",
		"LPORT",
		"RHOST",
		"RPORT",
		"LPORT_RAND",
		"SCHEME",
		"ACL",
	}

	if count > len(allHeaders) || count == 0 {
		count = len(allHeaders)
	}

	return allHeaders[0:count]
}

func (t *Tunnel) Row(count int) []string {
	allRowItems := []string{
		t.ID,
		t.Lhost,
		t.Lport,
		t.Rhost,
		t.Rport,
		fmt.Sprint(t.LportRandom),
		t.Scheme,
		t.ACL,
	}

	if count > len(allRowItems) || count == 0 {
		count = len(allRowItems)
	}

	return allRowItems[0:count]
}
