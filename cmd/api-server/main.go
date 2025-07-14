package main

import (
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"github.com/tomoffice/go-clean-architecture/config"
	"github.com/tomoffice/go-clean-architecture/internal/bootstrap"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("配置載入失敗: " + err.Error())
	}
	middlewareContainer := middleware.NewContainer()
	app := bootstrap.NewApp(cfg, middlewareContainer)
	app.Run()
}
