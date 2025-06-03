package mapper

import (
	gindto "module-clean/internal/framework/http/gin/dto"
	"module-clean/internal/modules/member/interface_adapter/dto"
)

func GinDTOToRegisterMemberDTO(ginDTO gindto.GinRegisterMemberRequestDTO) dto.RegisterMemberRequestDTO {
	return dto.RegisterMemberRequestDTO{
		Name:     ginDTO.Name,
		Email:    ginDTO.Email,
		Password: ginDTO.Password,
	}
}
func GinDTOToGetMemberByIDDTO(ginDTO gindto.GinGetMemberByIDURIRequestDTO) dto.GetMemberByIDRequestDTO {
	return dto.GetMemberByIDRequestDTO{
		ID: ginDTO.ID,
	}
}
func GinDTOToGetMemberByEmailDTO(ginDTO gindto.GinGetMemberByEmailQueryRequestDTO) dto.GetMemberByEmailRequestDTO {
	return dto.GetMemberByEmailRequestDTO{
		Email: ginDTO.Email,
	}
}
func GinDTOtoListMemberDTO(ginDTO gindto.GinListMemberQueryRequestDTO) dto.ListMemberRequestDTO {
	return dto.ListMemberRequestDTO{
		Page:    ginDTO.Page,
		Limit:   ginDTO.Limit,
		SortBy:  ginDTO.SortBy,
		OrderBy: ginDTO.OrderBy,
	}
}
func GinDTOToUpdateMemberDTO(ginURI gindto.GinUpdateMemberURIRequestDTO, ginBody gindto.GinUpdateMemberBodyRequestDTO) dto.UpdateMemberRequestDTO {
	return dto.UpdateMemberRequestDTO{
		ID:       ginURI.ID,
		Name:     ginBody.Name,
		Email:    ginBody.Email,
		Password: ginBody.Password,
	}
}
func GinDTOToDeleteMemberDTO(ginDTO gindto.GinDeleteMemberURIRequestDTO) dto.DeleteMemberRequestDTO {
	return dto.DeleteMemberRequestDTO{
		ID: ginDTO.ID,
	}
}
