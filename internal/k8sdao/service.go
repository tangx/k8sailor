package k8sdao

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceSpec struct {
	// 8080:80
	// !8080:80
	Port string
}

func CreateServcieByName(ctx context.Context, namespace string, name string) {
	opts := metav1.CreateOptions{}
	svc := v1.Service{}
	clientset.CoreV1().Services(namespace).Create(ctx, &svc, opts)
}

func GetServiceByName(ctx context.Context, namespace string, name string) (*v1.Service, error) {
	opts := metav1.GetOptions{}
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("Get Service Failed: %w", err)
	}

	return svc, nil
}
