package config

import (
	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

const (
	CreateTunnelRemoteParamDescr = "[required] the ports are defined from the servers' perspective. " +
		"'Remote' refers to the ports and interfaces of the client., e.g. '3389'" +
		"It's required unless -s uses a well-known scheme (SSH, RDP, VNC, HTTP, HTTPS)." +
		"Additionally if -b or -d parameters are provided and port is not provided, " +
		"a default corresponding value will be used 22 for ssh and 3389 for rdp"

	CreateTunnelLong = `creates a new tunnel, e.g.
rportcli tunnel create -l 0.0.0.0:3394 -r 22 -d bc0b705d-b5fb-4df5-84e3-82dba437bbef -s ssh --acl 10.1.2.3
this example opens port 3394 on the rport server and forwards to port 22 of the client bc0b705d-b5fb-4df5-84e3-82dba437bbef
with ssh url scheme and an IP address 10:1:2:3 allowed to access the tunnel
`
	CreateTunnelLocalDescr = `refers to the ports of the rport server address to use for a new tunnel, e.g. '3390' or '0.0.0.0:3390'.
If local is not specified, a random server port will be assigned automatically`

	CreateTunnelLaunchSSHDescr = `Start the ssh client after the tunnel is established and close tunnel on ssh exit.
Any parameter passed are append to the ssh command. i.e. -b "-l root"`
)

func GetCreateTunnelParamReqs(isRDPUserRequired bool) []ParameterRequirement {
	return []ParameterRequirement{
		GetNoPromptParamReq(),
		GetClientIDForTunnelParamReq(),
		GetClientNameForTunnelParamReq(),
		{
			Field:       Local,
			Description: CreateTunnelLocalDescr,
			ShortName:   "l",
		},
		{
			Field:       Remote,
			Description: CreateTunnelRemoteParamDescr,
			ShortName:   "r",
			IsRequired:  true,
			Validate:    RequiredValidate,
			Help:        "Enter a remote port value",
			IsEnabled:   isRemoteEnabled,
		},
		{
			Field:       Scheme,
			Description: "URI scheme to be used. For example, 'ssh', 'rdp', etc.",
			ShortName:   "s",
		},
		{
			Field:       ACL,
			Description: "ACL, IP addresses who is allowed to use the tunnel. For example, '142.78.90.8,201.98.123.0/24,'",
			Default:     DefaultACL,
			ShortName:   "a",
		},
		{
			Field:       CheckPort,
			Description: "A flag whether to check availability of a public port. By default check is disabled.",
			ShortName:   "p",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       LaunchSSH,
			Description: CreateTunnelLaunchSSHDescr,
			ShortName:   "b",
			Type:        StringRequirementType,
		},
		{
			Field: LaunchRDP,
			Description: `Start the default RDP client after the tunnel is established, e.g. -d 1
Optionally pass the rdp-width and rdp-height params for RDP window size`,
			ShortName: "d",
			Type:      BoolRequirementType,
			Default:   false,
		},
		{
			Field:       RDPWidth,
			Description: `RDP window width`,
			ShortName:   "w",
			Type:        IntRequirementType,
			Default:     1024,
		},
		{
			Field:       RDPHeight,
			Description: `RDP window height`,
			ShortName:   "i",
			Type:        IntRequirementType,
			Default:     768,
		},
		{
			Field:       RDPUser,
			Description: `username for a RDP session`,
			ShortName:   "u",
			Type:        StringRequirementType,
			Help:        "Enter a RDP user name",
			IsEnabled:   func(providedParams *options.ParameterBag) bool { return isRDPUserRequired },
		},
		{
			Field:       SkipIdleTimeout,
			Description: `if given, a tunnel will be created without an idle timeout`,
			ShortName:   "k",
			Type:        BoolRequirementType,
			Default:     false,
		},
		{
			Field:       IdleTimeoutMinutes,
			Description: `timeout in minutes for idle tunnels to be closed`,
			ShortName:   "m",
			Type:        IntRequirementType,
		},
		{
			Field:       UseHTTPProxy,
			Description: "Use https proxy (note: -s or --scheme must be either http or https)",
			Type:        BoolRequirementType,
			Default:     false,
		},
	}
}

func isRemoteEnabled(providedParams *options.ParameterBag) bool {
	scheme := providedParams.ReadString(Scheme, "")
	if scheme != "" && utils.GetPortByScheme(scheme) > 0 {
		return false
	}

	launchSSH := providedParams.ReadString(LaunchSSH, "")
	if launchSSH != "" {
		return false
	}

	launchRDP := providedParams.ReadBool(LaunchRDP, false)
	return !launchRDP
}

func GetDeleteTunnelParamReqs() []ParameterRequirement {
	return []ParameterRequirement{
		{
			Field:       ClientID,
			Description: "[conditionally required] client id, if not provided, client name should be given",
			Validate:    RequiredValidate,
			ShortName:   "c",
			IsRequired:  true,
			IsEnabled: func(providedParams *options.ParameterBag) bool {
				return providedParams.ReadString(ClientNameFlag, "") == ""
			},
			Help: "Enter a client id",
		},
		{
			Field:       ClientNameFlag,
			Description: `client name, if no client id is provided`,
			ShortName:   "n",
		},
		{
			Field:       ForceDeletion,
			ShortName:   "f",
			Default:     false,
			Description: `force tunnel deletion if it has active connections`,
			Type:        BoolRequirementType,
		},
		{
			Field:       TunnelID,
			Description: "[required]  tunnel id to delete",
			ShortName:   "u", // t is used for timeout
			IsRequired:  true,
			Validate:    RequiredValidate,
			Help:        "Enter a tunnel id",
		},
	}
}

func GetClientIDForTunnelParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       ClientID,
		Description: "[conditionally required] client id, if not provided, client name should be given",
		Validate:    RequiredValidate,
		ShortName:   "c",
		IsRequired:  true,
		IsEnabled: func(providedParams *options.ParameterBag) bool {
			return providedParams.ReadString(ClientNameFlag, "") == ""
		},
		Help: "Enter a client ID",
	}
}

func GetClientNameForTunnelParamReq() (paramReq ParameterRequirement) {
	return ParameterRequirement{
		Field:       ClientNameFlag,
		Description: `client name, if no client id is provided`,
		ShortName:   "n",
	}
}
