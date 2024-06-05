package common

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	config     *viper.Viper
	configOnce sync.Once
)

func NewViper() *viper.Viper {
	configOnce.Do(
		func() {
			config = viper.New()
		},
	)
	return config
}
