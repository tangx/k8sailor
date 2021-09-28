package k8sdao

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodByLabels(ctx context.Context, namespace string, labels map[string]string) (*corev1.PodList, error) {

	opts := metav1.ListOptions{
		LabelSelector: convertMapToSelector(labels),
	}

	return clientset.CoreV1().Pods(namespace).List(ctx, opts)
}
