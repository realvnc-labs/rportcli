// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package utils

import (
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() ([]byte, error) {
	return term.ReadPassword(syscall.Stdin)
}
