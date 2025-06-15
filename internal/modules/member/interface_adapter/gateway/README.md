## **目錄說明**

這層就是把 UseCase 跟 infra（DB、API 之類的）隔開。

所有資料轉換、錯誤 mapping 都在這裡集中處理。

UseCase 只認 interface，不會直接碰到底層 repository，這裡才知道怎麼跟 infra 打交道。

---

## **設計原則**

- 依賴反轉。UseCase 只依賴抽象，infra 換了不用動核心。
- 單一職責。只做 mapping，不寫業務邏輯。
- 禁止 infra 型別（sqlx、ent、api 回應）外流到 domain/entity。
- 錯誤要統一轉換，不能直接丟出 DB error。

---

## **命名規範**

- 具體 struct 跟檔名都加 Gateway 或 Adapter（例：MemberSQLXGateway）。
- interface（像 MemberRepository）在 usecase port/output。
- Gateway 本身會持有真正 infra repo 實例。

---

## **使用範例**

```
// UseCase 只會認這種 interface（定義在 usecase/port/output）
type MemberRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Member, error)
	Create(ctx context.Context, m *entity.Member) error
}

// Gateway 實作，這裡才知道底層是 SQLX、Ent 還是 API
type MemberSQLXGateway struct {
	repo sqlx.MemberSQLXRepository
}

func (g *MemberSQLXGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	m, err := g.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapSQLErrorToBusinessError(err)
	}
	// infra model 轉 domain entity
	return &entity.Member{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
	}, nil
}
```

---

## **避免事項**

- 不准把 infra 的型別（sqlx model、ent schema、第三方 response）直接往外丟。
- 必須 mapping，不能偷懶。
- 這層不負責任何商業邏輯或驗證。
- 錯誤一定要 mapping，不可以直接傳 infra error。

---