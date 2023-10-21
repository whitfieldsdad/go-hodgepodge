package hodgepodge

import (
	"net"
	"sort"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type NetworkInterface struct {
	Index                  int             `json:"index,omitempty"`
	Name                   string          `json:"name,omitempty"`
	MTU                    int             `json:"mtu,omitempty"`
	NetworkLocation        NetworkLocation `json:"network_location,omitempty"`
	NetworkInterfaceTraits `json:"network_interface_traits,omitempty"`
}

type NetworkInterfaceTraits struct {
	IsLoopback     bool `json:"is_loopback,omitempty"`
	IsUp           bool `json:"is_up,omitempty"`
	IsMulticast    bool `json:"is_multicast,omitempty"`
	IsBroadcast    bool `json:"is_broadcast,omitempty"`
	IsPointToPoint bool `json:"is_point_to_point,omitempty"`
}

func ListNetworkInterfaces() ([]NetworkInterface, error) {
	log.Info("Listing network interfaces...")
	nics, err := listNetworkInterfaces()
	if err != nil {
		log.Warnf("Failed to list network interfaces: %v\n", err)
		return nil, err
	}
	nicNames := make([]string, len(nics))
	for i, nic := range nics {
		nicNames[i] = nic.Name
	}
	sort.Strings(nicNames)

	log.Infof("Found %d network interfaces: %v\n", len(nicNames), nicNames)
	return nics, nil
}

func listNetworkInterfaces() ([]NetworkInterface, error) {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var nics []NetworkInterface
	for _, networkInterface := range networkInterfaces {
		nics = append(nics, parseNetworkInterfaceInfo(networkInterface))
	}
	return nics, nil
}

func getTraits(nic net.Interface) NetworkInterfaceTraits {
	return NetworkInterfaceTraits{
		IsLoopback:     nic.Flags&net.FlagLoopback != 0,
		IsUp:           nic.Flags&net.FlagUp != 0,
		IsMulticast:    nic.Flags&net.FlagMulticast != 0,
		IsBroadcast:    nic.Flags&net.FlagBroadcast != 0,
		IsPointToPoint: nic.Flags&net.FlagPointToPoint != 0,
	}
}

func GetNetworkInterface(name string) (*NetworkInterface, error) {
	nic, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	o := parseNetworkInterfaceInfo(*nic)
	return &o, nil
}

func parseNetworkInterfaceInfo(nic net.Interface) NetworkInterface {
	networkLocation, _ := GetNetworkLocationFromNetworkInterface(nic.Name)
	return NetworkInterface{
		Index:                  nic.Index,
		Name:                   nic.Name,
		MTU:                    nic.MTU,
		NetworkInterfaceTraits: getTraits(nic),
		NetworkLocation:        *networkLocation,
	}
}

func GetNetworkLocationFromNetworkInterface(name string) (*NetworkLocation, error) {
	nic, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	ipv4Addresses := []string{}
	ipv6Addresses := []string{}
	addrs, err := nic.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		if ip.To4() != nil {
			ipv4Addresses = append(ipv4Addresses, ip.String())
		} else {
			ipv6Addresses = append(ipv6Addresses, ip.String())
		}
	}
	return &NetworkLocation{
		IPv4Addresses: ipv4Addresses,
		IPv6Addresses: ipv6Addresses,
		MACAddress:    nic.HardwareAddr.String(),
	}, nil
}

func GetPrimaryNetworkInterface() (*NetworkInterface, error) {
	primaryIP, err := GetPrimaryIPAddress()
	if err != nil {
		return nil, err
	}
	a := primaryIP
	networkInterfaces, err := ListNetworkInterfaces()
	if err != nil {
		return nil, err
	}
	for _, networkInterface := range networkInterfaces {
		if networkInterface.IsLoopback {
			continue
		}
		if networkInterface.IsUp {
			for _, b := range networkInterface.NetworkLocation.IPv4Addresses {
				if a == b {
					return &networkInterface, nil
				}
			}
			for _, b := range networkInterface.NetworkLocation.IPv6Addresses {
				if a == b {
					return &networkInterface, nil
				}
			}
		}
	}
	return nil, errors.New("no primary network interface found")
}
