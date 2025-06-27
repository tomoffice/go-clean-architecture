package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"module-clean/config"
	"module-clean/internal/framework/database"

	"module-clean/internal/modules"
	"module-clean/internal/modules/member"
)

type App struct {
	Config  *config.Config
	Modules []modules.Module
}

func NewApp(cfg *config.Config) *App {
	return &App{
		Config:  cfg,
		Modules: []modules.Module{},
	}
}

func (a *App) Run() {
	db, err := database.InitSQLiteDB(a.Config.Database.DSN)
	if err != nil {
		log.Fatalf("DB 初始化失敗: %v", err)
	}
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	apiRouterGroup := engine.Group("/api/v1")
	// module
	memberModule := member.NewMemberModule(db, apiRouterGroup)
	a.Modules = append(a.Modules, memberModule)
	a.registerALlModules(a.Modules)

	addr := fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port)
	fmt.Printf("Starting server on %s ...\n", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("啟動服務失敗: %v", err)
	}
}
func (a *App) registerALlModules(modules []modules.Module) {
	for _, m := range modules {
		m.Assemble()
	}
}
