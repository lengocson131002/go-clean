package config

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
}

// SQLX
func GetDatabaseSqlx(p *PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.Username, p.Password, p.Database, p.SslMode)

	db, err := connectDatabase("postgres", dsn)

	if err != nil {
		return nil, err
	}

	// config connections
	db.SetMaxIdleConns(p.IdleConnection)
	db.SetMaxOpenConns(p.MaxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(p.MaxLifeTimeConnection))

	return db, nil
}

func connectDatabase(driverName string, dsn string) (*sqlx.DB, error) {
	db := sqlx.MustConnect(driverName, dsn)
	err := db.Ping()
	if err != nil {
		if db != nil {
			err = db.Close()
		}
		return nil, err
	}
	return db, err
}
