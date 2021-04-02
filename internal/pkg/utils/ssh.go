package utils

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func RunSSH(args []string) error {
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
