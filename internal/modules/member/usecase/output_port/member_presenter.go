package output_port

//go:generate mockgen -source=member_presenter.go -destination=../../interface_adapter/controller/mock/mock_member_presenter.go -package=mock
import (
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/dto"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
)

// mockgen 尚未支援泛型，所以用別名展開所有泛型返回類型
type CreateMemberResp = sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO]
type GetMemberByIDResp = sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO]
type GetMemberByEmailResp = sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO]
type ListMemberResp = sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO]
type UpdateMemberResp = sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO]
type DeleteMemberResp = sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO]

// 為 any 的情況也必須別名化
type AnyResp = sharedviewmodel.HTTPResponse[any]

type MemberPresenter interface {
	PresentCreateMember(member *entity.Member) CreateMemberResp
	PresentGetMemberByID(member *entity.Member) GetMemberByIDResp
	PresentGetMemberByEmail(member *entity.Member) GetMemberByEmailResp
	PresentListMembers(members []*entity.Member, total int) ListMemberResp
	PresentUpdateMember(member *entity.Member) UpdateMemberResp
	PresentDeleteMember(member *entity.Member) DeleteMemberResp

	// UseCase 層業務錯誤處理
	PresentUseCaseError(err error) (int, AnyResp)
	// Input 驗證錯誤處理
	PresentValidationError(err error) (int, AnyResp)
	// Binding error 處理
	PresentBindingError(err error) (int, AnyResp)
}
