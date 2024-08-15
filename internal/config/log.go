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
