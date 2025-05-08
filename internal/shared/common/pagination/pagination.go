package pagination

import (
	"module-clean/internal/shared/common/enum"
)

type Pagination struct {
	Limit   int          `json:"limit"`    // 一次最多拿幾筆
	Offset  int          `json:"offset"`   // 從某一筆開始拿
	SortBy  string       `json:"sort_by"`  // 排序的欄位
	OrderBy enum.OrderBy `json:"order_by"` // 升冪或降冪
	Total   int          `json:"total"`    // 總筆數
}
