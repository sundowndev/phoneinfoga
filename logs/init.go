package logs

import (
	"github.com/sirupsen/logrus"
)

func Init() {
	config := getConfig()
	logrus.SetLevel(config.Level)
	logrus.SetReportCaller(config.ReportCaller)
}
