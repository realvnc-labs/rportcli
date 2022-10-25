package rdp

import (
	"os"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestWriteRdpFile(t *testing.T) {
	fileInput := models.FileInput{
		Address:      "node1.rport.io:63231",
		ScreenHeight: 600,
		ScreenWidth:  800,
		UserName:     "Monster",
	}

	writer := &FileWriter{}

	filePath, err := writer.WriteRDPFile(fileInput)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		e := os.Remove(filePath)
		if e != nil {
			logrus.Error(e)
		}
	}()

	fileContents, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	expectedContent := `screen mode id:i:1
use multimon:i:0
desktopwidth:i:800
desktopheight:i:600
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
full address:s:node1.rport.io:63231
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
username:s:Monster
`
	assert.Equal(t, expectedContent, string(fileContents))
}
