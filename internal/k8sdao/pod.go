package k8sdao

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodByLabels(namespace string, labels map[string]string) (*corev1.PodList, error) {
	ctx := context.TODO()

	opts := metav1.ListOptions{
		LabelSelector: convertMapToSelector(labels),
	}

	return clientset.CoreV1().Pods(namespace).List(ctx, opts)
}

func convertMapToSelector(labels map[string]string) string {
	l := []string{}
	for k, v := range labels {
		l = append(l, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(l, ",")
}
