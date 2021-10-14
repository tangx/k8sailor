package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/ingress"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func IngressRouterGroup(base *gin.RouterGroup) {
	ing := base.Group("/ingresses")

	ing.GET("/:name/output", getNetIngressByName)
	ing.POST("/:name", createIngressByName)
}

func getNetIngressByName(c *gin.Context) {
	input := ingress.GetIngressByNameInput{}
	if err := ginbinder.ShouldBindRequest(c, &input); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	v1Ing, err := ingress.GetIngressByName(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, v1Ing)
}

func createIngressByName(c *gin.Context) {
	input := ingress.CreateIngressByNameInput{}
	if err := ginbinder.ShouldBindRequest(c, &input); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	v1Ing, err := ingress.CreateIngressByName(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, v1Ing)
}
