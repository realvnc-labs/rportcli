package config

import (
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/cli"
)

// Configuration constants
const (
	ServerURL        = "server_url"
	Login            = "login"
	Password         = "password"
	DefaultServerURL = "http://localhost:3000"
)

// GetParameterRequirements returns required configuration options
func GetParameterRequirements() []cli.ParameterRequirement {
	return []cli.ParameterRequirement{
		{
			Field:       ServerURL,
			Help:        "Enter Server Url",
			Default:     DefaultServerURL,
			Description: "Server address of rport to connect to",
			ShortName:   "s",
		},
		{
			Field:       Login,
			Help:        "Enter a valid login value",
			Validate:    cli.RequiredValidate,
			Description: "Login to the rport server",
			ShortName:   "l",
		},
		{
			Field:       Password,
			Help:        "Enter a valid password value",
			Validate:    cli.RequiredValidate,
			Description: "Password to the rport server",
			ShortName:   "p",
		},
	}
}
