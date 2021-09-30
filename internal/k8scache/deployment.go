package k8scache

import "sync"

// DeploymentCache 为本地 deployment cache
//   Tank 水箱， 为 Deployment 的容器
//   Tank 的数据结构为 map[namespace]map[name]*v1.Deployment
//   外部使用 sync.map 管理不同 namespace， 内部使用 map[name]*v1.Deployment 不同 deployment。
type DeploymentCache struct {
	rwmu sync.RWMutex
	Tank sync.Map
}

func (d *DeploymentCache) InformerKind() string {
	return "deployment"
}
func (d *DeploymentCache) OnAdd(obj interface{})         {}
func (d *DeploymentCache) OnUpdate(old, new interface{}) {}
func (d *DeploymentCache) OnDelete(obj interface{})      {}
