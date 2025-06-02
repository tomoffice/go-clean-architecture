package main

import (
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"module-clean/internal/framework/database"
	"module-clean/internal/framework/http/gin"
	ginrouter "module-clean/internal/framework/http/gin/router"
	memberrepo "module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
	"module-clean/internal/modules/member/interface_adapter/controller"
	membergateway "module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"module-clean/internal/modules/member/interface_adapter/presenter/http"
	"module-clean/internal/modules/member/interface_adapter/router"
	"module-clean/internal/modules/member/usecase"
)

func main() {
	// 初始化 DB
	db := database.InitSQLiteDB("./identifier.sqlite")

	// 初始化 Repository（Infrastructure → Outbound Adapter）
	memberRepo := memberrepo.NewSQLXMemberRepo(db)

	// 建立 Gateway（Domain → Inbound Adapter）
	memberGateway := membergateway.NewMemberSQLXGateway(memberRepo)

	// 建立 UseCase（Application Layer）
	memberUseCase := usecase.NewMemberUseCase(memberGateway)

	// 建立 Presenter（Interface Adapter）
	memberPresenter := http.NewMemberPresenter()

	// 建立 Controller（Interface Adapter）
	memberController := controller.NewMemberController(memberUseCase, memberPresenter)

	// 建立 Router（Interface Adapter）
	memberRouter := router.NewMemberRouter(memberController)

	// 初始化 Router（Gin）並綁定所有模組路由
	engine := ginrouter.NewGinEngine("/api/v1", memberRouter.RegisterRoutes)

	// 啟動 HTTP Server（或 gRPC / WebSocket 等其他協定）
	gin.StartHTTPServer(":81", engine)
}
