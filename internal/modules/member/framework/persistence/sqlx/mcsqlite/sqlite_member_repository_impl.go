package mcsqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	sqlx2 "github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"time"
)

type sqlxMemberRepo struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer tracer.Tracer
}

func NewSQLXMemberRepo(db *sqlx.DB, logger logger.Logger, tracer tracer.Tracer) sqlx2.MemberSQLXRepository {
	return &sqlxMemberRepo{
		db:     db,
		logger: logger,
		tracer: tracer,
	}
}
func (s sqlxMemberRepo) Create(ctx context.Context, m *sqlx2.MemberSQLXModel) error {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.Create")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 插入",
		logger.NewField("query", "INSERT INTO members"),
		logger.NewField("member_email", m.Email),
	)

	_, err := s.db.ExecContext(ctx, queryInsertMember, m.Name, m.Email, m.Password)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 插入失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", m.Email),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return mapSQLError(err)
	}

	contextLogger.Debug("SQL 插入成功",
		logger.NewField("member_email", m.Email),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)

	//id, err := result.LastInsertId()
	//if err != nil {
	//	return err
	//}
	//m.ID = int(id)
	return nil
}
func (s sqlxMemberRepo) GetByID(ctx context.Context, id int) (*sqlx2.MemberSQLXModel, error) {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.GetByID")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 查詢(ID)",
		logger.NewField("query", "SELECT * FROM members WHERE id = ?"),
		logger.NewField("member_id", id),
	)

	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(ctx, member, querySelectByID, id)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 查詢(ID)失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}

	contextLogger.Debug("SQL 查詢(ID)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return member, nil
}
func (s sqlxMemberRepo) GetByEmail(ctx context.Context, email string) (*sqlx2.MemberSQLXModel, error) {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.GetByEmail")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 查詢",
		logger.NewField("query", "SELECT * FROM members WHERE email = ?"),
		logger.NewField("member_email", email),
	)

	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(ctx, member, querySelectByEmail, email)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 查詢失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", email),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}

	contextLogger.Debug("SQL 查詢成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", email),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return member, nil
}
func (s sqlxMemberRepo) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*sqlx2.MemberSQLXModel, error) {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.GetAll")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	query := fmt.Sprintf(querySelectAllBase, pagination.SortBy, pagination.OrderBy)
	contextLogger.Debug("開始執行 SQL 列表查詢",
		logger.NewField("query", "SELECT * FROM members ORDER BY ... LIMIT ? OFFSET ?"),
		logger.NewField("limit", pagination.Limit),
		logger.NewField("offset", pagination.Offset),
	)

	members := make([]*sqlx2.MemberSQLXModel, 0)
	err := s.db.SelectContext(ctx, &members, query, pagination.Limit, pagination.Offset)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 列表查詢失敗",
			logger.NewField("error", err),
			logger.NewField("limit", pagination.Limit),
			logger.NewField("offset", pagination.Offset),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}

	contextLogger.Debug("SQL 列表查詢成功",
		logger.NewField("count", len(members)),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return members, nil
}
func (s sqlxMemberRepo) CountAll(ctx context.Context) (int, error) {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.CountAll")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 總數查詢",
		logger.NewField("query", "SELECT COUNT(*) FROM members"),
	)

	var count int
	err := s.db.GetContext(ctx, &count, queryCountMembers)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 總數查詢失敗",
			logger.NewField("error", err),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return 0, mapSQLError(err)
	}

	contextLogger.Debug("SQL 總數查詢成功",
		logger.NewField("count", count),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return count, nil
}
func (s sqlxMemberRepo) UpdateProfile(ctx context.Context, m *sqlx2.MemberSQLXModel) (*sqlx2.MemberSQLXModel, error) {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.UpdateProfile")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 資料更新",
		logger.NewField("query", "UPDATE members SET name = ?, email = ? WHERE id = ?"),
		logger.NewField("member_id", m.ID),
		logger.NewField("member_email", m.Email),
	)

	result, err := s.db.ExecContext(ctx, queryUpdateMemberProfile, m.Name, m.Email, m.ID)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 資料更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", m.ID),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		contextLogger.Error("SQL 資料更新結果檢查失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", m.ID),
		)
		return nil, err
	}

	if rowsAffected == 0 {
		contextLogger.Error("SQL 資料更新未影響任何行",
			logger.NewField("member_id", m.ID),
		)
		return nil, ErrDBNoEffect
	}

	contextLogger.Debug("SQL 資料更新成功",
		logger.NewField("member_id", m.ID),
		logger.NewField("member_email", m.Email),
		logger.NewField("rows_affected", rowsAffected),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return m, nil
}
func (s sqlxMemberRepo) UpdateEmail(ctx context.Context, id int, email string) error {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.UpdateEmail")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL Email 更新",
		logger.NewField("query", "UPDATE members SET email = ? WHERE id = ?"),
		logger.NewField("member_id", id),
		logger.NewField("new_email", email),
	)

	result, err := s.db.ExecContext(ctx, queryUpdateMemberEmail, email, id)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL Email 更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("new_email", email),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return mapSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		contextLogger.Error("SQL Email 更新結果檢查失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}

	if rowsAffected == 0 {
		contextLogger.Error("SQL Email 更新未影響任何行",
			logger.NewField("member_id", id),
		)
		return ErrDBNoEffect
	}

	contextLogger.Debug("SQL Email 更新成功",
		logger.NewField("member_id", id),
		logger.NewField("new_email", email),
		logger.NewField("rows_affected", rowsAffected),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return nil
}
func (s sqlxMemberRepo) UpdatePassword(ctx context.Context, id int, password string) error {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.UpdatePassword")
	defer span.End()

	// 創建帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 密碼更新",
		logger.NewField("query", "UPDATE members SET password = ? WHERE id = ?"),
		logger.NewField("member_id", id),
	)

	result, err := s.db.ExecContext(ctx, queryUpdateMemberPassword, password, id)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 密碼更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return mapSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		contextLogger.Error("SQL 密碼更新結果檢查失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}

	if rowsAffected == 0 {
		contextLogger.Error("SQL 密碼更新未影響任何行",
			logger.NewField("member_id", id),
		)
		return ErrDBNoEffect
	}

	contextLogger.Debug("SQL 密碼更新成功",
		logger.NewField("member_id", id),
		logger.NewField("rows_affected", rowsAffected),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return nil
}
func (s sqlxMemberRepo) Delete(ctx context.Context, id int) error {
	// 創建 Repository 層的子 span
	ctx, span := s.tracer.Start(ctx, "MemberRepository.Delete")
	defer span.End()

	// 創庺帶有 context 的 logger 用於追蹤
	contextLogger := s.logger.WithContext(ctx).With(logger.NewField("layer", "repository"))

	startTime := time.Now()
	contextLogger.Debug("開始執行 SQL 刪除",
		logger.NewField("query", "DELETE FROM members WHERE id = ?"),
		logger.NewField("member_id", id),
	)

	result, err := s.db.ExecContext(ctx, queryDeleteMember, id)
	duration := time.Since(startTime)

	if err != nil {
		contextLogger.Error("SQL 刪除失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return mapSQLError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		contextLogger.Error("SQL 刪除結果檢查失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return err
	}

	if rows != 1 {
		contextLogger.Error("SQL 刪除未影響預期行數",
			logger.NewField("member_id", id),
			logger.NewField("rows_affected", rows),
		)
		return ErrDBNoEffect
	}

	contextLogger.Debug("SQL 刪除成功",
		logger.NewField("member_id", id),
		logger.NewField("rows_affected", rows),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return nil
}
