package core

import (
	"reflect"
)

type Connection interface {
	Read() (*MessagePackage, error)
	Write(msgID MessageID, bytes []byte) error
	Close() error
}

type ServerCallbackFn func(conn Connection) error
type Server interface {
	Start() error  // DONT block
	Close() error
	SetCallback(cb ServerCallbackFn)
}

type Parser interface {
	Marshal(handler Handler, msgID MessageID, data interface{}) (bytes []byte, err error)
	Unmarshal(handler Handler, msgPkg *MessagePackage) (data interface{}, err error)
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
