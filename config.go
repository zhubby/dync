package dync

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

var defaultConfig = Config{
	Namespace:     v1.NamespaceDefault,
	CacheResync:   30 * time.Second,
	ConfigMapName: "dyn-config-map",
}

func Setting(c Config) {
	defaultConfig = c
}

type Config struct {
	Namespace     string `json:"namespace"`
	CacheResync   time.Duration
	ConfigMapName string `json:"config_map_name"`
}
