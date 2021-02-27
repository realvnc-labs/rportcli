package models

type Status struct {
	ClientsConnected    int    `json:"clients_connected"`
	ClientsDisconnected int    `json:"clients_disconnected"`
	Version             string `json:"version"`
	Fingerprint         string `json:"fingerprint"`
	ConnectURL          string `json:"connect_url"`
}
