package pod

import (
	"context"
	"fmt"
	"time"

	"github.com/tangx/k8sailor/internal/k8sdao"
	corev1 "k8s.io/api/core/v1"
)

type Pod struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Images     []string          `json:"images"`
	NodeName   string            `json:"nodeName"`
	NodeIp     string            `json:"nodeIp"`
	CreateTime time.Time         `json:"createTime"`
	PodIP      string            `json:"podIp"`
	Status     PodStatus         `json:"status"`
	Labels     map[string]string `json:"labels"`
}

type PodStatus struct {
	Phase   corev1.PodPhase `json:"phase"`
	Message string          `json:"message"`
	Reason  string          `json:"reason"`
}

type GetPodsByLabelsInput struct {
	Namespace string            `query:"namespace"`
	Labels    map[string]string `body:"" mime:"json"`
}

func GetPodsByLabels(ctx context.Context, input GetPodsByLabelsInput) ([]*Pod, error) {

	selector := k8sdao.ConvertMapToSelector(input.Labels)
	v1Pods, err := k8sdao.GetPodByLabels(ctx, input.Namespace, selector)
	if err != nil {
		return nil, err
	}

	pods := make([]*Pod, len(v1Pods.Items))

	for i, item := range v1Pods.Items {
		pods[i] = extractPod(item)
	}

	return pods, nil
}

type GetPodByNameInput struct {
	Name         string `uri:"name"`
	Namespace    string `query:"namespace"`
	OutputFormat string `query:"outputFormat,default=json"`
}

func GetCorePodByName(ctx context.Context, input GetPodByNameInput) (*corev1.Pod, error) {
	v1pod, err := k8sdao.GetPodByName(ctx, input.Namespace, input.Name)
	if err != nil {
		return nil, err
	}

	return v1pod, nil
}

type DeletePodByNameInput struct {
	Name      string `uri:"name"`
	Namespace string `query:"namespace"`
}

// DeletePodByName 根据名字 删除 pod
func DeletePodByName(ctx context.Context, input DeletePodByNameInput) error {
	err := k8sdao.DeletePodByName(ctx, input.Namespace, input.Name)
	if err != nil {
		return fmt.Errorf("k8s internal error: %w", err)
	}

	return nil
}

// extractPod 转换成业务本身的 Pod
func extractPod(item corev1.Pod) *Pod {
	return &Pod{
		Name:       item.Name,
		Namespace:  item.Namespace,
		Images:     PodImages(item.Spec),
		NodeName:   item.Spec.NodeName,
		NodeIp:     item.Status.HostIP,
		CreateTime: item.CreationTimestamp.Time,
		PodIP:      item.Status.PodIP,
		Status: PodStatus{
			Phase:   item.Status.Phase,
			Message: item.Status.Message,
			Reason:  item.Status.Reason,
		},
		Labels: item.Labels,
	}
}
