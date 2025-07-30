package member

import "github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/router"

// Module 會員模組 - 具體產品
type Module struct {
	router *router.MemberRouter
}

// NewModule 創建會員模組實例
func NewModule(router *router.MemberRouter) *Module {
	return &Module{
		router: router,
	}
}

// Name 實現 Module 接口
func (m *Module) Name() string {
	return "member"
}

// Setup 實現 Module 接口
func (m *Module) Setup() error {
	// 純粹的委派，不做任何組裝邏輯
	return m.router.Register()
}

// Shutdown 實現 Module 接口
func (m *Module) Shutdown() error {
	// 未來可能需要的清理邏輯
	return nil
}
