// Package input_port 定義應用層（UseCase）所需的輸入模型，
// 負責將外部資料（例如 HTTP 請求參數、JSON payload）轉換成對 UseCase 友好的結構。
// 它能：
//   1. 清楚表達 UseCase 的輸入意圖與範圍（必填、選填、預設值等）
//   2. 隔離外部傳入資料與核心 Domain Entity，避免直接依賴 Entity 結構
//   3. 支援部分欄位更新（Patch）等複雜場景，確保欄位修改的明確性
package input_port

// PatchUpdateMemberInputModel 表示「更新會員」UseCase 的輸入參數。
//   - ID      ：要更新的會員識別碼，必填
//   - Name    ：會員名稱，選填（nil 表示不變）
//   - Email   ：會員信箱，選填（nil 表示不變）
//   - Password：會員密碼，選填（nil 表示不變）
type PatchUpdateMemberInputModel struct {
	ID       int
	Name     *string
	Email    *string
	Password *string
}