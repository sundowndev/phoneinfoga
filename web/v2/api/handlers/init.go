package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"sync"
)

var once sync.Once
var RemoteLibrary *remote.Library

func Init(filterEngine filter.Filter) {
	once.Do(func() {
		RemoteLibrary = remote.NewLibrary(filterEngine)
		remote.InitScanners(RemoteLibrary)
		logrus.Debug("Scanners and plugins initialized")
	})
}
