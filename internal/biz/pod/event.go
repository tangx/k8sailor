package pod

import (
	"context"

	"github.com/tangx/k8sailor/internal/k8scache"
)

type PodEventStatus struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Message   string `json:"message"`
	Reason    string `json:"reason"`
}

type GetPodEventByNameInput struct {
	Namespace string `query:"namespace"`
	Name      string `uri:"name"`
}

func GetPodEventByName(ctx context.Context, input GetPodEventByNameInput) PodEventStatus {
	status := k8scache.EventTank.GetPodEvent(ctx, input.Namespace, input.Name)

	return PodEventStatus{
		Namespace: input.Namespace,
		Name:      input.Name,
		Message:   status.Message,
		Reason:    status.Reason,
	}
}
