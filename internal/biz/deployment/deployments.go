package deployment

import (
	"github.com/tangx/k8sailor/internal/k8sdao"
	corev1 "k8s.io/api/core/v1"
)

type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// Replicas 实际期望的 pod 数量
	Replicas int32 `json:"replicas"`

	// 镜像列表
	Images []string `json:"images"`

	Status DeploymentStatus `json:"status"`
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
func GetAllDeployments(input GetAllDeploymentsInput) ([]Deployment, error) {

	v1Deps, err := k8sdao.GetAllDeployments(input.Namespace)
	if err != nil {
		return nil, err
	}

	deps := make([]Deployment, len(v1Deps))
	for i, item := range v1Deps {
		deps[i] = Deployment{
			Name:      item.Name,
			Namespace: item.Namespace,
			Replicas:  *item.Spec.Replicas,
			Images:    podImages(item.Spec.Template.Spec),
			Status: DeploymentStatus{
				Replicas:            item.Status.Replicas,
				AvailableReplicas:   item.Status.AvailableReplicas,
				UnavailableReplicas: item.Status.UnavailableReplicas,
			},
		}
	}

	return deps, nil
}

// 返回 Pod 的镜像列表
func podImages(podSpec corev1.PodSpec) []string {
	images := containerImages(podSpec.Containers)
	initImages := containerImages(podSpec.InitContainers)

	return append(images, initImages...)
}

// containerImages 返回容器的镜像列表
func containerImages(containers []corev1.Container) []string {
	n := len(containers)
	images := make([]string, n)
	for i := 0; i < n; i++ {
		images[i] = containers[i].Image
	}
	return images
}
