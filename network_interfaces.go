package main

import (
	"net"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type NetworkInterface struct {
	Name          string   `json:"name"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	MACAddress    string   `json:"mac_address"`
}

func ListNetworkInterfaces() ([]*NetworkInterface, error) {
	log.Info("Listing network interfaces")
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	nics := []*NetworkInterface{}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		nic := &NetworkInterface{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr.String(),
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}
			if ip.To4() != nil {
				nic.IPv4Addresses = append(nic.IPv4Addresses, ip.String())
			} else {
				nic.IPv6Addresses = append(nic.IPv6Addresses, ip.String())
			}
		}
		nics = append(nics, nic)
	}
	log.Infof("Found %d network interfaces", len(nics))
	return nics, nil
}

func GetPrimaryNetworkInterface() (*NetworkInterface, error) {
	log.Infof("Identifying primary network interface")
	pnic, err := getPrimaryNetworkInterface()
	if err != nil {
		return nil, err
	}
	log.Infof("Primary network interface is %s", pnic.Name)
	return pnic, nil
}

func getPrimaryNetworkInterface() (*NetworkInterface, error) {
	primaryIPv4Address, _ := GetPrimaryIPv4Address()
	primaryIPv6Address, _ := GetPrimaryIPv6Address()

	nics, err := ListNetworkInterfaces()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list network interfaces")
	}
	for _, nic := range nics {
		if nic.Name == "lo" {
			continue
		}
		if len(nic.IPv4Addresses) > 0 || len(nic.IPv6Addresses) > 0 {
			return nic, nil
		}
		for _, ipv4Address := range nic.IPv4Addresses {
			if ipv4Address == primaryIPv4Address.String() {
				return nic, nil
			}
		}
		for _, ipv6Address := range nic.IPv6Addresses {
			if ipv6Address == primaryIPv6Address.String() {
				return nic, nil
			}
		}
	}
	return nil, errors.New("no network interfaces found")
}
