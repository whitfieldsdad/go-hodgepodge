package hodgepodge

import (
	"github.com/denisbrodbeck/machineid"
)

type Host struct {
	Id                string             `json:"id"`
	OperatingSystem   OperatingSystem    `json:"operating_system"`
	NetworkLocation   NetworkLocation    `json:"network_location"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces"`
}

func GetHost() (*Host, error) {
	id, err := GetHostId()
	if err != nil {
		return nil, err
	}
	os := GetOperatingSystem()
	networkLocation, _ := GetNetworkLocation()
	networkInterfaces, _ := ListNetworkInterfaces()

	host := &Host{
		Id:                id,
		OperatingSystem:   *os,
		NetworkLocation:   *networkLocation,
		NetworkInterfaces: networkInterfaces,
	}
	return host, nil
}

func GetHostId() (string, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", err
	}
	h := GetSHA256([]byte(id))
	return h, nil
}
