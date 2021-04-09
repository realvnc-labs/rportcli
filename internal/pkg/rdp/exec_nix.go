// +build linux

package rdp

func CommandProvider(filePath string) (cmd string, args []string) {
	return "remmina", []string{filePath}
}
