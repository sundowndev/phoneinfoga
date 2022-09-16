package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/sundowndev/phoneinfoga/v2/build"
)

type Config struct {
	Level        logrus.Level
	ReportCaller bool
}

func getConfig() Config {
	config := Config{
		Level:        logrus.WarnLevel,
		ReportCaller: false,
	}

	if !build.IsRelease() {
		config.Level = logrus.DebugLevel
	}

	return config
}
