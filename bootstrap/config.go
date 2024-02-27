package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/env"
)

func GetConfigure() env.Configure {
	var file env.ConfigFile = ".env"
	return env.NewViperConfig(&file)
}
