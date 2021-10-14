package k8sdao

import (
	"context"
	"fmt"
	"net/url"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngressByName(ctx context.Context, namespace string, name string) (*netv1.Ingress, error) {
	opts := metav1.GetOptions{}

	ing, err := clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("Get Ingress failed: %w", err)
	}

	return ing, nil
}

type IngressBackend struct {
	Endpoint string `json:"endpoint"`
	Service  string `json:"service"`
}

func CreateIngressByName(ctx context.Context, namespace string, name string, backends []IngressBackend) (*netv1.Ingress, error) {
	opts := metav1.CreateOptions{}
	ing := &netv1.Ingress{}
	ing, err := clientset.NetworkingV1().Ingresses(namespace).Create(ctx, ing, opts)
	if err != nil {
		return nil, fmt.Errorf("create Ingress failed: %w", err)
	}

	return ing, err
}

func parseEndpoint(endpoint string) {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Host)
	fmt.Println(u.RawPath)
}
