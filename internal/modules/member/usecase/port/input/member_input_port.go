package input

//go:generate mockgen -source=member_input_port.go -destination=../../../interface_adapter/controller/mock/mock_member_input_port.go -package=mock
import (
	"context"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/inputmodel"
	"module-clean/internal/shared/common/pagination"
)

type MemberInputPort interface {
	RegisterMember(ctx context.Context, member *entity.Member) (*entity.Member, error)
	GetMemberByID(ctx context.Context, id int) (*entity.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*entity.Member, error)
	ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, int, error)
	UpdateMember(ctx context.Context, patch *inputmodel.PatchUpdateMemberInputModel) (*entity.Member, error)
	DeleteMember(ctx context.Context, id int) (*entity.Member, error)
}
