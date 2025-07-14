package output

//go:generate mockgen -source=member_presenter.go -destination=../../../interface_adapter/controller/mock/mock_member_presenter.go -package=mock
import (
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/outputmodel"
)

type MemberPresenter interface {
	PresentRegisterMember(member *entity.Member) outputmodel.RegisterMemberResponse
	PresentGetMemberByID(member *entity.Member) outputmodel.GetMemberByIDResponse
	PresentGetMemberByEmail(member *entity.Member) outputmodel.GetMemberByEmailResponse
	PresentListMembers(members []*entity.Member, total int) outputmodel.ListMemberResponse
	PresentUpdateMemberProfile(member *entity.Member) outputmodel.UpdateMemberProfileResponse
	PresentUpdateMemberEmail() outputmodel.UpdateMemberEmailResponse
	PresentUpdateMemberPassword() outputmodel.UpdateMemberPasswordResponse
	PresentDeleteMember(member *entity.Member) outputmodel.DeleteMemberResponse
	// PresentBindingError 處理輸入綁定錯誤
	PresentBindingError(errCode int, message string) outputmodel.ErrorResponse
	// PresentValidationError 處理驗證錯誤
	PresentValidationError(err error) (int, outputmodel.ErrorResponse)
	// PresentUseCaseError 處理用例錯誤
	PresentUseCaseError(err error) (int, outputmodel.ErrorResponse)
}
