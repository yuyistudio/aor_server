package core

type MessageID uint16

type MessagePackage struct {
	Len       uint32
	MessageID MessageID
	Data      []byte
}
