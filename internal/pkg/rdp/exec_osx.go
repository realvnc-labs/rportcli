//go:build darwin

// TODO: check with TH about the above change

package rdp

func CommandProvider(filePath string) (cmd string, args []string) {
	return "open", []string{filePath}
}
