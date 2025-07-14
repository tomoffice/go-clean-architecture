package member

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/tomoffice/go-clean-architecture/internal/modules"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx/mcsqlite"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/controller"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/gateway/repository"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/presenter/http"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/router"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase"
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
