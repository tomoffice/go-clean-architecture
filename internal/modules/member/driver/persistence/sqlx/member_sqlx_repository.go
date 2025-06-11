package sqlx

import (
	"context"
	"module-clean/internal/shared/common/pagination"
)

type MemberSQLXRepository interface {
	Create(ctx context.Context, m *MemberSQLXModel) error
	GetByID(ctx context.Context, id int) (*MemberSQLXModel, error)
	GetByEmail(ctx context.Context, email string) (*MemberSQLXModel, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*MemberSQLXModel, error)
	CountAll(ctx context.Context) (int, error)
	UpdateProfile(ctx context.Context, m *MemberSQLXModel) (*MemberSQLXModel, error)
	UpdateEmail(ctx context.Context, id int, newEmail string) error
	UpdatePassword(ctx context.Context, id int, newPassword string) error
	Delete(ctx context.Context, id int) error
}
