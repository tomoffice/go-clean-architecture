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
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/inputmodel"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/input"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/output"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
)

type MemberUseCase struct {
	MemberGateway output.MemberPersistence
	logger        logger.Logger
	tracer        tracer.Tracer
}

func NewMemberUseCase(memberRepo output.MemberPersistence, logger logger.Logger, tracer tracer.Tracer) input.MemberInputPort {
	return &MemberUseCase{
		MemberGateway: memberRepo,
		logger:        logger,
		tracer:        tracer,
	}
}
func (m *MemberUseCase) RegisterMember(ctx context.Context, member *entity.Member) (*entity.Member, error) {
	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	// 創建 UseCase 層的子 span，自動記錄執行時間
	ctx, span := m.tracer.Start(ctx, "MemberController.GetByID")
	defer span.End()

	contextLogger.Info("開始會員註冊",
		logger.NewField("member_email", member.Email),
		logger.NewField("member_name", member.Name),
	)

	err := m.MemberGateway.Create(ctx, member)
	if err != nil {
		contextLogger.Error("會員註冊 Gateway 創建失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", member.Email),
		)
		return nil, err
	}
	// 為了通用 repository，無論底層是否會 mutate 傳入 entity，
	// 一律透過唯一欄位查詢回傳完整 entity，減少 infra 依賴。
	retrieveMember, err := m.MemberGateway.GetByEmail(ctx, member.Email)
	if err != nil {
		contextLogger.Error("會員註冊後查詢失敗",
			logger.NewField("error", err.Error()),
			logger.NewField("member_email", member.Email),
		)
		return nil, err
	}

	contextLogger.Info("會員註冊成功",
		logger.NewField("member_id", retrieveMember.ID),
		logger.NewField("member_email", retrieveMember.Email),
	)
	return retrieveMember, nil
}
func (m *MemberUseCase) GetMemberByID(ctx context.Context, id int) (*entity.Member, error) {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.GetMemberByID")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員查詢(ID)",
		logger.NewField("member_id", id),
	)

	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		contextLogger.Error("會員查詢(ID) Gateway 執行失敗",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", id),
		)
		return nil, err
	}

	contextLogger.Debug("會員查詢(ID)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
	)
	return member, nil
}
func (m *MemberUseCase) GetMemberByEmail(ctx context.Context, email string) (*entity.Member, error) {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.GetMemberByEmail")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員查詢(Email)",
		logger.NewField("member_email", email),
	)

	member, err := m.MemberGateway.GetByEmail(ctx, email)
	if err != nil {
		contextLogger.Error("會員查詢(Email) Gateway 執行失敗",
			logger.NewField("error", err.Error()),
			logger.NewField("member_email", email),
		)
		return nil, err
	}

	contextLogger.Debug("會員查詢(Email)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
	)
	return member, nil
}
func (m *MemberUseCase) ListMembers(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, int, error) {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.ListMembers")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員列表查詢",
		logger.NewField("limit", pagination.Limit),
		logger.NewField("offset", pagination.Offset),
	)

	members, err := m.MemberGateway.GetAll(ctx, pagination)
	if err != nil {
		contextLogger.Error("會員列表查詢 Gateway 執行失敗",
			logger.NewField("error", err),
			logger.NewField("limit", pagination.Limit),
			logger.NewField("offset", pagination.Offset),
		)
		return nil, 0, err
	}
	total, err := m.MemberGateway.CountAll(ctx)
	if err != nil {
		contextLogger.Error("會員總數查詢 Gateway 執行失敗",
			logger.NewField("error", err),
		)
		return nil, 0, err
	}

	contextLogger.Debug("會員列表查詢成功",
		logger.NewField("count", len(members)),
		logger.NewField("total", total),
	)
	return members, total, nil
}
func (m *MemberUseCase) UpdateMemberProfile(ctx context.Context, patch *inputmodel.PatchUpdateMemberProfileInputModel) (*entity.Member, error) {
	// 創庺 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.UpdateMemberProfile")
	defer span.End()

	// 創庺帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員資料更新",
		logger.NewField("member_id", patch.ID),
	)

	member, err := m.MemberGateway.GetByID(ctx, patch.ID)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		member.Name = *patch.Name
	}
	member, err = m.MemberGateway.UpdateProfile(ctx, member)
	if err != nil {
		contextLogger.Error("會員資料更新 Gateway 執行失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", patch.ID),
		)
		return nil, err
	}

	contextLogger.Debug("會員資料更新成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
	)
	return member, nil
}
func (m *MemberUseCase) UpdateMemberEmail(ctx context.Context, id int, newEmail, password string) error {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.UpdateMemberEmail")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員 Email 更新",
		logger.NewField("member_id", id),
		logger.NewField("new_email", newEmail),
	)

	// 先檢查新 email 是否被其他人使用
	existedMember, err := m.MemberGateway.GetByEmail(ctx, newEmail)
	if err == nil && existedMember.ID != id {
		// 新 email 已被其他人使用
		contextLogger.Error("會員 Email 更新失敗：新 Email 已被使用",
			logger.NewField("member_id", id),
			logger.NewField("new_email", newEmail),
			logger.NewField("existing_member_id", existedMember.ID),
		)
		return ErrMemberEmailAlreadyExists
	} else if err == nil && existedMember.ID == id {
		// 新 email 與舊 email 相同，直接返回錯誤
		contextLogger.Error("會員 Email 更新失敗：新舊 Email 相同",
			logger.NewField("member_id", id),
			logger.NewField("email", newEmail),
		)
		return ErrMemberUpdateSameEmail
	}
	if err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			// 正常情境：新 email 不存在，可以繼續檢查密碼與更新
			// 不 return，繼續往下走
		} else {
			// 異常情境：DB 或其它技術錯誤
			contextLogger.Error("會員 Email 更新檢查失敗",
				logger.NewField("error", err),
				logger.NewField("member_id", id),
				logger.NewField("new_email", newEmail),
			)
			return err
		}
	}
	// 驗證是否存在 member
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		contextLogger.Error("會員 Email 更新檢查會員失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}
	// 驗證密碼與更新 email
	// 確認密碼是否正確
	if member.Password != password {
		contextLogger.Error("會員 Email 更新失敗：密碼錯誤",
			logger.NewField("member_id", id),
			logger.NewField("member_email", member.Email),
		)
		return ErrMemberPasswordIncorrect
	}
	// 執行 email 更新
	err = m.MemberGateway.UpdateEmail(ctx, id, newEmail)
	if err != nil {
		contextLogger.Error("會員 Email 更新 Gateway 執行失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("new_email", newEmail),
		)
		return err
	}

	contextLogger.Debug("會員 Email 更新成功",
		logger.NewField("member_id", id),
		logger.NewField("old_email", member.Email),
		logger.NewField("new_email", newEmail),
	)
	return nil
}
func (m *MemberUseCase) UpdateMemberPassword(ctx context.Context, id int, newPassword, oldPassword string) error {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.UpdateMemberPassword")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員密碼更新",
		logger.NewField("member_id", id),
	)

	// 先檢查新舊密碼是否一樣
	if newPassword == oldPassword {
		contextLogger.Error("會員密碼更新失敗：新舊密碼相同",
			logger.NewField("member_id", id),
		)
		return ErrMemberUpdateSamePassword
	}
	// 取得目前 member，驗證密碼
	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		contextLogger.Error("會員密碼更新檢查會員失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}
	// 確認舊密碼是否正確
	if member.Password != oldPassword {
		contextLogger.Error("會員密碼更新失敗：舊密碼錯誤",
			logger.NewField("member_id", id),
			logger.NewField("member_email", member.Email),
		)
		return ErrMemberPasswordIncorrect
	}
	// 執行密碼更新
	err = m.MemberGateway.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		contextLogger.Error("會員密碼更新 Gateway 執行失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}

	contextLogger.Debug("會員密碼更新成功",
		logger.NewField("member_id", id),
		logger.NewField("member_email", member.Email),
	)
	return nil
}
func (m *MemberUseCase) DeleteMember(ctx context.Context, id int) (*entity.Member, error) {
	// 創建 UseCase 層的子 span
	ctx, span := m.tracer.Start(ctx, "MemberUseCase.DeleteMember")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := m.logger.WithContext(ctx).With(logger.NewField("layer", "usecase"))

	contextLogger.Debug("開始會員刪除",
		logger.NewField("member_id", id),
	)

	member, err := m.MemberGateway.GetByID(ctx, id)
	if err != nil {
		contextLogger.Error("會員刪除檢查會員失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return nil, err
	}

	err = m.MemberGateway.Delete(ctx, id)
	if err != nil {
		contextLogger.Error("會員刪除 Gateway 執行失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("member_email", member.Email),
		)
		return nil, err
	}

	contextLogger.Debug("會員刪除成功",
		logger.NewField("member_id", id),
		logger.NewField("member_email", member.Email),
	)
	return member, nil
}
