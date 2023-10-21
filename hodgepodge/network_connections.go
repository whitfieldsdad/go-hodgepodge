package hodgepodge

import (
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/net"
)

type NetworkConnection struct {
	LocalHost     string `json:"local_host"`
	LocalPort     int64  `json:"local_port"`
	RemoteHost    string `json:"remote_host"`
	RemotePort    int64  `json:"remote_port"`
	Pid           int64  `json:"pid"`
	Status        string `json:"status"`
	AddressFamily string `json:"address_family"`
	SocketType    string `json:"socket_type"`
}

// TODO: extract address family and socket type
func ListNetworkConnections() ([]*NetworkConnection, error) {
	log.Info("Listing network connections")
	conns, err := net.Connections("all")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list network connections")
	}
	connections := []*NetworkConnection{}
	for _, conn := range conns {
		connections = append(connections, &NetworkConnection{
			LocalHost:  conn.Laddr.IP,
			LocalPort:  int64(conn.Laddr.Port),
			RemoteHost: conn.Raddr.IP,
			RemotePort: int64(conn.Raddr.Port),
			Pid:        int64(conn.Pid),
			Status:     conn.Status,
			//AddressFamily: "",
			//SocketType:    "",
		})
	}
	log.Infof("Found %d network connections", len(connections))
	return connections, nil
}
