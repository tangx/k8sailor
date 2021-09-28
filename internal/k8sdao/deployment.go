package k8sdao

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListDeployments(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	opts := metav1.ListOptions{}
	return clientset.AppsV1().Deployments(namespace).List(ctx, opts)

}

func GetDeploymentByName(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	opts := metav1.GetOptions{}
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}
