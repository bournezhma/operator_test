package pkg

import (
	net "k8s.io/api/networking/v1"
	api "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	coreInformer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	netLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
)

type controller struct {
	client        kubernetes.Interface
	serviceLister coreLister.ServiceLister
	ingressLister netLister.IngressLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) updateService(oldObj, newObj interface{}) {
	//TODO: compare annotations
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

// 只需要入队一个key就行，worker可以根据key去Indexer里查找obj
func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	c.queue.Add(key)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*net.Ingress)
	service := api.GetControllerOf(ingress)
	if service != nil {
		return
	}
	if service.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) Run(stopCh chan struct{}) {
	<-stopCh
}

func NewController(client kubernetes.Interface, serviceInformer coreInformer.ServiceInformer,
	ingressInformer netInformer.IngressInformer) controller {
	c := controller{
		client:        client,                   // 与 APIServer 交互
		serviceLister: serviceInformer.Lister(), // Indexer，缓解 APIServer 压力
		ingressLister: ingressInformer.Lister(), // Indexer，缓解 APIServer 压力
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(),
			"ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c
}
