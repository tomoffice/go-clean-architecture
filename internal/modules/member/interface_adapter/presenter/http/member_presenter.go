package http

import (
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/mapper"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/outputmodel"
	sharedenum "github.com/tomoffice/go-clean-architecture/internal/shared/enum"
	sharedviewmodel "github.com/tomoffice/go-clean-architecture/internal/shared/viewmodel/http"
	"strconv"
)

type MemberPresenter struct{}

func NewMemberPresenter() *MemberPresenter {
	return &MemberPresenter{}
}

func (p *MemberPresenter) PresentRegisterMember(member *entity.Member) outputmodel.RegisterMemberResponse {
	respDTO := mapper.EntityToCreateMemberResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentGetMemberByID(member *entity.Member) outputmodel.GetMemberByIDResponse {
	respDTO := mapper.EntityToGetMemberByIDResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentGetMemberByEmail(member *entity.Member) outputmodel.GetMemberByEmailResponse {
	respDTO := mapper.EntityToGetMemberByEmailResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentListMembers(members []*entity.Member, total int) outputmodel.ListMemberResponse {
	respDTO := mapper.EntityToListMemberResponseDTO(members)
	meta := &sharedviewmodel.MetaPayload{
		Total: total,
	}
	return buildSuccessResponseWithMeta(respDTO, meta)
}

func (p *MemberPresenter) PresentUpdateMemberProfile(member *entity.Member) outputmodel.UpdateMemberProfileResponse {
	respDTO := mapper.EntityToUpdateMemberProfileResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentUpdateMemberEmail() outputmodel.UpdateMemberEmailResponse {
	respDTO := mapper.EntityToUpdateMemberEmailResponseDTO()
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentUpdateMemberPassword() outputmodel.UpdateMemberPasswordResponse {
	respDTO := mapper.EntityToUpdateMemberPasswordResponseDTO()
	return buildSuccessResponse(respDTO)

}

func (p *MemberPresenter) PresentDeleteMember(member *entity.Member) outputmodel.DeleteMemberResponse {
	respDTO := mapper.EntityToDeleteMemberResponseDTO(member)
	return buildSuccessResponse(respDTO)
}

func (p *MemberPresenter) PresentBindingError(errCode int, message string) outputmodel.ErrorResponse {
	return buildFailedResponse(errCode, message)
}

func (p *MemberPresenter) PresentValidationError(err error) (int, outputmodel.ErrorResponse) {
	errCode, message := MapMemberValidationError(err)
	return errCode, buildFailedResponse(errCode, message)
}

func (p *MemberPresenter) PresentUseCaseError(err error) (int, outputmodel.ErrorResponse) {
	errCode, message := MapMemberUseCaseToPresenterError(err)
	return errCode, buildFailedResponse(errCode, message)
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
func buildFailedResponse(code int, message string) outputmodel.ErrorResponse {
	return outputmodel.ErrorResponse{
		Error: &sharedviewmodel.ErrorPayload{
			Code:    strconv.Itoa(code),
			Message: message,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusFailed),
	}
}
