package apis

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
	"github.com/tangx/k8sailor/internal/biz/deployment"
	"github.com/tangx/k8sailor/pkg/confgin/httpresponse"
)

func DeploymentRouterGroup(base *gin.RouterGroup) {
	// 创建 deployment 路由组
	deployment := base.Group("/deployments")
	{

		/*
			如果没有这个 handler， 那么 gin 默认生成 301 规则
			[GIN-debug] redirecting request 301: /k8sailor/v0/deployments --> /k8sailor/v0/deployments?namespace=default

			在没有跨域的情况下， 例如 `curl -L` 没有什么问题， 或者浏览器直接访问
			但是对于 axios 的时候，将会出现 cors 错误
			axios redirect No 'Access-Control-Allow-Origin' header is present

			但是如果手动实现这个 301 的 handler， 一切正常。

			从浏览器 network 请求记录上来看，
				1. 手动实现的情况下， 总共触发了 **2次** 请求。
				2. gin 自动实现情况， 只有 **1次** 请求。

			gin issue 有相同问题:
				https://github.com/gin-gonic/gin/issues/1985
				https://github.com/gin-gonic/gin/issues/2413#issuecomment-645768561
			目前， gin 先执行 middleware, 后找路由并执行 handler。
		*/
		deployment.GET("/", func(c *gin.Context) {
			// fmt.Println(c.Request.URL)  // /k8sailor/v0/deployments/?namespace=default
			// fmt.Println(c.Request.URL.Path) // /k8sailor/v0/deployments/
			// fmt.Println(c.Request.URL.RawQuery)  // namespace=default
			_url := strings.TrimRight(c.Request.URL.Path, "/") + "?" + c.Request.URL.RawQuery
			c.Redirect(301, _url)
		})

		// 针对 所有 deployment 操作
		deployment.GET("", handlerGetAllDeployments)

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
