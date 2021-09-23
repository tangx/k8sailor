package k8s

import (
	"context"
	"fmt"

	"github.com/tangx/k8sailor/cmd/k8sailor/global"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Connent() *kubernetes.Clientset {
	// var kubeconfig *string
	kubeconfig := &global.Flags.Config
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

	return clientset
}
