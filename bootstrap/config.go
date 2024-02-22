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
	httpPort := cfg.GetInt("APP_HTTP_PORT")
	grpcPort := cfg.GetInt("APP_GRPC_PORT")
	baseUrl := cfg.GetString("APP_BASE_URL")
	grRunningThreshold := cfg.GetInt("APP_GR_RUNNING_THRESHOLD")
	gcMaxPauseThresholdms := cfg.GetInt("APP_GC_PAUSE_THRESHOLD_MS")

	return &ServerConfig{
		Name:               name,
		AppVersion:         version,
		HttpPort:           httpPort,
		GrpcPort:           grpcPort,
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

func GetYugabyteConfig(cfg env.Configure) *YugabyteConfig {
	username := cfg.GetString("DB_YUGABYTE_USER")
	password := cfg.GetString("DB_YUGABYTE_PASSWORD")
	host := cfg.GetString("DB_YUGABYTE_HOST")
	port := cfg.GetInt("DB_YUGABYTE_PORT")
	sslmode := "disable"
	database := cfg.GetString("DB_YUGABYTE_DBNAME")
	idleConnection := cfg.GetInt("DB_YUGABYTE_POOL_IDLE_CONNECTION")
	maxConnection := cfg.GetInt("DB_YUGABYTE_MAX_POOL_SIZE")
	maxLifeTimeConnection := cfg.GetInt("DB_YUGABYTE_MAX_LIFE_TIME")
	maxLifeIdleConnection := cfg.GetInt("DB_YUGABYTE_IDLE_TIMEOUT")

	return &YugabyteConfig{
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

func GetT24MqConfig(cfg env.Configure) *T24Config {
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
