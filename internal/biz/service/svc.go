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
