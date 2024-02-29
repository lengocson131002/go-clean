package bootstrap

import (
	"fmt"
	"time"

	"github.com/lengocson131002/go-clean/infras/data"
	"github.com/lengocson131002/go-clean/pkg/config"
	"github.com/lengocson131002/go-clean/pkg/database"
	_ "github.com/lib/pq"
)

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

func GetYugabyteConfig(cfg config.Configure) *YugabyteConfig {
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
