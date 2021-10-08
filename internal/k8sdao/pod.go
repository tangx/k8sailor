package k8sdao

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodByLabels(ctx context.Context, namespace string, labelSelector string) (*corev1.PodList, error) {

	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	return clientset.CoreV1().Pods(namespace).List(ctx, opts)
}

func GetPodByName(ctx context.Context, namespace string, name string) (*corev1.Pod, error) {
	opts := metav1.GetOptions{}
	return clientset.CoreV1().Pods(namespace).Get(ctx, name, opts)
}

func DeletePodByName(ctx context.Context, namespace string, name string) error {
	opts := metav1.DeleteOptions{}
	return clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}
