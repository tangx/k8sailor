# 使用 gin 初始化一个 API Server

> tag: https://github.com/tangx/k8sailor/tree/feat/04-httpserver-initial


```bash
cd cmd/k8sailor && go run . httpserver

启动 web 服务器

Usage:
  k8sailor httpserver [flags]

Flags:
  -h, --help   help for httpserver

Global Flags:
      --config string   k8s 配置授权文件 (default "./k8sconfig/config.yml")

2021/09/24 07:56:51 open config/local.yml: no such file or directory
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /k8sailor/v0/ping         --> github.com/tangx/k8sailor/cmd/k8sailor/apis.RootGroup.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8088
```
## 创建 pkg/confgin 初始化配置文件

为了方便服务配置管理， 将使用 [使用 jarvis 初始化配置](#使用 jarvis 初始化配置)

需要对 httpserver [pkg/confgin/gin.go](/pkg/confgin/gin.go) 进行一些初始化配置 

```go
// Server 定义一个 gin httpserver 的关键字段
// `env:""` 可以通过 `github.com/go-jarvis/jarvis` 库渲染成配置文件
type Server struct {
	Host    string `env:""`
	Port    int    `env:""`
	Appname string `env:""`
	engine  *gin.Engine
}
```

其中需要额外强调的是, httpserver **必须** 有自己的 **应用名** 的前缀路由。 **应该** 有自己的版本路由。

```
http://127.0.0.1:8088/appname/v0/ping

```

在服务容器化后， 具有自己 **应用名** 路由的服务对 **各家** `ingress` 规则都是友好的。 如果没有， 上线之后要强行使用 `rewrite` 实现的话， 那就必须依赖 ingress controller 的实现了。 

> 如果 `rewrite` 规则里面 **正则表达式** ， 那就让运维哭去吧

目前所知

+ `nginx ingress controller` 支持 rewrite， 支持正则表达式
+ `traefik` 好像可以使用 `middleware` 实现， 不支持正则表达式
+ `istio` 有自己的规则， 但不支持正则表达式 

```go
// RegisterRoute 注册
func (s *Server) RegisterRoute(registerFunc func(rg *gin.RouterGroup)) {

	// 注册以服务名为根的路由信息，方便在 k8s ingress 中做转发
	base := s.engine.Group(s.Appname)

	// 注册业务子路由
	registerFunc(base)
}
```

## 使用 jarvis 初始化配置

`jarvis` 是一个对 **配置解析** 和 **配置加载** 操作封装的库。 

1. 可以方便的通过 `config struct` 解析出对应的配置参数。
2. 在启动时， 支持通过 **配置文件** 和 **环境变量** 加载配置参数， 对 k8s 容器应用还算友好。
3. 支持 giltab 分支配置特点， 可以在不同分支使用不同的变量值。

具体使用案例， 可以参考 github 的 demo [github.com/go-jarvis/jarvis](https://github.com/go-jarvis/jarvis)


```go
// 定义服务相关信息
var (
	HttpServer = &confgin.Server{}

	app = jarvis.App{
		Name: "k8sailor",
	}
)

// 使用 jarvis 初始化配置文件
func init() {
	config := &struct {
		HttpServer *confgin.Server
	}{
		HttpServer: HttpServer,
	}
	app.Conf(config)
}
```

在运行的时候， 会在 **运行** 根目录生成 `config/default.yml` 文件。 该文件不要直接修改， 每次运行将被覆盖。 

加载顺序 `default.yml -> config.yml -> local.yml / config.branch.yml -> env`。

如果本地开发， 可以把一些关键的敏感配置放在 `local.yml` 中并 `.gitignore` 忽略。


## 为命令行添加 httpserver 子命令

初始化 [cmd/httpserver.go](/cmd/k8sailor/cmd/httpserver.go) 子命令

并设置启动命令

```go
// runHttpserver 启动 http server
func runHttpserver() {
	// 1. 将 apis 注册到 httpserver 中
	global.HttpServer.RegisterRoute(apis.RootGroup)

	// 2. 启动服务
	if err := global.HttpServer.Run(); err != nil {
		logrus.Fatalf("start httpserver failed: %v", err)
	}
}
```

完成之后， 在 `cmd/root.go` 命令中添加子命令

```go
func init() {
	cobrautils.BindFlags(rootCmd, global.Flags)

	// 添加子命令
	rootCmd.AddCommand(cmdHttpserver)
}
```


## 创建并注册路由

## 自定义启动参数

由于是用了 **jarvis** 库， 在程序启动的时候， 会在运行目录生成 `config/defualt.yml` 配置文件。

复制并重命名为 `config.yml` 覆盖默认值。

```yaml
k8sailor__HttpServer_Appname: k8sailor
k8sailor__HttpServer_Host: ""
k8sailor__HttpServer_Port: 8088
```

## 启动

如开头所示， 可以看到, 配置项已成功被应用

+ httpserver 的根路由为 `/k8sailor/v0/xxxxx`
+ httpserver 的监听端口为 `8088`

如果要在容器中运行， 只需要在容器中注入相同变量名的变量

```bash
export k8sailor__HttpServer_Appname=k8sailor
```
