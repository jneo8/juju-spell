package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	AppName          = "jujuspell"
	JujuSpellLogFile = "jujuspell.log"
)

var (
	AppLogFile string
)

func InitLogLoc() error {
	var err error
	appLogDir, err := xdg.StateFile(AppName)
	if err != nil {
		return err
	}
	if err := EnsureFullPath(appLogDir, DefaultDirMod); err != nil {
		return err
	}
	AppLogFile = filepath.Join(appLogDir, JujuSpellLogFile)
	return nil
}
