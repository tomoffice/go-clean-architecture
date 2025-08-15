package dao

//go:generate mockgen -source=member_dao.go -destination=../../interface_adapter/gateway/mock/mock_member_dao.go -package=mock

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"time"
)

type MemberRecord struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type MemberDAO interface {
	Create(ctx context.Context, m *MemberRecord) error
	GetByID(ctx context.Context, id int) (*MemberRecord, error)
	GetByEmail(ctx context.Context, email string) (*MemberRecord, error)
	GetAll(ctx context.Context, p pagination.Pagination) ([]*MemberRecord, error)
	UpdateProfile(ctx context.Context, m *MemberRecord) (*MemberRecord, error)
	UpdateEmail(ctx context.Context, id int, newEmail string) error
	UpdatePassword(ctx context.Context, id int, newPassword string) error
	Delete(ctx context.Context, id int) error
	CountAll(ctx context.Context) (int, error)
}
