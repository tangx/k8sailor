package ingress

import (
	"context"

	"github.com/tangx/k8sailor/internal/k8sdao"
	netv1 "k8s.io/api/networking/v1"
)

type GetIngressByNameInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
}

func GetIngressByName(ctx context.Context, input GetIngressByNameInput) (*netv1.Ingress, error) {
	return k8sdao.GetIngressByName(ctx, input.Namespace, input.Name)
}

type CreateIngressByNameInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
	Body      Rules  `body:"" mime:"json"`
}

type Rules struct {
	Endpoints []string `json:"endpoints"`
}

func CreateIngressByName(ctx context.Context, input CreateIngressByNameInput) (*netv1.Ingress, error) {
	// input = CreateIngressByNameInput{
	// 	Namespace: "default",
	// 	Name:      "my-nginx-123",
	// 	Body: []string{
	// 		"http://www.baidu.com/api?backend=my-nginx-service:8080&tls=secret1",
	// 		"http://www.baidu.com/v0/api*",
	// 	},
	// }

	return k8sdao.CreateIngressByName(ctx, input.Namespace, input.Name, input.Body.Endpoints)
}
