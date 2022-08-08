package launcher

import (
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/exec"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/rdp"
)

func LaunchRDPTunnel(tunnelCreated *models.TunnelCreated, user string, height, width int) error {
	clientName := tunnelCreated.ClientName
	if clientName == "" {
		clientName = "client-id-" + tunnelCreated.ClientID
	}
	rdpFileInput := models.FileInput{
		Address:      fmt.Sprintf("%s:%s", tunnelCreated.RportServer, tunnelCreated.Lport),
		ScreenHeight: height,
		ScreenWidth:  width,
		UserName:     user,
		FileName:     fmt.Sprintf("%s.rdp", clientName),
	}
	fw := rdp.FileWriter{}
	filePath, err := fw.WriteRDPFile(rdpFileInput)
	if err != nil {
		return err
	}

	return exec.StartDefaultApp(filePath)
}
