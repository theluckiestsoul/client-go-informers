package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	path.Join()
	defaultKubeConfigPath := filepath.Join(homeDir, ".kube", "config")
	cfgFile := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig file")
	cfg, err := clientcmd.BuildConfigFromFlags("", *cfgFile)
	if err != nil {
		fmt.Println("Using in-cluster config")
		cfg, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("failed to build in-cluster config: %v", err)
		}
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("failed to create clientset: %v", err)
	}

	fmt.Println("Clientset created", cs)

	informerFactory := informers.NewSharedInformerFactory(cs, 30*time.Second)

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Pod added")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Pod updated")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Pod deleted")
		},
	})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podInformer.Lister().Pods("default").Get("nginx")
	if err != nil {
		//log.Fatalf("failed to get pod: %v", err)
	}
	fmt.Println("Pod found", pod.Name)
	sigchan := make(chan struct{})

	<-sigchan
}
