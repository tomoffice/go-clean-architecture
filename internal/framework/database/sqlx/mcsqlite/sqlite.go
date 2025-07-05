package mcsqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// ConnConfig SQLite 連接配置
type ConnConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConnConfig 創建預設連接配置
func DefaultConnConfig(dsn string) *ConnConfig {
	return &ConnConfig{
		DSN:             dsn,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}
}

// NewDB 使用預設配置創建 SQLite 資料庫連接
func NewDB(dsn string) (*sqlx.DB, error) {
	cfg := DefaultConnConfig(dsn)
	return NewDBWithConfig(cfg)
}

// NewDBWithConfig 使用自訂配置創建 SQLite 資料庫連接
func NewDBWithConfig(cfg *ConnConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
