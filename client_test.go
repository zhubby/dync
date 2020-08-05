package dync_test

import (
	"github.com/spf13/viper"
	_ "github.com/zhubby/dync"
	"log"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	var stop = make(chan int)

	tick := time.NewTicker(5 * time.Second)

	for range tick.C {
		log.Printf("%+v", viper.AllSettings())
	}

	<-stop
}
