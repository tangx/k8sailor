package deployment

import (
	"context"
	"fmt"

	"github.com/tangx/k8sailor/internal/biz/pod"
	"github.com/tangx/k8sailor/internal/biz/replicaset"
	"github.com/tangx/k8sailor/internal/k8scache"
	"github.com/tangx/k8sailor/internal/k8sdao"
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// Replicas 实际期望的 pod 数量
	Replicas int32 `json:"replicas"`

	// 镜像列表
	Images []string `json:"images"`

	Status DeploymentStatus `json:"status"`

	Labels map[string]string `json:"labelSelector"`
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
// 业务层，可以对接不同来源的数据。
func ListDeployments(ctx context.Context, input ListDeploymentsInput) ([]*Deployment, error) {

	/* k8s api 返回的数据 */
	// v1Deps, err := k8sdao.ListDeployments(ctx, input.Namespace)

	/* 使用 informer 保存在本地的 cache 数据 */
	v1Deps, err := k8scache.DepTank.ListDeployments(ctx, input.Namespace)

	if err != nil {
		return nil, err
	}

	deps := make([]*Deployment, len(v1Deps))
	for i, item := range v1Deps {
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

type GetPodsByDeploymentInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
}

// GetPodsByDeployment 根据 Deployment name 获取所有 Pod
func GetPodsByDeployment(ctx context.Context, input GetPodsByDeploymentInput) ([]*pod.Pod, error) {
	// get deployment
	dInput := GetDeploymentByNameInput{
		Namespace: input.Namespace,
		Name:      input.Name,
	}
	dep, err := GetDeploymentByName(ctx, dInput)
	if err != nil {
		return nil, err
	}

	// get active replica set
	rsInput := replicaset.ListReplicaSetInput{
		Namespace: dep.Namespace,
		Labels:    dep.Labels,
	}
	rsList, err := replicaset.ListReplicaSet(ctx, rsInput)
	if err != nil {
		return nil, err
	}

	// get pods
	allPods := []*pod.Pod{}
	for _, rs := range rsList {
		pInput := pod.GetPodsByLabelsInput{
			Namespace: rs.Namespace,
			Labels:    rs.Labels,
		}

		pods, err := pod.GetPodsByLabels(ctx, pInput)
		if err != nil {
			return nil, err
		}
		fmt.Println(len(pods), pods)
		allPods = append(allPods, pods...)
	}

	return allPods, nil
}

// SetDeploymentReplicasInput 调整 deployment pod 数量参数
// Replicas 为了避免 **0值** 影响。
//   1. 使用为 *int 指针对象， 自行在业务逻辑中进行校验
//   2. 另外也可以使用， `binding` tag， 由 gin 框架的 valicator 帮忙校验。 https://github.com/go-playground/validator
// Namespace 设置了默认值， 如果请求不提供将由 gin 框架自己填充。
type SetDeploymentReplicasInput struct {
	Namespace string `query:"namespace,default=default"`
	Name      string `uri:"name"`
	Replicas  *int   `query:"replicas" binding:"required"`
}

// SetDeploymentReplicas 设置 deployment 的 pod 副本数量
func SetDeploymentReplicas(ctx context.Context, input SetDeploymentReplicasInput) (bool, error) {

	// 参数验证
	if input.Replicas == nil {
		err := fmt.Errorf("replicas must be provide")
		return false, err
	}

	err := k8sdao.SetDeploymentReplicas(ctx, input.Namespace, input.Name, *input.Replicas)

	// err==nil -> true 表示设置成功
	// 这里选择不在 api 层进行判断是为了保持各层的业务分工一致性。
	return err == nil, err
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
		Labels: item.Spec.Selector.MatchLabels,
	}
}
