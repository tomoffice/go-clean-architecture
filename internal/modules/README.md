# internal/modules

## 用途與職責

業務模組目錄，專門存放各個領域模組（如 member、order、product 等）。每個模組遵循 Clean Architecture 四層分離原則，包含 entity、usecase、interface_adapter、framework 等分層結構，實現高內聚、低耦合的模組化架構。

## 設計原則

- **模組自治**: 每個模組自成一格，彼此間低耦合
- **依賴注入**: 所有依賴在 bootstrap 階段組裝注入，模組本身僅做聚合
- **介面規範**: 模組必須實作 Module 介面，統一路由註冊方式
- **Clean Architecture**: 嚴格遵循四層分離，確保依賴方向正確

## 目錄結構

```
internal/modules/
├── {module}/                    # 各業務模組目錄
│   ├── entity/                  # 領域實體層
│   ├── usecase/                 # 應用邏輯層
│   ├── interface_adapter/       # 介面適配器層
│   ├── framework/              # 框架驅動層
│   └── module.go               # 模組聚合定義
├── modules.go                   # 模組註冊管理
└── README.md                   # 本說明檔
```

## 使用方式

### 定義新模組

```go
// 1. 建立模組結構
type Module struct {
    router router.Router  // 聚合依賴，不自行建構
}

// 2. 建構函數（依賴由 bootstrap 注入）
func NewModule(router router.Router) *Module {
    return &Module{router: router}
}

// 3. 實作 Module 介面
func (m *Module) RegisterRoutes(routerGroup *gin.RouterGroup) {
    m.router.RegisterRoutes(routerGroup)
}
```

### Bootstrap 組裝

```go
// 在 cmd/bootstrap.go 完成依賴組裝
repo := persistence.NewRepository(db)
gateway := repository.NewGateway(repo)
usecase := usecase.NewUseCase(gateway)
presenter := presenter.NewPresenter()
controller := controller.NewController(usecase, presenter)
router := router.NewRouter(controller)
module := module.NewModule(router)
```

## 注意事項

- **嚴禁在模組內自動組裝依賴**，所有 new 操作都在 bootstrap 階段
- **一個業務領域一個模組**，不要將多個領域混合
- **所有路由註冊必須透過 RegisterRoutes 方法**，不可散落各處
- **依賴方向必須符合 Clean Architecture**，外層依賴內層