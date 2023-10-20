package main

import (
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/elastic/go-sysinfo"
)

func GetOperatingSystem() *OperatingSystem {
	os := &OperatingSystem{
		Type:         runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
	log.Infof("Identifying OS (type: %s, architecture: %s)", os.Type, os.Architecture)
	host, err := sysinfo.Host()
	if err != nil {
		log.Errorf("Failed to lookup host info: %s", err)
	} else {
		info := host.Info()
		os.Name = info.OS.Name
		os.Version = info.OS.Version
		os.KernelVersion = info.KernelVersion
		log.Infof("Identified OS (type: %s, architecture: %s, name: %s, version: %s, kernel version: %s)", os.Type, os.Architecture, os.Name, os.Version, os.KernelVersion)
	}
	return os
}

type OperatingSystem struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	KernelVersion string `json:"kernel_version"`
	Architecture  string `json:"architecture"`
}
