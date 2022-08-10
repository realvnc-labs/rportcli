package launcher

import (
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/exec"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
)

func LaunchURITunnel(tunnelCreated *models.TunnelCreated, scheme string) error {
	// Create a URI to be opened by the default app of OS
	uri := fmt.Sprintf("%s://%s:%s", utils.GetHandlerByScheme(scheme), tunnelCreated.RportServer, tunnelCreated.Lport)
	return exec.StartDefaultApp(uri)
}
