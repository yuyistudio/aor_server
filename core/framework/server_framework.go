package framework

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"io"
)

type ServerFramework struct {
	Handler core.Handler
	Server  core.Server
	parser  core.Parser
}

func NewServerFramework(registry core.Handler, server core.Server, parser core.Parser) *ServerFramework {
	st := new(ServerFramework)
	st.Handler = registry
	st.Server = server
	st.Server.SetCallback(st.callback)
	st.parser = parser
	return st
}

func (f *ServerFramework) callback(conn core.Connection) error {
	for {
		pkg, err := conn.Read()
		if err != nil {
			if err == io.EOF {
				log.Debug("connection closed")
				return nil
			} else {
				log.Error("failed to read message, error %v", err)
				return err
			}
		}
		data, err := f.parser.Unmarshal(f.Handler, pkg)
		if err != nil {
			log.Error("failed to parse data, error `%v`", err)
			return err
		}
		resp, err := f.Handler.Process(conn, pkg.MessageID, data)
		if err != nil {
			log.Error("failed to process data, error `%v`", err)
			return err
		}
		bytes, err := f.parser.Marshal(f.Handler, pkg.MessageID, resp)
		if err != nil {
			log.Error("failed to marshal data, error `%v`", err)
		}
		err = conn.Write(pkg.MessageID, bytes)
		if err != nil {
			log.Error("failed to write response, error `%v`", err)
		}
	}
}

func (f *ServerFramework) Start() error {
	return f.Server.Start()
}
