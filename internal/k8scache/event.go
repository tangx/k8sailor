package k8scache

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
)

type PodEventCache struct {
	// once 用于初始化
	once sync.Once

	rwmu sync.RWMutex
	tank map[string]PodStatus
}
type PodStatus struct {
	Reason  string
	Message string
}

func (e *PodEventCache) InformerKind() string {
	return "event"
}

func (e *PodEventCache) OnAdd(obj interface{}) {
	event, ok := e.validEvent(obj)
	if !ok {
		return
	}

	e.onAdd(event)

	logrus.Debugln("pod event on add")
}

func (e *PodEventCache) OnDelete(obj interface{}) {
	event, ok := e.validEvent(obj)
	if !ok {
		return
	}

	e.onDelete(event)
	logrus.Debugln("pod event on delete")

}

func (e *PodEventCache) OnUpdate(old, new interface{}) {
	_, ok := e.validEvent(old)
	newEvent, ok2 := e.validEvent(new)
	if !ok || !ok2 {
		return
	}

	e.onUpdate(newEvent)

	logrus.Debugln("pod event on update")
}

func (e *PodEventCache) onAdd(event *corev1.Event) {
	e.initialOnce()

	e.rwmu.Lock()
	defer e.rwmu.Unlock()

	// key := "pod-namespace-podname"
	key := e.eventKey(event)
	e.tank[key] = PodStatus{
		Message: event.Message,
		Reason:  event.Reason,
	}
}

func (e *PodEventCache) onDelete(event *corev1.Event) {

	e.rwmu.Lock()
	defer e.rwmu.Unlock()

	key := e.eventKey(event)
	delete(e.tank, key)
}

func (e *PodEventCache) onUpdate(event *corev1.Event) {

	e.rwmu.Lock()
	defer e.rwmu.Unlock()

	key := e.eventKey(event)
	e.tank[key] = PodStatus{
		Message: event.Message,
		Reason:  event.Reason,
	}
}

// initialOnce 初始化
func (e *PodEventCache) initialOnce() {

	initialFunc := func() {
		if e.tank == nil {
			e.tank = make(map[string]PodStatus)
		}
	}

	e.once.Do(initialFunc)
}

// validEvent 检查是否为 pod 事件
func (e *PodEventCache) validEvent(obj interface{}) (*corev1.Event, bool) {
	event, ok := obj.(*corev1.Event)
	if !ok {
		return nil, false
	}

	// 非 pod 事件， 丢弃
	if event.InvolvedObject.Kind != "Pod" {
		return nil, false
	}

	return event, true
}

func (e *PodEventCache) eventKey(event *corev1.Event) string {
	return e.keyname(event.InvolvedObject.Kind, event.InvolvedObject.Namespace, event.InvolvedObject.Name)
}

func (e *PodEventCache) keyname(kind string, namespace string, name string) string {
	return fmt.Sprintf("%s-%s-%s", kind, namespace, name)
}

func (e *PodEventCache) GetPodEvent(ctx context.Context, namespace string, name string) PodStatus {
	key := e.keyname("Pod", namespace, name)

	e.rwmu.RLock()
	defer e.rwmu.RUnlock()

	status := e.tank[key]
	return status
}
