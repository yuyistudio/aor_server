package main

import (
	"fmt"
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/framework"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/examples/project/rpc"
	"github.com/yuyistudio/aor_server/implement/network"
	"github.com/yuyistudio/aor_server/utility"
)

func RegisterHandler(msgID core.MessageID, req *rpc.RegisterReq) (resp *rpc.RegisterResp, err error) {
	log.Info("received %v", req.Name)
	resp = new(rpc.RegisterResp)
	resp.Code = 0
	resp.Message = fmt.Sprintf("this is response message, req.Name is `%s`", req.Name)
	return
}

func main() {
	// 新建一个服务
	service := framework.NewServiceFramework()

	// 新增一个TCP服务器。一个服务可以同时包含多个端口的服务器。
	tcpConf := new(network.TcpServerConfig)
	tcpConf.Addr = "127.0.0.1:8316"
	tcpConf.MaxConnNum = 100
	tcpConf.PendingWriteNum = 100
	handler := service.AddServer(utility.NewTcpJsonServer(tcpConf))

	// 定义TCP服务的handler函数。
	handler.RegisterCallback(rpc.MidRegister, RegisterHandler)

	// 开启服务。这里会阻塞住，直到收到退出消息。
	service.Start()
}
