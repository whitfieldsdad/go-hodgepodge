package main

import (
	"github.com/charmbracelet/log"
)

type NetworkLocation struct {
	Hostname      string   `json:"hostname"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	MACAddress    string   `json:"mac_address"`
}

// GetNetworkLocation returns the system's network location (i.e. where it resides on the local network).
func GetNetworkLocation() (*NetworkLocation, error) {
	log.Info("Resolving network location...")
	pnic, err := GetPrimaryNetworkInterface()
	if err != nil {
		return nil, err
	}
	loc := &pnic.NetworkLocation
	loc.Hostname, _ = GetHostname()
	log.Info("Resolved network location")
	return loc, nil
}
