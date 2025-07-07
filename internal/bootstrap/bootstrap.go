package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"module-clean/internal/framework/database/sqlx/mcsqlite"
	"module-clean/internal/framework/http/gin/middleware"
	"module-clean/internal/modules/member"

	"log"
	"module-clean/config"
)

type App struct {
	Config              *config.Config
	MiddlewareContainer *middleware.Container
}

func NewApp(cfg *config.Config, middlewareContainer *middleware.Container) *App {
	return &App{
		Config:              cfg,
		MiddlewareContainer: middlewareContainer,
	}
}

func (a *App) Run() {
	db, err := mcsqlite.NewDB(a.Config.Database.DSN)
	if err != nil {
		log.Fatalf("DB 初始化失敗: %v", err)
	}
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	apiRouterGroup := engine.Group("/api/v1")
	// middleware
	engine.Use(a.MiddlewareContainer.CORS())
	// 創建會員模組
	memberModuleFactory := member.NewModuleFactory()
	memberModule, err := memberModuleFactory.CreateModule(db, apiRouterGroup)
	if err != nil {
		log.Fatalf("創建會員模組失敗: %v", err)
	}
	// 初始化會員模組
	if err := memberModule.Setup(); err != nil {
		log.Fatalf("初始化會員模組失敗: %v", err)
	}
	log.Printf("模組 %s 初始化成功", memberModule.Name())

	addr := fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port)
	fmt.Printf("Starting server on %s ...\n", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("啟動服務失敗: %v", err)
	}
}
