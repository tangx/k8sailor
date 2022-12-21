# 使用 api/biz/dao 分层结构管理数据请求，获取 deployment 数据

> tag: https://github.com/tangx/k8sailor/tree/feat/06-get-all-deployments

```
client ->
apis ->
    biz ->
    dao ->
```

将业务逻辑部分分为经典三层，想法是这样的，可能实现有错误。

1. `apis` 接入层: 只用于管理 http 请求与交互。 
2. `biz` 业务层: 用于处理 api 层来的请求， 封装原始数据
3. `dao` 数据访问层: 与数据库, cluster 等交互。 存取数据。


## 重新调整目录结构

1. 新创建 `/internal` 目录用于存放业务信息
2. 在 `/internal` 目录下新创建 **业务层** [/internal/biz](/internal/biz)  和 **k8s dao 层** [/internal/k8sdao](/internal/k8sdao)
3. 将 `apis` **接入层** 从原来的 `/cmd/k8sailor/apis` 移动到了 [/internal/apis](/internal/apis)


### 使用 jarvis, 删除 cobra flags

1. 统一使用 `jarvis` 进行变量管理， 因此删除了 `global/config.go` 和 `k8s/cluster.go` 原来的 `flag` 相关变量声明使用的代码
2. 但 `cobra cmd` 命令相关代码依旧保留， 感觉以后会用到。
3. 将原来的 `/pkg/k8s` 重命名为 [/pkg/confk8s](/pkg/confk8s)， 命名风格更加统一。

## 获取 deployments 信息

1. api 处理用户请求参数， 请求 biz Operator 方法
2. biz Operator， 请求 k8sdao Operator， 并 **处理/过滤** 原始数据
3. k8sdao 与 cluster 交互， 返回原始数据。

> 有点问题， 三个模块， 三次同名方法。 有点麻烦。

### kubernete.ClientSet 客户端


在 [/cmd/k8sailor/global/config.go](/cmd/k8sailor/global/config.go) 中声明 `KubeClient`

并在 [/internal/k8sdao/clientset.go](/internal/k8sdao/clientset.go) 包中赋值给新变量名保存， 使用短名字方便后续调用。

```go
var clientset = global.KubeClient.Client()
```

### 在 dao 层获取 deployment 数据

在 [/internal/k8sdao/deployments.go](/internal/k8sdao/deployments.go) 中， 封装了一个 **获取指定 namespace** 所有 Pod 的方法。 并返回给下游

```go
func GetAllDeployments(namespace string) ([]appsv1.Deployment, error) {
	ctx := context.TODO()
	opts := metav1.ListOptions{}
	v1Deps, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return v1Deps.Items, nil
}
```

### 在 biz 层提取 dao 层原始数据

在 **业务层** ， 调用 dao 层的 api 获取 k8s cluster 的原始数据， 并根据业务世纪需求提取必要信息形成 **新的业务层的 Deployment 对象**， 并返回给用户。 

之前还真不知道， deployment 中有如此多的 `replicas` 字段（这里还没列举完）

```go
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// Replicas 实际期望的 pod 数量
	Replicas int32 `json:"replicas"`

	// 镜像列表
	Images []string `json:"images"`

	Status DeploymentStatus `json:"status"`
}
type DeploymentStatus struct {
	// 标签匹配的 Pod 数量
	Replicas int32 `json:"replicas"`
	// 可用 pod 数量
	AvailableReplicas int32 `json:"availableReplicas"`
	// 不可用数量
	UnavailableReplicas int32 `json:"unavailableReplicas"`
}
```

同样是在这里， 定义了业务层中 **每个方法** 的请求参数。

```go
type GetAllDeploymentsInput struct {
	Namespace string `query:"namespace"`
}
// GetAllDeployments 获取 namespace 下的所有 deployments
func GetAllDeployments(input GetAllDeploymentsInput) ([]Deployment, error) {
	v1Deps, err := k8sdao.GetAllDeployments(input.Namespace)
// ... 省略
}
```

### apis 接入层处理用户请求， 返回用户需要的数据

在 apis 接入层中， 定义了各个请求的 **方法、路由和处理器（hanlder）**。

同时也将用户请求绑定到 `biz` 中的 **方法请求参数** 上。

```go
func handlerGetAllDeployments(c *gin.Context) {
	params := &deployment.GetAllDeploymentsInput{}
	// 绑定用户请求参数
	err := ginbinder.ShouldBindRequest(c, params)
	if err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}
	// ...省略
}
```

[ginbinder](https://github.com/tangx/ginbinder) 是针对 gin 框架封装的一个请求数据绑定库， 可以方便的将 `http request` 中的请求参数 **一次性全部** 绑定到 **接收者(params)** 中


## 跑起来

使用 `make httpserver` 命令启动 server 服务

```go
[GIN-debug] GET    /k8sailor/v0/ping         --> github.com/tangx/k8sailor/internal/apis.RootGroup.func1 (3 handlers)
[GIN-debug] GET    /k8sailor/v0/deployments/ --> github.com/tangx/k8sailor/internal/apis.handlerGetAllDeployments (3 handlers)
[GIN-debug] GET    /k8sailor/v0/deployments/:name --> github.com/tangx/k8sailor/internal/apis.DeploymentRouterGroup.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8088
```

使用 `vscode REST client` 请求 `/k8sailor/v0/deployments/` 接口

```bash
### GET all deployments
GET http://127.0.0.1:8088/k8sailor/v0/deployments?namespace=default

```

结果与期望一致

```json5
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 24 Sep 2021 18:05:39 GMT
Content-Length: 336
Connection: close

{
  "code": 0,
  "data": [
    {
      "name": "my-nginx-1",
      "namespace": "default",
      "replicas": 1,
      "images": [
        "nginx:alpine"
      ],
      "status": {
        "replicas": 1,
        "availableReplicas": 1,
        "unavailableReplicas": 0
      }
    },
    {
      "name": "my-nginx-2",
      "namespace": "default",
      "replicas": 2,
      "images": [
        "nginx:alpine"
      ],
      "status": {
        "replicas": 2,
        "availableReplicas": 2,
        "unavailableReplicas": 0
      }
    }
  ],
  "error": ""
}
```

## Next

现在数据有了， 接下来要使用 `vue3+typescript` 在前端展示了。 想想都头疼。
