package deployment

import (
	"github.com/tangx/k8sailor/internal/biz/pod"
	"github.com/tangx/k8sailor/internal/k8sdao"
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

type GetAllDeploymentsInput struct {
	Namespace string `query:"namespace"`
}

// GetAllDeployments 获取 namespace 下的所有 deployments
func GetAllDeployments(input GetAllDeploymentsInput) ([]*Deployment, error) {

	v1Deps, err := k8sdao.GetAllDeployments(input.Namespace)
	if err != nil {
		return nil, err
	}

	deps := make([]*Deployment, len(v1Deps))
	for i, item := range v1Deps {
		deps[i] = &Deployment{
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

	return deps, nil
}
