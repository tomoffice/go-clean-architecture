package main

import (
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"module-clean/internal/framework/database"
	"module-clean/internal/framework/gin"
	router2 "module-clean/internal/framework/gin/router"
	memberrepo "module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
	"module-clean/internal/modules/member/interface_adapter/controller"
	"module-clean/internal/modules/member/interface_adapter/router"
	"module-clean/internal/modules/member/usecase"
)

func main() {
	// Step 1: 初始化 DB
	db := database.InitSQLiteDB("./identifier.sqlite")
	// Step 2: 初始化 Repository（Infrastructure → Outbound Adapter）
	memberRepo := memberrepo.NewSQLXMemberRepo(db)

	// Step 3: 建立 UseCase（Application Layer）
	memberUseCase := usecase.NewMemberUseCase(memberRepo)

	// Step 4: 建立 Controller（Interface Adapter）
	memberController := controller.NewMemberController(memberUseCase)

	// Step 5: 建立 Router（Interface Adapter）
	memberRouter := router.NewMemberRouter(memberController)

	// Step 6: 初始化 Router（Gin）並綁定所有模組路由
	engine := router2.NewGinEngine("/api/v1", memberRouter.RegisterRoutes)

	// Step 7: 啟動 HTTP Server（或 gRPC / WebSocket 等其他協定）
	gin.StartHTTPServer(":81", engine)
}
