package usecase

import (
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/shared/interface_adapter/viewmodel/http"
)

type MemberOutputPort interface {
	PresentCreateMember(member *entities.Member) http.HTTPResponse[dto.CreateMemberResponseDTO]
	PresentGetMemberByID(member *entities.Member) http.HTTPResponse[dto.GetMemberByIDResponseDTO]
	PresentGetMemberByEmail(member *entities.Member) http.HTTPResponse[dto.GetMemberByEmailResponseDTO]
	PresentListMembers(members []*entities.Member, total int) http.HTTPResponse[dto.ListMemberResponseDTO]
	PresentUpdateMember(member *entities.Member) http.HTTPResponse[dto.UpdateMemberResponseDTO]
	PresentDeleteMember(member *entities.Member) http.HTTPResponse[dto.DeleteMemberResponseDTO]
	PresentError(err error) (int, http.HTTPResponse[any])
}
