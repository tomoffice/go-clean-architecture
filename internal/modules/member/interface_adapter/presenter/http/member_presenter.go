package http

import (
	"module-clean/internal/framework/http/gin/errordefs"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/interface_adapter/mapper"
	sharedenum "module-clean/internal/shared/common/enum"
	"module-clean/internal/shared/common/errorcode"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
	"strconv"
)

type MemberPresenter struct{}

func NewMemberPresenter() *MemberPresenter {
	return &MemberPresenter{}
}

func (p *MemberPresenter) PresentCreateMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO] {
	respDTO := mapper.EntityToCreateMemberResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentGetMemberByID(member *entity.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO] {
	respDTO := mapper.EntityToGetMemberByIDResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentGetMemberByEmail(member *entity.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO] {
	respDTO := mapper.EntityToGetMemberByEmailResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentListMembers(members []*entity.Member, total int) sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO] {
	respDTO := mapper.EntityToListMemberResponseDTO(members)
	meta := &sharedviewmodel.MetaPayload{
		Total: total,
	}
	return buildSuccessResponseWithMeta(respDTO, meta)
}

func (p *MemberPresenter) PresentUpdateMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO] {
	respDTO := mapper.EntityToUpdateMemberResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentDeleteMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO] {
	respDTO := mapper.EntityToDeleteMemberResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentBindingError(err error) (int, sharedviewmodel.HTTPResponse[any]) {
	errCode, message := errordefs.MapGinBindingError(err)
	if errCode == errorcode.ErrInternalServer {
		message = "Unclassified binding error"
	}
	httpStatus := MapErrorCodeToHTTPStatus(errCode)
	return httpStatus, buildFailedResponse(errCode, message)
}

func (p *MemberPresenter) PresentValidationError(err error) (int, sharedviewmodel.HTTPResponse[any]) {
	errCode := errorcode.ErrValidationFailed
	message := err.Error()
	httpStatus := MapErrorCodeToHTTPStatus(errCode)
	return httpStatus, buildFailedResponse(errCode, message)
}

func (p *MemberPresenter) PresentUseCaseError(err error) (int, sharedviewmodel.HTTPResponse[any]) {
	errCode, message := MapMemberUseCaseError(err)
	httpStatus := MapErrorCodeToHTTPStatus(errCode)
	return httpStatus, buildFailedResponse(errCode, message)
}

func buildSuccessResponse[T any](data T) sharedviewmodel.HTTPResponse[T] {
	return sharedviewmodel.HTTPResponse[T]{
		Data:             data,
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}
func buildSuccessResponseWithMeta[T any](data T, meta *sharedviewmodel.MetaPayload) sharedviewmodel.HTTPResponse[T] {
	return sharedviewmodel.HTTPResponse[T]{
		Data:             data,
		Meta:             meta,
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}
func buildFailedResponse(code int, message string) sharedviewmodel.HTTPResponse[any] {
	return sharedviewmodel.HTTPResponse[any]{
		Error: &sharedviewmodel.ErrorPayload{
			Code:    strconv.Itoa(code),
			Message: message,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusFailed),
	}
}
