package bootstrap

import (
	"fmt"
	"log"
	"module-clean/config"
	"module-clean/internal/framework/database"
	"module-clean/internal/framework/http/gin"
	ginrouter "module-clean/internal/framework/http/gin/router"
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
	engine := ginrouter.NewGinEngine("")
	apiRouterGroup := engine.Group("/api/v1")
	// member module
	memberModule := member.NewMemberModule(db, apiRouterGroup)
	a.Modules = append(a.Modules, memberModule)
	a.registerALlModules(a.Modules)
	addr := fmt.Sprintf("%s:%s", a.Config.Server.HTTP.Host, a.Config.Server.HTTP.Port)
	fmt.Printf("Starting server on %s ...\n", addr)
	gin.StartHTTPServer(addr, engine)
}
func (a *App) registerALlModules(modules []modules.Module) {
	for _, m := range modules {
		m.Assemble()
	}
}
