package dto

// GinBindingRegisterMemberRequestDTO (POST /api/v1/members)
type GinBindingRegisterMemberRequestDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GinBindingGetMemberByIDURIRequestDTO (GET /api/v1/members/:id)
type GinBindingGetMemberByIDURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}

// GinBindingGetMemberByEmailQueryRequestDTO (GET /api/v1/members/email/:email)
type GinBindingGetMemberByEmailQueryRequestDTO struct {
	Email string `form:"email" binding:"required"`
}

// GinBindingListMemberQueryRequestDTO (GET /api/v1/members?page=&limit=&sort_by=&order_by=)
type GinBindingListMemberQueryRequestDTO struct {
	Page    int    `form:"page" binding:"required"`
	Limit   int    `form:"limit" binding:"required"`
	SortBy  string `form:"sort_by" binding:"omitempty"`
	OrderBy string `form:"order_by" binding:"omitempty"`
}

// GinBindingUpdateMemberURIRequestDTO (PATCH /api/v1/members/:id)
type GinBindingUpdateMemberURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}

// GinBindingUpdateMemberProfileBodyRequestDTO (PATCH /api/v1/members/:id)
type GinBindingUpdateMemberProfileBodyRequestDTO struct {
	Name *string `json:"name,omitempty" binding:"omitempty"`
}

// GinBindingUpdateMemberEmailBodyRequestDTO (PATCH /api/v1/members/:id/email)
type GinBindingUpdateMemberEmailBodyRequestDTO struct {
	NewEmail string `json:"new_email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GinBindingUpdateMemberPasswordBodyRequestDTO (PATCH /api/v1/members/:id/password)
type GinBindingUpdateMemberPasswordBodyRequestDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// GinBindingDeleteMemberURIRequestDTO (DELETE /api/v1/members/:id)
type GinBindingDeleteMemberURIRequestDTO struct {
	ID int `uri:"id" binding:"required"`
}
