package utility

import (
	"github.com/yuyistudio/aor_server/core/framework"
	"github.com/yuyistudio/aor_server/implement/handler"
	"github.com/yuyistudio/aor_server/implement/network"
	"github.com/yuyistudio/aor_server/implement/parser"
	"github.com/yuyistudio/aor_server/implement/message_package"
)

func NewUdpJsonServer(udpConf *network.UdpServerConfig) *framework.ServerFramework {
	f := framework.NewServerFramework(false, handler.NewJsonHandler(), network.NewUdpServer(udpConf), parser.NewJsonParser(), message_package.NewDiagramPackager())
	return f
}
