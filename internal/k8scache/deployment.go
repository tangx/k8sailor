package k8scache

import (
	"sync"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
)

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
	set, dep, ok := d.extract(obj)
	if !ok {
		return
	}

	d.onAdd(set, dep)
}

func (d *DeploymentCache) OnUpdate(old, new interface{}) {
	set, _, ok := d.extract(old)
	if !ok {
		return
	}
	_, dep, ok := d.extract(new)
	if !ok {
		return
	}

	d.onUpdate(set, dep)
}

func (d *DeploymentCache) OnDelete(obj interface{}) {
	set, dep, ok := d.extract(obj)
	if !ok {
		return
	}

	d.onDelete(set, dep)
}

// onAdd 添加对象
func (d *DeploymentCache) onAdd(set map[string]*appsv1.Deployment, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	set[dep.Name] = dep
	logrus.Debugf("add deployment %s of %s", dep.Name, dep.Namespace)
}

func (d *DeploymentCache) onDelete(set map[string]*appsv1.Deployment, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	delete(set, dep.Name)
	logrus.Debugf("delete deployment %s of %s", dep.Name, dep.Namespace)
}

func (d *DeploymentCache) onUpdate(set map[string]*appsv1.Deployment, dep *appsv1.Deployment) {

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	set[dep.Name] = dep
	logrus.Debugf("update deployment %s of %s", dep.Name, dep.Namespace)
}

func newDeploymentSet() map[string]*appsv1.Deployment {
	return make(map[string]*appsv1.Deployment)
}

// extract 提取信息
func (d *DeploymentCache) extract(obj interface{}) (set map[string]*appsv1.Deployment, dep *appsv1.Deployment, ok bool) {

	dep, ok = obj.(*appsv1.Deployment)
	// 如果 obj 不是 appsv1 对象， 则推出
	if !ok {
		return nil, nil, false
	}

	namespace := dep.Namespace

	// 提取 sync.Map 中的对象
	objSet, ok := d.tank.Load(namespace)
	// 如果当前 namesapce 不存在， 则新建
	if !ok {
		objSet = newDeploymentSet()
		d.tank.Store(namespace, objSet)
	}

	// 检查 objSet 类型检查
	set, ok = objSet.(map[string]*appsv1.Deployment)
	if !ok {
		return nil, dep, false
	}

	return set, dep, true
}
