package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/env"
)

func GetConfigure() env.Configure {
	var file env.ConfigFile = ".env"
	return env.NewViperConfig(&file)
}

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
	maxLifeIdleConnection := cfg.GetInt("DB_POOL_MAX_IDLE_TIME")

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
		MaxIdleTimeConnection: maxLifeIdleConnection,
	}
}

func GetServerConfig(cfg env.Configure) *ServerConfig {
	name := cfg.GetString("APP_NAME")
	version := cfg.GetString("APP_VERSION")
	port := cfg.GetInt("APP_PORT")
	baseUrl := cfg.GetString("APP_BASE_URL")
	grRunningThreshold := cfg.GetInt("APP_GR_RUNNING_THRESHOLD")
	gcMaxPauseThresholdms := cfg.GetInt("APP_GC_PAUSE_THRESHOLD_MS")

	return &ServerConfig{
		Name:               name,
		AppVersion:         version,
		Port:               port,
		BaseURI:            baseUrl,
		GrRunningThreshold: grRunningThreshold,
		GcPauseThresholdMs: gcMaxPauseThresholdms,
		EnvFilePath:        "./.env",
	}
}

func GetTracingConfig(cfg env.Configure) *TraceConfig {
	serviceName := cfg.GetString("TRACE_SERVICE_NAME")
	endpoint := cfg.GetString("TRACE_ENDPOINT")

	return &TraceConfig{
		ServiceName: serviceName,
		Endpoint:    endpoint,
	}
}
