package pod

import (
	"context"
	"time"

	"github.com/tangx/k8sailor/internal/k8sdao"
	corev1 "k8s.io/api/core/v1"
)

type Pod struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Images     []string  `json:"images"`
	NodeName   string    `json:"nodeName"`
	CreateTime time.Time `json:"createTime"`
	PodIP      string    `json:"podIp"`
	Status     PodStatus `json:"status"`
	// Status2    corev1.PodStatus
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

	v1Pods, err := k8sdao.GetPodByLabels(ctx, input.Namespace, input.Labels)
	if err != nil {
		return nil, err
	}

	pods := make([]*Pod, len(v1Pods.Items))

	for i, v1pod := range v1Pods.Items {
		pods[i] = &Pod{
			Name:       v1pod.Name,
			Namespace:  v1pod.Namespace,
			Images:     PodImages(v1pod.Spec),
			NodeName:   v1pod.Spec.NodeName,
			CreateTime: v1pod.CreationTimestamp.Time,
			PodIP:      v1pod.Status.PodIP,
			Status: PodStatus{
				Phase:   v1pod.Status.Phase,
				Message: v1pod.Status.Message,
				Reason:  v1pod.Status.Reason,
			},
			// Status2: v1pod.Status,
		}
	}

	return pods, nil
}
