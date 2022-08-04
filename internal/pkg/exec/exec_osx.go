//go:build darwin
// +build darwin

package exec

const OpenCmd = "open"

func CommandProvider(filePath string) (cmd string, args []string) {
	return OpenCmd, []string{filePath}
}
