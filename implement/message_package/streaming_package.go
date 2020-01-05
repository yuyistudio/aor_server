package message_package

/*
Package Structure:
4-bits: data length, maximum 4GB data
2-bits: message type, maximum 64K message types
N-bits: data
*/

import (
	"encoding/binary"
	"io"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/core"
	"errors"
)

const lenFieldByteCount = 4
const idFieldByteCount = 2

type StreamingPackager struct {
	maxMsgLen int
	minMsgLen int
}

func NewStreamingPackager() *StreamingPackager {
	c := new(StreamingPackager)
	c.minMsgLen = 1
	c.maxMsgLen = 16 * 1024 * 1024 // limit to 16MB manually
	return c
}

var GlobalStreamingPackager = NewStreamingPackager()

// goroutine safe
func (p *StreamingPackager) ReadPackage(conn io.Reader) (*core.MessagePackage, error) {
	pkg := new(core.MessagePackage)

	// read len field
	{
		lenFieldBytes := make([]byte, lenFieldByteCount)
		if _, err := io.ReadFull(conn, lenFieldBytes); err != nil {
			if err != io.EOF {
				log.Error("failed to read len-field, err %v", err)
			}
			return nil, err
		}
		var msgLen = binary.LittleEndian.Uint32(lenFieldBytes)
		if int(msgLen) > p.maxMsgLen {
			return nil, errors.New("message too long")
		} else if int(msgLen) < p.minMsgLen {
			return nil, errors.New("message too short")
		}
		pkg.Len = msgLen
	}

	// read id field
	{
		idFieldBytes := make([]byte, idFieldByteCount)
		if _, err := io.ReadFull(conn, idFieldBytes); err != nil {
			if err != io.EOF {
				log.Error("failed to read id-field, err %v", err)
			}
			return nil, err
		}
		pkg.MessageID = core.MessageID(binary.LittleEndian.Uint16(idFieldBytes))
	}

	// data
	{
		msgData := make([]byte, int(pkg.Len))
		if _, err := io.ReadFull(conn, msgData); err != nil {
			return nil, err
		}
		pkg.Data = msgData
	}
	return pkg, nil
}

// goroutine safe
func (p *StreamingPackager) Write(conn io.Writer, msgID core.MessageID, bytes []byte) error {
	// check len
	if len(bytes) > p.maxMsgLen {
		return errors.New("message too long")
	} else if len(bytes) < p.minMsgLen {
		return errors.New("message too short")
	}

	// construct package
	msg := make([]byte, lenFieldByteCount+idFieldByteCount+len(bytes))
	binary.LittleEndian.PutUint32(msg, uint32(len(bytes)))
	binary.LittleEndian.PutUint16(msg[lenFieldByteCount:], uint16(msgID))
	copy(msg[lenFieldByteCount+idFieldByteCount:], bytes)

	// write data
	_, err := conn.Write(msg)
	return err
}

// goroutine safe
func (p *StreamingPackager) ReadIdPackage(conn io.Reader) (*core.MessagePackage, error) {
	pkg := new(core.MessagePackage)

	// read len field
	{
		lenFieldBytes := make([]byte, lenFieldByteCount)
		if _, err := io.ReadFull(conn, lenFieldBytes); err != nil {
			if err != io.EOF {
				log.Error("failed to read len-field, err %v", err)
			}
			return nil, err
		}
		var msgLen = binary.LittleEndian.Uint32(lenFieldBytes)
		if int(msgLen) > p.maxMsgLen {
			return nil, errors.New("message too long")
		} else if int(msgLen) < p.minMsgLen {
			return nil, errors.New("message too short")
		}
		pkg.Len = msgLen
	}

	// read id field
	{
		idFieldBytes := make([]byte, idFieldByteCount)
		if _, err := io.ReadFull(conn, idFieldBytes); err != nil {
			if err != io.EOF {
				log.Error("failed to read id-field, err %v", err)
			}
			return nil, err
		}
		pkg.MessageID = core.MessageID(binary.LittleEndian.Uint16(idFieldBytes))
	}

	// data
	{
		msgData := make([]byte, int(pkg.Len))
		if _, err := io.ReadFull(conn, msgData); err != nil {
			return nil, err
		}
		pkg.Data = msgData
	}
	return pkg, nil
}

// goroutine safe
func (p *StreamingPackager) WriteIdPackage(conn io.Writer, msgID core.MessageID, bytes []byte) error {
	// check len
	if len(bytes) > p.maxMsgLen {
		return errors.New("message too long")
	} else if len(bytes) < p.minMsgLen {
		return errors.New("message too short")
	}

	// construct package
	msg := make([]byte, lenFieldByteCount+idFieldByteCount+len(bytes))
	binary.LittleEndian.PutUint32(msg, uint32(len(bytes)))
	binary.LittleEndian.PutUint16(msg[lenFieldByteCount:], uint16(msgID))
	copy(msg[lenFieldByteCount+idFieldByteCount:], bytes)

	// write data
	_, err := conn.Write(msg)
	return err
}
