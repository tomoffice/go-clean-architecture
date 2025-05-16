package persistence

import (
	"context"
	"module-clean/internal/shared/common/pagination"
)

type MemberRepository interface {
	Create(ctx context.Context, m *MemberRepoModel) error
	GetByID(ctx context.Context, id int) (*MemberRepoModel, error)
	GetByEmail(ctx context.Context, email string) (*MemberRepoModel, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*MemberRepoModel, error)
	CountAll(ctx context.Context) (int, error)
	Update(ctx context.Context, m *MemberRepoModel) (*MemberRepoModel, error)
	Delete(ctx context.Context, id int) error
}
