package main

import (
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"module-clean/config"
	"module-clean/internal/bootstrap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("配置載入失敗: " + err.Error())
	}
	app := bootstrap.NewApp(cfg)
	app.Run()
}
