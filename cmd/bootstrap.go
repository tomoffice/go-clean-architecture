package main

import (
	"fmt"
	"log"
	"module-clean/config"
	"module-clean/internal/framework/database"
	"module-clean/internal/framework/http/gin"
	ginrouter "module-clean/internal/framework/http/gin/router"
	"module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
	"module-clean/internal/modules/member/interface_adapter/controller"
	"module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"module-clean/internal/modules/member/interface_adapter/presenter/http"
	"module-clean/internal/modules/member/interface_adapter/router"
	"module-clean/internal/modules/member/usecase"

	"module-clean/internal/modules/member"
)

type App struct {
	Config *config.Config
}

func NewApp(cfg *config.Config) *App { return &App{Config: cfg} }

func (a *App) Run() {
	engine := ginrouter.NewGinEngine("")
	apiRouterGroup := engine.Group("/api/v1")

	// 組裝 infra
	db, err := database.InitSQLiteDB(a.Config.Database.DSN, "sqlite_init.sql", "sqlite_seed.sql")
	if err != nil {
		log.Fatalf("DB初始化失敗: %v", err)
	}

	// 組裝member模組
	memberInfraRepo := sqlite.NewSQLXMemberRepo(db)
	memberGateway := repository.NewMemberSQLXGateway(memberInfraRepo)
	uc := usecase.NewMemberUseCase(memberGateway)
	presenter := http.NewMemberPresenter()
	ctrl := controller.NewMemberController(uc, presenter)
	memberRouter := router.NewMemberRouter(ctrl)
	memberModule := member.NewModule(memberRouter)
	memberModule.RegisterRoutes(apiRouterGroup)

	fmt.Printf("Starting server on %s ...\n", a.Config.Server.HTTP.Port)
	gin.StartHTTPServer(fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port), engine)

}
