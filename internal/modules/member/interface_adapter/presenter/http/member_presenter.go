package http

import (
	"module-clean/internal/framework/http/gin/errordefs"
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/modules/member/interface_adapter/dto"
	sharedenum "module-clean/internal/shared/common/enum"
	"module-clean/internal/shared/common/errorcode"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
	"strconv"
	"time"
)

type MemberPresenter struct{}

func NewMemberPresenter() *MemberPresenter {
	return &MemberPresenter{}
}

func (p *MemberPresenter) PresentCreateMember(member *entities.Member) sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO] {
	return sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO]{
		Data: dto.CreateMemberResponseDTO{
			ID:    member.ID,
			Name:  member.Name,
			Email: member.Email,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

func (p *MemberPresenter) PresentGetMemberByID(member *entities.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO] {
	return sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO]{
		Data: dto.GetMemberByIDResponseDTO{
			ID:        member.ID,
			Name:      member.Name,
			Email:     member.Email,
			CreatedAt: member.CreatedAt.Format(time.RFC3339),
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

func (p *MemberPresenter) PresentGetMemberByEmail(member *entities.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO] {
	return sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO]{
		Data: dto.GetMemberByEmailResponseDTO{
			ID:        member.ID,
			Name:      member.Name,
			Email:     member.Email,
			CreatedAt: member.CreatedAt.Format(time.RFC3339),
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

func (p *MemberPresenter) PresentListMembers(members []*entities.Member, total int) sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO] {
	items := make([]dto.MemberListItemDTO, len(members))
	for i, m := range members {
		items[i] = dto.MemberListItemDTO{
			ID:    m.ID,
			Name:  m.Name,
			Email: m.Email,
		}
	}
	return sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO]{
		Data: dto.ListMemberResponseDTO{
			Members: items,
		},
		Meta: &sharedviewmodel.MetaPayload{
			Total: total,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

func (p *MemberPresenter) PresentUpdateMember(member *entities.Member) sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO] {
	return sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO]{
		Data: dto.UpdateMemberResponseDTO{
			ID:    member.ID,
			Name:  &member.Name,
			Email: &member.Email,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

func (p *MemberPresenter) PresentDeleteMember(member *entities.Member) sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO] {
	return sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO]{
		Data: dto.DeleteMemberResponseDTO{
			ID:    member.ID,
			Name:  &member.Name,
			Email: &member.Email,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusSuccess),
	}
}

//	func (p *MemberPresenter) PresentError(err error) (int, sharedviewmodel.HTTPResponse[any]) {
//		errCode, message := MapMemberUseCaseError(err)
//		httpStatus := MapErrorCodeToHTTPStatus(errCode)
//		return httpStatus, sharedviewmodel.HTTPResponse[any]{
//			Error: &sharedviewmodel.ErrorPayload{
//				Code:    strconv.Itoa(errCode),
//				Message: message,
//			},
//			BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusFailed),
//		}
//	}
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
	return httpStatus, sharedviewmodel.HTTPResponse[any]{
		Error: &sharedviewmodel.ErrorPayload{
			Code:    strconv.Itoa(errCode),
			Message: message,
		},
		BaseHTTPResponse: sharedviewmodel.NewBaseHTTPResponse(sharedenum.APIStatusFailed),
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
