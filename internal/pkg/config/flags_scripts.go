package config

import (
	"strconv"

	options "github.com/breathbath/go_utils/v2/pkg/config"
)

const (
	scriptsClientIDsDescription = "[required] Comma separated client ids on which the script should be executed. " +
		"Alternatively use -n to execute a script by client name(s), or use --search flag."
)

func GetScriptParamReqs() (paramReqs []ParameterRequirement) {
	return []ParameterRequirement{
		GetNoPromptParamReq(),
		GetReadYAMLParamReq(),
		GetClientIDsParamReq(scriptsClientIDsDescription),
		{
			Field:       ClientNameFlag,
			Description: "[deprecated] Comma separated client names on which the script should be executed",
			ShortName:   "",
		},
		{
			Field:       ClientNamesFlag,
			Description: `Comma separated client names for which the command should be executed"`,
			ShortName:   "n",
		},
		{
			Field:       ClientSearchFlag,
			Description: "Search clients on all fields, supports wildcards (*).",
		},
		{
			Field:       Script,
			Help:        "Enter script path",
			Validate:    RequiredValidate,
			Description: "Path to the script file",
			ShortName:   "s",
			IsRequired:  true,
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(EmbeddedScript, "") == ""
			},
		},
		{
			Field:       EmbeddedScript,
			Help:        "Enter script content",
			Validate:    RequiredValidate,
			Description: "Script content to be executed on the clients",
			ShortName:   "c",
			IsRequired:  true,
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(Script, "") == ""
			},
		},
		{
			Field:       Timeout,
			Help:        "Enter timeout in seconds",
			Description: "timeout in seconds that was used to observe the script execution",
			Default:     strconv.Itoa(DefaultCmdTimeoutSeconds),
			ShortName:   "t",
		},
		{
			Field:       GroupIDs,
			Help:        "Enter comma separated group IDs",
			Description: "Comma separated client group IDs",
			ShortName:   "g",
		},
		{
			Field:       ExecConcurrently,
			Help:        "execute the script concurrently on multiple clients",
			Description: "execute the script concurrently on multiple clients",
			ShortName:   "r",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       IsFullOutput,
			Help:        "output detailed information of a script execution",
			Description: "output detailed information of a script execution",
			ShortName:   "f",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       IsSudo,
			Help:        "execute script as sudo",
			Description: "execute script as sudo",
			ShortName:   "u",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       Interpreter,
			Help:        "enter interpreter/shell name for the script execution",
			Description: "interpreter/shell name for the script execution",
			ShortName:   "i",
			Type:        StringRequirementType,
		},
		{
			Field:       Cwd,
			Help:        "enter current working directory",
			Description: "current working directory",
			ShortName:   "w",
			Type:        StringRequirementType,
			Default:     "",
		},
	}
}
