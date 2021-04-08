// +build windows

package rdp

func CommandProvider(filePath string) (cmd string, args []string) {
	return "start", []string{filePath}
}
