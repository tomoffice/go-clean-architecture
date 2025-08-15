# internal/framework/http/gin/dto

## 用途與職責

Gin 框架專用 DTO 目錄，負責處理 HTTP 請求的資料綁定和驗證。作為 Framework 層的一部分，專門處理 Gin 特有的 Request/Response 資料結構，包含 JSON、URI、Query 參數的綁定邏輯。

## 設計原則

- **框架隔離**: 僅供 Controller 層和 Gin Mapper 使用
- **邊界清晰**: 不參與業務邏輯，用完即轉換為業務 DTO
- **職責單一**: 專注於 HTTP 資料綁定和基礎驗證
- **依賴方向**: 絕不流入 UseCase、Entity、Repository 等內層

## 目錄結構

```
internal/framework/http/gin/dto/
├── member_dto.go               # Member 相關 Gin DTO
├── common_dto.go              # 共用 DTO 結構
└── README.md                  # 本說明檔
```

## 使用方式

### Controller 資料綁定

```go
// JSON 綁定
var req gindto.GinCreateMemberRequestDTO
if err := ctx.ShouldBindJSON(&req); err != nil {
    return c.presenter.ValidationError(ctx, err)
}

// URI 參數綁定
var uriReq gindto.GinGetMemberByIDURIRequestDTO
if err := ctx.ShouldBindUri(&uriReq); err != nil {
    return c.presenter.ValidationError(ctx, err)
}

// Query 參數綁定
var queryReq gindto.GinGetMemberByEmailQueryRequestDTO
if err := ctx.ShouldBindQuery(&queryReq); err != nil {
    return c.presenter.ValidationError(ctx, err)
}
```

### DTO 轉換

```go
// 轉換為業務 DTO 後傳入 UseCase
businessReq := ginmapper.GinCreateMemberRequestDTOToCreateMemberDTO(req)
result, err := c.usecase.CreateMember(ctx, businessReq)
```

### DTO 定義範例

```go
// GinCreateMemberRequestDTO - JSON 請求
type GinCreateMemberRequestDTO struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// GinGetMemberByIDURIRequestDTO - URI 參數
type GinGetMemberByIDURIRequestDTO struct {
    ID int `uri:"id" binding:"required,min=1"`
}
```

## 注意事項

- **嚴禁直接傳入 UseCase 層**，必須先轉換為業務 DTO
- **命名必須加 Gin 前綴**，並標明資料來源（JSON/URI/Query）
- **僅處理資料綁定和基礎驗證**，不包含業務邏輯驗證
- **不可在業務層引用**，保持 Framework 層邊界清晰