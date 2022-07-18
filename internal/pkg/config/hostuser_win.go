//go:build windows
// +build windows

package config

import (
	"errors"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

func userCurrent() (string, error) {
	pw_name := make([]uint16, 256)
	pwname_size := uint32(len(pw_name)) - 1
	err := windows.GetUserNameEx(windows.NameSamCompatible, &pw_name[0], &pwname_size)
	if err != nil {
		return "", errors.New("unable to get windows user name")
	}
	s := syscall.UTF16ToString(pw_name)
	u := filepath.Base(s)
	return u, nil
}
