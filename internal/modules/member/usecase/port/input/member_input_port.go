package input

//go:generate mockgen -source=member_input_port.go -destination=../../../interface_adapter/controller/mock/mock_member_input_port.go -package=mock
import (
	"context"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/inputmodel"
	"module-clean/internal/shared/pagination"
)

type MemberInputPort interface {
	RegisterMember(ctx context.Context, member *entity.Member) (*entity.Member, error)
	GetMemberByID(ctx context.Context, id int) (*entity.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*entity.Member, error)
	ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, int, error)
	UpdateMemberProfile(ctx context.Context, patch *inputmodel.PatchUpdateMemberProfileInputModel) (*entity.Member, error)
	UpdateMemberEmail(ctx context.Context, id int, newEmail, password string) error
	UpdateMemberPassword(ctx context.Context, id int, oldPassword, newPassword string) error
	DeleteMember(ctx context.Context, id int) (*entity.Member, error)
}
