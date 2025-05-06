package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InitSQLiteDB(dataSourceName string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("failed to open sqlite database: %v", err)
	}

	// 可選擇性做連線檢查
	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping sqlite database: %v", err)
	}

	return db
}
