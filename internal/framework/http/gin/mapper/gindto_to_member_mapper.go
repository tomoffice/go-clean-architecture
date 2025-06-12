package mapper

import (
	gindto "module-clean/internal/framework/http/gin/dto"
	"module-clean/internal/modules/member/interface_adapter/dto"
)

func GinDTOToRegisterMemberDTO(ginDTO gindto.GinBindingRegisterMemberRequestDTO) dto.RegisterMemberRequestDTO {
	return dto.RegisterMemberRequestDTO{
		Name:     ginDTO.Name,
		Email:    ginDTO.Email,
		Password: ginDTO.Password,
	}
}
func GinDTOToGetMemberByIDDTO(ginDTO gindto.GinBindingGetMemberByIDURIRequestDTO) dto.GetMemberByIDRequestDTO {
	return dto.GetMemberByIDRequestDTO{
		ID: ginDTO.ID,
	}
}
func GinDTOToGetMemberByEmailDTO(ginDTO gindto.GinBindingGetMemberByEmailQueryRequestDTO) dto.GetMemberByEmailRequestDTO {
	return dto.GetMemberByEmailRequestDTO{
		Email: ginDTO.Email,
	}
}
func GinDTOtoListMemberDTO(ginDTO gindto.GinBindingListMemberQueryRequestDTO) dto.ListMemberRequestDTO {
	return dto.ListMemberRequestDTO{
		Page:    ginDTO.Page,
		Limit:   ginDTO.Limit,
		SortBy:  ginDTO.SortBy,
		OrderBy: ginDTO.OrderBy,
	}
}
func GinDTOToUpdateMemberProfileDTO(ginURI gindto.GinBindingUpdateMemberURIRequestDTO, ginBody gindto.GinBindingUpdateMemberProfileBodyRequestDTO) dto.UpdateMemberProfileRequestDTO {
	return dto.UpdateMemberProfileRequestDTO{
		ID:   ginURI.ID,
		Name: ginBody.Name,
	}
}
func GinDTOToUpdateMemberEmailDTO(ginURI gindto.GinBindingUpdateMemberURIRequestDTO, ginBody gindto.GinBindingUpdateMemberEmailBodyRequestDTO) dto.UpdateMemberEmailRequestDTO {
	return dto.UpdateMemberEmailRequestDTO{
		ID:       ginURI.ID,
		NewEmail: ginBody.NewEmail,
		Password: ginBody.Password,
	}
}
func GinDTOToUpdateMemberPasswordDTO(ginURI gindto.GinBindingUpdateMemberURIRequestDTO, ginBody gindto.GinBindingUpdateMemberPasswordBodyRequestDTO) dto.UpdateMemberPasswordRequestDTO {
	return dto.UpdateMemberPasswordRequestDTO{
		ID:          ginURI.ID,
		OldPassword: ginBody.OldPassword,
		NewPassword: ginBody.NewPassword,
	}
}
func GinDTOToDeleteMemberDTO(ginDTO gindto.GinBindingDeleteMemberURIRequestDTO) dto.DeleteMemberRequestDTO {
	return dto.DeleteMemberRequestDTO{
		ID: ginDTO.ID,
	}
}
