package launcher

import (
	"fmt"
	"strings"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

type TunnelLauncher struct {
	SSHParamsFlat    string
	LaunchRDP        bool
	RDPHeight        int
	RDPWidth         int
	RDPUser          string
	LaunchURIHandler bool
	Scheme           string
	params           *options.ParameterBag
}

func NewTunnelLauncher(params *options.ParameterBag) (*TunnelLauncher, error) {
	tl := &TunnelLauncher{
		SSHParamsFlat:    params.ReadString(config.LaunchSSH, ""),
		LaunchRDP:        params.ReadBool(config.LaunchRDP, false),
		LaunchURIHandler: params.ReadBool(config.LaunchURIHandler, false),
		Scheme:           params.ReadString(config.Scheme, ""),
		RDPWidth:         params.ReadInt(config.RDPWidth, 0),  // Default already set by GetCreateTunnelParamReqs()
		RDPHeight:        params.ReadInt(config.RDPHeight, 0), // Default already set by GetCreateTunnelParamReqs()
		RDPUser:          params.ReadString(config.RDPUser, ""),
	}
	tl.params = params
	// Set some obvious defaults
	switch {
	case tl.LaunchRDP && tl.Scheme == "":
		tl.Scheme = utils.RDP
	case tl.SSHParamsFlat != "" && tl.Scheme == "":
		tl.Scheme = utils.SSH
	}
	// Validate the combination of parameters so tunnels that can't be launch won't be created.
	// tunnel controller create() exists if an error is returned here.
	switch {
	case utils.Empty2int(tl.SSHParamsFlat)+utils.Bool2int(tl.LaunchURIHandler)+utils.Bool2int(tl.LaunchRDP) > 1:
		// Catch more than one launcher
		return nil, fmt.Errorf(
			"conflict: only one launch parameter of '--launch-uri, --launch-rdp, --launch-ssh' allowed",
		)
	case tl.SSHParamsFlat != "" && tl.Scheme != utils.SSH:
		// Catch mismatch of launcher and scheme
		return nil, fmt.Errorf(
			"launching the ssh client on scheme '%s' is not supported", tl.Scheme,
		)
	case strings.HasPrefix(tl.SSHParamsFlat, "-p") || strings.Contains(tl.SSHParamsFlat, " -p"):
		// Reject ssh params that already contain the -p (port) parameter
		return nil, fmt.Errorf(
			"do not pass a port with '-p': port will be set dynamically",
		)
	case tl.LaunchRDP && tl.Scheme != utils.RDP:
		// Catch mismatch of launcher and scheme
		return nil, fmt.Errorf(
			"launching the remote desktop client on scheme '%s' is not supported", tl.Scheme,
		)
	case tl.LaunchURIHandler:
		// Catch unsupported schemes
		ok, supported := utils.IsSupportedHandlerScheme(tl.Scheme)
		if !ok {
			return nil, fmt.Errorf(
				"scheme '%s' is not supported to be handeled by a default app. supported schemes: %s",
				tl.Scheme,
				strings.Join(supported, ", "),
			)
		}
	}
	return tl, nil
}

// Execute will execute a tunnel launcher if applicable
func (tl *TunnelLauncher) Execute(tunnelCreated *models.TunnelCreated) (deleteAfter bool, launch error) {
	switch {
	case tl.SSHParamsFlat != "":
		// Launch SSH
		deleteAfter = true
		launch = LaunchSSHTunnel(tunnelCreated, tl.SSHParamsFlat)
	case tl.LaunchRDP:
		// Launch the remote desktop app
		deleteAfter = false
		launch = LaunchRDPTunnel(tunnelCreated, tl.RDPUser, tl.RDPHeight, tl.RDPWidth)
		return
	case tl.LaunchURIHandler:
		// Launch the default app by scheme
		deleteAfter = false
		launch = LaunchURITunnel(tunnelCreated, tl.Scheme)
		return
	default:
		// Do nothing after the tunnel is created
		deleteAfter = false
		launch = nil
	}
	return
}
