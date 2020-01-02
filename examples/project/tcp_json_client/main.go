package main

import (
	"fmt"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/examples/project/rpc"
	"github.com/yuyistudio/aor_server/utility"
	"sync"
)

func GoRoutineTest() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		// notice: use one client per goroutine, as client isn't goroutine-safe.
		client := utility.NewClient("127.0.0.1:8316")
		if err := client.Connect(); err != nil {
			log.Error("connection failed, err %v", err)
		}
		req := new(rpc.RegisterReq)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				req.Name = fmt.Sprintf("yy name %d", i)
				req.Password = "yy password"
				var resp rpc.RegisterResp
				err := client.SyncSend(rpc.MidRegister, req, &resp)
				if err != nil {
					log.Error("failed to send rpc, error %v", err)
					return
				}
				log.Info("response %v, %v", resp.Code, resp.Message)
			}
		}()
	}
	log.Info("waiting")
	wg.Wait()
	log.Info("done")
}

func SingleCallTest() {
	// 首先创建一个client，并连接到服务端
	client := utility.NewClient("127.0.0.1:8316")
	if err := client.Connect(); err != nil {
		log.Error("connection failed, err %v", err)
	}

	// 新建一个请求
	req := new(rpc.RegisterReq)
	req.Name = "yy name"
	req.Password = "yy password"

	// 请求服务端，获取返回结果
	var resp rpc.RegisterResp
	err := client.SyncSend(rpc.MidRegister, req, &resp)
	if err != nil {
		log.Error("failed to send rpc, error %v", err)
		return
	}

	// 处理返回结果
	log.Info("response %v, %v", resp.Code, resp.Message)
}

func main() {
	SingleCallTest()
	log.Info("done")
}
