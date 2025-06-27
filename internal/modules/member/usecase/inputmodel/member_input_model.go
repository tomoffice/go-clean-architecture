// Package inputmodel 定義應用層（UseCase）專用的輸入資料模型。
// 負責將外部請求（如 HTTP URI、JSON、表單等）轉換為 UseCase 可用且語意明確的結構。
//
// 職責：
// 1. 明確描述 UseCase 需要的輸入資料（包含必填、選填、部分更新等需求）
// 2. 將外部資料來源（如 DTO、表單）與核心 Domain Entity 解耦
// 3. 支援 PATCH、局部更新等複雜場景，確保欄位變動意圖明確
// 4. 僅作為 UseCase 的 input，嚴禁混用於 Domain/Entity 層
package inputmodel

// PatchUpdateMemberProfileInputModel 為「更新會員資訊」UseCase 的輸入模型。
//   - 僅用於 UseCase 內部，不對外暴露。
//   - 支援 PATCH 部分欄位更新，欄位為 nil 表示不更新該欄位。
//   - 嚴禁直接對應 Domain Entity，需經 Mapper 轉換。
type PatchUpdateMemberProfileInputModel struct {
	ID   int
	Name *string
	// NickName *string
	// Avatar   *string
}
type PatchUpdateMemberPasswordInputModel struct {
	ID          int
	OldPassword string
	NewPassword string
}
