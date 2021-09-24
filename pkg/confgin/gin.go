package confgin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Server 定义一个 gin httpserver 的关键字段
// `env:""` 可以通过 `github.com/go-jarvis/jarvis` 库渲染成配置文件
type Server struct {
	Host    string `env:""`
	Port    int    `env:""`
	Appname string `env:""`
	engine  *gin.Engine
}

// SetDefaults 设置默认值
func (s *Server) SetDefaults() {
	if s.Port == 0 {
		s.Port = 80
	}
	if s.Appname == "" {
		s.Appname = "app"
	}
}

// Init 初始化关键信息
func (s *Server) Init() {

	if s.engine == nil {
		s.SetDefaults()
		s.engine = gin.Default()
	}

}

// Run 启动业务
func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	return s.engine.Run(addr)
}

// RegisterRoute 注册
func (s *Server) RegisterRoute(registerFunc func(rg *gin.RouterGroup)) {

	// 注册以服务名为根的路由信息，方便在 k8s ingress 中做转发
	base := s.engine.Group(s.Appname)

	// 注册业务子路由
	registerFunc(base)
}
