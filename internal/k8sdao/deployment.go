package k8sdao

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAllDeployments(namespace string) ([]appsv1.Deployment, error) {
	ctx := context.TODO()
	opts := metav1.ListOptions{}
	v1Deps, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return v1Deps.Items, nil
}
