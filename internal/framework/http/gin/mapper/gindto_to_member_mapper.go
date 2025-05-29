package mapper

import (
	gindto "module-clean/internal/framework/http/gin/dto"
	"module-clean/internal/modules/member/interface_adapter/dto"
)

func GinDTOToCreateMemberDTO(ginDTO gindto.GinCreateMemberRequestDTO) dto.CreateMemberRequestDTO {
	return dto.CreateMemberRequestDTO{
		Name:     ginDTO.Name,
		Email:    ginDTO.Email,
		Password: ginDTO.Password,
	}
}
func GinDTOToGetMemberByIDDTO(ginDTO gindto.GinGetMemberByIDRequestDTO) dto.GetMemberByIDRequestDTO {
	return dto.GetMemberByIDRequestDTO{
		ID: ginDTO.ID,
	}
}
func GinDTOToGetMemberByEmailDTO(ginDTO gindto.GinGetMemberByEmailRequestDTO) dto.GetMemberByEmailRequestDTO {
	return dto.GetMemberByEmailRequestDTO{
		Email: ginDTO.Email,
	}
}
func GinDTOtoListMemberDTO(ginDTO gindto.GinListMemberRequestDTO) dto.ListMemberRequestDTO {
	return dto.ListMemberRequestDTO{
		Page:    ginDTO.Page,
		Limit:   ginDTO.Limit,
		SortBy:  ginDTO.SortBy,
		OrderBy: ginDTO.OrderBy,
	}
}
func GinDTOToUpdateMemberDTO(ginDTO gindto.GinUpdateMemberRequestDTO) dto.UpdateMemberRequestDTO {
	return dto.UpdateMemberRequestDTO{
		ID:       ginDTO.ID,
		Name:     ginDTO.Name,
		Email:    ginDTO.Email,
		Password: ginDTO.Password,
	}
}
func GinDTOToDeleteMemberDTO(ginDTO gindto.GinDeleteMemberRequestDTO) dto.DeleteMemberRequestDTO {
	return dto.DeleteMemberRequestDTO{
		ID: ginDTO.ID,
	}
}
