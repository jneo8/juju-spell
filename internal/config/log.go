package config

import (
	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	log.AddHook(filename.NewHook())
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
		ForceColors:     true,
	})
}

func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
