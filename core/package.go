package core

/*
Package Structure:
4-bits: data length, maximum 4GB data
2-bits: message type, maximum 64K message types
N-bits: data
*/

import (
	"encoding/binary"
	"errors"
	"github.com/yuyistudio/aor_server/core/log"
	"io"
)

type MessageID uint16

type MessagePackage struct {
	Len       uint32
	MessageID MessageID
	Data      []byte
}

const lenFieldByteCount = 4
const idFieldByteCount = 2

type MessagePackageHandler struct {
	maxMsgLen int
	minMsgLen int
}

func NewMessagePackageHandler() *MessagePackageHandler {
	c := new(MessagePackageHandler)
	c.minMsgLen = 1
	c.maxMsgLen = 16 * 1024 * 1024 // limit to 16MB manually
	return c
}

// goroutine safe
func (p *MessagePackageHandler) ReadPackage(conn io.Reader) (*MessagePackage, error) {
	pkg := new(MessagePackage)

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
		pkg.MessageID = MessageID(binary.LittleEndian.Uint16(idFieldBytes))
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
func (p *MessagePackageHandler) Write(conn io.Writer, msgID MessageID, bytes []byte) error {
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
