package models

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
