package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/lengocson131002/go-clean/pkg/logger"
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

var db *gorm.DB
var onceDb sync.Once

func GetDatabase(p *PostgresConfig, logger logger.LoggerInterface) *gorm.DB {

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
	Logger logger.LoggerInterface
}

func (l *logWriter) Printf(message string, args ...interface{}) {
	l.Logger.Trace(message, args...)
}
