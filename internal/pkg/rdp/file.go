package rdp

import (
	"io"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ScreenWidthPlaceholder  = "SCREEN_WIDTH"
	ScreenHeightPlaceholder = "SCREEN_HEIGHT"
	AddressPlaceholder      = "ADDRESS"
	UserNamePlaceholder     = "USER_NAME"
	defaultScreenWidth      = 1024
	defaultScreenHeight     = 768
)

const template = `screen mode id:i:1
use multimon:i:0
desktopwidth:i:{{SCREEN_WIDTH}}
desktopheight:i:{{SCREEN_HEIGHT}}
client bpp:i:32
winposstr:s:0,3,0,0,800,600a
compression:i:1
keyboardhook:i:2
audiocapturemode:i:0
videoplaybackmode:i:1
connection type:i:7
networkautodetect:i:1
bandwidthautodetect:i:1
displayconnectionbar:i:1
enableworkspacereconnect:i:0
disable wallpaper:i:0
allow font smoothing:i:0
allow desktop composition:i:0
disable full window drag:i:1
disable menu anims:i:1
disable themes:i:0
disable cursor setting:i:0
bitmapcachepersistenable:i:1
full address:s:{{ADDRESS}}
audiomode:i:2
redirectprinters:i:0
redirectcomports:i:0
redirectsmartcards:i:0
redirectclipboard:i:1
redirectposdevices:i:0
drivestoredirect:s:
autoreconnection enabled:i:1
authentication level:i:2
prompt for credentials:i:0
negotiate security layer:i:1
remoteapplicationmode:i:0
alternate shell:s:
shell working directory:s:
gatewayhostname:s:
gatewayusagemethod:i:4
gatewaycredentialssource:i:4
gatewayprofileusagemethod:i:0
promptcredentialonce:i:0
gatewaybrokeringtype:i:0
use redirection server name:i:0
rdgiskdcproxy:i:0
kdcproxyname:s:
username:s:{{USER_NAME}}
`

type FileInput struct {
	Address      string
	ScreenHeight int
	ScreenWidth  int
	UserName     string
}

func WriteRdpFile(fi FileInput, w io.Writer) error {
	if fi.ScreenWidth == 0 {
		fi.ScreenWidth = defaultScreenWidth
	}

	if fi.ScreenHeight == 0 {
		fi.ScreenHeight = defaultScreenHeight
	}

	placeholderValues := map[string]string{
		ScreenWidthPlaceholder:  strconv.Itoa(fi.ScreenWidth),
		ScreenHeightPlaceholder: strconv.Itoa(fi.ScreenHeight),
		AddressPlaceholder:      fi.Address,
		UserNamePlaceholder:     fi.UserName,
	}

	content := template
	for k, v := range placeholderValues {
		content = strings.ReplaceAll(content, "{{"+k+"}}", v)
	}

	logrus.Debugf("created a rdp file")
	logrus.Debug(content)

	_, err := w.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
