package parser

import (
	"encoding/json"
	"fmt"
	"github.com/yuyistudio/aor_server/core"
	"reflect"
)

type JsonParser struct {
}

func NewJsonParser() core.Parser {
	return new(JsonParser)
}

func (p *JsonParser) Unmarshal(handler core.Handler, msgPkg *core.MessagePackage) (data interface{}, err error) {
	entry := handler.GetCallback(msgPkg.MessageID)
	if entry == nil {
		return nil, fmt.Errorf("unrecognized message id `%v`", msgPkg.MessageID)
	}
	data = reflect.New(entry.ArgType).Interface()
	err = json.Unmarshal(msgPkg.Data, data)
	return
}

func (p *JsonParser) Marshal(handler core.Handler, msgID core.MessageID, data interface{}) (bytes []byte, err error) {
	return json.Marshal(data)
}
