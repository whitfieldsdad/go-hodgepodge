package main

import (
	"net"
	"os"

	"github.com/charmbracelet/log"
)

func GetHostname() (string, error) {
	return os.Hostname()
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
