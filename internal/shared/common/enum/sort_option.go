package enum

type OrderBy string

const (
	OrderByAsc  OrderBy = "asc"
	OrderByDesc OrderBy = "desc"
)

func (s OrderBy) IsValid() bool {
	if s == OrderByAsc || s == OrderByDesc {
		return true
	}
	return false
}
