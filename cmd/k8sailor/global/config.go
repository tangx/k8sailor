package global

import (
	"github.com/go-jarvis/jarvis"
	"github.com/tangx/k8sailor/pkg/confgin"
)

// 定义命令行相关参数
type CmdFlags struct {
	Config string `flag:"config" usage:"k8s 配置授权文件" persistent:"true"`
}

var Flags = &CmdFlags{
	Config: "./k8sconfig/config.yml",
}

// 定义服务相关信息
var (
	HttpServer = &confgin.Server{}

	app = jarvis.App{
		Name: "k8sailor",
	}
)

// 使用 jarvis 初始化配置文件
func init() {
	config := &struct {
		HttpServer *confgin.Server
	}{
		HttpServer: HttpServer,
	}
	app.Conf(config)
}
