# 连接 k3s 集群

> tag: https://github.com/tangx/k8sailor/tree/feat/02-connect-cluster


使用 sdk 链接 k3s cluster 并获取 deployment 信息

```bash
cd cmd/k8sailor && go run .

 * my-nginx-1 (1 replicas)
 * my-nginx-2 (2 replicas)
```

## 下载 client-go sdk


之前在安装 k3s 集群的时候，版本是 v0.21.4。 因此。 这里选择 client-go sdk 的版本也是 v0.21.4

如果还有其他环境， 可以使用 `go mod edit` 命令锁定 client-go 的版本

```bash
go get k8s.io/client-go@v0.21.4

go mod edit -replace=k8s.io/client-go=k8s.io/client-go@v0.21.4
```

## 连接集群并获取 deployment

> https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go


连接到 cluster 的鉴权方式有多种， 后面可以根据 cobra 传递的参数值， 选择不同的鉴权方式。
这里直接参考官方 demo 使用配置文件方式鉴权。

修改一下 kubeconfig 配置来源地址。

[pkg/k8s/cluster.go](/pkg/k8s/cluster.go)

```go

// 从 cobra 配置中获取地址
kubeconfig := &global.Flags.Config

// 其他一样
config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
if err != nil {
    panic(err)
}
clientset, err := kubernetes.NewForConfig(config)
if err != nil {
    panic(err)
}


/* clientset 测试开始， 打印 default namespace 下的所有 deployment */
depClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
list, err := depClient.List(context.TODO(), metav1.ListOptions{})
if err != nil {
	panic(err)
}
for _, d := range list.Items {
	fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
}
/* clienset 测试结束 */

```

## 运行

在 `cmd/root.go` 中调用 k8s cluster 的连接函数

```go
var rootCmd = &cobra.Command{
	// .. 省略
	Run: func(cmd *cobra.Command, args []string) {
		// 连接 k3s
		k8s.Connent()
	},
}
```

运行结果如开篇所示。
