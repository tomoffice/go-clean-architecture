# internal/framework/http/gin/dto

## 目錄說明

本目錄專門放置 **Gin 專用的 Request/Response DTO**，其唯一用途為：

- 搭配 Gin 進行參數綁定（binding）、資料驗證（binding tag）、JSON 解析（json tag）、URI/Query 解析（uri、form、query tag）等。
- 不包含任何業務邏輯，也不參與 UseCase 層或 Entity 層資料流。

## 設計原則

- **只允許在 Controller 層或 Gin 相關的 Mapper 進行引用**。
- **嚴禁在 UseCase、Entity、Domain、Repository、Presenter 等核心層級出現**。
- 如需傳遞資料至 UseCase，請轉換為 `internal/modules/**/interface_adapter/dto` 目錄下的 DTO。

## 命名建議

- 統一加上 `Gin` 前綴，例如：
    - `GinCreateMemberRequestDTO`
    - `GinUpdateMemberRequestDTO`
    - `GinListMemberRequestDTO`

## 使用範例

1. 在 Controller 層進行綁定
    ```go
    var req gindto.GinCreateMemberRequestDTO
    if err := ctx.ShouldBindJSON(&req); err != nil {
        // ...
    }
    ```

2. 轉換為業務 DTO
    ```go
    businessReq := ginmapper.GinDTOToCreateMemberDTO(req)
    ```

## 避免事項

- 請勿將 gin 專用 DTO 直接傳入 UseCase 或 Entity。
- 請勿將業務 DTO 放入本目錄。

---