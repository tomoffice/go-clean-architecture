# Gateway 層

## 核心職責

Gateway 層是 `UseCase`（應用程式核心邏輯）與 `Infrastructure`（外部基礎設施，如資料庫、第三方 API）之間的橋樑。它的存在是為了**徹底隔離**這兩者。

它主要負責兩件事：

1.  **隔離 (Isolation):** `UseCase` 只依賴於 Gateway 提供的抽象介面（定義在 `usecase/port/output`），完全不知道底層是 SQL 資料庫、NoSQL 資料庫還是某個 RESTful API。
2.  **轉譯 (Translation):**
    *   **資料模型轉換 (Data Model Mapping):** 將從 Infrastructure 層獲取的資料模型（例如 `dao.MemberRecord`）轉換為 `UseCase` 能理解的領域實體（`entity.Member`）。
    *   **錯誤轉換 (Error Mapping):** 將 Infrastructure 層的專有錯誤（例如 `sql.ErrNoRows`）轉換為 `UseCase` 定義的、具有業務意義的錯誤。

## 設計原則

*   **依賴反轉原則 (Dependency Inversion):** `UseCase` 不依賴具體的 Gateway 實作，而是依賴抽象介面。這使得我們可以輕易地抽換底層實作（例如，將資料庫從 `SQLx` 換成 `GORM`），而無需修改任何核心業務邏輯。
*   **單一職責原則 (Single Responsibility):** Gateway 的職責僅限於「轉譯」。它不應該包含任何業務流程、資料驗證或流程控制的邏輯。

## 結構與命名

*   **Gateway 實作** (例如 `MemberRepoGateway`) 位於 `interface_adapter/gateway/repository` 或 `interface_adapter/gateway/api` 目錄下。
*   它所實作的**介面** (例如 `MemberPersistence`) 定義在 `usecase/port/output` 中。
*   Gateway struct 會依賴一個 `DAO` (Data Access Object) 介面或 API 客戶端，而不是直接依賴具體的資料庫驅動。

## 實作流程範例

1.  `UseCase` 呼叫 `MemberPersistence` 介面的 `GetByID` 方法。
2.  `MemberRepoGateway` 作為介面的實作，接收到這個呼叫。
3.  `MemberRepoGateway` 呼叫其持有的 `memberDAO` 的 `GetByID` 方法，與資料庫進行互動。
4.  `memberDAO` 回傳 `dao.MemberRecord` 和一個可能的 infrastructure error。
5.  `MemberRepoGateway` 透過 `MapInfraErrorToUsecaseError` 將 infrastructure error 轉換為 usecase error。
6.  `MemberRepoGateway` 將 `dao.MemberRecord` 轉換為 `entity.Member`。
7.  最後，將 `entity.Member` 和轉換後的 usecase error 回傳給 `UseCase`。

## 程式碼範例

**1. UseCase 定義的 Port (介面)**
```go
// --- 位於 usecase/port/output/member_port.go ---
package output

// MemberPersistence 定義了 UseCase 期望的持久化操作
type MemberPersistence interface {
	GetByID(ctx context.Context, id int) (*entity.Member, error)
	Create(ctx context.Context, m *entity.Member) error
    // ... 其他方法
}
```

**2. Gateway 的具體實作**
```go
// --- 位於 interface_adapter/gateway/repository/member_repository_gateway.go ---
package repository

// MemberRepoGateway 實作了 MemberPersistence 介面
type MemberRepoGateway struct {
	dao    dao.MemberDAO // 依賴 DAO 抽象，而非具體實作
	logger logger.Logger
	tracer tracer.Tracer
}

// GetByID 實作了資料獲取、轉譯的過程
func (g MemberRepoGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	// 呼叫 DAO 層來獲取資料庫紀錄
	record, err := g.dao.GetByID(ctx, id)
	if err != nil {
		// 將底層錯誤轉譯為 UseCase 能理解的錯誤
		return nil, MapInfraErrorToUsecaseError(err)
	}

	// 將 DAO model (dao.MemberRecord) 轉譯為 domain entity (entity.Member)
	member := &entity.Member{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		Password:  "", // 密碼等敏感資訊不應外傳
		CreatedAt: record.CreatedAt,
	}

	return member, nil
}
```

## 開發規範

*   **嚴禁洩漏 (No Leaks):** 絕不將 `DAO` 層或任何 `Infrastructure` 層的資料模型（如 `dao.MemberRecord`）或原始錯誤（如 `sql.ErrNoRows`）回傳給 `UseCase`。
*   **專注轉譯 (Focus on Translation):** 不要在 Gateway 中實作任何業務邏輯、資料驗證或流程控制。這些都屬於 `UseCase` 的職責。
*   **徹底轉換 (Complete Mapping):** 所有從 `DAO` 或 API 客戶端獲得的資料模型和潛在錯誤，都必須被完整地處理和轉譯。
