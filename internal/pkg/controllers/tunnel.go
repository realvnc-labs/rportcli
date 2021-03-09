package controllers

import (
	"context"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	options "github.com/breathbath/go_utils/utils/config"
	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

const (
	ClientID   = "clid"
	Local      = "local"
	Remote     = "remote"
	Scheme     = "scheme"
	ACL        = "acl"
	CheckPort  = "checkp"
	DefaultACL = "<<YOU CURRENT PUBLIC IP>>"
)

func GetCreateTunnelRequirements() []config.ParameterRequirement {
	return []config.ParameterRequirement{
		{
			Field:       ClientID,
			Description: "[required] unique client id retrieved previously",
			Validate:    config.RequiredValidate,
			ShortName:   "d",
			IsRequired:  true,
		},
		{
			Field: Local,
			Description: `refers to the ports of the rport server address to use for a new tunnel, e.g. '3390' or '0.0.0.0:3390'. 
If local is not specified, a random server port will be assigned automatically`,
			ShortName: "l",
		},
		{
			Field: Remote,
			Description: "[required] the ports are defined from the servers' perspective. " +
				"'Remote' refers to the ports and interfaces of the client., e.g. '3389'",
			ShortName:  "r",
			IsRequired: true,
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
			Description: "A flag whether to check availability of a public port. By default check is enabled. To disable it specify '0 or false'.",
			ShortName:   "p",
		},
	}
}

type TunnelRenderer interface {
	RenderTunnels(rw io.Writer, tunnels []*models.Tunnel) error
	RenderTunnel(rw io.Writer, t *models.Tunnel) error
}

type IPProvider interface {
	GetIP(ctx context.Context) (string, error)
}

type TunnelController struct {
	Rport          *api.Rport
	TunnelRenderer TunnelRenderer
	IPProvider     IPProvider
}

func (cc *TunnelController) Tunnels(ctx context.Context, rw io.Writer) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	tunnels := make([]*models.Tunnel, 0)
	for _, cl := range clResp.Data {
		for _, t := range cl.Tunnels {
			t.Client = cl.ID
			tunnels = append(tunnels, t)
		}
	}

	return cc.TunnelRenderer.RenderTunnels(rw, tunnels)
}

func (cc *TunnelController) Delete(ctx context.Context, rw io.Writer, clientID, tunnelID string) error {
	err := cc.Rport.DeleteTunnel(ctx, clientID, tunnelID)
	if err != nil {
		return err
	}

	err = output.RenderHeader(rw, "OK")
	if err != nil {
		return err
	}

	return nil
}

func (cc *TunnelController) Create(ctx context.Context, rw io.Writer, params *options.ParameterBag) error {
	clientID := params.ReadString(ClientID, "")
	local := params.ReadString(Local, "")
	remote := params.ReadString(Remote, "")
	scheme := params.ReadString(Scheme, "")
	acl := params.ReadString(ACL, "")
	if acl == "" || acl == DefaultACL {
		ip, e := cc.IPProvider.GetIP(context.Background())
		if e != nil {
			logrus.Errorf("failed to fetch IP: %v", e)
		} else {
			acl = ip
		}
	}

	checkPort := params.ReadString(CheckPort, "")
	tun, err := cc.Rport.CreateTunnel(ctx, clientID, local, remote, scheme, acl, checkPort)
	if err != nil {
		return err
	}

	return cc.TunnelRenderer.RenderTunnel(rw, tun.Data)
}
