package message_package

import (
	"encoding/binary"
	"io"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/core"
)

type DiagramPackager struct {
}

func NewDiagramPackager() *DiagramPackager {
	c := new(DiagramPackager)
	return c
}

var GlobalDiagramPackager = NewDiagramPackager()

// goroutine safe
func (p *DiagramPackager) ReadPackage(conn io.Reader) (*core.MessagePackage, error) {
	pkg := new(core.MessagePackage)
	log.Debug("reading from udp reader")
	diagram := make([]byte, 4096)
	diagramSize, err := conn.Read(diagram)
	diagram = diagram[:diagramSize]
	log.Debug("result `%v` error `%v`", string(diagram), err)
	if err != nil {
		if len(diagram) == 0 {
			return nil, err
		} else {
			// ok
		}
	}

	// read id field
	idFieldBytes := diagram[:idFieldByteCount]
	pkg.MessageID = core.MessageID(binary.LittleEndian.Uint16(idFieldBytes))
	pkg.Len = uint32(len(diagram) - idFieldByteCount)
	pkg.Data = diagram[idFieldByteCount:]
	return pkg, nil
}

// goroutine safe
func (p *DiagramPackager) Write(conn io.Writer, msgID core.MessageID, bytes []byte) error {
	log.Debug("writing to udp")
	// construct package
	msg := make([]byte, idFieldByteCount+len(bytes))
	binary.LittleEndian.PutUint16(msg[:idFieldByteCount], uint16(msgID))
	copy(msg[idFieldByteCount:], bytes)

	// write data
	log.Debug("writing data `%s`", string(msg))
	_, err := conn.Write(msg)
	return err
}
