package persistence

import (
	"context"
	"module-clean/internal/shared/common/pagination"
)

type MemberRepository interface {
	Create(ctx context.Context, m *MemberSQLXModel) error
	GetByID(ctx context.Context, id int) (*MemberSQLXModel, error)
	GetByEmail(ctx context.Context, email string) (*MemberSQLXModel, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*MemberSQLXModel, error)
	CountAll(ctx context.Context) (int, error)
	Update(ctx context.Context, m *MemberSQLXModel) (*MemberSQLXModel, error)
	Delete(ctx context.Context, id int) error
}
