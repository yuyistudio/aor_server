package network

import (
	"github.com/yuyistudio/aor_server/core"
	"github.com/yuyistudio/aor_server/core/log"
	"github.com/yuyistudio/aor_server/implement/connection"
	"net"
	"strings"
	"sync"
	"time"
)

type TcpServerConfig struct {
	Addr            string
	MaxConnNum      int
	PendingWriteNum int
}

func (c *TcpServerConfig) Validate() {
	if c.Addr == "" {
		panic("[TcpServerConfig] invalid addr")
	}
	if !strings.Contains(c.Addr, ":") {
		panic("[TcpServerConfig] port missing")
	}
	if c.MaxConnNum < 1 {
		panic("[TcpServerConfig] invalid MaxConnNum")
	}
	if c.PendingWriteNum < 1 {
		panic("[TcpServerConfig] invalid PendingWriteNum")
	}
}

type TCPServer struct {
	conf             *TcpServerConfig
	cb               core.ServerCallbackFn
	listener         net.Listener
	connections      map[net.Conn]core.Connection
	connectionsMutex sync.Mutex
	listenerWg       sync.WaitGroup
	wgForConnections sync.WaitGroup

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
}

func NewTcpServer(conf *TcpServerConfig) core.Server {
	conf.Validate()
	s := new(TCPServer)
	s.conf = conf
	s.connections = make(map[net.Conn]core.Connection)
	return s
}

func (s *TCPServer) SetCallback(cb core.ServerCallbackFn) {
	s.cb = cb
}

func (server *TCPServer) Start() error {
	if err := server.init(); err != nil {
		return err
	}
	go server.run()
	return nil
}

func (server *TCPServer) init() error {
	listener, err := net.Listen("tcp", server.conf.Addr)
	if err != nil {
		return err
		return err
	}
	server.listener = listener
	return nil
}

func (server *TCPServer) run() {
	server.listenerWg.Add(1)
	defer server.listenerWg.Done()

	var retryDelay time.Duration
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if retryDelay == 0 {
					retryDelay = 5 * time.Millisecond
				} else {
					retryDelay *= 2
				}
				if max := 1 * time.Second; retryDelay > max {
					retryDelay = max
				}
				log.Error("error %v, retry after %v", err.Error(), retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return
		}
		retryDelay = 0
		log.Debug("connection received")

		if len(server.connections) >= server.conf.MaxConnNum {
			log.Error("max connections count reached, conf.count %d", server.conf.MaxConnNum)
			conn.Close()
			continue
		}
		server.connectionsMutex.Lock()
		tcpConn := connection.NewTcpConnection(conn)
		server.connections[conn] = tcpConn
		server.connectionsMutex.Unlock()

		server.wgForConnections.Add(1)
		go func() {
			defer server.wgForConnections.Done()
			server.cb(tcpConn)

			// cleanup
			tcpConn.Close()
			server.connectionsMutex.Lock()
			delete(server.connections, conn)
			server.connectionsMutex.Unlock()
		}()
	}
}

func (server *TCPServer) Close() error {
	log.Debug("tcp server closing")
	server.listener.Close()        // trigger an error to stop run() loop
	server.listenerWg.Wait()       // waiting for loop done
	server.wgForConnections.Wait() // waiting for all agents to return
	return nil
}
