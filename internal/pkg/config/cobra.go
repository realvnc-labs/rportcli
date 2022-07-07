package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DefineCommandInputs(c *cobra.Command, reqs []ParameterRequirement) {
	for _, req := range reqs {
		defaultStr := ""
		if req.Default != nil {
			defaultStr = fmt.Sprint(req.Default)
		}

		switch req.Type {
		case BoolRequirementType:
			boolValDefault := true
			if defaultStr == "" || defaultStr == "0" || defaultStr == "false" {
				boolValDefault = false
			}
			c.Flags().BoolP(req.Field, req.ShortName, boolValDefault, req.Description)
		case IntRequirementType:
			defaultInt, err := strconv.Atoi(defaultStr)
			if err == nil {
				c.Flags().IntP(req.Field, req.ShortName, defaultInt, req.Description)
			} else {
				c.Flags().IntP(req.Field, req.ShortName, 0, req.Description)
			}
		case StringSliceRequirementType:
			c.Flags().StringSliceP(req.Field, req.ShortName, nil, req.Description)
		default:
			c.Flags().StringP(req.Field, req.ShortName, defaultStr, req.Description)
		}
	}
}

type FlagValuesProvider struct {
	flags *pflag.FlagSet
}

func CreateFlagValuesProvider(flags *pflag.FlagSet) options.ValuesProvider {
	return &FlagValuesProvider{flags: flags}
}

func (fvp *FlagValuesProvider) Dump(w io.Writer) (err error) {
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(fvp.ToKeyValues())
	return
}

func (fvp *FlagValuesProvider) ToKeyValues() map[string]interface{} {
	res := make(map[string]interface{})
	fvp.flags.VisitAll(func(flag *pflag.Flag) {
		res[flag.Name] = flag.Value.String()
	})

	return res
}

func (fvp *FlagValuesProvider) Read(name string) (val interface{}, found bool) {
	fl := fvp.flags.Lookup(name)
	if fl == nil {
		return nil, false
	}

	return fl.Value.String(), true
}

func (fvp *FlagValuesProvider) ReadFlag(reqField, reqType string) (result interface{}, isFound bool, err error) {
	flags := fvp.flags
	switch reqType {
	case BoolRequirementType:
		boolVal, e := flags.GetBool(reqField)
		if e != nil {
			return nil, false, e
		}
		return boolVal, true, nil
	case IntRequirementType:
		intVal, e := flags.GetInt(reqField)
		if e != nil {
			return nil, false, e
		}
		return intVal, true, nil
	case StringSliceRequirementType:
		sliceVal, e := flags.GetStringSlice(reqField)
		if e != nil {
			return nil, false, e
		}
		return sliceVal, true, nil
	default:
		strVal, e := flags.GetString(reqField)
		if e != nil {
			return nil, false, e
		}
		return strVal, true, nil
	}
}

func (fvp *FlagValuesProvider) ChangedFlag(reqField string) (isFound bool) {
	return fvp.flags.Changed(reqField)
}
