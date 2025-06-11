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
func GinDTOToUpdateMemberProfileDTO(ginURI gindto.GinUpdateMemberURIRequestDTO, ginBody gindto.GinUpdateMemberProfileBodyRequestDTO) dto.UpdateMemberProfileRequestDTO {
	return dto.UpdateMemberProfileRequestDTO{
		ID:   ginURI.ID,
		Name: ginBody.Name,
	}
}
func GinDTOToUpdateMemberEmailDTO(ginURI gindto.GinUpdateMemberURIRequestDTO, ginBody gindto.GinUpdateMemberEmailBodyRequestDTO) dto.UpdateMemberEmailRequestDTO {
	return dto.UpdateMemberEmailRequestDTO{
		ID:       ginURI.ID,
		NewEmail: ginBody.NewEmail,
		Password: ginBody.Password,
	}
}
func GinDTOToUpdateMemberPasswordDTO(ginURI gindto.GinUpdateMemberURIRequestDTO, ginBody gindto.GinUpdateMemberPasswordBodyRequestDTO) dto.UpdateMemberPasswordRequestDTO {
	return dto.UpdateMemberPasswordRequestDTO{
		ID:          ginURI.ID,
		OldPassword: ginBody.OldPassword,
		NewPassword: ginBody.NewPassword,
	}
}
func GinDTOToDeleteMemberDTO(ginDTO gindto.GinDeleteMemberURIRequestDTO) dto.DeleteMemberRequestDTO {
	return dto.DeleteMemberRequestDTO{
		ID: ginDTO.ID,
	}
}
