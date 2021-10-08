package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/pod"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func PodRouterGroup(base *gin.RouterGroup) {
	pod := base.Group("/pods")

	pod.GET("/:name/event", getPodEventByName)
	pod.GET("/:name/output", getCorePodByName)
}

func getPodEventByName(c *gin.Context) {

	input := pod.GetPodEventByNameInput{}

	err := ginbinder.ShouldBindRequest(c, &input)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	event := pod.GetPodEventByName(c, input)

	httpresponse.OK(c, event)
}

func getCorePodByName(c *gin.Context) {
	input := pod.GetPodByNameInput{}
	err := ginbinder.ShouldBindRequest(c, &input)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	v1pod, err := pod.GetCorePodByName(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	if input.OutputFormat == "yaml" {
		c.YAML(http.StatusOK, v1pod)
		return
	}

	c.JSON(http.StatusOK, v1pod)
}
