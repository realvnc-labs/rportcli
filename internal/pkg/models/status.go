package models

type Status struct {
	SessionsCount int    `json:"sessions_count"`
	Version       string `json:"version"`
	Fingerprint   string `json:"fingerprint"`
	ConnectURL    string `json:"connect_url"`
}
