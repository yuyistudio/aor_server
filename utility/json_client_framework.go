package utility

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"io"
	"github.com/yuyistudio/aor_server/implement/parser"
)

type ClientFramework struct {
	conn core.Connection
	resp chan interface{}
	parser core.Parser
}

func NewClient(conn core.Connection) *ClientFramework {
	st := new(ClientFramework)
	st.conn = conn
	st.parser = parser.NewJsonParser()
	return st
}

func (client *ClientFramework) Call(msgID core.MessageID, data interface{}, resp interface{}) error {
	err := client.Send(msgID, data)
	if err != nil {
		log.Error("send failed, error %v", err)
		return err
	}
	return client.GetResponse(resp)
}

func (client *ClientFramework) Send(msgID core.MessageID, data interface{}) error {
	bytes, err := client.parser.Marshal(msgID, data)
	if err != nil {
		return err
	}
	return client.conn.WriteMessage(msgID, bytes)
}

func (client *ClientFramework) GetResponse(resp interface{}) error {
	pkg, err := client.conn.ReadMessage()
	if err != nil {
		if err == io.EOF {
			log.Error("failed to read response, error `connection closed`")
		} else {
			log.Error("failed to read response, error %v", err)
		}
		return err
	}
	return client.parser.Unmarshal(pkg, resp)
}
