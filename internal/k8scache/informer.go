package k8scache

import (
	"github.com/tangx/k8sailor/pkg/confk8s"
)

// Handler Group
var (
	DepTank = &DeploymentCache{}
)

func RegisterHandlers(informer *confk8s.Informer) {
	informer.WithEventHandlers(DepTank)
}
