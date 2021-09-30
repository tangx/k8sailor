package k8scache

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentMapper = map[string]*appsv1.Deployment

// DeploymentCache 为本地 deployment cache
//   Tank 水箱， 为 Deployment 的容器
//   Tank 的数据结构为 map[namespace]map[name]*v1.Deployment
//   外部使用 sync.map 管理不同 namespace， 内部使用 map[name]*v1.Deployment 不同 deployment。
type DeploymentCache struct {
	rwmu sync.RWMutex
	tank sync.Map
}

func (d *DeploymentCache) InformerKind() string {
	return "deployment"
}
func (d *DeploymentCache) OnAdd(obj interface{}) {
	mapper, dep, ok := d.extract(obj)
	if !ok {
		return
	}

	d.onAdd(mapper, dep)
}

func (d *DeploymentCache) OnUpdate(old, new interface{}) {
	mapper, _, ok := d.extract(old)
	if !ok {
		return
	}
	_, dep, ok := d.extract(new)
	if !ok {
		return
	}

	d.onUpdate(mapper, dep)
}

func (d *DeploymentCache) OnDelete(obj interface{}) {
	mapper, dep, ok := d.extract(obj)
	if !ok {
		return
	}

	d.onDelete(mapper, dep)
}

// onAdd 添加对象
func (d *DeploymentCache) onAdd(depMap DeploymentMapper, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	depMap[dep.Name] = dep
	logrus.Debugf("add deployment %s of %s", dep.Name, dep.Namespace)
}

func (d *DeploymentCache) onDelete(mapper DeploymentMapper, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	delete(mapper, dep.Name)
	logrus.Debugf("delete deployment %s of %s", dep.Name, dep.Namespace)
}

func (d *DeploymentCache) onUpdate(mapper DeploymentMapper, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	mapper[dep.Name] = dep
	logrus.Debugf("update deployment %s of %s", dep.Name, dep.Namespace)
}

func newDeploymentMap() DeploymentMapper {
	return make(DeploymentMapper)
}

// extract 提取信息
func (d *DeploymentCache) extract(obj interface{}) (mapper DeploymentMapper, dep *appsv1.Deployment, ok bool) {

	dep, ok = obj.(*appsv1.Deployment)
	// 如果 obj 不是 appsv1 对象， 则推出
	if !ok {
		return nil, nil, false
	}

	namespace := dep.Namespace
	mapper, ok = d.getDeploymentMapper(namespace)
	if !ok {
		return nil, dep, false
	}

	return mapper, dep, true
}

func (d *DeploymentCache) getDeploymentMapper(namespace string) (mapper DeploymentMapper, ok bool) {
	obj, ok := d.tank.Load(namespace)
	if !ok {
		obj := newDeploymentMap()
		d.tank.Store(namespace, obj)
	}

	mapper, ok = obj.(DeploymentMapper)
	return mapper, ok
}

/* 业务功能 */

// ListDeployments 返回 namespace 下的所有 deployments
func (d *DeploymentCache) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	mapper, _ := d.getDeploymentMapper(namespace)

	// 返回非指针对象，兼容 appsv1.DeploymentList.Items()
	depList := []appsv1.Deployment{}
	for key := range mapper {
		depList = append(depList, *mapper[key])
	}

	return depList, nil
}

func (d *DeploymentCache) GetDeploymentByName(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	mapper, _ := d.getDeploymentMapper(namespace)

	dep, ok := mapper[name]
	if !ok {
		return nil, fmt.Errorf("no deployment named by %s\n", name)
	}
	return dep, nil
}
