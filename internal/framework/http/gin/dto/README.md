# **internal/framework/http/gin/dto**

## **目錄說明**

本目錄**僅用於放置 Gin 專用的 Request/Response DTO**，其唯一用途如下：

- 僅供 Gin Controller 進行參數綁定（binding tag）、資料驗證、JSON 解析（json tag）、URI/Query 綁定（uri、form、query tag）等。
- **不包含任何業務邏輯，不參與 UseCase、Entity 等核心資料流。**

---

## **設計原則**

- **僅允許在 Controller 層或 Gin 專用的 Mapper 被引用。**
- **嚴禁在 UseCase、Entity、Domain、Repository、Presenter 等核心層級出現。**
- 如需傳遞資料至 UseCase，必須先轉換為 internal/modules/**/interface_adapter/dto 目錄下的業務 DTO。

---

## **命名規範**

- 統一採用 Gin 前綴，明確型別及來源。
    - 範例：
        - GinCreateMemberRequestDTO
        - GinUpdateMemberRequestDTO
        - GinListMemberRequestDTO
        - GinGetMemberByIDURIRequestDTO
        - GinGetMemberByEmailQueryRequestDTO
- **請明確標註型別來源**（如 URI / Query / JSON），避免混淆。

---

## **使用範例**

### **1. Controller 層綁定與驗證**

```
// 綁定 JSON
var req gindto.GinCreateMemberRequestDTO
if err := ctx.ShouldBindJSON(&req); err != nil {
    // 處理綁定錯誤
}

// 綁定 URI
var uriReq gindto.GinGetMemberByIDURIRequestDTO
if err := ctx.ShouldBindUri(&uriReq); err != nil {
    // 處理綁定錯誤
}

// 綁定 Query
var queryReq gindto.GinGetMemberByEmailQueryRequestDTO
if err := ctx.ShouldBindQuery(&queryReq); err != nil {
    // 處理綁定錯誤
}
```

### **2. DTO 轉換**

```
// 轉換為業務 DTO
businessReq := ginmapper.GinCreateMemberRequestDTOToCreateMemberDTO(req)
```

---

## **避免事項**

- 請勿將 Gin 專用 DTO 直接傳入 UseCase、Entity 或核心業務流程。
- 請勿將任何業務 DTO 放入本目錄（如 interface_adapter/dto 之 DTO 不可放這裡）。

---

## **參考範例**

```
// GinCreateMemberRequestDTO - 綁定 JSON
type GinCreateMemberRequestDTO struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// GinGetMemberByIDURIRequestDTO - 綁定 URI
type GinGetMemberByIDURIRequestDTO struct {
    ID int `uri:"id" binding:"required,min=1"`
}

// GinGetMemberByEmailQueryRequestDTO - 綁定 Query
type GinGetMemberByEmailQueryRequestDTO struct {
    Email string `form:"email" binding:"required,email"`
}
```