package repository

import (
	"context"
	"module-clean/internal/modules/member/driver/persistence"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/shared/common/pagination"
	"time"
)

type MemberRepositoryGateway struct {
	infraRepo persistence.MemberRepository
}

func NewMemberRepositoryGateway(infraRepo persistence.MemberRepository) MemberRepositoryGateway {
	return MemberRepositoryGateway{
		infraRepo: infraRepo,
	}
}
func (m2 MemberRepositoryGateway) Create(ctx context.Context, m *entity.Member) error {
	repoModel := &persistence.MemberRepoModel{
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	err := m2.infraRepo.Create(ctx, repoModel)
	if err != nil {
		return MapInfraErrorToGatewayError(err)
	}
	return nil
}

func (m2 MemberRepositoryGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	repoModel, err := m2.infraRepo.GetByID(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToGatewayError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		return nil, ErrGatewayMemberMappingFailed
	}
	return &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}, nil
}

func (m2 MemberRepositoryGateway) GetByEmail(ctx context.Context, email string) (*entity.Member, error) {
	repoModel, err := m2.infraRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, MapInfraErrorToGatewayError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		return nil, ErrGatewayMemberMappingFailed
	}
	return &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}, nil
}

func (m2 MemberRepositoryGateway) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, error) {
	repoModels, err := m2.infraRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, MapInfraErrorToGatewayError(err)
	}
	var members []*entity.Member
	for _, repoModel := range repoModels {
		createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
		if err != nil {
			return nil, ErrGatewayMemberMappingFailed
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

func (m2 MemberRepositoryGateway) Update(ctx context.Context, m *entity.Member) (*entity.Member, error) {
	repoModel := &persistence.MemberRepoModel{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	//讀已寫所以不用回傳
	_, err := m2.infraRepo.Update(ctx, repoModel)
	if err != nil {
		return nil, MapInfraErrorToGatewayError(err)
	}
	return m, nil
}

func (m2 MemberRepositoryGateway) Delete(ctx context.Context, id int) error {
	err := m2.infraRepo.Delete(ctx, id)
	if err != nil {
		return MapInfraErrorToGatewayError(err)
	}
	return nil
}

func (m2 MemberRepositoryGateway) CountAll(ctx context.Context) (int, error) {
	count, err := m2.infraRepo.CountAll(ctx)
	if err != nil {
		return 0, MapInfraErrorToGatewayError(err)
	}
	return count, nil
}
