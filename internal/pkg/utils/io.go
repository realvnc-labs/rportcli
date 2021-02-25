package utils

import (
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/sirupsen/logrus"
)

const (
	InterruptMessage = "command interrupted, exiting"
)

func CalcTerminalColumnsCount() int {
	actualTerminalWidth, _ := consolesize.GetConsoleSize()

	logrus.Debugf("actual terminal width is %d", actualTerminalWidth)

	return actualTerminalWidth
}
