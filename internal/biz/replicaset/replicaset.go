package replicaset

import (
	"context"

	"github.com/tangx/k8sailor/internal/k8sdao"
)

type ReplicaSet struct {
	Namespace string
	Name      string
	Labels    map[string]string
}

type ListReplicaSetInput struct {
	Namespace string
	Labels    map[string]string
}

// ListReplicaSet 当前有 Pod 的关联的 ReplicaSet
func ListReplicaSet(ctx context.Context, input ListReplicaSetInput) ([]ReplicaSet, error) {
	v1Rs, err := k8sdao.ListReplicaSet(ctx, input.Namespace, input.Labels)
	if err != nil {
		return nil, err
	}

	rs := []ReplicaSet{}
	for _, item := range v1Rs.Items {
		if item.Status.FullyLabeledReplicas != 0 {
			rs = append(rs, ReplicaSet{
				Namespace: item.Namespace,
				Name:      item.Name,
				Labels:    item.Spec.Selector.MatchLabels,
			})
		}
	}

	return rs, nil
}
