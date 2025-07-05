package member

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"module-clean/internal/modules"
	"module-clean/internal/modules/member/framework/persistence/sqlx/mcsqlite"
	"module-clean/internal/modules/member/interface_adapter/controller"
	"module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"module-clean/internal/modules/member/interface_adapter/presenter/http"
	"module-clean/internal/modules/member/interface_adapter/router"
	"module-clean/internal/modules/member/usecase"
)

// Factory 會員模組工廠
type Factory struct{}

// NewModuleFactory 創建會員模組工廠
func NewModuleFactory() modules.ModuleFactory {
	return &Factory{}
}

// CreateModule 創建會員模組
func (f *Factory) CreateModule(db *sqlx.DB, rg *gin.RouterGroup) (modules.Module, error) {
	// 組裝所有組件
	repo := mcsqlite.NewSQLXMemberRepo(db)
	gateway := repository.NewMemberSQLXGateway(repo)
	useCase := usecase.NewMemberUseCase(gateway)
	presenter := http.NewMemberPresenter()
	controller := controller.NewMemberController(useCase, presenter)
	router := router.NewMemberRouter(controller, rg)

	// 創建並返回模組實例
	return New(router), nil
}
