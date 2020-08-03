package dync

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	toolscache "k8s.io/client-go/tools/cache"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func init() {
	var (
		restCfg *rest.Config
		err     error
	)
	restCfg, err = config.GetConfig()

	if err != nil {
		panic(err)
	}
	caches, err := cache.New(restCfg, cache.Options{Resync: &defaultConfig.CacheResync, Namespace: defaultConfig.Namespace})
	if err != nil {
		panic(err)
	}
	cmInformer, err := caches.GetInformer(&v1.ConfigMap{})
	if err != nil {
		panic(err)
	}
	cmInformer.AddEventHandler(toolscache.ResourceEventHandler(&handler{}))
	go caches.Start(nil)
}

type handler struct {
}

func (h handler) OnAdd(obj interface{}) {
	log.Print(obj)
}

func (h handler) OnUpdate(oldObj, newObj interface{}) {
	log.Print(newObj)
}

func (h handler) OnDelete(obj interface{}) {
	log.Print(obj)
}
