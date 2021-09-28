package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/deployment"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func DeploymentRouterGroup(base *gin.RouterGroup) {
	// 创建 dep 路由组
	dep := base.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		dep.GET("", handlerListDeployments)

		// 针对特定的命名资源操作
		dep.GET("/:name", hanlderGetDeploymentByName)

		dep.GET("/:name/pods", handlerGetPodsByDeployment)
	}
}

// handlerListDeployments 获取所有 deployments
func handlerListDeployments(c *gin.Context) {
	params := &deployment.ListDeploymentsInput{}
	err := ginbinder.ShouldBindRequest(c, params)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	deps, err := deployment.ListDeployments(c, *params)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, deps)
}

// hanlderGetDeploymentByName 根据 name 获取 deployment
func hanlderGetDeploymentByName(c *gin.Context) {
	input := deployment.GetDeploymentByNameInput{}
	err := ginbinder.ShouldBindRequest(c, &input)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}
	dep, err := deployment.GetDeploymentByName(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}
	httpresponse.OK(c, dep)
}

// handlerGetPodsByDeployment 根据 deployment 获取 pods
func handlerGetPodsByDeployment(c *gin.Context) {
	// get deployment
	input := deployment.GetPodsByDeploymentInput{}
	err := ginbinder.ShouldBindRequest(c, &input)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	pods, err := deployment.GetPodsByDeployment(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, pods)
}
