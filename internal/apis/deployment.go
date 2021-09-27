package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/deployment"
	"github.com/tangx/k8sailor/internal/biz/pod"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func DeploymentRouterGroup(base *gin.RouterGroup) {
	// 创建 d 路由组
	d := base.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		d.GET("", handlerGetAllDeployments)

		// 针对特定的命名资源操作
		d.GET("/:name", handlerGetPodsByDeployment)
	}
}

// handlerGetAllDeployments 获取所有 deployments
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

func handlerGetPodsByDeployment(c *gin.Context) {
	params := &pod.GetPodsByLabelsInput{}
	err := ginbinder.ShouldBindRequest(c, params)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	pods, err := pod.GetPodsByLabels(*params)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, pods)
}
