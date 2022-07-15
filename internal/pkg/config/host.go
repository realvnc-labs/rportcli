package config

import (
	"os"
	"os/user"
	"time"
)

type HostInfo struct {
	FetchedAt time.Time
	Username  string
	Machine   string
}

func GetHostInfo() (hostInfo *HostInfo, err error) {
	userInfo, err := user.Current()
	if err != nil {
		return nil, err
	}
	machine, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hostInfo = &HostInfo{
		FetchedAt: time.Now(),
		Username:  userInfo.Name,
		Machine:   machine,
	}
	return hostInfo, nil
}
