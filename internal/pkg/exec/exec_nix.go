//go:build linux
// +build linux

package exec

const OpenCmd = "xdg-open"

func CommandProvider(filePath string) (cmd string, args []string) {
	return OpenCmd, []string{filePath}
}
