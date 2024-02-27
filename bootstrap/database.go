package bootstrap

import (
	"fmt"
	"time"

	"github.com/lengocson131002/go-clean/infras/data"
	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/env"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	SslMode               string
	IdleConnection        int
	MaxConnection         int
	MaxLifeTimeConnection int //seconds
	MaxIdleTimeConnection int // seconds
}

type YugabyteConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	SslMode               string
	IdleConnection        int
	MaxConnection         int
	MaxLifeTimeConnection int //seconds
	MaxIdleTimeConnection int // seconds
}

func GetUserDatabaseConfig(cfg env.Configure) *PostgresConfig {
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

func GetDatabaseConnector() database.DatabaseConnector {
	return database.NewSqlxDatabaseConnector()
}

func GetUserDatabase(p *PostgresConfig, conn database.DatabaseConnector) *data.UserDatabase {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.Username, p.Password, p.Database, p.SslMode)

	db, err := conn.Connect("postgres", dsn, &database.PoolOptions{
		MaxIdleCount: p.IdleConnection,
		MaxOpen:      p.MaxConnection,
		MaxLifetime:  time.Duration(p.MaxLifeTimeConnection) * time.Second,
		MaxIdleTime:  time.Duration(p.MaxIdleTimeConnection) * time.Second,
	})

	if err != nil {
		panic(err)
	}

	return &data.UserDatabase{
		DB: db,
	}
}

func GetMasterDataDatabase(y *YugabyteConfig, conn database.DatabaseConnector) *data.MasterDataDatabase {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s", y.Username, y.Password, y.Host, y.Port, y.Database, fmt.Sprintf("sslmode=%s", y.SslMode))

	db, err := conn.Connect("postgres", dsn, &database.PoolOptions{
		MaxIdleCount: y.IdleConnection,
		MaxOpen:      y.MaxConnection,
		MaxLifetime:  time.Duration(y.MaxLifeTimeConnection) * time.Second,
		MaxIdleTime:  time.Duration(y.MaxIdleTimeConnection) * time.Second,
	})

	if err != nil {
		panic(err)
	}

	return &data.MasterDataDatabase{
		DB: db,
	}
}
