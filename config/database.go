package config

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lengocson131002/go-clean/pkg/logger"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

// GORM

func GetDatabase(p *PostgresConfig, logger logger.Logger) *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.Username, p.Password, p.Database, p.SslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(&logWriter{Logger: logger}, gormlogger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  gormlogger.Info,
		}),
	})

	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
	}

	// config connections
	connection.SetMaxIdleConns(p.IdleConnection)
	connection.SetMaxOpenConns(p.MaxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(p.MaxLifeTimeConnection))

	return db
}

type logWriter struct {
	Logger logger.Logger
}

func (l *logWriter) Printf(message string, args ...interface{}) {
	l.Logger.Trace(message, args...)
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
