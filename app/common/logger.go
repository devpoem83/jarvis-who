package common

import (
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var Logger logrus.Logger

func LoggerInit(level string) {
	logLevel := logrus.DebugLevel
	timeFormat := "2006-01-02 15:04:05"
	logFormat := "[%lvl%] %time% --> %msg% \n"

	if level == "info" {
		logLevel = logrus.InfoLevel
	} else if level == "warn" {
		logLevel = logrus.WarnLevel
	} else if level == "error" {
		logLevel = logrus.ErrorLevel
	}

	Logger = logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: timeFormat,
			LogFormat:       logFormat,
		},
	}

	Logger.Info("Log Level : ", logLevel)
}
