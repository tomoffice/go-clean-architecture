package database

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql

var fs embed.FS

func InitSQLiteDB(dsn string, migrationFiles ...string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	for _, f := range migrationFiles {
		if err := ExecMigration(db, f); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func ExecMigration(db *sqlx.DB, file string) error {
	sql, err := fs.ReadFile("migrations/" + file)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	//fmt.Printf("=== Migration [%s] ===\n%s\n", file, string(sql)) // Debug 用
	_, err = db.Exec(string(sql))
	return err
}
