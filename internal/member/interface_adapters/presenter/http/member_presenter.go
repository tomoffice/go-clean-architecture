package http

import (
	"module-clean/internal/member/domain/entities"
	"module-clean/internal/member/interface_adapters/dto"
	"time"
)

func PresentCreateMemberDTO(member *entities.Member) dto.CreateMemberResponseDTO {
	return dto.CreateMemberResponseDTO{
		ID:    member.ID,
		Name:  member.Name,
		Email: member.Email,
	}
}
func PresentGetMemberByIDDTO(member *entities.Member) dto.GetMemberByIDResponseDTO {
	return dto.GetMemberByIDResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func PresentGetMemberByEmailDTO(member *entities.Member) dto.GetMemberByEmailResponseDTO {
	return dto.GetMemberByEmailResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func PresentListMemberDTO(members []*entities.Member) dto.ListMemberResponseDTO {
	items := make([]dto.MemberListItemDTO, len(members))
	for i, m := range members {
		items[i] = dto.MemberListItemDTO{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		}
	}
	return dto.ListMemberResponseDTO{Members: items}
}
func PresentUpdateMemberDTO(member *entities.Member) dto.UpdateMemberResponseDTO {
	return dto.UpdateMemberResponseDTO{
		ID:    member.ID,
		Name:  &member.Name,
		Email: &member.Email,
	}
}
func PresentDeleteMemberDTO(member *entities.Member) dto.DeleteMemberResponseDTO {
	return dto.DeleteMemberResponseDTO{
		ID:    member.ID,
		Name:  &member.Name,
		Email: &member.Email,
	}
}
