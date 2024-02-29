package bootstrap

import "github.com/lengocson131002/go-clean/pkg/config"

type T24Config struct {
	Username   string
	MqHost     string
	MqPort     int
	MqChannel  string
	MqManager  string
	MqNameIn   string
	MqNameOut  string
	MqUsername string
	MqPassword string
	MqTimeout  int
}

func GetT24MqConfig(cfg config.Configure) *T24Config {
	return &T24Config{
		Username:   cfg.GetString("T24_USERNAME"),
		MqHost:     cfg.GetString("T24_MQ_HOST"),
		MqPort:     cfg.GetInt("T24_MQ_PORT"),
		MqChannel:  cfg.GetString("T24_MQ_CHANNEL"),
		MqManager:  cfg.GetString("T24_MQ_MANAGER"),
		MqNameIn:   cfg.GetString("T24_MQ_NAME_IN"),
		MqNameOut:  cfg.GetString("T24_MQ_NAME_OUT"),
		MqTimeout:  cfg.GetInt("T24_MQ_TIMEOUT_MS"),
		MqUsername: cfg.GetString("T24_MQ_USERNAME"),
		MqPassword: cfg.GetString("T24_MQ_PASSWORD"),
	}
}
