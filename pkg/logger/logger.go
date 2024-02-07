package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func NewLogger(env string, file string) *logrus.Logger {

	var log = logrus.New()

	src, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err", err)
	}
	log.Out = src

	switch env {
	case EnvLocal:
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05"
		customFormatter.FullTimestamp = true
		log.Formatter = customFormatter
	case EnvProd:
		customFormatter2 := new(logrus.JSONFormatter)
		customFormatter2.TimestampFormat = "2006-01-02 15:04:05"
		log.Formatter = customFormatter2
		log.Debug("logger debug mode enabled")
	}

	return log
}
