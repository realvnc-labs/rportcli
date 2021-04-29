package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/rdp"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

const (
	ClientID      = "client"
	TunnelID      = "tunnel"
	Local         = "local"
	Remote        = "remote"
	Scheme        = "scheme"
	ACL           = "acl"
	CheckPort     = "checkp"
	LaunchSSH     = "launch-ssh"
	LaunchRDP     = "launch-rdp"
	RDPWidth      = "rdp-width"
	RDPHeight     = "rdp-height"
	RDPUser       = "rdp-user"
	DefaultACL    = "<<YOU CURRENT PUBLIC IP>>"
	ForceDeletion = "force"
)

type TunnelRenderer interface {
	RenderTunnels(tunnels []*models.Tunnel) error
	RenderTunnel(t output.KvProvider) error
	RenderDelete(s output.KvProvider) error
}

type IPProvider interface {
	GetIP(ctx context.Context) (string, error)
}

type RDPFileWriter interface {
	WriteRDPFile(fi models.FileInput) (filePath string, err error)
}

type TunnelController struct {
	Rport          *api.Rport
	TunnelRenderer TunnelRenderer
	IPProvider     IPProvider
	ClientSearch   ClientSearch
	SSHFunc        func(sshParams []string) error
	RDPWriter      RDPFileWriter
	RDPExecutor    *rdp.Executor
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

	err := tc.Rport.DeleteTunnel(ctx, clientID, tunnelID, params.ReadBool(ForceDeletion, false))
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
	clientID = params.ReadString(ClientID, "")
	clientName = params.ReadString(ClientNameFlag, "")
	if clientID == "" && clientName == "" {
		err = errors.New("no client id nor name provided")
		return
	}

	if clientID != "" {
		return
	}

	client, err := tc.ClientSearch.FindOne(ctx, clientName, params)
	if err != nil {
		return
	}

	return client.ID, clientName, nil
}

func (tc *TunnelController) Create(ctx context.Context, params *options.ParameterBag) error {
	clientID, clientName, err := tc.getClientIDAndClientName(ctx, params)
	if err != nil {
		return err
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

	remotePortAndHostStr, scheme, err := tc.resolveRemoteAddrAndScheme(params)
	if err != nil {
		return err
	}

	local := params.ReadString(Local, "")
	checkPort := params.ReadString(CheckPort, "")
	tunResp, err := tc.Rport.CreateTunnel(ctx, clientID, local, remotePortAndHostStr, scheme, acl, checkPort)
	if err != nil {
		return err
	}

	tunnelCreated := tunResp.Data
	tunnelCreated.Usage = tc.generateUsage(tunnelCreated, params)
	if tunnelCreated.ClientID == "" {
		tunnelCreated.ClientID = clientID
	}
	if tunnelCreated.ClientName == "" && clientName != "" {
		tunnelCreated.ClientID = clientName
	}

	if clientName == "" && tunnelCreated.ClientName != "" {
		clientName = tunnelCreated.ClientName
	}

	err = tc.TunnelRenderer.RenderTunnel(tunnelCreated)
	if err != nil {
		return err
	}

	launchSSHStr := params.ReadString(LaunchSSH, "")
	shouldLaunchRDP := params.ReadBool(LaunchRDP, false)

	return tc.launchHelperFlowIfNeeded(ctx, launchSSHStr, clientID, clientName, shouldLaunchRDP, tunnelCreated, params)
}

func (tc *TunnelController) resolveRemoteAddrAndScheme(params *options.ParameterBag) (remotePortAndHostStr, scheme string, err error) {
	remotePortAndHostStr = params.ReadString(Remote, "")
	remotePortInt, _ := utils.ExtractPortAndHost(remotePortAndHostStr)

	scheme = params.ReadString(Scheme, "")
	if scheme == "" && remotePortInt > 0 {
		scheme = utils.GetSchemeByPort(remotePortInt)
	}

	if scheme != "" && remotePortAndHostStr == "" {
		remotePortInt = utils.GetPortByScheme(scheme)
	}

	launchSSHStr := params.ReadString(LaunchSSH, "")
	if launchSSHStr != "" {
		if scheme == "" {
			scheme = utils.SSH
		}
		if scheme != utils.SSH {
			err = fmt.Errorf("scheme %s is not compatible with the %s option", scheme, LaunchSSH)
			return
		}
		if remotePortInt == 0 {
			remotePortInt = utils.GetPortByScheme(scheme)
		}
	}

	shouldLaunchRDP := params.ReadBool(LaunchRDP, false)
	if shouldLaunchRDP {
		if scheme == "" {
			scheme = utils.RDP
		}
		if scheme != utils.RDP {
			err = fmt.Errorf("scheme %s is not compatible with the %s option", scheme, LaunchRDP)
			return
		}
		if remotePortInt == 0 {
			remotePortInt = utils.GetPortByScheme(scheme)
		}
	}

	if remotePortAndHostStr == "" && remotePortInt > 0 {
		remotePortAndHostStr = strconv.Itoa(remotePortInt)
	}

	return remotePortAndHostStr, scheme, err
}

func (tc *TunnelController) launchHelperFlowIfNeeded(
	ctx context.Context,
	launchSSHStr, clientID, clientName string,
	shouldLaunchRDP bool,
	tunnelCreated *models.TunnelCreated,
	params *options.ParameterBag,
) error {
	if launchSSHStr == "" && !shouldLaunchRDP {
		return nil
	}

	if launchSSHStr != "" {
		deleteTunnelParams := options.New(options.NewMapValuesProvider(map[string]interface{}{
			ClientID: clientID,
			TunnelID: tunnelCreated.ID,
		}))
		return tc.startSSHFlow(ctx, tunnelCreated, params, deleteTunnelParams)
	}

	return tc.startRDPFlow(ctx, tunnelCreated, params, clientName, clientID)
}

func (tc *TunnelController) finishSSHFlow(ctx context.Context, deleteTunnelParams *options.ParameterBag, prevErr error) error {
	logrus.Debugf("will delete tunnel with params: %+v", deleteTunnelParams)
	deleteTunnelErr := tc.Delete(ctx, deleteTunnelParams)
	if prevErr == nil {
		return deleteTunnelErr
	}

	if deleteTunnelErr == nil {
		return prevErr
	}

	return fmt.Errorf("%v, %v", prevErr, deleteTunnelErr)
}

func (tc *TunnelController) startSSHFlow(
	ctx context.Context,
	tunnelCreated *models.TunnelCreated,
	params, deleteTunnelParams *options.ParameterBag,
) error {
	sshParamsFlat := params.ReadString(LaunchSSH, "")
	logrus.Debugf("ssh arguments are provided: '%s', will start an ssh session", sshParamsFlat)
	port, host, err := tc.extractPortAndHost(tunnelCreated, params)
	if err != nil {
		prevErr := fmt.Errorf("failed to parse rport URL '%s': %v", params.ReadString(config.ServerURL, ""), err)
		return tc.finishSSHFlow(ctx, deleteTunnelParams, prevErr)
	}

	if host == "" {
		return tc.finishSSHFlow(ctx, deleteTunnelParams, errors.New("failed to retrieve rport URL"))
	}
	sshStr := host

	if port != "" && !strings.Contains(sshParamsFlat, "-p") {
		sshStr += " -p " + port
	}
	sshStr += " " + strings.TrimSpace(sshParamsFlat)
	sshParams := strings.Split(sshStr, " ")

	logrus.Debugf("will execute ssh %s", sshStr)
	err = tc.SSHFunc(sshParams)

	return tc.finishSSHFlow(ctx, deleteTunnelParams, err)
}

func (tc *TunnelController) generateUsage(tunnelCreated *models.TunnelCreated, params *options.ParameterBag) string {
	port, host, err := tc.extractPortAndHost(tunnelCreated, params)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	if host == "" {
		return ""
	}

	if port != "" {
		return fmt.Sprintf("ssh -p %s %s -l ${USER}", port, host)
	}

	return fmt.Sprintf("ssh %s -l ${USER}", host)
}

func (tc *TunnelController) extractPortAndHost(
	tunnelCreated *models.TunnelCreated,
	params *options.ParameterBag,
) (port, host string, err error) {
	rportHost := params.ReadString(config.ServerURL, "")
	if rportHost == "" {
		return
	}

	var rportURL *url.URL
	rportURL, err = url.Parse(rportHost)
	if err != nil {
		return
	}

	host = rportURL.Hostname()

	if tunnelCreated.Lport != "" {
		port = tunnelCreated.Lport
	}

	return
}

func (tc *TunnelController) startRDPFlow(
	ctx context.Context,
	tunnelCreated *models.TunnelCreated,
	params *options.ParameterBag,
	clientName, clientID string,
) error {
	port, host, err := tc.extractPortAndHost(tunnelCreated, params)
	if err != nil {
		return err
	}

	if clientName == "" {
		logrus.Debug("since client name is not provided, will try to find a client by id " + clientID)
		clients, e := tc.ClientSearch.Search(ctx, clientID, params)
		if e != nil {
			return e
		}
		if len(clients) == 0 || clients[0].Name == "" {
			clientName = fmt.Sprint(time.Now().Unix())
		} else {
			clientName = clients[0].Name
		}

		logrus.Debugf("found client name %s", clientName)
	}

	rdpFileInput := models.FileInput{
		Address:      fmt.Sprintf("%s:%s", host, port),
		ScreenHeight: params.ReadInt(RDPHeight, 0),
		ScreenWidth:  params.ReadInt(RDPWidth, 0),
		UserName:     params.ReadString(RDPUser, ""),
		FileName:     fmt.Sprintf("%s.rdp", clientName),
	}

	filePath, err := tc.RDPWriter.WriteRDPFile(rdpFileInput)
	if err != nil {
		return err
	}

	return tc.RDPExecutor.StartRdp(filePath)
}
