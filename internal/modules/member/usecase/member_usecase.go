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
	"errors"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/usecase/inputmodel"
	"module-clean/internal/modules/member/usecase/port/input"
	"module-clean/internal/modules/member/usecase/port/output"
	"module-clean/internal/shared/pagination"
)

type MemberUseCase struct {
	MemberGateway output.MemberPersistence
}

func NewMemberUseCase(memberRepo output.MemberPersistence) input.MemberInputPort {
	return &MemberUseCase{
		MemberGateway: memberRepo,
	}
}
func (m *MemberUseCase) RegisterMember(ctx context.Context, member *entity.Member) (*entity.Member, error) {
	err := m.MemberGateway.Create(ctx, member)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	// 為了通用 repository，無論底層是否會 mutate 傳入 entity，
	// 一律透過唯一欄位查詢回傳完整 entity，減少 infra 依賴。
	retrieveMember, err := m.MemberGateway.GetByEmail(ctx, member.Email)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	return retrieveMember, nil
}
func (m *MemberUseCase) GetMemberByID(ctx context.Context, id int) (*entity.Member, error) {
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	return member, nil
}
func (m *MemberUseCase) GetMemberByEmail(ctx context.Context, email string) (*entity.Member, error) {
	member, err := m.MemberGateway.GetByEmail(ctx, email)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	return member, nil
}
func (m *MemberUseCase) ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, int, error) {
	members, err := m.MemberGateway.GetAll(ctx, pagination)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, 0, ucErr
	}
	total, err := m.MemberGateway.CountAll(ctx)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, 0, ucErr
	}
	return members, total, nil
}
func (m *MemberUseCase) UpdateMemberProfile(ctx context.Context, patch *inputmodel.PatchUpdateMemberProfileInputModel) (*entity.Member, error) {
	member, err := m.MemberGateway.GetByID(ctx, patch.ID)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	if patch.Name != nil {
		member.Name = *patch.Name
	}
	member, err = m.MemberGateway.UpdateProfile(ctx, member)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	return member, nil
}
func (m *MemberUseCase) UpdateMemberEmail(ctx context.Context, id int, newEmail, password string) error {
	// 先檢查新 email 是否被其他人使用
	existedMember, err := m.MemberGateway.GetByEmail(ctx, newEmail)
	if err == nil && existedMember.ID != id {
		// 新 email 已被其他人使用
		return ErrUseCaseMemberEmailAlreadyExists
	} else if err == nil && existedMember.ID == id {
		// 新 email 與舊 email 相同，直接返回錯誤
		return ErrUseCaseMemberUpdateSameEmail
	}
	if err != nil {
		err = MapGatewayErrorToUseCaseError(err)
		if errors.Is(err, ErrUseCaseMemberNotFound) {
			// 正常情境：新 email 不存在，可以繼續檢查密碼與更新
			// 不 return，繼續往下走
		} else {
			// 異常情境：DB 或其它技術錯誤
			return err
		}
	}
	// 驗證密碼與更新 email
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		return MapGatewayErrorToUseCaseError(err)
	}
	// 確認密碼是否正確
	if member.Password != password {
		return ErrUseCaseMemberPasswordIncorrect
	}
	// 執行 email 更新
	err = m.MemberGateway.UpdateEmail(ctx, id, newEmail)
	if err != nil {
		return MapGatewayErrorToUseCaseError(err)
	}
	return nil
}
func (m *MemberUseCase) UpdateMemberPassword(ctx context.Context, id int, newPassword, oldPassword string) error {
	// 先檢查新舊密碼是否一樣
	if newPassword == oldPassword {
		return ErrUseCaseMemberUpdateSamePassword
	}
	// 取得目前 member，驗證密碼
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		return MapGatewayErrorToUseCaseError(err)
	}
	// 確認舊密碼是否正確
	if member.Password != oldPassword {
		return ErrUseCaseMemberPasswordIncorrect
	}
	// 執行密碼更新
	err = m.MemberGateway.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		return MapGatewayErrorToUseCaseError(err)
	}
	return nil
}
func (m *MemberUseCase) DeleteMember(ctx context.Context, id int) (*entity.Member, error) {
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}

	err = m.MemberGateway.Delete(ctx, id)
	if err != nil {
		ucErr := MapGatewayErrorToUseCaseError(err)
		return nil, ucErr
	}
	return member, nil
}
