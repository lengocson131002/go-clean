package config

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var web *fiber.App
var onceWeb sync.Once

type ServerConfig struct {
	Name       string
	AppVersion string
	Port       int
	BaseURI    string
}
