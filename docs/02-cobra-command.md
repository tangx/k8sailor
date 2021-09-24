# 使用 cobra 管理命令与参数

> tag: https://github.com/tangx/k8sailor/tree/feat/01-cobra-command

为了更加方便的管理配置文件的来源， 这里使用 [cobra](https://github.com/spf13/cobra) 进行命令行构建

效果如下

```bash
cd cmd/k8sailor && go run .
k8s 管理平台

Usage:
  k8sailor [flags]

Flags:
      --config string   k8s 配置授权文件 (default "./k8sconfig/config.yml")
  -h, --help            help for k8sailor
```

## 编码

### 变量管理

在 `cmd/k8sailor/global` 目录中管理 **全局** 变量。

其中，定义一个 `CmdFlag` 结构体管理所有 cobra flags。

```go
type CmdFlags struct {
	Config string `flag:"config" usage:"k8s 配置授权文件" persistent:"true"`
}

var Flags = &CmdFlags{
	Config: "./k8sconfig/config.yml",
}
```

### cobra

在 `cmd/k8sailor/cmd` 中管理所有 cobra 命令。 [root.go](/cmd/k8sailor/cmd/root.go)

在代码中使用了 [cobrautils](https://github.com/go-jarvis/cobrautils) 库帮助管理 flag 绑定。

```go
func init() {
	cobrautils.BindFlags(rootCmd, global.Flags)
}
```

### 启动

在 [main.go](/cmd/k8sailor/main.go) 调用 `cmd/root.go` 的启动函数。 运行结果如上所示。


## 目录结构

```bash
# tree 
.
├── README.md
├── cmd
│   └── k8sailor
│       ├── cmd
│       │   └── root.go
│       ├── global
│       │   └── config.go
│       ├── k8sconfig
│       │   └── config.yml
│       └── main.go
├── go.mod
└── go.sum

6 directories, 9 files
```