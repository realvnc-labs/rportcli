package config

import (
	"strconv"

	options "github.com/breathbath/go_utils/v2/pkg/config"
)

const (
	ClientNameFlag   = "name"
	ClientSearchFlag = "search"

	ClientIDs        = "cids"
	Command          = "command"
	Script           = "script"
	GroupIDs         = "gids"
	Timeout          = "timeout"
	ExecConcurrently = "conc"
	AbortOnError     = "abort"
	Cwd              = "cwd"
	IsSudo           = "is_sudo"
	Interpreter      = "interpreter"
	IsFullOutput     = "full-command-response"

	DefaultCmdTimeoutSeconds = 30
)

func GetNoPromptFlagSpec() (flagSpec ParameterRequirement) {
	return ParameterRequirement{
		Field: NoPrompt,
		Help:  "Flag to disable prompting when missing values",
		// TODO: is it just authentication parameters?
		Description: "Never prompt when missing parameters",
		ShortName:   "q",
		Type:        BoolRequirementType,
		Default:     false,
	}
}

func GetReadYAMLFlagSpec() (flagSpec ParameterRequirement) {
	return ParameterRequirement{
		Field:       ReadYAML,
		Description: "Read parameters from a YAML file",
		ShortName:   "y",
		Type:        StringSliceRequirementType,
	}
}

func GetCommandFlagSpecs() (flagSpecs []ParameterRequirement) {
	return []ParameterRequirement{
		GetNoPromptFlagSpec(),
		GetReadYAMLFlagSpec(),
		{
			Field:    ClientIDs,
			Help:     "Enter comma separated client IDs",
			Validate: RequiredValidate,
			Description: "Comma separated client ids for which the command should be executed. " +
				"Alternatively use -n to execute a command by client name(s), or use --search flag.",
			ShortName: "d",
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(ClientNameFlag, "") == "" && providedParams.ReadString(ClientSearchFlag, "") == ""
			},
		},
		{
			Field:       ClientNameFlag,
			Description: "Comma separated client names for which the command should be executed",
			ShortName:   "n",
		},
		{
			Field:       ClientSearchFlag,
			Description: "Search clients on all fields, supports wildcards (*).",
		},
		{
			Field:       Command,
			Help:        "Enter command",
			Description: "Command which should be executed on the clients",
			ShortName:   "c",
			IsRequired:  true,
		},
		{
			Field:       Timeout,
			Help:        "Enter timeout in seconds",
			Description: "timeout in seconds that was used to observe the command execution",
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
			Help:        "execute the command concurrently on multiple clients",
			Description: "execute the command concurrently on multiple clients",
			ShortName:   "r",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       IsFullOutput,
			Help:        "output detailed information of a job execution",
			Description: "output detailed information of a job execution",
			ShortName:   "f",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       IsSudo,
			Help:        "should execute command as sudo",
			Description: "execute command as sudo",
			ShortName:   "u",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       Interpreter,
			Help:        "enter interpreter/shell name for the command execution",
			Description: "interpreter/shell name for the command execution",
			ShortName:   "i",
			Type:        StringRequirementType,
		},
		{
			Field:       AbortOnError,
			Description: "if true and command fails on one client, it's not executed on others",
			Help:        "should abort command if it fails on any client",
			ShortName:   "a",
			Type:        BoolRequirementType,
			Default:     false,
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

func GetScriptFlagSpecs() (flagSpecs []ParameterRequirement) {
	return []ParameterRequirement{
		GetNoPromptFlagSpec(),
		GetReadYAMLFlagSpec(),
		{
			Field:    ClientIDs,
			Help:     "Enter comma separated client IDs",
			Validate: RequiredValidate,
			Description: "Comma separated client ids on which the script should be executed. " +
				"Alternatively use -n to execute a script by client name(s), or use --search flag.",
			ShortName: "d",
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(ClientNameFlag, "") == "" && providedParams.ReadString(ClientSearchFlag, "") == ""
			},
			IsRequired: true,
		},
		{
			Field:       ClientNameFlag,
			Description: "Comma separated client names on which the script should be executed",
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
