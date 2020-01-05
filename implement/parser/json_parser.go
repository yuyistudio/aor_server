package parser

import (
	"encoding/json"
	"github.com/yuyistudio/aor_server/core"
)

type JsonParser struct {
}

func NewJsonParser() core.Parser {
	return new(JsonParser)
}

func (p *JsonParser) Unmarshal(msgPkg *core.MessagePackage, data interface{}) error {
	return json.Unmarshal(msgPkg.Data, data)
}

func (p *JsonParser) Marshal(msgID core.MessageID, data interface{}) (bytes []byte, err error) {
	return json.Marshal(data)
}
