package connection

import (
	"github.com/yuyistudio/aor_server/core"
	"net"
	"github.com/yuyistudio/aor_server/implement/message_package"
)

type TcpConnection struct {
	conn     net.Conn
}

func NewTcpConnection(conn net.Conn) *TcpConnection {
	tcpConn := new(TcpConnection)
	tcpConn.conn = conn
	return tcpConn
}

func (p *TcpConnection) ReadMessage() (*core.MessagePackage, error) {
	return message_package.GlobalStreamingPackager.ReadPackage(p.conn)
}

func (p *TcpConnection) WriteMessage(msgID core.MessageID, bytes []byte) error {
	return message_package.GlobalStreamingPackager.Write(p.conn, msgID, bytes)
}

func (p *TcpConnection) Close() error {
	return p.conn.Close()
}
