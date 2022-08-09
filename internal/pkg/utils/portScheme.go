package utils

import (
	"fmt"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/exec"
)

type PortScheme struct {
	Port          int
	Scheme        string
	HandlerScheme string
}

const (
	SSH     = "ssh"
	RDP     = "rdp"
	VNC     = "vnc"
	HTTP    = "http"
	HTTPS   = "https"
	REALVNC = "realvnc"
)

var SupportedURIHandlers = []string{HTTP, HTTPS, VNC, REALVNC}
var PortSchemesMap = []PortScheme{
	{
		Port:   22,
		Scheme: SSH,
	},
	{
		Port:   3389,
		Scheme: RDP,
	},
	{
		Port:          5900,
		Scheme:        VNC,
		HandlerScheme: VNC,
	},
	{
		Port:          80,
		Scheme:        HTTP,
		HandlerScheme: HTTP,
	},
	{
		Port:          443,
		Scheme:        HTTPS,
		HandlerScheme: HTTPS,
	},
	{
		Port:          5900,
		Scheme:        REALVNC,
		HandlerScheme: "com.realvnc.vncviewer.connect",
	},
}

func GetPortByScheme(scheme string) int {
	for i := range PortSchemesMap {
		portScheme := PortSchemesMap[i]
		if portScheme.Scheme == scheme {
			return portScheme.Port
		}
	}

	return 0
}

func GetHandlerByScheme(scheme string) string {
	for _, i := range PortSchemesMap {
		if i.Scheme == scheme {
			return i.HandlerScheme
		}
	}

	return ""
}

// IsSupportedHandlerScheme check if a scheme can be opened by a default app.
// If not, return the list of supported schemes
func IsSupportedHandlerScheme(scheme string) (ok bool, supported []string) {
	var r []string
	for _, i := range PortSchemesMap {
		if i.HandlerScheme != "" {
			if i.Scheme == scheme {
				return true, []string{}
			}
			r = append(r, i.Scheme)
		}
	}

	return false, r
}

func GetUsageByScheme(scheme, host, port string) string {
	switch scheme {
	case "ssh":
		return fmt.Sprintf("ssh %s -p %s <...more ssh options>", host, port)
	case "rdp":
		return fmt.Sprintf("Connect remote desktop to remote pc '%s:%s'", host, port)
	case "realvnc":
		return fmt.Sprintf("Connect VNCViewer to VNCServer address '%s:%s'", host, port)
	case "vnc":
		return fmt.Sprintf("Connect a vnc viewer to server address '%s:%s'", host, port)
	case "http":
		return fmt.Sprintf("Open the following address with a browser 'http://%s:%s'", host, port)
	case "https":
		return fmt.Sprintf("Open the following address with a browser 'https://%s:%s'", host, port)
	}
	for _, i := range PortSchemesMap {
		if i.Scheme == scheme && i.HandlerScheme != "" {
			return fmt.Sprintf("%s %s://%s:%s", exec.OpenCmd, i.HandlerScheme, host, port)
		}
	}
	return fmt.Sprintf("Connect to '%s' on port '%s", host, port)
}
