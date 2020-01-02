Aor Server
====

概述
-----

通用的网络服务框架，基于net包的网络封装，用于网络业务的开发。

如何使用
---------

目录含义：
* core  定义接口，不用关心
* core/framework  定义接口的组合关系，不用关心
* implement  预置的TcpServer/JsonParser的实现，暂时不用关心
* utility  少量的工具代码，示例依赖此代码
* examples/project  一个完整的示例

一个JSON编解码的TCP服务示例：
```golang
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

// 定义接口处理函数
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
```

设计思路
----

框架将网络服务抽象为：
* server 负责网络交互，提供读写字节流的功能
* parser 负责字节流和结构体直接的转换
* handler 负责处理输入的结构体，并返回一个结构体作为返回结果


Licensing
---------

Aor Server is licensed under the Apache License.
