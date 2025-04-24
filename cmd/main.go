package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // or mysql, pgx, etc.
	"log"
	"module-clean/internal/framework/gin"

	"module-clean/internal/framework/server"
	memberrepo "module-clean/internal/member/infrastructure/persistence/sqlx"
	"module-clean/internal/member/interface_adapters/controller"
	"module-clean/internal/member/usecase"
)

func main() {
	// Step 1: 初始化 DB（可抽離至 infrastructure/init 資源載入）
	db, err := sqlx.Open("sqlite3", "./identifier.sqlite")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Step 2: 初始化 Repository（Infrastructure → Outbound Adapter）
	memberRepo := memberrepo.NewSQLXMemberRepo(db)

	// Step 3: 建立 UseCase（Application Layer）
	memberUseCase := usecase.NewMemberUseCase(memberRepo)

	// Step 4: 建立 Controller（Interface Adapter）
	memberController := controller.NewMemberController(memberUseCase)

	// Step 5: 初始化 Router（Gin）並綁定所有模組路由
	engine := gin.InitGinRouter(memberController)

	// Step 6: 啟動 HTTP Server（或 gRPC / WebSocket 等其他協定）
	server.StartHTTPServer(":81", engine)
}
