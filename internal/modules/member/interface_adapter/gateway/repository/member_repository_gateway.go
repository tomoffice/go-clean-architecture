package repository

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/output"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"time"
)

type MemberSQLXGateway struct {
	infraRepo sqlx.MemberSQLXRepository
}

func NewMemberSQLXGateway(infraRepo sqlx.MemberSQLXRepository) output.MemberPersistence {
	return MemberSQLXGateway{
		infraRepo: infraRepo,
	}
}
func (g MemberSQLXGateway) Create(ctx context.Context, m *entity.Member) error {
	repoModel := &sqlx.MemberSQLXModel{
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	err := g.infraRepo.Create(ctx, repoModel)
	if err != nil {
		return MapInfraErrorToUsecaseError(err)
	}
	return nil
}

func (g MemberSQLXGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	repoModel, err := g.infraRepo.GetByID(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	return &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}, nil
}

func (g MemberSQLXGateway) GetByEmail(ctx context.Context, email string) (*entity.Member, error) {
	repoModel, err := g.infraRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	return &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}, nil
}

func (g MemberSQLXGateway) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, error) {
	repoModels, err := g.infraRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	var members []*entity.Member
	for _, repoModel := range repoModels {
		createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
		if err != nil {
			return nil, MapInfraErrorToUsecaseError(err)
		}
		members = append(members, &entity.Member{
			ID:        repoModel.ID,
			Name:      repoModel.Name,
			Email:     repoModel.Email,
			Password:  "",
			CreatedAt: createdAt,
		})
	}
	return members, nil
}

func (g MemberSQLXGateway) UpdateProfile(ctx context.Context, m *entity.Member) (*entity.Member, error) {
	repoModel := &sqlx.MemberSQLXModel{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	//讀已寫所以不用回傳
	_, err := g.infraRepo.UpdateProfile(ctx, repoModel)
	if err != nil {
		return nil, MapInfraErrorToUsecaseError(err)
	}
	return m, nil
}

func (g MemberSQLXGateway) UpdateEmail(ctx context.Context, id int, newEmail string) error {
	err := g.infraRepo.UpdateEmail(ctx, id, newEmail)
	if err != nil {
		return MapInfraErrorToUsecaseError(err)
	}
	return nil
}

func (g MemberSQLXGateway) UpdatePassword(ctx context.Context, id int, newPassword string) error {
	err := g.infraRepo.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		return MapInfraErrorToUsecaseError(err)
	}
	return nil
}

func (g MemberSQLXGateway) Delete(ctx context.Context, id int) error {
	err := g.infraRepo.Delete(ctx, id)
	if err != nil {
		return MapInfraErrorToUsecaseError(err)
	}
	return nil
}

func (g MemberSQLXGateway) CountAll(ctx context.Context) (int, error) {
	count, err := g.infraRepo.CountAll(ctx)
	if err != nil {
		return 0, MapInfraErrorToUsecaseError(err)
	}
	return count, nil
}
