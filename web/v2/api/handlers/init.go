package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"sync"
)

var once sync.Once
var remoteLibrary *remote.Library

func init() {
	once.Do(func() {
		remoteLibrary = remote.NewLibrary(filter.NewEngine())
		remote.InitScanners(remoteLibrary)
		logrus.Debug("Scanners and plugins initialized")
	})
}
