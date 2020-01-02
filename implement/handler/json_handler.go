package handler

import (
	"fmt"
	"github.com/yuyistudio/aor_server/core"
	"reflect"
)

type JsonHandler struct {
	// func fn(msgID MessageID, request *UserStruct) (response *UserStruct, error)
	entries map[core.MessageID]*core.Callback // msg_id -> fn
}

func NewJsonHandler() *JsonHandler {
	h := new(JsonHandler)
	h.entries = make(map[core.MessageID]*core.Callback)
	return h
}

func check(fn interface{}) error {
	// check request
	t := reflect.TypeOf(fn)
	if t.NumIn() != 2 {
		panic(fmt.Sprintf("invalid argument count %d", t.NumIn()))
	}
	arg := t.In(0)
	if arg.Kind() != reflect.Uint16 {
		panic(fmt.Sprintf("invalid first argment type, type `core.MessageID` expected"))
	}
	arg = t.In(1)
	if arg.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("invalid seconds argument type, [pointer] of struct expected"))
	} else if arg.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("invalid seconds argument type, pointer of [struct] expected"))
	}

	// check response
	if t.NumOut() != 2 {
		panic(fmt.Sprintf("invalid return values count %d", t.NumIn()))
	}
	arg = t.Out(0)
	if arg.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("invalid seconds return type, [pointer] of struct expected"))
	} else if arg.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("invalid seconds return type, pointer of [struct] expected"))
	}
	arg = t.Out(1)
	if arg.Kind() != reflect.Interface {
		panic(fmt.Sprintf("invalid second return type, error expected"))
	}
	return nil
}

func (r *JsonHandler) RegisterCallback(msgID core.MessageID, fn interface{}) {
	check(fn)
	if _, ok := r.entries[msgID]; ok {
		panic(fmt.Sprintf("registered more than once, id %d", msgID))
	}
	check(fn)
	entry := new(core.Callback)
	entry.Fn = reflect.ValueOf(fn)
	entry.ArgType = reflect.TypeOf(fn).In(1).Elem()
	entry.ID = msgID
	r.entries[msgID] = entry
}

func (r *JsonHandler) GetCallback(msgID core.MessageID) *core.Callback {
	return r.entries[msgID]
}

func (r *JsonHandler) Process(conn core.Connection, msgID core.MessageID, req interface{}) (interface{}, error) {
	entry := r.GetCallback(msgID)
	results := entry.Fn.Call([]reflect.Value{reflect.ValueOf(msgID), reflect.ValueOf(req)})
	err := results[1].Interface()
	if err == nil {
		return results[0].Interface(), nil
	} else {
		return nil, err.(error)
	}
}
