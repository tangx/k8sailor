package service

import (
	"context"

	"github.com/tangx/k8sailor/internal/k8sdao"
	corev1 "k8s.io/api/core/v1"
)

type GetCoreServerByNameInput struct {
	Namespace    string `query:"namespace"`
	Name         string `uri:"name"`
	OutputFormat string `query:"outputFormat"`
}

func GetCoreServerByName(ctx context.Context, input GetCoreServerByNameInput) (*corev1.Service, error) {
	return k8sdao.GetServiceByName(ctx, input.Namespace, input.Name)
}

// CreateServcieByNameInput 参数
// port example:
// 	port:= port:targetPort            clusterIp
// 	port:= !nodePort:port:targetPort  指定 nodePort 值
// 	port:= !port:targetPort           随机 nodePort 值
//
// symbal:
// 	_blank: clusterIp
// 	!: nodePort
type CreateServcieByNameInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
	Body      struct {
		Services []string `json:"services"`
	} `body:"" mime:"json"`
}

func CreateServiceByName(ctx context.Context, input CreateServcieByNameInput) (*corev1.Service, error) {
	v1svc, err := k8sdao.CreateServcieByName(ctx, input.Namespace, input.Name, input.Body.Services)
	if err != nil {
		return nil, err
	}

	return v1svc, nil
}
