package exec

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func StartDefaultApp(filePath string) error {
	c := exec.Command(OpenCmd, filePath)

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
