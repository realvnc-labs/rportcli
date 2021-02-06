package models

type User struct {
	User   string   `json:"user"`
	Groups []string `json:"groups"`
}
