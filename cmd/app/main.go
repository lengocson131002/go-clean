package main

import (
	"github.com/lengocson131002/go-clean/config"
	"github.com/lengocson131002/go-clean/presenter/http/v1"
)

func main() {
	// config all dependencies
	cfg := config.GetConfigure()
	log := config.GetLogger()
	validate := config.GetValidator()

	bootstrapConfig := config.BootstrapConfig{
		Config:    cfg,
		Logger:    log,
		Validator: validate,
	}

	// Start http
	err := http.RunServer(&bootstrapConfig)
	if err != nil {
		panic(err)
	}

	// Start GRPC

	// Start ...
}
