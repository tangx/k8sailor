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

	// 针对 appname 下的路由，允许跨域
	base.Use(MiddleCors())

	// 注册业务子路由
	registerFunc(base)
}

func AppendGroup(base *gin.RouterGroup, register func(base *gin.RouterGroup)) {
	register(base)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
