package k8scache

import "github.com/tangx/k8sailor/cmd/k8sailor/global"

// Handler Group
var (
	DepTank = &DeploymentCache{}
)

func init() {
	clientset := global.KubeClient.Client()
	informer := global.KubeInformer

	informer.WithClientset(clientset).WithEventHandler(DepTank)
	informer.Start()
}
