// +build windows

package utils

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	passBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil && strings.Contains(err.Error(), "The handle is invalid") {
		err = fmt.Errorf("your terminal does not support password promting, please use PowerShell or CMD or specify -p parameter explicitly")
	}

	return passBytes, err
}
