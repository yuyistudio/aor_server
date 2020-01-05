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

func HeartBeatHandler(msgID core.MessageID, req *rpc.HeartBeatReq) (resp *rpc.HeartBeatResp, err error) {
	resp = new(rpc.HeartBeatResp)
	resp.Code = 0
	return
}

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
	tcp := service.AddServer(utility.NewTcpJsonServer(tcpConf))

	// 定义TCP服务的handler函数。
	tcp.RegisterCallback(rpc.MidRegister, RegisterHandler)
	tcp.RegisterCallback(rpc.MidHeartBeat, HeartBeatHandler)

	// 新增一个TCP服务器。一个服务可以同时包含多个端口的服务器。
	udpConf := new(network.UdpServerConfig)
	udpConf.Addr = "127.0.0.1:8316"
	udpConf.ConcurrentCount = 2
	udpConf.MaxDiagramSize = 4096
	udp := service.AddServer(utility.NewUdpJsonServer(udpConf))

	// 定义UDP服务的handler函数
	udp.RegisterCallback(rpc.MidRegister, RegisterHandler)

	// 开启服务。这里会阻塞住，直到收到退出消息。
	service.Start()
}
