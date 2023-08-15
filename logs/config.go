package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/sundowndev/phoneinfoga/v2/build"
	"os"
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

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		loglevel, _ := logrus.ParseLevel(lvl)
		config.Level = loglevel
	}

	return config
}
