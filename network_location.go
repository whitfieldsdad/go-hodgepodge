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

func GetNetworkLocation() (*NetworkLocation, error) {
	log.Info("Resolving network location...")
	hostname, _ := GetHostname()
	loc := &NetworkLocation{
		Hostname: hostname,
	}
	nic, err := GetPrimaryNetworkInterface()
	if err != nil {
		log.Warn("Failed to get PNIC while resolving network location")
		log.Infof("Resolved network location (hostname: %s)", hostname)
	} else {
		loc.IPv4Addresses = nic.IPv4Addresses
		loc.IPv6Addresses = nic.IPv6Addresses
		loc.MACAddress = nic.MACAddress
		log.Infof("Resolved network location (hostname: %s, IPv4 addresses: %s, IPv6 addresses: %s, MAC address: %s)", hostname, nic.IPv4Addresses, nic.IPv6Addresses, nic.MACAddress)
	}
	return loc, nil
}
