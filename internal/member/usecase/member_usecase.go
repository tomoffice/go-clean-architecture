// Package usecase 定義應用層業務邏輯（Application Business Rules），
// 用於實作實際的使用情境流程（Use Case），例如 CreateMember、ListMember 等。
//
// 職責:
// - 接收 input_model 作為輸入
// - 調用 domain service / repository 執行邏輯
// - 回傳 entity、output_model 或 error 結果
// - 不依賴外部框架（如 HTTP、DB）
package usecase

import (
	"context"
	"module-clean/internal/common/pagination"
	"module-clean/internal/member/domain/entities"
	"module-clean/internal/member/domain/repository"
)

type UseCase struct {
	MemberRepo repository.MemberRepository
}

func NewUseCase(memberRepo repository.MemberRepository) *UseCase {
	return &UseCase{
		MemberRepo: memberRepo,
	}
}
func (m *UseCase) RegisterMember(ctx context.Context, member *entities.Member) error {
	return m.MemberRepo.Create(ctx, member)
}
func (m *UseCase) GetMemberByID(ctx context.Context, id int) (*entities.Member, error) {
	return m.MemberRepo.GetByID(ctx, id)
}
func (m *UseCase) GetMemberByEmail(ctx context.Context, email string) (*entities.Member, error) {
	return m.MemberRepo.GetByEmail(ctx, email)
}
func (m *UseCase) ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, error) {
	return m.MemberRepo.GetAll(ctx, pagination)
}
func (m *UseCase) UpdateMember(ctx context.Context, patch *PatchUpdateMemberInput) (*entities.Member, error) {
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
	return m.MemberRepo.Update(ctx, member)
}
func (m *UseCase) DeleteMember(ctx context.Context, id int) (*entities.Member, error) {
	member, err := m.MemberRepo.GetByID(ctx, id)
	if err != nil {

	}
	err = m.MemberRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return member, nil
}
