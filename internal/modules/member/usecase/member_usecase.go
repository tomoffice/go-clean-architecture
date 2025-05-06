// Package usecase 定義應用層業務邏輯（Application Business Rules），
// 用於實作實際的使用情境流程（Use Case），例如 CreateMember、ListMember 等。
//
// 職責:
// - 接收 input_model 作為輸入
// - 調用 domain service / repository 執行邏輯
// - 回傳 entity、output_model 或 errordefs 結果
// - 不依賴外部框架（如 HTTP、DB）
package usecase

import (
	"context"
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/modules/member/domain/repository"
	"module-clean/internal/shared/pagination"
)

type MemberUseCase struct {
	MemberRepo repository.MemberRepository
}

func NewMemberUseCase(memberRepo repository.MemberRepository) *MemberUseCase {
	return &MemberUseCase{
		MemberRepo: memberRepo,
	}
}
func (m *MemberUseCase) RegisterMember(ctx context.Context, member *entities.Member) error {
	err := m.MemberRepo.Create(ctx, member)
	if err != nil {
		return MapInfraErrorToUseCaseError(err)
	}
	return nil
}
func (m *MemberUseCase) GetMemberByID(ctx context.Context, id int) (*entities.Member, error) {
	member, err := m.MemberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}
	return member, nil
}
func (m *MemberUseCase) GetMemberByEmail(ctx context.Context, email string) (*entities.Member, error) {
	member, err := m.MemberRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}
	return member, nil
}
func (m *MemberUseCase) ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, int, error) {
	members, err := m.MemberRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, 0, MapInfraErrorToUseCaseError(err)
	}
	total, err := m.MemberRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, MapInfraErrorToUseCaseError(err)
	}
	return members, total, nil
}
func (m *MemberUseCase) UpdateMember(ctx context.Context, patch *PatchUpdateMemberInput) (*entities.Member, error) {
	member, err := m.MemberRepo.GetByID(ctx, patch.ID)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		member.Name = *patch.Name
	}
	if patch.Email != nil {
		member.Email = *patch.Email
	}
	if patch.Password != nil {
		member.Password = *patch.Password
	}
	member, err = m.MemberRepo.Update(ctx, member)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}
	return member, nil

}
func (m *MemberUseCase) DeleteMember(ctx context.Context, id int) (*entities.Member, error) {
	member, err := m.MemberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}

	err = m.MemberRepo.Delete(ctx, id)
	if err != nil {
		return nil, MapInfraErrorToUseCaseError(err)
	}
	return member, nil
}
