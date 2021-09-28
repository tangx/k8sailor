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
	// 创建 dep 路由组
	dep := base.Group("/deployments")
	{
		// 针对 所有 deployment 操作
		dep.GET("", handlerGetAllDeployments)

		// 针对特定的命名资源操作
		dep.GET("/:name", handlerGetPodsByDeployment)
	}
}

// handlerGetAllDeployments 获取所有 deployments
func handlerGetAllDeployments(c *gin.Context) {
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

func handlerGetPodsByDeployment(c *gin.Context) {
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

	pInput := pod.GetPodsByLabelsInput{
		Namespace: dep.Namespace,
		Labels:    dep.LabelSelector.MatchLabels,
	}
	pods, err := pod.GetPodsByLabels(c, pInput)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
	}

	httpresponse.OK(c, pods)
}
