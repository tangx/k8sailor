package k8sdao

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListReplicaSet(ctx context.Context, namespace string, labelSelector string) (*appsv1.ReplicaSetList, error) {
	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	return clientset.AppsV1().ReplicaSets(namespace).List(ctx, opts)
}
