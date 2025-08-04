package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tomoffice/go-clean-architecture/internal/framework/database/sqlx/mcsqlite"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/middleware"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"

	"github.com/tomoffice/go-clean-architecture/config"
	"log"
)

type App struct {
	Config              *config.Config
	MiddlewareContainer *middleware.Container
	Logger              logger.Logger
	Tracer              tracer.Tracer
}

func NewApp(cfg *config.Config, logger logger.Logger, tracer tracer.Tracer) *App {
	return &App{
		Config:              cfg,
		Logger:              logger,
		Tracer:              tracer,
		MiddlewareContainer: middleware.NewContainer(logger, tracer),
	}
}

func (a *App) Run() {
	a.Logger.Debug("設定值", logger.NewField("config", a.Config))
	// 初始化數據庫
	db, err := mcsqlite.NewDB(a.Config.Database.DSN)
	if err != nil {
		log.Fatalf("DB 初始化失敗: %v", err)
	}
	a.Logger.Debug("DB 連接成功", logger.NewField("dsn", a.Config.Database.DSN))

	// 設置 Gin 引擎
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	// middleware - logging middleware should be early in the chain
	engine.Use(a.MiddlewareContainer.Logging())
	engine.Use(a.MiddlewareContainer.CORS())

	// 設置 API 路由組
	apiRouterGroup := engine.Group("/api/v1")

	// 創建會員模組
	memberModuleFactory := member.NewModuleFactory()
	memberModule, err := memberModuleFactory.CreateModule(db, apiRouterGroup, a.Logger, a.Tracer)
	if err != nil {
		//log.Fatalf("創建會員模組失敗: %v", err)
		a.Logger.Error("創建會員模組失敗", logger.NewField("error", err))
	}
	a.Logger.Info("會員模組創建成功", logger.NewField("module", memberModule.Name()))

	// 初始化會員模組
	if err := memberModule.Setup(); err != nil {
		//log.Fatalf("初始化會員模組失敗: %v", err)
		a.Logger.Error("初始化會員模組失敗", logger.NewField("error", err))
	}
	//log.Printf("模組 %s 初始化成功", memberModule.Name())
	a.Logger.Info("模組初始化成功", logger.NewField("module", memberModule.Name()))

	// 啟動服務器
	addr := fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port)
	//fmt.Printf("Starting server on %s ...\n", addr)
	a.Logger.Info("啟動服務", logger.NewField("address", addr))
	if err := engine.Run(addr); err != nil {
		//log.Fatalf("啟動服務失敗: %v", err)
		a.Logger.Error("啟動服務失敗", logger.NewField("error", err))
	}
}
