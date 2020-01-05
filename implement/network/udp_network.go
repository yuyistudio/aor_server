package network

import (
	"strings"
	"sync"
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/implement/connection"
	"net"
)

type UdpServerConfig struct {
	Addr            string
	ConcurrentCount int
	MaxDiagramSize  int
}

func (c *UdpServerConfig) Validate() {
	if c.Addr == "" {
		panic("[UdpServerConfig] invalid addr")
	}
	if !strings.Contains(c.Addr, ":") {
		panic("[UdpServerConfig] port missing")
	}
	if c.ConcurrentCount <= 0 {
		panic("[UdpServerConfig] invalid concurrent count")
	}
	if c.MaxDiagramSize <= 64 {
		panic("[UdpServerConfig] invalid MaxDiagramSize")
	}
}

type UdpServer struct {
	conf             *UdpServerConfig
	cb               core.ServerCallbackFn
	listener *net.UDPConn
	localAddr *net.UDPAddr
	listenerWg sync.WaitGroup
}

func NewUdpServer(conf *UdpServerConfig) core.Server {
	conf.Validate()
	s := new(UdpServer)
	s.conf = conf
	return s
}

func (s *UdpServer) SetCallback(cb core.ServerCallbackFn) {
	s.cb = cb
}

func (server *UdpServer) Start() error {
	if err := server.init(); err != nil {
		return err
	}
	go server.run()
	return nil
}

func (server *UdpServer) init() error {
	addr, err := net.ResolveUDPAddr("udp", server.conf.Addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	server.localAddr = addr
	server.listener = conn
	return nil
}

func (server *UdpServer) run() {
	for i := 0; i < server.conf.ConcurrentCount; i++ {
		server.listenerWg.Add(1)
		go func() {
			defer server.listenerWg.Done()
			var buf = make([]byte, server.conf.MaxDiagramSize)
			for {
				readBytesCount, addr, err := server.listener.ReadFromUDP(buf)
				if err != nil {
					log.Error("now break, failed to read from udp, error %v", err)
					break
				}
				data := buf[:readBytesCount]
				log.Debug("udp data `%v`", string(data))
				udpConn := connection.NewUdpConnection(data, server.localAddr, addr, server.listener)
				server.cb(udpConn)
			}
		}()
	}
}

func (server *UdpServer) Close() error {
	log.Debug("udp server closing")
	server.listener.Close()        // trigger an error to stop run() loop
	server.listenerWg.Wait()       // waiting for loop done
	return nil
}

