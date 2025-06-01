package dto

type GinRegisterMemberRequestDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GinGetMemberByIDRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}

type GinGetMemberByEmailRequestDTO struct {
	Email string `form:"email" binding:"required"`
}

type GinListMemberRequestDTO struct {
	Page    int    `form:"page" binding:"required"`
	Limit   int    `form:"limit" binding:"required"`
	SortBy  string `form:"sort_by" binding:"omitempty"`
	OrderBy string `form:"order_by" binding:"omitempty"`
}

type GinUpdateMemberRequestDTO struct {
	ID       int     `json:"id" binding:"required"`
	Name     *string `json:"name,omitempty" binding:"omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty"`
	Password *string `json:"password,omitempty" binding:"omitempty"`
}

type GinDeleteMemberRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}
