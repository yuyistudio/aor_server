package utility

import (
	"github.com/yuyistudio/aor_server/core/framework"
	"github.com/yuyistudio/aor_server/implement/handler"
	"github.com/yuyistudio/aor_server/implement/network"
	"github.com/yuyistudio/aor_server/implement/parser"
	"github.com/yuyistudio/aor_server/implement/message_package"
)

func NewTcpJsonServer(tcpConf *network.TcpServerConfig) *framework.ServerFramework {
	f := framework.NewServerFramework(true, handler.NewJsonHandler(), network.NewTcpServer(tcpConf), parser.NewJsonParser(), message_package.NewStreamingPackager())
	return f
}
