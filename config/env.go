package config

import (
	"sync"

	"github.com/lengocson131002/go-clean/pkg/env"
)

var cfg env.ConfigInterface
var onceConfig sync.Once

func GetConfigure() env.ConfigInterface {
	if cfg == nil {
		onceConfig.Do(func() {
			cfg = env.NewViperConfig(".env")
		})
	}

	return cfg
}
