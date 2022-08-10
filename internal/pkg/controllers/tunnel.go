package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/launcher"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
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
}

func (tc *TunnelController) Tunnels(ctx context.Context, params *options.ParameterBag) error {
	clResp, err := tc.Rport.Clients(
		ctx,
		api.NewPaginationFromParams(params),
		api.NewFilters(
			"id", params.ReadString(config.ClientID, ""),
			"name", params.ReadString(config.ClientNameFlag, ""),
			"*", params.ReadString(config.ClientSearchFlag, ""),
		),
	)
	if err != nil {
		return err
	}
	clients := clResp.Data

	tunnels := make([]*models.Tunnel, 0)
	for _, cl := range clients {
		for _, t := range cl.Tunnels {
			t.ClientID = cl.ID
			t.ClientName = cl.Name
			tunnels = append(tunnels, t)
		}
	}

	return tc.TunnelRenderer.RenderTunnels(tunnels)
}

func (tc *TunnelController) Delete(ctx context.Context, params *options.ParameterBag) error {
	clientID, _, err := tc.getClientIDAndClientName(ctx, params)
	if err != nil {
		return err
	}

	tunnelID := params.ReadString(config.TunnelID, "")
	err = tc.Rport.DeleteTunnel(ctx, clientID, tunnelID, params.ReadBool(config.ForceDeletion, false))
	if err != nil {
		if strings.Contains(err.Error(), "tunnel is still active") {
			return fmt.Errorf("%v, use -f to delete it anyway", err)
		}
		return err
	}

	err = tc.TunnelRenderer.RenderDelete(&models.OperationStatus{Status: "Tunnel successfully deleted"})
	if err != nil {
		return err
	}

	return nil
}

func (tc *TunnelController) getClientIDAndClientName(
	ctx context.Context,
	params *options.ParameterBag,
) (clientID, clientName string, err error) {
	clientID = params.ReadString(config.ClientID, "")
	clientName = params.ReadString(config.ClientNameFlag, "")
	if clientID == "" && clientName == "" {
		err = errors.New("no client id or name provided")
		return
	}
	if clientID != "" && clientName != "" {
		err = errors.New("both client id and name provided. Please provide one or the other")
		return
	}

	if clientID != "" {
		return
	}

	clients, err := tc.Rport.Clients(ctx, api.NewPaginationWithLimit(25), api.NewFilters("name", clientName))
	if err != nil {
		return
	}
	var client *models.Client
	numClients := len(clients.Data)
	maxClientsForSelection := 15
	switch {
	case numClients == 1:
		client = clients.Data[0]
	case numClients < 1:
		return "", "", fmt.Errorf("unknown client with name %q", clientName)
	case numClients > maxClientsForSelection:
		return "", "", fmt.Errorf(
			"client with name %q is ambiguous, use a more precise name or use the client id", clientName,
		)
	default:
		names := []string{}
		for _, client := range clients.Data {
			names = append(names, "- "+client.Name)
		}
		return "", "", fmt.Errorf(
			"client with name %q is ambiguous, use a more precise name or use the client id.\ndo you mean:\n%s",
			clientName,
			strings.Join(names, "\n"),
		)
	}
	return client.ID, client.Name, nil
}

func (tc *TunnelController) Create(ctx context.Context, params *options.ParameterBag) error {
	TunnelLauncher, err := launcher.NewTunnelLauncher(params)
	if err != nil {
		return err
	}
	clientID, clientName, err := tc.getClientIDAndClientName(ctx, params)
	if err != nil {
		return err
	}

	acl := params.ReadString(config.ACL, "")
	if (acl == "" || acl == config.DefaultACL) && tc.IPProvider != nil {
		ip, e := tc.IPProvider.GetIP(ctx)
		if e != nil {
			logrus.Errorf("failed to fetch IP: %v", e)
		} else {
			acl = ip
		}
	}

	// deconstruct the values of '-r, --remote' using either <IP address>:<PORT> e.g. 127.0.0.1:22
	// or just <PORT> e.g. 22.
	remotePortAndHostStr := params.ReadString(config.Remote, "")
	if TunnelLauncher.Scheme != "" && remotePortAndHostStr == "" {
		// if '-r, --remote' is not given, try to get the port from the scheme
		// If we have just a port convert back to string.
		// For the RPort server API a port without a host is sufficient
		// to create a tunnel to this port on localhost
		remotePortAndHostStr = strconv.Itoa(utils.GetPortByScheme(TunnelLauncher.Scheme))
	}

	local := params.ReadString(config.Local, "")
	checkPort := params.ReadString(config.CheckPort, "")
	skipIdleTimeout := params.ReadBool(config.SkipIdleTimeout, false)
	idleTimeoutMinutes := 0
	useHTTPProxy := params.ReadBool(config.UseHTTPProxy, false)
	if !skipIdleTimeout {
		idleTimeoutMinutes = params.ReadInt(config.IdleTimeoutMinutes, 0)
	}
	tunResp, err := tc.Rport.CreateTunnel(
		ctx,
		clientID,
		local,
		remotePortAndHostStr,
		TunnelLauncher.Scheme,
		acl,
		checkPort,
		idleTimeoutMinutes,
		skipIdleTimeout,
		useHTTPProxy,
	)
	if err != nil {
		return err
	}
	// Map the API response to the struct
	tunnelCreated := tunResp.Data
	// Enrich with local data
	tunnelCreated.RportServer = tc.Rport.BaseURL
	tunnelCreated.ClientID = clientID
	if clientName != "" {
		tunnelCreated.ClientName = clientName
	}
	tc.getRportServerName(tunnelCreated)
	tunnelCreated.Usage = utils.GetUsageByScheme(tunnelCreated.Scheme, tunnelCreated.RportServer, tunnelCreated.Lport)

	err = tc.TunnelRenderer.RenderTunnel(tunnelCreated)
	if err != nil {
		return err
	}
	del, err := TunnelLauncher.Execute(tunnelCreated)
	if del {
		deleteTunnelParams := options.New(options.NewMapValuesProvider(map[string]interface{}{
			config.ClientID: tunnelCreated.ClientID,
			config.TunnelID: tunnelCreated.ID,
		}))
		return tc.Delete(ctx, deleteTunnelParams)
	}
	return err
}

// getRportServerName extracts just the server name from the Rport API URL
func (tc *TunnelController) getRportServerName(tunnelCreated *models.TunnelCreated) {
	var rportURL *url.URL
	rportURL, err := url.Parse(tunnelCreated.RportServer)
	if err != nil {
		return
	}
	tunnelCreated.RportServer = rportURL.Hostname()
	// @todo: Get the tunnel host from the API. Tunnel host can differ from API host
}
