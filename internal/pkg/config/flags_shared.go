package config

import options "github.com/breathbath/go_utils/v2/pkg/config"

const (
	ClientNameFlag   = "name"
	ClientNamesFlag  = "names"
	ClientSearchFlag = "search"

	ClientIDs        = "cids"
	Command          = "command"
	Script           = "script"
	EmbeddedScript   = "exec"
	GroupIDs         = "gids"
	Timeout          = "timeout"
	ExecConcurrently = "conc"
	AbortOnError     = "abort"
	Cwd              = "cwd"
	IsSudo           = "is_sudo"
	Interpreter      = "interpreter"
	IsFullOutput     = "full-command-response"
	WriteExecLog     = "write-execlog"
	ReadExecLog      = "read-execlog"

	ClientID           = "client"
	TunnelID           = "tunnel"
	Local              = "local"
	Remote             = "remote"
	Scheme             = "scheme"
	ACL                = "acl"
	CheckPort          = "checkp"
	IdleTimeoutMinutes = "idle-timeout-minutes"
	SkipIdleTimeout    = "skip-idle-timeout"
	LaunchSSH          = "launch-ssh"
	LaunchRDP          = "launch-rdp"
	RDPWidth           = "rdp-width"
	RDPHeight          = "rdp-height"
	RDPUser            = "rdp-user"
	DefaultACL         = "<<YOU CURRENT PUBLIC IP>>"
	ForceDeletion      = "force"

	DefaultCmdTimeoutSeconds = 30
)

func GetNoPromptParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       NoPrompt,
		Help:        "Flag to disable prompting when missing values or confirmations (will answer all y/n questions with y)",
		Description: "No prompting when missing values or confirmations (will answer all y/n questions with y)",
		ShortName:   "q",
		Type:        BoolRequirementType,
		Default:     false,
	}
}

func GetReadYAMLParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       ReadYAML,
		Description: "Read parameters from a YAML file (parameters from the command line will have priority)",
		ShortName:   "y",
		Type:        StringSliceRequirementType,
	}
}

func GetClientIDsParamReq(desc string) (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       ClientIDs,
		Help:        "Enter comma separated client IDs",
		Validate:    RequiredValidate,
		Description: desc,
		ShortName:   "d",
		IsEnabled: func(providedParams *options.ParameterBag) bool {
			return providedParams.ReadString(ClientNameFlag, "") == "" &&
				providedParams.ReadString(ClientNamesFlag, "") == "" &&
				providedParams.ReadString(ClientSearchFlag, "") == ""
		},
		IsRequired: true,
	}
}

func GetWriteExecutionLogParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       WriteExecLog,
		Help:        "keep a log of the current execution",
		Description: "write a log of the execution output",
		ShortName:   "",
		Type:        StringRequirementType,
		Default:     "",
	}
}

func GetReadExecutionLogParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field: ReadExecLog,
		Help:  "use client ids with failed runs from specified execution log",
		Description: "read execution log from which to extract failed client ids, " +
			"which will be used to target clients for the next run",
		ShortName: "",
		Type:      StringRequirementType,
		Default:   "",
	}
}
