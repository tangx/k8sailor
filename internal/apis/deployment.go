package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/deployment"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func DeploymentRouterGroup(base *gin.RouterGroup) {
	// 创建 deployment 路由组
	deployment := base.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		deployment.GET("/", handlerGetAllDeployments)

		// 针对特定的命名资源操作
		deployment.GET("/:name", func(c *gin.Context) {
			err := errors.New("deployment not found")
			httpresponse.Error(c, http.StatusNotFound, err)
		})
	}
}

func handlerGetAllDeployments(c *gin.Context) {
	params := &deployment.GetAllDeploymentsInput{}
	err := ginbinder.ShouldBindRequest(c, params)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	deps, err := deployment.GetAllDeployments(*params)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, deps)
}
