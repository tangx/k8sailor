package deployment

import (
	"context"

	"github.com/tangx/k8sailor/internal/biz/pod"
	"github.com/tangx/k8sailor/internal/k8sdao"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// Replicas 实际期望的 pod 数量
	Replicas int32 `json:"replicas"`

	// 镜像列表
	Images []string `json:"images"`

	Status DeploymentStatus `json:"status"`

	LabelSelector *metav1.LabelSelector
}

type DeploymentStatus struct {
	// 标签匹配的 Pod 数量
	Replicas int32 `json:"replicas"`
	// 可用 pod 数量
	AvailableReplicas int32 `json:"availableReplicas"`
	// 不可用数量
	UnavailableReplicas int32 `json:"unavailableReplicas"`
}

type ListDeploymentsInput struct {
	Namespace string `query:"namespace"`
}

// ListDeployments 获取 namespace 下的所有 deployments
func ListDeployments(input ListDeploymentsInput) ([]*Deployment, error) {

	v1Deps, err := k8sdao.ListDeployments(input.Namespace)
	if err != nil {
		return nil, err
	}

	deps := make([]*Deployment, len(v1Deps.Items))
	for i, item := range v1Deps.Items {
		deps[i] = extractDeployment(item)
	}

	return deps, nil
}

type GetDeploymentByNameInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
}

// GetDeploymentByName 通过名称获取 deployment
func GetDeploymentByName(ctx context.Context, input GetDeploymentByNameInput) (*Deployment, error) {
	v1dep, err := k8sdao.GetDeploymentByName(ctx, input.Namespace, input.Name)
	if err != nil {
		return nil, err
	}

	dep := extractDeployment(*v1dep)
	return dep, nil
}

// extractDeployment 转换成业务本身的 Deployment
func extractDeployment(item appsv1.Deployment) *Deployment {
	return &Deployment{
		Name:      item.Name,
		Namespace: item.Namespace,
		Replicas:  *item.Spec.Replicas,
		Images:    pod.PodImages(item.Spec.Template.Spec),
		Status: DeploymentStatus{
			Replicas:            item.Status.Replicas,
			AvailableReplicas:   item.Status.AvailableReplicas,
			UnavailableReplicas: item.Status.UnavailableReplicas,
		},
		LabelSelector: item.Spec.Selector,
	}
}
