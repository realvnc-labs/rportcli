package config

import (
	"fmt"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/spf13/cobra"
)

func DefineCommandInputs(c *cobra.Command, reqs []ParameterRequirement) {
	for _, req := range reqs {
		if req.Type == BoolRequirementType {
			boolValDefault := true
			if req.Default == "" || req.Default == "0" || req.Default == "false" {
				boolValDefault = false
			}
			c.Flags().BoolP(req.Field, req.ShortName, boolValDefault, req.Description)
		} else {
			c.Flags().StringP(req.Field, req.ShortName, req.Default, req.Description)
		}
	}
}

func LoadAllParams(
	c *cobra.Command,
	reqs []ParameterRequirement,
	promptReader PromptReader,
) (params *options.ParameterBag, err error) {
	configFromFile, err := LoadParamsFromFileAndEnv()
	if err != nil {
		return params, err
	}
	configFromCommandAndPrompt, err := CollectParamsFromCommandAndPrompt(c, reqs, promptReader)
	if err != nil {
		return params, err
	}

	configFromFile.MergeParameterBag(configFromCommandAndPrompt)

	return configFromFile, nil
}

func CollectParamsFromCommandAndPrompt(
	c *cobra.Command,
	reqs []ParameterRequirement,
	promptReader PromptReader,
) (params *options.ParameterBag, err error) {
	paramsRaw := make(map[string]string, len(reqs))
	for _, req := range reqs {
		if req.Type == BoolRequirementType {
			boolVal, e := c.Flags().GetBool(req.Field)
			if e != nil {
				return params, e
			}
			paramsRaw[req.Field] = fmt.Sprint(boolVal)
		} else {
			strVal, e := c.Flags().GetString(req.Field)
			if e != nil {
				return params, e
			}
			paramsRaw[req.Field] = strVal
		}
	}
	paramsFromInput := FromValues(paramsRaw)
	missedRequirements := CheckRequirements(paramsFromInput, reqs)
	if len(missedRequirements) == 0 {
		return paramsFromInput, nil
	}
	err = PromptRequiredValues(missedRequirements, paramsRaw, promptReader)
	if err != nil {
		return
	}
	params = FromValues(paramsRaw)

	return params, nil
}
