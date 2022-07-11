package config

import (
	"strconv"
)

const (
	commandClientIDsDescription = "[required] Comma separated client ids for which the command should be executed. " +
		"Alternatively use -n to execute a command by client name(s), or use --search flag."
)

func GetCommandParamReqs() (paramReqs []ParameterRequirement) {
	return []ParameterRequirement{
		GetNoPromptParamReq(),
		GetReadYAMLParamReq(),
		GetClientIDsParamReq(commandClientIDsDescription),
		{
			Field:       ClientNameFlag,
			Description: "[deprecated] Comma separated client names for which the command should be executed",
			ShortName:   "n",
		},
		{
			Field:       ClientNamesFlag,
			Description: `Comma separated client names for which the command should be executed"`,
			ShortName:   "",
		},
		{
			Field:       ClientSearchFlag,
			Description: "Search clients on all fields, supports wildcards (*).",
		},
		{
			Field:       Command,
			Help:        "Enter command",
			Description: "[required] Command which should be executed on the clients",
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
