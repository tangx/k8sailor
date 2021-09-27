# 获取并展示 Deployments 信息


![display-deployments.png](/docs/assets/img/08/display-deployments.png)

## 使用 Axios 请求 Deployments 数据

安装 axios 客户端

```bash
# 安装 axios
yarn add axios
```

创建 [/webapp/src/apis](/webapp/src/apis) 目录， 用于存放所有针对 k8sailor 后端的数据请求


## server 端允许 cors 跨域

跨域在 gin 中的实现其实就是 `gin.HandlerFunc`， 可以理解成一种中间件。

以下就是跨域规则，规则比较暴力， 几乎允许所了所有。

```go
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
```

接下来是允许跨域



```go
// RegisterRoute 注册
func (s *Server) RegisterRoute(registerFunc func(rg *gin.RouterGroup)) {

	// 注册以服务名为根的路由信息，方便在 k8s ingress 中做转发
	base := s.engine.Group(s.Appname)

	// 针对 appname 下的路由，允许跨域
	base.Use(cors())

	// 注册业务子路由
	registerFunc(base)
}
```

这里并没有在 **根路由** 下允许， 而是在 `/:appname` 下允许。 

也就是说如下

```bash
# 允许跨域
/appname/deployments   
/appname/pods/:podname  

# 不允许
/ping
```

## vue3 展示数据

### 使用 onMounted 加载数据

### 使用 v-for 显示数据

### 使用 v-if 进行条件渲染

### 使用 v-click 输入查询条件


## 问题遗留

### 301 重定向遇到跨域问题。