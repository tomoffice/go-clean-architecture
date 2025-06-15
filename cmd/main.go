package main

import (
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"module-clean/config"
)

func main() {
	cfg := config.LoadConfig()
	app := NewApp(cfg)
	app.Run()
}
