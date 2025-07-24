package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"github.com/tomoffice/go-clean-architecture/config"
	"github.com/tomoffice/go-clean-architecture/internal/bootstrap"
	"github.com/tomoffice/go-clean-architecture/internal/framework/logger"
)

func main() {
	// 1. 載入配置
	cfg, err := config.Load()
	if err != nil {
		panic("配置載入失敗: " + err.Error())
	}

	// 2. 創建 logger
	appLogger, cleanup, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("創建 Logger 失敗: " + err.Error())
	}
	defer func() {
		if err := cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}()

	// 3. 創建並啟動應用
	app := bootstrap.NewApp(cfg, appLogger)
	app.Run()
}
