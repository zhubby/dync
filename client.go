package dync

import (
	"k8s.io/client-go/rest"
)

var	restCfg  *rest.Config

func init() {
	var err error
	restCfg, err = rest.InClusterConfig()
	if err !=nil {
		panic(err)
	}
}