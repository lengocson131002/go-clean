package config

import (
	"github.com/lengocson131002/go-clean/pkg/env"
)

func GetDatabaseConfig(cfg env.Configure) *PostgresConfig {
	username := cfg.GetString("DB_USERNAME")
	password := cfg.GetString("DB_PASSWORD")
	host := cfg.GetString("DB_HOST")
	port := cfg.GetInt("DB_PORT")
	sslmode := cfg.GetString("DB_SSL_MODE")
	database := cfg.GetString("DB_NAME")
	idleConnection := cfg.GetInt("DB_POOL_IDLE_CONNECTION")
	maxConnection := cfg.GetInt("DB_POOL_MAX_CONNECTION")
	maxLifeTimeConnection := cfg.GetInt("DB_POOL_MAX_LIFE_TIME")

	return &PostgresConfig{
		Username:              username,
		Password:              password,
		Host:                  host,
		Port:                  port,
		Database:              database,
		SslMode:               sslmode,
		IdleConnection:        idleConnection,
		MaxConnection:         maxConnection,
		MaxLifeTimeConnection: maxLifeTimeConnection,
	}
}

func GetServerConfig(cfg env.Configure) *ServerConfig {
	name := cfg.GetString("APP_NAME")
	version := cfg.GetString("APP_VERSION")
	port := cfg.GetInt("APP_PORT")
	baseUrl := cfg.GetString("APP_BASE_URL")

	return &ServerConfig{
		Name:       name,
		AppVersion: version,
		Port:       port,
		BaseURI:    baseUrl,
	}

}
