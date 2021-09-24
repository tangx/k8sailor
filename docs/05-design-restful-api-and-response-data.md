# RESTful-API 与 http-response-data

强烈建议使用 **RESTful** 风格来设计 API 文档。


## RESTful api

```yaml
# kubectl create deployment nginx-tools --image nginx:alpine --output=yaml --dry-run=client
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: nginx-tools
  name: nginx-tools
# ... 省略

# kubectl create namespace  hello --dry-run=client -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: null
  name: hello
# ... 省略
```

可以看到， k8s api 中都有一个对应的 `kind` 描述资源类型， 这个正好符合 RESTful 中资源定位的需求。

大概就是这样。

```bash
# 所有资源操作
GET     /appname/v0/:resources

## 特定志愿操作
GET     /appname/v0/:resources/:name?params
POST    /appname/v0/:resources/:name?params
DELETE  /appname/v0/:resources/:name
```

```bash
# 获取所有 deployemnt 信息, 默认会设计一些限定条件， 比如说 namespace=default
GET /k8sailor/v0/deployments

# 针对特定名称资源的 deployment 操作
GET /k8sailor/v0/deployments/my-nginx-01?namespace=kube-system
DELETE /k8sailor/v0/deployments/my-nginx-01?namespace=default
```

回到代码中 [/cmd/k8sailor/apis/route.go](/cmd/k8sailor/apis/route.go)

```go
// RootGroup 向 httpserver 注册根路由
func RootGroup(base *gin.RouterGroup) {
// ... 省略

	// 创建 deployment 路由组
	deployment := v0.Group("/deployments")
	{
		// 针对 所有 deployment 操作， 这里还没有绑定 handler
		deployment.GET("/")

		// 针对特定的命名资源操作
		// 直接返回 找不到
		deployment.GET("/:name", func(c *gin.Context) {
			err := errors.New("deployment not found")
			httpresponse.Error(c, http.StatusNotFound, err)
		})
	}
}
```

## http response

对于应答消息， 不建议将 **成功** 和 **失败** 内容分成两个不同的 **结构体** 发送给客户端， 否则客户端在使用的时候还需要在判断应答的结构体属于哪种。 如果服务端一旦修改了应答结构体，客户端可能就崩掉了。

```json5
// 成功
{
    "data":"success data"
}

// 失败
{
    "error":"error message"
}
```

因此需要对 http 相应进行一些简单的封装。

1. 把应答消息封装成一个标准结构， 具体消息信息用某个字段占有。
    + data 表示成功消息
    + error 表示失败消息

2. `http status code` 本身就对 **行为和资源** 的有了一个明确的描述， 并且是通用的。 因此最好能将 `response code` 和 `http status code` 之间建立一个映射关系， 这样通过 code 也快速的判断 response 状态和内容。
    + 这里只是简单的将 http status code 用作 response code 。
    + 如果 http code 是 200， 则 response code 强制设置成 0。 一般情况下，非 0 表示异常。

[/pkg/confgin/httpresponse/response.go](/pkg/confgin/httpresponse/response.go)

```go
func Common(c *gin.Context, code int, data interface{}, err error) {
	_err := ""
	if err != nil {
		_err = err.Error()
	}

	// 强制设置
	if code == 200 {
		code = 0
	}

	resp := Response{
		Code:  code,
		Data:  data,
		Error: _err,
	}

	c.JSON(code, resp)
}
```

## 文献

+ **gitlab** RESTful API: https://docs.gitlab.com/ee/api/api_resources.html
+ **github** RESTful API: https://docs.github.com/en/rest/overview/resources-in-the-rest-api
+ **kalaserach** resp api 最佳实践: https://kalasearch.cn/blog/rest-api-best-practices/


## 运行起来

### 启动服务
```bash
cd cmd/k8sailor && go run . httpserver

[GIN-debug] GET    /k8sailor/v0/ping         --> github.com/tangx/k8sailor/cmd/k8sailor/apis.RootGroup.func1 (3 handlers)
[GIN-debug] GET    /k8sailor/v0/deployments/ --> github.com/gin-gonic/gin.CustomRecoveryWithWriter.func1 (2 handlers)
[GIN-debug] GET    /k8sailor/v0/deployments/:name --> github.com/tangx/k8sailor/cmd/k8sailor/apis.RootGroup.func2 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8088

```

### 请求接口

这里推荐一个 vscode 下比较好用的 http client `REST client`, 类似 postman

> https://marketplace.visualstudio.com/items?itemName=humao.rest-client

```bash
### GET deployment by name
GET http://127.0.0.1:8088/k8sailor/v0/deployments/my-nginx-01
```

请求结果， 资源找不到， http status code, data code 等都符合预期。

```json
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Date: Fri, 24 Sep 2021 03:18:34 GMT
Content-Length: 55
Connection: close

{
  "code": 404,
  "data": null,
  "error": "deployment not found"
}
```

