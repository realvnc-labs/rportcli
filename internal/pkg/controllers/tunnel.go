package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

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
	RenderTunnel(t output.KvProvider) error
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

func (tc *TunnelController) Tunnels(ctx context.Context) error {
	clResp, err := tc.Rport.Clients(ctx)
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

	return tc.TunnelRenderer.RenderTunnels(tunnels)
}

func (tc *TunnelController) Delete(ctx context.Context, params *options.ParameterBag) error {
	clientID := params.ReadString(ClientID, "")
	tunnelID := params.ReadString(TunnelID, "")
	clientName := params.ReadString(ClientNameFlag, "")

	if clientID == "" && clientName == "" {
		return errors.New("no client id nor name provided")
	}

	if clientID == "" {
		clients, err := tc.ClientSearch.Search(ctx, clientName, params)
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

	err := tc.Rport.DeleteTunnel(ctx, clientID, tunnelID)
	if err != nil {
		return err
	}

	err = tc.TunnelRenderer.RenderDelete(&models.OperationStatus{Status: "OK"})
	if err != nil {
		return err
	}

	return nil
}

func (tc *TunnelController) Create(ctx context.Context, params *options.ParameterBag) error {
	var err error
	clientID := params.ReadString(ClientID, "")
	clientName := params.ReadString(ClientNameFlag, "")
	if clientID == "" && clientName == "" {
		return errors.New("no client id nor name provided")
	}

	if clientID == "" {
		clientID, err = tc.findClientID(ctx, clientName, params)
		if err != nil {
			return err
		}
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
	if (acl == "" || acl == DefaultACL) && tc.IPProvider != nil {
		ip, e := tc.IPProvider.GetIP(context.Background())
		if e != nil {
			logrus.Errorf("failed to fetch IP: %v", e)
		} else {
			acl = ip
		}
	}

	checkPort := params.ReadString(CheckPort, "")
	tunResp, err := tc.Rport.CreateTunnel(ctx, clientID, local, remote, scheme, acl, checkPort)
	if err != nil {
		return err
	}
	tunnelCreated := tunResp.Data
	tunnelCreated.Usage = tc.generateUsage(tunnelCreated, params)

	return tc.TunnelRenderer.RenderTunnel(tunnelCreated)
}

func (tc *TunnelController) generateUsage(tunnelCreated *models.TunnelCreated, params *options.ParameterBag) string {
	rportHost := params.ReadString(config.ServerURL, "")
	if rportHost == "" {
		return ""
	}

	rportURL, err := url.Parse(rportHost)
	if err != nil {
		logrus.Error(err)
	} else {
		rportHost = rportURL.Hostname()
	}

	if tunnelCreated.Lport != "" {
		return fmt.Sprintf("ssh -p %s %s -l ${USER}", tunnelCreated.Lport, rportHost)
	}

	return fmt.Sprintf("ssh %s -l ${USER}", rportHost)
}

func (tc *TunnelController) findClientID(ctx context.Context, clientName string, params *options.ParameterBag) (string, error) {
	clients, err := tc.ClientSearch.Search(ctx, clientName, params)
	if err != nil {
		return "", err
	}

	if len(clients) == 0 {
		return "", fmt.Errorf("unknown client '%s'", clientName)
	}

	if len(clients) != 1 {
		return "", fmt.Errorf("client identified by '%s' is ambiguous, use a more precise name or use the client id", clientName)
	}
	return clients[0].ID, nil
}
