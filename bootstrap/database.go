package bootstrap

import (
	"fmt"
	"time"

	"github.com/lengocson131002/go-clean/pkg/database"
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

func GetDatabaseConnector() database.DatabaseConnector {
	return database.NewSqlxDatabaseConnector()
}

func GetDatabase(p *PostgresConfig, conn database.DatabaseConnector) (*database.Gdbc, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.Username, p.Password, p.Database, p.SslMode)

	db, err := conn.Connect("postgres", dsn, &database.PoolOptions{
		MaxIdleCount: p.IdleConnection,
		MaxOpen:      p.MaxConnection,
		MaxLifetime:  time.Duration(p.MaxLifeTimeConnection) * time.Second,
		MaxIdleTime:  time.Duration(p.MaxIdleTimeConnection) * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return db, err
}
