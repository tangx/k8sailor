package global

import (
	"github.com/go-jarvis/jarvis"
	"github.com/tangx/k8sailor/pkg/confgin"
	"github.com/tangx/k8sailor/pkg/confk8s"
)

// 定义服务相关信息
var (
	HttpServer = &confgin.Server{}
	KubeClient = &confk8s.Client{}

	app = jarvis.App{
		Name: "k8sailor",
	}
)

// 使用 jarvis 初始化配置文件
func init() {
	config := &struct {
		HttpServer *confgin.Server
		KubeClient *confk8s.Client
	}{
		HttpServer: HttpServer,
		KubeClient: KubeClient,
	}
	app.Conf(config)
}
