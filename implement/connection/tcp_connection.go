package connection

import (
	"github.com/yuyistudio/aor_server/core"
	"net"
)

type TcpConnection struct {
	conn           net.Conn
	packageHandler *core.MessagePackageHandler
}

func NewTcpConnection(conn net.Conn) *TcpConnection {
	c := new(TcpConnection)
	c.conn = conn
	c.packageHandler = core.NewMessagePackageHandler()
	return c
}

func (p *TcpConnection) Read() (*core.MessagePackage, error) {
	return p.packageHandler.ReadPackage(p.conn)
}

func (p *TcpConnection) Write(msgID core.MessageID, bytes []byte) error {
	return p.packageHandler.Write(p.conn, msgID, bytes)
}

func (p *TcpConnection) Close() error {
	return p.conn.Close()
}
