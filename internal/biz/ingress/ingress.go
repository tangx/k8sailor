package ingress

import (
	"context"

	"github.com/tangx/k8sailor/internal/k8sdao"
	netv1 "k8s.io/api/networking/v1"
)

type CreateIngressByNameInput struct {
	Namespace string                  `query:"namespace"`
	Name      string                  `uri:"name"`
	Body      []k8sdao.IngressBackend `body:"" mime:"json"`
}

func CreateIngressByName(ctx context.Context, input CreateIngressByNameInput) (*netv1.Ingress, error) {
	input = CreateIngressByNameInput{
		Namespace: "default",
		Name:      "my-nginx-123",
		Body: []k8sdao.IngressBackend{
			{
				Endpoint: "http://www.baidu.com/api",
				Service:  "my-nginx-service:8080",
			},
			{
				Endpoint: "http://www.baidu.com/v0/api*",
				Service:  "my-nginx-service-v0:8080",
			},
		},
	}

	return k8sdao.CreateIngressByName(ctx, input.Namespace, input.Name, input.Body)
}
