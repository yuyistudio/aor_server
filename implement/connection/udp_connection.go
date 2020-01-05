package connection

import (
	"github.com/yuyistudio/aor_server/core"
	"net"
	"strings"
	"github.com/yuyistudio/aor_server/implement/message_package"
	"github.com/pkg/errors"
	"github.com/yuyistudio/aor_server/core/log"
)

type UdpConnection struct {
	diagram           []byte
	localAddr *net.UDPAddr
	remoteAddr *net.UDPAddr
	conn *net.UDPConn
}

func NewUdpConnection(data []byte, localAddr, remoteAddr *net.UDPAddr, conn *net.UDPConn) *UdpConnection {
	c := new(UdpConnection)
	c.diagram = data
	c.localAddr = localAddr
	c.remoteAddr = remoteAddr
	c.conn = conn
	return c
}

func (p *UdpConnection) ReadMessage() (*core.MessagePackage, error) {
	if p.diagram != nil {
		log.Debug("reading diagram `%s`", p.diagram)
		reader := strings.NewReader(string(p.diagram))
		p.diagram = nil
		return message_package.GlobalDiagramPackager.ReadPackage(reader)
	} else if p.conn != nil {
		log.Debug("reading connection")
		return message_package.GlobalDiagramPackager.ReadPackage(p.conn)
	}
	return nil, errors.New("invalid diagram or connection")
}

func (p *UdpConnection) WriteMessage(msgID core.MessageID, bytes []byte) error {
	if p.conn == nil {
		conn, err := net.DialUDP("udp", p.localAddr, p.remoteAddr)
		if err != nil {
			return err
		}
		p.conn = conn
	}
	if p.conn.RemoteAddr() == nil {
		return message_package.GlobalDiagramPackager.Write(p, msgID, bytes)
	} else {
		return message_package.GlobalDiagramPackager.Write(p.conn, msgID, bytes)
	}
}

func (p *UdpConnection) Write(data []byte) (n int, err error) {
	return p.conn.WriteToUDP(data, p.remoteAddr)
}
/*
func (p *UdpConnection) Read(data []byte) (n int, err error) {
	return p.conn.ReadFromUDP(data, p.remoteAddr)
}
*/

func (p *UdpConnection) Close() error {
	// do nothing
	return nil
}
