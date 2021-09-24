package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

// RootGroup 向 httpserver 注册根路由
func RootGroup(base *gin.RouterGroup) {
	v0 := base.Group("v0")
	v0.GET("/ping", func(c *gin.Context) {
		c.String(200, "poing")
	})

	// 创建 deployment 路由组
	deployment := v0.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		deployment.GET("/")

		// 针对特定的命名资源操作
		deployment.GET("/:name", func(c *gin.Context) {
			err := errors.New("deployment not found")
			httpresponse.Error(c, http.StatusNotFound, err)
		})
	}
}
