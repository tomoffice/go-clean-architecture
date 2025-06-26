package member

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"module-clean/internal/modules"
	"module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
	"module-clean/internal/modules/member/interface_adapter/controller"
	"module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"module-clean/internal/modules/member/interface_adapter/presenter/http"
	"module-clean/internal/modules/member/interface_adapter/router"
	"module-clean/internal/modules/member/usecase"
)

type MemberModule struct {
	db *sqlx.DB
	rg *gin.RouterGroup
}

func NewMemberModule(db *sqlx.DB, rg *gin.RouterGroup) modules.Module {
	return &MemberModule{
		db: db,
		rg: rg,
	}
}

// Assemble 組裝 MemberModule 的所有組件
func (m *MemberModule) Assemble() {
	infraRepo := sqlite.NewSQLXMemberRepo(m.db)
	gateway := repository.NewMemberSQLXGateway(infraRepo)
	uc := usecase.NewMemberUseCase(gateway)
	presenter := http.NewMemberPresenter()
	ctrl := controller.NewMemberController(uc, presenter)
	memberRouter := router.NewMemberRouter(ctrl, m.rg)
	memberRouter.Register()
}
