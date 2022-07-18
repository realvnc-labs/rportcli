package main

import (
	"fmt"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/cmd"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		if err == controllers.ErrNoClientIDsToUse {
			// no client ids means no work, so exit with a regular message and a non-error code.
			msg := err.Error()
			displayMsg := strings.ToUpper(msg[:1]) + msg[1:]
			fmt.Println(displayMsg)
		} else {
			logrus.Fatal(err)
		}
	}
}
