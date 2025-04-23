package repository

import (
	"context"
	"module-clean/internal/common/pagination"
	"module-clean/internal/member/domain/entities"
)

type MemberRepository interface {
	Create(ctx context.Context, m *entities.Member) error
	GetByID(ctx context.Context, id int) (*entities.Member, error)
	GetByEmail(ctx context.Context, email string) (*entities.Member, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, error)
	Update(ctx context.Context, m *entities.Member) (*entities.Member, error)
	Delete(ctx context.Context, id int) error
}
