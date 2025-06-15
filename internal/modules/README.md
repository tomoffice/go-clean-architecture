# **internal/modules**

## **目錄說明**

這個目錄專門用來放各個領域模組（如：member、order、product 等），每個模組底下都會有自己獨立的分層結構，遵循 Clean Architecture，包含 driver（infra/persistence）、entity、interface_adapter、usecase 等。

---

## **設計原則**

- **每個模組都要自成一格，彼此間低耦合。**
- **所有依賴（repository、gateway、usecase、presenter、controller、router）** 都在 cmd/bootstrap.go 組裝、注入，**Module 本身只做聚合，不做建構。**
- **模組必須實作 Module 介面**（目前只規範 RegisterRoutes，後續可擴充）。

---

## **命名規範**

- 模組資料夾直接用業務英文單字小寫（如 member、order）。
- 各層、元件、檔名依職責命名，不要搞 ambiguous 的縮寫。
- interface_adapter 下如有 gateway、dto、mapper、presenter、controller 等，務必清楚標示用途。

---

## **使用範例**

### **定義一個新模組**

1. 建立模組資料夾（如：internal/modules/product/）
2. 依照 Clean Architecture 切 driver、entity、interface_adapter、usecase 等子資料夾。
3. 新增 module.go，結構如下：

```GO
type Module struct {
    // 這裡聚合該模組所有依賴（如 repo、gateway、usecase、presenter、controller、router）
}

func NewModule() *Module {
    // 不做任何 new，所有依賴都由 bootstrap 傳進來
    return &Module{...}
}

func (m *Module) RegisterRoutes(routerGroup *gin.RouterGroup) {
    // 呼叫 router 內註冊路由
}
```

1. 在 cmd/bootstrap.go 完成該模組的依賴組裝與註冊：

```GO
otherInfraRepo := otherinfra.NewSQLXOtherRepo(db)
otherGateway := othergateway.NewOtherSQLXGateway(otherInfraRepo)
otherUc := usecase.NewOtherUseCase(otherGateway)
otherPresenter := http.NewOtherPresenter()
otherCtrl := controller.NewOtherController(otherUc, otherPresenter)
otherRouter := router.NewOtherRouter(otherCtrl)
otherModule := other.NewModule(apiRouterGroup, otherRouter)
otherModule.RegisterRoutes()
```

---

## **避免事項**

- **嚴禁** 在模組內自動組裝依賴（不允許 module.go 自行 new 各層）。
- **不要** 把多個領域混在同一模組下，每個業務就一個獨立 module。
- **不要** 在模組裡藏 DI 或初始化邏輯，所有東西都應該由 bootstrap 明確注入。
- **不要** 把 controller/router 等註冊行為散落在外，全部走 RegisterRoutes。

---