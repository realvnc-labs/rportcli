package config

import (
	"fmt"

	options "github.com/breathbath/go_utils/utils/config"
)

// Configuration constants
const (
	ServerURL        = "server_url"
	Login            = "login"
	Password         = "password"
	DefaultServerURL = "http://localhost:3000"
)

// Validate validation callback
type Validate func(fieldName string, val interface{}) error

// RequiredValidate validation logic for required field
var RequiredValidate = func(fieldName string, val interface{}) error {
	valStr, ok := val.(string)
	if val == nil || (ok && valStr == "") || fmt.Sprint(val) == "" {
		return fmt.Errorf("value '%s' is required and cannot be empty", fieldName)
	}
	return nil
}

// ParameterRequirement contains information about a parameter requirement
type ParameterRequirement struct {
	Field       string
	Help        string
	Validate    Validate
	Default     string
	Description string
}

// GetParameterRequirements returns required configuration options
func GetParameterRequirements() []ParameterRequirement {
	return []ParameterRequirement{
		{
			Field:       ServerURL,
			Help:        "Enter Server Url",
			Default:     DefaultServerURL,
			Description: "Server address of rport to connect to",
		},
		{
			Field:       Login,
			Help:        "Enter a valid login value",
			Validate:    RequiredValidate,
			Description: "Login to the rport server",
		},
		{
			Field:       Password,
			Help:        "Enter a valid password value",
			Validate:    RequiredValidate,
			Description: "Password to the rport server",
		},
	}
}

// GetNotMatchedRequirements reads parameters which are missing in the configuration or having a default value
func GetNotMatchedRequirements(params *options.ParameterBag) []ParameterRequirement {
	missedRequirements := make([]ParameterRequirement, 0)
	for _, req := range GetParameterRequirements() {
		paramInConfig, _ := params.Read(req.Field, nil)
		if req.Default != "" && paramInConfig == req.Default {
			missedRequirements = append(missedRequirements, req)
			continue
		}

		if req.Validate == nil {
			continue
		}
		err := req.Validate(req.Field, paramInConfig)
		if err != nil {
			missedRequirements = append(missedRequirements, req)
		}
	}

	return missedRequirements
}
