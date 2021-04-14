package utils

import (
	"strconv"
	"strings"
)

type PortScheme struct {
	Port   int
	Scheme string
}

const (
	SSH   = "ssh"
	RDP   = "rdp"
	VNC   = "vnc"
	HTTP  = "http"
	HTTPS = "https"
)

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
		Port:   5900,
		Scheme: VNC,
	},
	{
		Port:   80,
		Scheme: HTTP,
	},
	{
		Port:   443,
		Scheme: HTTPS,
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

func GetSchemeByPort(port int) string {
	for i := range PortSchemesMap {
		portScheme := PortSchemesMap[i]
		if portScheme.Port == port {
			return portScheme.Scheme
		}
	}

	return ""
}

func ExtractPortAndHost(input string) (port int, host string) {
	var portStr string
	var err error
	if strings.Contains(input, ":") {
		hostPortParts := strings.Split(input, ":")
		host = hostPortParts[0]
		portStr = hostPortParts[1]
	} else {
		portStr = input
	}

	port, err = strconv.Atoi(portStr)
	if err != nil {
		port = 0
	}

	return
}
