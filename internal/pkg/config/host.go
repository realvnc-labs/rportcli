package config

import (
	"os"
	"time"
)

type HostInfo struct {
	FetchedAt time.Time
	Username  string
	Machine   string
}

func GetHostInfo() (hostInfo *HostInfo, err error) {
	var name string
	name, err = userCurrent()
	if err != nil {
		return nil, err
	}
	machine, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hostInfo = &HostInfo{
		FetchedAt: time.Now(),
		Username:  name,
		Machine:   machine,
	}
	return hostInfo, nil
}
