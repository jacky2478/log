package main

import (
	"github.com/jacky2478/mlog"
)

func initLog() {
	mlog.SetDepth(4)
	mlog.SetFlags(mlog.LstdFlags)
	mlog.SetHighlighting(false)
	mlog.SetLevel(mlog.LOG_LEVEL_ALL)
}

func main() {
	initLog()

	mlog.Info("test Info")
	mlog.Infof("test Infof: %s", "hello mlog")

	mlog.Debug("test Debug")
	mlog.Debugf("test Debugf: %s", "hello mlog")

	mlog.Warning("test Warning")
	mlog.Warningf("test Warningf: %s", "hello mlog")

	mlog.Error("test Error")
	mlog.Errorf("test Errorf: %s", "hello mlog")
}