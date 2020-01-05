package core

import (
	"reflect"
	"io"
)

type Connection interface {
	ReadMessage() (*MessagePackage, error)
	WriteMessage(msgID MessageID, bytes []byte) error
	Close() error
}

type ServerCallbackFn func(conn Connection) error
type Server interface {
	Start() error  // DONT block
	Close() error
	SetCallback(cb ServerCallbackFn)
}

type Parser interface {
	// 通过handler，获取callback，进而获取
	Marshal(msgID MessageID, data interface{}) (bytes []byte, err error)
	Unmarshal(msgPkg *MessagePackage, data interface{}) error
}

type Callback struct {
	ID      MessageID
	ArgType reflect.Type
	Fn      reflect.Value
}

type Handler interface {
	RegisterCallback(msgID MessageID, fn interface{})
	GetCallback(msgID MessageID) *Callback
	Process(conn Connection, msgID MessageID, req interface{}) (response interface{}, err error)
}

type MessagePackager interface {
	ReadPackage(conn io.Reader) (*MessagePackage, error)
	Write(conn io.Writer, msgID MessageID, bytes []byte) error
}
