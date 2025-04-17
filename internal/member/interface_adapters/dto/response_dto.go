package dto

type MemberResponseDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
type MemberListResponseDTO struct {
	Total int                 `json:"total"`
	Limit int                 `json:"limit"`
	Page  int                 `json:"page"`
	Data  []MemberResponseDTO `json:"data"`
}
