package config

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/adrg/xdg"
)

func InitLogLoc(logFile string) error {
	logDir := filepath.Dir(logFile)
	if err := EnsureFullPath(logDir, DefaultDirMod); err != nil {
		return err
	}
	return nil
}

func GetXDGStateFile() string {
	dir, err := xdg.StateFile(AppName)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
