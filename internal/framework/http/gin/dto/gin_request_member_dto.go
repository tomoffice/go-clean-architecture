package dto

// GinRegisterMemberRequestDTO (POST /api/v1/members)
type GinRegisterMemberRequestDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GinGetMemberByIDURIRequestDTO (GET /api/v1/members/:id)
type GinGetMemberByIDURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}

// GinGetMemberByEmailQueryRequestDTO (GET /api/v1/members/email/:email)
type GinGetMemberByEmailQueryRequestDTO struct {
	Email string `form:"email" binding:"required"`
}

// GinListMemberQueryRequestDTO (GET /api/v1/members?page=&limit=&sort_by=&order_by=)
type GinListMemberQueryRequestDTO struct {
	Page    int    `form:"page" binding:"required"`
	Limit   int    `form:"limit" binding:"required"`
	SortBy  string `form:"sort_by" binding:"omitempty"`
	OrderBy string `form:"order_by" binding:"omitempty"`
}

// GinUpdateMemberURIRequestDTO (PATCH /api/v1/members/:id)
type GinUpdateMemberURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}
type GinUpdateMemberBodyRequestDTO struct {
	Name     *string `json:"name,omitempty" binding:"omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty"`
	Password *string `json:"password,omitempty" binding:"omitempty"`
}

// GinDeleteMemberURIRequestDTO (DELETE /api/v1/members/:id)
type GinDeleteMemberURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}
