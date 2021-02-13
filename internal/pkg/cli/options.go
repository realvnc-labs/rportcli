package cli

import (
	"errors"
	"fmt"
	"strings"

	options "github.com/breathbath/go_utils/utils/config"
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
	ShortName   string
	Help        string
	Validate    Validate
	Default     string
	Description string
}

// CheckRequirements reads parameters which are missing in the configuration or having a default value
func CheckRequirements(params *options.ParameterBag, requirementsToCheck []ParameterRequirement) []ParameterRequirement {
	missedRequirements := make([]ParameterRequirement, 0)
	for _, req := range requirementsToCheck {
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

func CheckRequirementsError(params *options.ParameterBag, requirementsToCheck []ParameterRequirement) error {
	missedRequirements := make([]ParameterRequirement, 0, len(requirementsToCheck))
	for _, req := range requirementsToCheck {
		paramInConfig, _ := params.Read(req.Field, nil)
		if req.Validate == nil {
			continue
		}
		err := req.Validate(req.Field, paramInConfig)
		if err != nil {
			missedRequirements = append(missedRequirements, req)
		}
	}
	if len(missedRequirements) == 0 {
		return nil
	}
	errorStrs := make([]string, 0, len(missedRequirements))
	for _, missedRequirement := range missedRequirements {
		errorStrs = append(errorStrs, fmt.Sprintf("missing value for %s: %s", missedRequirement.Field, missedRequirement.Description))
	}

	return errors.New(strings.Join(errorStrs, "\n"))
}

// FromValues creates a parameter bag from provided values
func FromValues(inputParams map[string]string) (params *options.ParameterBag) {
	inputParamsI := map[string]interface{}{}
	for k, v := range inputParams {
		inputParamsI[k] = v
	}
	vp := options.NewMapValuesProvider(inputParamsI)

	return &options.ParameterBag{BaseValuesProvider: vp}
}
