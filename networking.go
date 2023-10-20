package main

import (
	"errors"
	"net"
	"os"

	"github.com/charmbracelet/log"
)

func GetHostname() (string, error) {
	return os.Hostname()
}

func GetPrimaryIPAddress() (string, error) {
	ipv4, err := GetPrimaryIPv4Address()
	if err == nil {
		return ipv4.String(), nil
	}
	ipv6, err := GetPrimaryIPv6Address()
	if err == nil {
		return ipv6.String(), nil
	}
	return "", errors.New("failed to resolve primary IP address")
}

func GetPrimaryIPv4Address() (*net.IP, error) {
	log.Infof("Resolving primary IPv4 address...")
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	log.Infof("Resolved primary IPv4 address: %s", localAddr.IP)
	return &localAddr.IP, nil
}

func GetPrimaryIPv6Address() (*net.IP, error) {
	log.Infof("Resolving primary IPv6 address...")
	conn, err := net.Dial("udp6", "[2001:4860:4860::8888]:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	log.Infof("Resolved primary IPv6 address: %s", localAddr.IP)
	return &localAddr.IP, nil
}

func IsLoopbackAddress(ip string) bool {
	v := net.ParseIP(ip)
	return v != nil && v.IsLoopback()
}

// GetHostnameFromIPAddress returns the hostname associated with the given IP address.
func GetHostnameFromIPAddress(ipAddress string) (string, error) {
	addresses, err := net.LookupAddr(ipAddress)
	if err != nil {
		return "", err
	}
	if len(addresses) == 0 {
		return "", nil
	} else if len(addresses) > 1 {
		log.Infof("Multiple hostnames found for IP address %s: %v - blindly returning the first one\n", ipAddress, addresses)
	}
	return addresses[0], nil
}

// GetIPv4AddressFromHostname returns the IPv4 address associated with the given hostname.
func GetIPv4AddressFromHostname(hostname string) (*net.IP, error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	for _, address := range addresses {
		if address.To4() != nil {
			return &address, nil
		}
	}
	return nil, nil
}

// GetIPv6AddressFromHostname returns the IPv6 address associated with the given hostname.
func GetIPv6AddressFromHostname(hostname string) (*net.IP, error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	for _, address := range addresses {
		if address.To16() == nil {
			return &address, nil
		}
	}
	return nil, nil
}
