package utility

import (
	"github.com/yuyistudio/aor_server/core/framework"
	"github.com/yuyistudio/aor_server/implement/handler"
	"github.com/yuyistudio/aor_server/implement/network"
	"github.com/yuyistudio/aor_server/implement/parser"
)

func NewTcpJsonServer(tcpConf *network.TcpServerConfig) *framework.ServerFramework {
	f := framework.NewServerFramework(handler.NewJsonHandler(), network.NewTcpServer(tcpConf), parser.NewJsonParser())
	return f
}
