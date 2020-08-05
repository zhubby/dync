package dync

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const PackageName = "dync"

func TempConfigDir() string {
	return os.TempDir() + PackageName
}

func init() {
	if globalConfig.Debug {
		log.Printf("%+v", globalConfig)
	}
	var (
		restCfg *rest.Config
		err     error
	)
	restCfg, err = config.GetConfig()
	if err != nil {
		panic(err)
	}
	caches, err := cache.New(restCfg, cache.Options{Resync: &globalConfig.CacheResync, Namespace: globalConfig.Namespace})
	if err != nil {
		panic(err)
	}
	cmInformer, err := caches.GetInformer(&v1.ConfigMap{})
	if err != nil {
		panic(err)
	}
	handler := newHandler()
	cmInformer.AddEventHandler(handler)
	go caches.Start(nil)
	if caches.WaitForCacheSync(nil) {
		if err := handler.loadConfig(); err != nil {
			panic(err)
		}
		handler.watch()
	}
}

type handler struct {
	fs    afero.Fs
	cmMap map[string][]string
}

func (h handler) loadConfig() error {
	for k, v := range h.cmMap {
		for _, s := range v {
			file := fmt.Sprintf("%s/%s/%s", TempConfigDir(), k, s)
			viper.SetConfigFile(file)
			viper.SetConfigType(globalConfig.FileType)
			if err := viper.MergeInConfig(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h handler) watch() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
}

type ResourceEventHandler interface {
	toolscache.ResourceEventHandler
	loadConfig() error
	watch()
}

func newHandler() ResourceEventHandler {
	return &handler{fs: afero.NewOsFs(), cmMap: globalConfig.ConfigMapNameAndKey}
}

func (h handler) OnAdd(obj interface{}) {
	if cm, ok := obj.(*v1.ConfigMap); ok {
		var fileDir = fmt.Sprintf("%s/%s", TempConfigDir(), cm.Name)

		for _, v := range h.cmMap[cm.Name] {
			if err := h.fs.MkdirAll(fileDir, 0755); err == nil {
				if globalConfig.Debug {
					log.Printf("Create config temp dir %s", fileDir)
				}
			}
			var fileName = fmt.Sprintf("%s/%s", fileDir, v)
			if err := afero.WriteFile(h.fs, fileName, []byte(cm.Data[v]), 0644); err != nil {
				if globalConfig.Debug {
					log.Printf("Write config file fail %v", err)
				}
				continue
			}
		}
	}
}

func (h handler) OnUpdate(oldObj, newObj interface{}) {
	if cm, ok := newObj.(*v1.ConfigMap); ok {
		var fileDir = fmt.Sprintf("%s/%s", TempConfigDir(), cm.Name)

		for _, v := range h.cmMap[cm.Name] {
			if err := h.fs.MkdirAll(fileDir, 0755); err == nil {
				if globalConfig.Debug {
					log.Printf("Create config temp dir %s", fileDir)
				}
			}
			var fileName = fmt.Sprintf("%s/%s", fileDir, v)
			if err := afero.WriteFile(h.fs, fileName, []byte(cm.Data[v]), 0644); err != nil {
				if globalConfig.Debug {
					log.Printf("Write config file fail %v", err)
				}
				continue
			}
		}
	}
}

func (h handler) OnDelete(obj interface{}) {
	if cm, ok := obj.(*v1.ConfigMap); ok {
		fileDir := fmt.Sprintf("%s/%s", TempConfigDir(), cm.Name)
		err := os.RemoveAll(fileDir)
		if err != nil {
			log.Printf("Delete Config Dir %s error %v", fileDir, err)
		}
	}
}
