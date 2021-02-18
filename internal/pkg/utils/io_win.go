// +build windows

package utils

import (
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	return term.ReadPassword(int(syscall.Stdin))
}
