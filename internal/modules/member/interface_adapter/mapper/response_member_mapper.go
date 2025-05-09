package mapper

import (
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"time"
)

func EntityToCreateMemberResponseDTO(member *entities.Member) dto.CreateMemberResponseDTO {
	return dto.CreateMemberResponseDTO{
		ID:    member.ID,
		Name:  member.Name,
		Email: member.Email,
	}
}
func EntityToGetMemberByIDResponseDTO(member *entities.Member) dto.GetMemberByIDResponseDTO {
	return dto.GetMemberByIDResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func EntityToGetMemberByEmailResponseDTO(member *entities.Member) dto.GetMemberByEmailResponseDTO {
	return dto.GetMemberByEmailResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func EntityToListMemberResponseDTO(members []*entities.Member) dto.ListMemberResponseDTO {
	items := make([]dto.ListMemberItemDTO, len(members))
	for i, m := range members {
		items[i] = dto.ListMemberItemDTO{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		}
	}
	return dto.ListMemberResponseDTO{
		Members: items,
	}
}
func EntityToUpdateMemberResponseDTO(member *entities.Member) dto.UpdateMemberResponseDTO {
	return dto.UpdateMemberResponseDTO{
		ID:    member.ID,
		Name:  &member.Name,
		Email: &member.Email,
	}
}
func EntityToDeleteMemberResponseDTO(member *entities.Member) dto.DeleteMemberResponseDTO {
	return dto.DeleteMemberResponseDTO{
		ID:    member.ID,
		Name:  &member.Name,
		Email: &member.Email,
	}
}
