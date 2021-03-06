package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/service"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func ServiceRouterGroup(base *gin.RouterGroup) {
	svc := base.Group("/services")

	svc.POST("/:name", createServiceByName)
	svc.GET("/:name/output", getCoreServiceByName)
}

func getCoreServiceByName(c *gin.Context) {
	input := service.GetCoreServerByNameInput{}
	if err := ginbinder.ShouldBindRequest(c, &input); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	v1svc, err := service.GetCoreServerByName(c, input)
	if err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, *v1svc)
}

func createServiceByName(c *gin.Context) {
	input := service.CreateServcieByNameInput{}
	if err := ginbinder.ShouldBindRequest(c, &input); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	svc, err := service.CreateServiceByName(c, input)
	if err != nil {
		logrus.Errorln(err)
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, svc)
}
