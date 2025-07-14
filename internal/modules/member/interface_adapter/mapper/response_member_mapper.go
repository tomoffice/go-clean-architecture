package mapper

import (
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dto"
	"time"
)

func EntityToCreateMemberResponseDTO(member *entity.Member) dto.RegisterMemberResponseDTO {
	return dto.RegisterMemberResponseDTO{
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
func EntityToUpdateMemberProfileResponseDTO(member *entity.Member) dto.UpdateMemberProfileResponseDTO {
	return dto.UpdateMemberProfileResponseDTO{
		ID:   member.ID,
		Name: member.Name,
	}
}
func EntityToUpdateMemberEmailResponseDTO() dto.UpdateMemberEmailResponseDTO {
	return dto.UpdateMemberEmailResponseDTO{}
}
func EntityToUpdateMemberPasswordResponseDTO() dto.UpdateMemberPasswordResponseDTO {
	return dto.UpdateMemberPasswordResponseDTO{}
}
func EntityToDeleteMemberResponseDTO(member *entity.Member) dto.DeleteMemberResponseDTO {
	return dto.DeleteMemberResponseDTO{
		ID:        member.ID,
		Name:      member.Name,
		Email:     member.Email,
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
}
