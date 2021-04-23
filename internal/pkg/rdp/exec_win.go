// +build windows

package rdp

func CommandProvider(filePath string) (cmd string, args []string) {
	return "cmd.exe", []string{"/C", "start", filePath}
}
