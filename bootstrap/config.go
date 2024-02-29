package bootstrap

import "github.com/lengocson131002/go-clean/pkg/config"

func GetConfigure() config.Configure {
	var file config.ConfigFile = ".env"
	return config.NewViperConfig(&file)
}
