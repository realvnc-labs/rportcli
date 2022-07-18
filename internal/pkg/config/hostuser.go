//go:build darwin || linux
// +build darwin linux

package config

import "os/user"

func userCurrent() (string, error) {
	userInfo, err := user.Current()
	if err != nil {
		return "", err
	}
	return userInfo.Name, nil
}
