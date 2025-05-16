package output_port

import (
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/dto"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
)

type MemberOutputPort interface {
	PresentCreateMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO]
	PresentGetMemberByID(member *entity.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO]
	PresentGetMemberByEmail(member *entity.Member) sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO]
	PresentListMembers(members []*entity.Member, total int) sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO]
	PresentUpdateMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO]
	PresentDeleteMember(member *entity.Member) sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO]

	// UseCase 層業務錯誤處理
	PresentUseCaseError(err error) (int, sharedviewmodel.HTTPResponse[any])
	// Input 驗證錯誤處理
	PresentValidationError(err error) (int, sharedviewmodel.HTTPResponse[any])
	// Binding error 處理
	PresentBindingError(err error) (int, sharedviewmodel.HTTPResponse[any])
}
