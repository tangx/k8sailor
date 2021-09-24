package deployment

import "github.com/tangx/k8sailor/internal/k8sdao"

type Deployment struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Replicas  int32  `json:"replicas,omitempty"`
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
		}
	}

	return deps, nil
}
