package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tomoffice/go-clean-architecture/internal/framework/database/sqlx/mcsqlite"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/middleware"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"

	"github.com/tomoffice/go-clean-architecture/config"
	"log"
)

type App struct {
	Config              *config.Config
	MiddlewareContainer *middleware.Container
	Logger              logger.Logger
}

func NewApp(cfg *config.Config, logger logger.Logger) *App {
	return &App{
		Config:              cfg,
		Logger:              logger,
		MiddlewareContainer: middleware.NewContainer(logger),
	}
}

func (a *App) Run() {
	// 初始化數據庫
	db, err := mcsqlite.NewDB(a.Config.Database.DSN)
	if err != nil {
		log.Fatalf("DB 初始化失敗: %v", err)
	}

	// 設置 Gin 引擎
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	
	// middleware - logging middleware should be early in the chain
	engine.Use(a.MiddlewareContainer.Logging())
	engine.Use(a.MiddlewareContainer.CORS())
	
	// 設置 API 路由組
	apiRouterGroup := engine.Group("/api/v1")
	
	// 創建會員模組
	memberModuleFactory := member.NewModuleFactory(a.Logger)
	memberModule, err := memberModuleFactory.CreateModule(db, apiRouterGroup)
	if err != nil {
		log.Fatalf("創建會員模組失敗: %v", err)
	}
	
	// 初始化會員模組
	if err := memberModule.Setup(); err != nil {
		log.Fatalf("初始化會員模組失敗: %v", err)
	}
	log.Printf("模組 %s 初始化成功", memberModule.Name())

	// 啟動服務器
	addr := fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port)
	fmt.Printf("Starting server on %s ...\n", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("啟動服務失敗: %v", err)
	}
}
