package apis

import "github.com/gin-gonic/gin"

// RootGroup 向 httpserver 注册根路由
func RootGroup(base *gin.RouterGroup) {
	v0 := base.Group("v0")
	v0.GET("/ping", func(c *gin.Context) {
		c.String(200, "poing")
	})
}
