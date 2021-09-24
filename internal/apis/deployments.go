package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func DeploymentRouterGroup(base *gin.RouterGroup) {
	// 创建 deployment 路由组
	deployment := base.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		deployment.GET("/", GetAllDeployments)

		// 针对特定的命名资源操作
		deployment.GET("/:name", func(c *gin.Context) {
			err := errors.New("deployment not found")
			httpresponse.Error(c, http.StatusNotFound, err)
		})
	}
}

func GetAllDeployments(c *gin.Context) {

}
