package utility

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/implement/connection"
	"net"
)

func NewUdpConnection(addr string) (core.Connection, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return nil, err
	}
	return connection.NewUdpConnection(nil, nil, remoteAddr, conn), nil
}

func NewTcpConnection(addr string) (core.Connection, error) {
	if conn, err := net.Dial("tcp", addr); err != nil {
		return nil, err
	} else {
		return connection.NewTcpConnection(conn), nil
	}
}
