package pod

import (
	"time"

	"github.com/tangx/k8sailor/internal/k8sdao"
)

type Pod struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Images    []string `json:"images"`
	NodeName  string   `json:"nodeName"`

	CreateTime time.Time `json:"createTime"`
}

type GetPodsByLabelsInput struct {
	Namespace string            `query:"namespace"`
	Labels    map[string]string `body:"" mime:"json"`
}

func GetPodsByLabels(input GetPodsByLabelsInput) ([]*Pod, error) {

	v1Pods, err := k8sdao.GetPodByLabels(input.Namespace, input.Labels)
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
		}
	}

	return pods, nil
}
