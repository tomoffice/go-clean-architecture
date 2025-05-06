// Package usecase 定義應用層邏輯使用的輸入資料結構，
// 用於隔離外部輸入（如 DTO）與內部的業務邏輯，並避免直接耦合 domain entity。
// 特別是在需要處理「部分欄位更新」或「可選輸入」等情境時，
// 使用 input struct 能提供更具語意、符合意圖的資料格式。
//
// 職責:
// - 定義 UseCase 所需的輸入資料結構
// - 提供與外部資料格式（DTO）及內部 entity 之間的緩衝層
// - 表達 UseCase 的輸入意圖與可選欄位
package usecase

type PatchUpdateMemberInput struct {
	ID       int
	Name     *string
	Email    *string
	Password *string
}
