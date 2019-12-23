package main

import (
	"fmt"
	"github.com/jacky2478/mlog/adapter"
)
var log = adapter.GetLog("example_adapter_sys")

func main() {
	useSys()
	useMlog()
}

func useSys() {
	log.Info("test Info")
	log.Infof("test Infof: %s", "hello mlog")

	log.Debug("test Debug")
	log.Debugf("test Debugf: %s", "hello mlog")

	log.Warning("test Warning")
	log.Warningf("test Warningf: %s", "hello mlog")

	log.Error("test Error")
	log.Errorf("test Errorf: %s", "hello mlog")
}

func useMlog() {
	adapter.UseMlog()
	log = adapter.GetLog("example_adapter_mlog")
	fmt.Println("")

	log.Info("test Info")
	log.Infof("test Infof: %s", "hello mlog")

	log.Debug("test Debug")
	log.Debugf("test Debugf: %s", "hello mlog")

	log.Warning("test Warning")
	log.Warningf("test Warningf: %s", "hello mlog")

	log.Error("test Error")
	log.Errorf("test Errorf: %s", "hello mlog")
}