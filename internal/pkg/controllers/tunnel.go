package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

const (
	ClientID   = "client"
	TunnelID   = "tunnel"
	Local      = "local"
	Remote     = "remote"
	Scheme     = "scheme"
	ACL        = "acl"
	CheckPort  = "checkp"
	DefaultACL = "<<YOU CURRENT PUBLIC IP>>"
)

type TunnelRenderer interface {
	RenderTunnels(tunnels []*models.Tunnel) error
	RenderTunnel(t *models.Tunnel) error
	RenderDelete(s output.KvProvider) error
}

type IPProvider interface {
	GetIP(ctx context.Context) (string, error)
}

type TunnelController struct {
	Rport          *api.Rport
	TunnelRenderer TunnelRenderer
	IPProvider     IPProvider
	ClientSearch   ClientSearch
}

func (cc *TunnelController) Tunnels(ctx context.Context) error {
	clResp, err := cc.Rport.Clients(ctx)
	if err != nil {
		return err
	}

	tunnels := make([]*models.Tunnel, 0)
	for _, cl := range clResp.Data {
		for _, t := range cl.Tunnels {
			t.ClientID = cl.ID
			t.ClientName = cl.Name
			tunnels = append(tunnels, t)
		}
	}

	return cc.TunnelRenderer.RenderTunnels(tunnels)
}

func (cc *TunnelController) Delete(ctx context.Context, clientID, clientName, tunnelID string) error {
	if clientID == "" && clientName == "" {
		return errors.New("no client id nor name provided")
	}

	if clientID == "" {
		clients, err := cc.ClientSearch.Search(ctx, clientName)
		if err != nil {
			return err
		}

		if len(clients) == 0 {
			return fmt.Errorf("unknown client '%s'", clientName)
		}

		if len(clients) != 1 {
			return fmt.Errorf("client identified by '%s' is ambiguous, use a more precise name or use the client id", clientName)
		}
		clientID = clients[0].ID
	}

	err := cc.Rport.DeleteTunnel(ctx, clientID, tunnelID)
	if err != nil {
		return err
	}

	err = cc.TunnelRenderer.RenderDelete(&models.OperationStatus{Status: "OK"})
	if err != nil {
		return err
	}

	return nil
}

func (cc *TunnelController) Create(ctx context.Context, params *options.ParameterBag) error {
	clientID := params.ReadString(ClientID, "")
	clientName := params.ReadString(ClientNameFlag, "")
	if clientID == "" && clientName == "" {
		return errors.New("no client id nor name provided")
	}

	if clientID == "" {
		clients, err := cc.ClientSearch.Search(ctx, clientName)
		if err != nil {
			return err
		}

		if len(clients) == 0 {
			return fmt.Errorf("unknown client '%s'", clientName)
		}

		if len(clients) != 1 {
			return fmt.Errorf("client identified by '%s' is ambiguous, use a more precise name or use the client id", clientName)
		}
		clientID = clients[0].ID
	}

	local := params.ReadString(Local, "")

	remote := params.ReadString(Remote, "")
	scheme := params.ReadString(Scheme, "")
	if scheme == "" {
		port, _ := utils.ExtractPortAndHost(remote)
		scheme = utils.GetSchemeByPort(port)
	}

	if remote == "" {
		remotePort := utils.GetPortByScheme(scheme)
		remote = strconv.Itoa(remotePort)
	}

	acl := params.ReadString(ACL, "")
	if (acl == "" || acl == DefaultACL) && cc.IPProvider != nil {
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

	return cc.TunnelRenderer.RenderTunnel(tun.Data)
}
