package main

import (
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"time"
)

func main() {

	// RESTClient
	/*	//config
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile) // ~/.kube/config
		if err != nil {
			panic(err)
		}
		config.GroupVersion = &v1.SchemeGroupVersion
		config.NegotiatedSerializer = scheme.Codecs
		config.APIPath = "/api"

		//client
		restClient, err := rest.RESTClientFor(config)
		if err != nil {
			panic(err)
		}

		//get data
		pod := v1.Pod{}
		err = restClient.Get().Namespace("default").Resource("pods").Name("samplepod").Do(context.TODO()).Into(&pod)
		if err != nil {
			println(err)
		} else {
			println(pod.Name)
		}*/

	// ClientSet
	/*	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile) // ~/.kube/config
		if err != nil {
			panic(err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		coreV1 := clientset.CoreV1()
		pod, err := coreV1.Pods("default").Get(context.TODO(), "samplepod", v1.GetOptions{})
		if err != nil {
			println(err)
		} else {
			println(pod.Name)
		}*/

	// Informer 08
	// create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	// create client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// get informer
	//factory := informers.NewSharedInformerFactory(clientset, 0)
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("default"))
	informer := factory.Core().V1().Pods().Informer()

	// add workqueue
	rateLimitingQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller")

	// add event handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("@%v\tEvent ADD\n", time.Now())
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitingQueue.AddRateLimited(key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("@%v\tEvent UPDATE\n", time.Now())
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitingQueue.AddRateLimited(key)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("@%v\tEvent DELETE\n", time.Now())
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Println("can't get key")
			}
			rateLimitingQueue.AddRateLimited(key)
		},
	})

	// start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	<-stopCh
}
