package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangx/k8sailor/cmd/k8sailor/global"
	"github.com/tangx/k8sailor/internal/apis"
	"github.com/tangx/k8sailor/internal/k8scache"
)

var cmdHttpserver = &cobra.Command{
	Use:  "httpserver",
	Long: "启动 web 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		// 启动 informer
		runInformer()

		// 启动服务
		runHttpserver()
	},
}

// runHttpserver 启动 http server
func runHttpserver() {
	// 1. 将 apis 注册到 httpserver 中
	global.HttpServer.RegisterRoute(apis.RootGroup)

	// 2. 启动服务
	if err := global.HttpServer.Run(); err != nil {
		logrus.Fatalf("start httpserver failed: %v", err)
	}
}

func runInformer() {

	clientset := global.KubeClient.Client()
	informer := global.KubeInformer.WithClientset(clientset)

	k8scache.RegisterHandlers(informer)

	informer.Start()
}
