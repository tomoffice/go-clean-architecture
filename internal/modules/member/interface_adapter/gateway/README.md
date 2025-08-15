# internal/modules/member/interface_adapter/gateway

## 用途與職責

Gateway 層作為 UseCase（應用邏輯層）與 Infrastructure（基礎設施層）之間的橋樑，負責徹底隔離核心業務邏輯與外部依賴。主要職責包括資料模型轉換、錯誤映射，以及協議轉譯，確保 UseCase 層完全不感知底層實現細節。

## 設計原則

- **依賴反轉**: UseCase 僅依賴抽象介面，不依賴具體 Gateway 實現
- **單一職責**: 專注於資料轉譯，不包含業務邏輯或驗證
- **完全隔離**: 絕不洩漏 Infrastructure 層的資料模型或錯誤到 UseCase
- **協議轉譯**: 負責不同層級間的資料格式轉換

## 目錄結構

```
internal/modules/member/interface_adapter/gateway/
├── repository/                          # 資料庫相關 Gateway
│   └── member_repository_gateway.go     # Member 資料庫 Gateway 實現
├── api/                                 # 第三方 API Gateway （如有需要）
└── README.md                           # 本說明檔
```

## 使用方式

### Gateway 實現

```go
// MemberRepoGateway 實作 UseCase 定義的介面
type MemberRepoGateway struct {
    dao    dao.MemberDAO  // 依賴 DAO 抽象
    logger logger.Logger
    tracer tracer.Tracer
}

// 實作 UseCase 期望的方法
func (g *MemberRepoGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
    // 1. 呼叫 DAO 層獲取資料
    record, err := g.dao.GetByID(ctx, id)
    if err != nil {
        // 2. 錯誤轉譯：Infrastructure -> UseCase
        return nil, MapInfraErrorToUsecaseError(err)
    }

    // 3. 資料轉譯：DAO Record -> Domain Entity
    member := &entity.Member{
        ID:        record.ID,
        Name:      record.Name,
        Email:     record.Email,
        CreatedAt: record.CreatedAt,
        // 注意：敏感資訊如密碼不應外傳
    }

    return member, nil
}
```

### UseCase Port 定義

```go
// 在 usecase/port/output/member_port.go 中定義
type MemberPersistence interface {
    GetByID(ctx context.Context, id int) (*entity.Member, error)
    Create(ctx context.Context, member *entity.Member) error
    Update(ctx context.Context, member *entity.Member) error
}
```

## 注意事項

- **嚴禁資料洩漏**: 絕不將 DAO 層資料模型或 Infrastructure 錯誤直接傳給 UseCase
- **專注轉譯職責**: 不可包含業務邏輯、資料驗證或流程控制
- **完整錯誤映射**: 所有 Infrastructure 錯誤都必須轉換為具業務意義的 UseCase 錯誤
- **敏感資訊處理**: 密碼等敏感資訊在轉譯過程中應適當處理或過濾
