package input_port

import (
	"context"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/shared/common/pagination"
)

type MemberInputPort interface {
	RegisterMember(ctx context.Context, member *entity.Member) (*entity.Member, error)
	GetMemberByID(ctx context.Context, id int) (*entity.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*entity.Member, error)
	ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, int, error)
	UpdateMember(ctx context.Context, patch *PatchUpdateMemberInput) (*entity.Member, error)
	DeleteMember(ctx context.Context, id int) (*entity.Member, error)
}
