package dync

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

var globalConfig = Config{
	Debug:               true,
	Namespace:           v1.NamespaceDefault,
	CacheResync:         30 * time.Second,
	ConfigMapNameAndKey: map[string][]string{"dyn-config": []string{"config"}},
	FileType:            "yaml",
}

func Setting(c Config) {
	globalConfig = c
}

type Config struct {
	Debug               bool
	Namespace           string
	CacheResync         time.Duration
	ConfigMapNameAndKey map[string][]string
	FileType            string
}
