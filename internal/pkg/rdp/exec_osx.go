// +build darwin

package rdp

func CommandProvider(filePath string) (cmd string, args []string) {
	return "open", []string{filePath}
}
