package mapper

import (
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"time"
)

func EntityToCreateMemberResponseDTO(member *entity.Member) dto.CreateMemberResponseDTO {
	return dto.CreateMemberResponseDTO{
		ID:    member.ID,
		Name:  member.Name,
		Email: member.Email,
	}
}
func EntityToGetMemberByIDResponseDTO(member *entity.Member) dto.GetMemberByIDResponseDTO {
	return dto.GetMemberByIDResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func EntityToGetMemberByEmailResponseDTO(member *entity.Member) dto.GetMemberByEmailResponseDTO {
	return dto.GetMemberByEmailResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
func EntityToListMemberResponseDTO(members []*entity.Member) dto.ListMemberResponseDTO {
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
func EntityToUpdateMemberResponseDTO(member *entity.Member) dto.UpdateMemberResponseDTO {
	return dto.UpdateMemberResponseDTO{
		ID:    member.ID,
		Name:  &member.Name,
		Email: &member.Email,
	}
}
func EntityToDeleteMemberResponseDTO(member *entity.Member) dto.DeleteMemberResponseDTO {
	return dto.DeleteMemberResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
