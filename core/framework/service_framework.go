package framework

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"os"
	"os/signal"
)

func init() {
	log.SetLogger(new(DefaultLogger))
}

type ServiceFramework struct {
	servers []*ServerFramework
}

func NewServiceFramework() *ServiceFramework {
	f := new(ServiceFramework)
	return f
}

func (frame *ServiceFramework) SetLogger(logger log.Logger) {
	log.SetLogger(logger)
}
func (frame *ServiceFramework) AddServer(server *ServerFramework) core.Handler {
	frame.servers = append(frame.servers, server)
	return server.Handler
}

func (frame *ServiceFramework) Start() {
	log.Info("starting %d server(server)", len(frame.servers))
	for index, server := range frame.servers {
		log.Info("starting server-%d", index)
		server.Start()
	}
	log.Info("service started, servers count %d", len(frame.servers))

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Info("waiting for signal")
	sig := <-c
	log.Info("service closing down (signal: %v)", sig)

	for _, s := range frame.servers {
		s.Server.Close()
	}
	log.Info("service closed")
}
