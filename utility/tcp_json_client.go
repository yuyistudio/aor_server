package utility

import (
	"encoding/json"
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/implement/connection"
	"io"
	"net"
)

type ClientFramework struct {
	conn core.Connection
	addr string
	resp chan interface{}
}

func NewClient(addr string) *ClientFramework {
	st := new(ClientFramework)
	st.addr = addr
	return st
}

func (client *ClientFramework) Connect() error {
	if conn, err := net.Dial("tcp", client.addr); err != nil {
		return err
	} else {
		client.conn = connection.NewTcpConnection(conn)
		return nil
	}
}

func (client *ClientFramework) SyncSend(msgID core.MessageID, data interface{}, resp interface{}) error {
	err := client.Send(msgID, data)
	if err != nil {
		return err
	}
	return client.GetResponse(resp)
}

func (client *ClientFramework) Send(msgID core.MessageID, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return client.conn.Write(msgID, bytes)
}

func (f *ClientFramework) GetResponse(resp interface{}) error {
	pkg, err := f.conn.Read()
	if err != nil {
		if err == io.EOF {
			log.Error("connection closed")
		} else {
			log.Error("failed to read message, error %v", err)
		}
		return err
	}
	return json.Unmarshal(pkg.Data, resp)
}
