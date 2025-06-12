package output

//go:generate mockgen -source=member_repository.go -destination=../../mock/mock_member_repository.go -package=mock
import (
	"context"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/shared/pagination"
)

type MemberRepository interface {
	Create(ctx context.Context, m *entity.Member) error
	GetByID(ctx context.Context, id int) (*entity.Member, error)
	GetByEmail(ctx context.Context, email string) (*entity.Member, error)
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, error)
	UpdateProfile(ctx context.Context, m *entity.Member) (*entity.Member, error)
	UpdateEmail(ctx context.Context, id int, newEmail string) error
	UpdatePassword(ctx context.Context, id int, newPassword string) error
	Delete(ctx context.Context, id int) error
	CountAll(ctx context.Context) (int, error)
}
