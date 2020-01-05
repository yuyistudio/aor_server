package main

import (
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/examples/project/rpc"
	"github.com/yuyistudio/aor_server/utility"
	"github.com/yuyistudio/aor_server/core"
	"time"
	"sync"
)

func SingleCallTest(useTcp bool, concurrency int, reqCount int, duration time.Duration) {
	var wg sync.WaitGroup
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 首先创建一个client，并连接到服务端
			var conn core.Connection
			var err error
			if useTcp {
				conn, err = utility.NewTcpConnection("127.0.0.1:8316")
				if err != nil {
					log.Error("connection to tcp server failed, err %v", err)
					return
				}
			} else {
				conn, err = utility.NewUdpConnection("127.0.0.1:8316")
				if err != nil {
					log.Error("connection to udp server failed, err %v", err)
				}
			}
			client := utility.NewClient(conn)

			// 心跳
			var doneChan = make(chan int, 1)
			if useTcp {
				go func() {
					for {
						time.Sleep(200*time.Millisecond)
						req := new(rpc.HeartBeatReq)
						var resp rpc.HeartBeatResp
						err := client.Call(rpc.MidHeartBeat, req, &resp)
						if err == nil {
							log.Info("heart beat ok")
						} else {
							log.Error("heart beat failed, error %v", err)
						}

						select {
						case <-doneChan:
							break
						default:
							continue
						}
					}
					doneChan <- 1
				}()
			}

			// 业务逻辑
			for i := 0; i < reqCount; i ++ {
				// 新建一个请求
				req := new(rpc.RegisterReq)
				req.Name = "yy name"
				req.Password = "yy password"

				// 请求服务端，获取返回结果
				var resp rpc.RegisterResp
				err = client.Call(rpc.MidRegister, req, &resp)
				if err != nil {
					log.Error("failed to send rpc, error %v", err)
					return
				}
				// 处理返回结果
				log.Info("response %v, %v", resp.Code, resp.Message)
				if duration > 0 {
					time.Sleep(duration)
				}
			}

			if useTcp {
				doneChan <- 1
				select {
				case <-doneChan:
				}
			}
		}()
	}
	wg.Wait()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func(){
		defer wg.Done()
		SingleCallTest(false, 1, 5, 1*time.Second)
	}()
	go func() {
		defer wg.Done()
		go SingleCallTest(true, 1, 5, 1*time.Second)
	}()
	wg.Wait()
	log.Info("done")
}
