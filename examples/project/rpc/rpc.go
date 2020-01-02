package rpc

// 这里定义 消息ID、消息Request、消息Response

import "github.com/yuyistudio/aor_server/core"

const (
	MidRegister core.MessageID = 1
)

type RegisterReq struct {
	Name     string
	Password string
}

type RegisterResp struct {
	Code    int
	Message string
}
