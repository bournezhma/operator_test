package main

import (
	"github.com/bournezhma/client-go-demo/11/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// Stage 1. create config

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Println("Can't get config from RecommendedHomeFIle. Try using InClusterConfig instead.")
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("Can't get config.")
		}
		config = inClusterConfig
	}

	// Stage 2. create client

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Can't create client.")
	}

	// Stage 3. create informer by factory	（这里还没有真正创建 informer，真正实例化在 controller 中）

	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()

	// Stage 4. register event handler （抽离到 controller 中去）

	// Stage 5. start informer
	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)
}
