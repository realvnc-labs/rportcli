package launcher

import (
	"os"
	"os/exec"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"
)

func LaunchSSHTunnel(tunnelCreated *models.TunnelCreated, sshParamsFlat string) error {
	logrus.Debugf("ssh arguments are provided: '%s', will start an ssh session", sshParamsFlat)
	sshStr := tunnelCreated.RportServer + " " + strings.TrimSpace(sshParamsFlat) + " -p " + tunnelCreated.Lport
	sshParams := strings.Split(sshStr, " ")
	logrus.Debugf("will execute 'ssh %s'", sshStr)

	return runSSH(sshParams)
}

var ExecCommand = exec.Command

func runSSH(args []string) error {
	c := exec.Command("ssh", args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	err := c.Run()
	logrus.Debugf("will run %s", c.String())
	if err != nil {
		return err
	}
	logrus.Debugf("finished run %s", c.String())

	return nil
}
