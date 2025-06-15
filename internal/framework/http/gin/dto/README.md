# **internal/framework/http/gin/dto**

## **目錄說明**

這層只放 Gin 專用的 Request/Response DTO。

用途就兩個：

- Controller 拿來 binding（binding tag）、驗證、JSON/URI/Query 解析。
- 這邊的 DTO 不做業務、不參與核心流程，用完就轉業務 DTO。

---

## **設計原則**

- 只能 Controller 或 Gin 的 Mapper 用到，其他地方都別碰。
- 禁止用在 UseCase、Entity、Domain、Repository、Presenter 等任何業務/核心層。
- 要給 UseCase，自己轉到 internal/modules/**/interface_adapter/dto 下的 DTO，再丟進去。

---

## **命名規範**

- 統一加 Gin 前綴，看到就知道是哪層來的。
    - 範例：
        - GinCreateMemberRequestDTO
        - GinUpdateMemberRequestDTO
        - GinListMemberRequestDTO
        - GinGetMemberByIDURIRequestDTO
        - GinGetMemberByEmailQueryRequestDTO
- 來源（URI/Query/JSON）一定要標明，別搞混。

---

## **使用範例**

### **1. Controller 綁定/驗證**

```
// 綁 JSON
var req gindto.GinCreateMemberRequestDTO
if err := ctx.ShouldBindJSON(&req); err != nil {
    // 處理綁定錯誤
}

// 綁 URI
var uriReq gindto.GinGetMemberByIDURIRequestDTO
if err := ctx.ShouldBindUri(&uriReq); err != nil {
    // 處理綁定錯誤
}

// 綁 Query
var queryReq gindto.GinGetMemberByEmailQueryRequestDTO
if err := ctx.ShouldBindQuery(&queryReq); err != nil {
    // 處理綁定錯誤
}
```

### **2. DTO 轉換給 UseCase**

```
// 一律轉成業務用 DTO
businessReq := ginmapper.GinCreateMemberRequestDTOToCreateMemberDTO(req)
```

---

## **避免事項**

- Gin DTO 不准直接塞進 UseCase、Entity 或其他業務流程。
- 業務 DTO 也不要亂丟進來（interface_adapter/dto 的東西別放這）。

---

## **參考範例**

```
// GinCreateMemberRequestDTO - 綁 JSON
type GinCreateMemberRequestDTO struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// GinGetMemberByIDURIRequestDTO - 綁 URI
type GinGetMemberByIDURIRequestDTO struct {
    ID int `uri:"id" binding:"required,min=1"`
}

// GinGetMemberByEmailQueryRequestDTO - 綁 Query
type GinGetMemberByEmailQueryRequestDTO struct {
    Email string `form:"email" binding:"required,email"`
}
```