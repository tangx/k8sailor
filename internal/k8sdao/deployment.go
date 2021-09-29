package k8sdao

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployments 返回 namespace 下的所有 deployments
func ListDeployments(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	opts := metav1.ListOptions{}
	return clientset.AppsV1().Deployments(namespace).List(ctx, opts)

}

// GetDeploymentByName 根据 name 获取 deployment
func GetDeploymentByName(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	opts := metav1.GetOptions{}
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

// SetDeploymentReplicas 设置 namespace 大小
func SetDeploymentReplicas(ctx context.Context, namespace string, name string, replicas int) error {

	opts := metav1.GetOptions{}
	v1Scale, err := clientset.AppsV1().Deployments(namespace).GetScale(ctx, name, opts)
	if err != nil {
		return err
	}

	if replicas < 0 || replicas > 10 {
		return errors.New("invalid replicas number")
	}
	// set new replicas
	v1Scale.Spec.Replicas = int32(replicas)

	upOpts := metav1.UpdateOptions{}
	v1Scale, err = clientset.AppsV1().Deployments(namespace).UpdateScale(ctx, name, v1Scale, upOpts)
	if err != nil {
		return err
	}

	return nil
}
