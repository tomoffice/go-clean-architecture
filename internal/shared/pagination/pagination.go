package pagination

import (
	"module-clean/internal/shared/enum"
)

// Pagination 分頁結構
//   - Limit : 一次最多拿幾筆
//   - Offset : 從某一筆開始拿
//   - SortBy : 排序的欄位 id name email created_at
//   - OrderBy : 升冪或降冪 asc desc
//   - Total : 總筆數
type Pagination struct {
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
	SortBy  string       `json:"sort_by"`
	OrderBy enum.OrderBy `json:"order_by"`
	Total   int          `json:"total"`
}
