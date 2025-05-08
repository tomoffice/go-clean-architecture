package usecase

import (
	"context"
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/shared/common/pagination"
)

type MemberInputPort interface {
	RegisterMember(ctx context.Context, member *entities.Member) (*entities.Member, error)
	GetMemberByID(ctx context.Context, id int) (*entities.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*entities.Member, error)
	ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, int, error)
	UpdateMember(ctx context.Context, patch *PatchUpdateMemberInput) (*entities.Member, error)
	DeleteMember(ctx context.Context, id int) (*entities.Member, error)
}
