package domain

import (
	"context"
	"module-clean/internal/common/pagination"
)

type MemberRepository interface {
	Create(ctx context.Context, m *Member) error
	GetByID(ctx context.Context, id int) (*Member, error)
	GetByEmail(ctx context.Context, email string) (*Member, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*Member, error)
	Update(ctx context.Context, m *Member) error
	Delete(ctx context.Context, id int) error
}
