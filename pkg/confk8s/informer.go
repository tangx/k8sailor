package confk8s

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type InformerEventHandler interface {
	InformerKind() string
	cache.ResourceEventHandler
}

type Informer struct {
	factory       informers.SharedInformerFactory
	DefaultResync time.Duration `env:""`
	clientset     *kubernetes.Clientset
}

func (inf *Informer) SetDefaults() {
	if inf.DefaultResync == 0 {
		inf.DefaultResync = 5 * time.Second
	}
}

func (inf *Informer) initial() {
	if inf.factory != nil {
		return
	}

	inf.SetDefaults()
	inf.factory = informers.NewSharedInformerFactory(inf.clientset, inf.DefaultResync)
}

// WithEventHandlers 注册 handler
func (inf *Informer) WithEventHandlers(handlers ...InformerEventHandler) *Informer {
	for _, handler := range handlers {
		kind := handler.InformerKind()
		switch kind {
		case "deployment":
			inf.factory.Apps().V1().Deployments().Informer().AddEventHandler(handler)
		case "event":
			inf.factory.Core().V1().Events().Informer().AddEventHandler(handler)
		}

	}

	return inf
}

// WithClientset 注册 kube clienset， 并创建初始化
func (inf *Informer) WithClientset(clientset *kubernetes.Clientset) *Informer {
	inf.clientset = clientset
	inf.initial()

	return inf
}

func (inf *Informer) Start() {
	fmt.Println("starting informer factory")
	inf.factory.Start(wait.NeverStop)
}
