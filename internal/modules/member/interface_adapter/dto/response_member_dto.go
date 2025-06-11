package dto

type RegisterMemberResponseDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type GetMemberByIDResponseDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
type GetMemberByEmailResponseDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
type ListMemberItemDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type ListMemberResponseDTO struct {
	Members []ListMemberItemDTO `json:"members"`
}
type UpdateMemberProfileResponseDTO struct {
	ID   int     `json:"id"`
	Name string `json:"name,omitempty"`
}
type UpdateMemberEmailResponseDTO struct{}
type UpdateMemberPasswordResponseDTO struct{}
type DeleteMemberResponseDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
