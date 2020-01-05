package framework

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"io"
	"reflect"
	"fmt"
)

type ServerFramework struct {
	Handler core.Handler
	Server  core.Server
	parser  core.Parser
	packager core.MessagePackager
	loopRead bool
}

func NewServerFramework(loopRead bool, registry core.Handler, server core.Server, parser core.Parser, packager core.MessagePackager) *ServerFramework {
	sf := new(ServerFramework)
	sf.Handler = registry
	sf.Server = server
	if loopRead {
		sf.Server.SetCallback(sf.callback)
	} else {
		sf.Server.SetCallback(sf.readOnce)
	}
	sf.parser = parser
	sf.packager = packager
	sf.loopRead = loopRead
	return sf
}

func (f *ServerFramework) callback(conn core.Connection) error {
	for {
		err := f.readOnce(conn)
		if err != nil {
			return err
		}
	}
}

func (f *ServerFramework) readOnce(conn core.Connection) error {
	// network -> bytes -> package
	pkg, err := conn.ReadMessage()
	if err != nil {
		if err == io.EOF {
			log.Debug("client connection closed")
		} else {
			log.Error("close connection now, failed to read message, error %v", err)
			return err
		}
		return err
	}

	// raw_package -> req
	entry := f.Handler.GetCallback(pkg.MessageID)
	if entry == nil {
		return fmt.Errorf("unrecognized message id `%v`", pkg.MessageID)
	}
	data := reflect.New(entry.ArgType).Interface()
	err = f.parser.Unmarshal(pkg, data)
	if err != nil {
		log.Error("failed to parse data, error `%v`", err)
		return err
	}

	// req -> resp
	resp, err := f.Handler.Process(conn, pkg.MessageID, data)
	if err != nil {
		log.Error("failed to process data, error `%v`", err)
		return err
	}

	// resp -> bytes -> network
	bytes, err := f.parser.Marshal(pkg.MessageID, resp)
	if err != nil {
		log.Error("failed to marshal data, error `%v`", err)
	}
	err = conn.WriteMessage(pkg.MessageID, bytes)
	if err != nil {
		log.Error("failed to write response, error `%v`", err)
	} else {
		log.Debug("write response done")
	}
	return nil
}

func (f *ServerFramework) Start() error {
	return f.Server.Start()
}
