package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lengocson131002/go-clean/pkg/database"
)

func NewSqlxTxDataSql(db *sqlx.DB) *database.TxDataSql {
	return &database.TxDataSql{
		DB: NewSqlxDBGdbc(db),
	}
}
