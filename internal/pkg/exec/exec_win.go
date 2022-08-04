//go:build windows
// +build windows

package exec

import "strings"

const OpenCmd = "cmd.exe /C start"

func CommandProvider(filePath string) (cmd string, args []string) {
	parts := strings.Split(OpenCmd, " ")
	return parts[0], []string{parts[1], parts[2], filePath}
}
