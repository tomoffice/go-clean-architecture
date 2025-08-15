package mcsqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	sqlx2 "github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dao"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"time"
)

// sqlxMemberRepo 實作 dao.MemberDAO
type sqlxMemberSqlite struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer tracer.Tracer
}

func NewSqlxMemberSqlite(db *sqlx.DB, log logger.Logger, tracer tracer.Tracer) dao.MemberDAO {
	baseLogger := log.With(logger.NewField("layer", "repository"))
	return &sqlxMemberSqlite{
		db:     db,
		logger: baseLogger,
		tracer: tracer,
	}
}
func (s sqlxMemberSqlite) Create(ctx context.Context, m *dao.MemberRecord) error {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.Create")
	defer span.End()

	startTime := time.Now()

	_, err := s.db.ExecContext(repoCtx, queryInsertMember, m.Name, m.Email, m.Password)
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
func (s sqlxMemberSqlite) GetByID(ctx context.Context, id int) (*dao.MemberRecord, error) {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.GetByID")
	defer span.End()
	startTime := time.Now()

	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(repoCtx, member, querySelectByID, id)
	duration := time.Since(startTime)
	if err != nil {
		contextLogger.Error("SQL 查詢(ID)失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}
	record, err := sqlxModelToDTO(member)
	if err != nil {
		contextLogger.Error("SQL 查詢(ID) DTO 轉換失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, err
	}
	contextLogger.Debug("SQL 查詢(ID)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return record, nil
}
func (s sqlxMemberSqlite) GetByEmail(ctx context.Context, email string) (*dao.MemberRecord, error) {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.GetByEmail")
	defer span.End()
	startTime := time.Now()

	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(repoCtx, member, querySelectByEmail, email)
	duration := time.Since(startTime)
	if err != nil {
		contextLogger.Error("SQL 查詢失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", email),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, mapSQLError(err)
	}
	record, err := sqlxModelToDTO(member)
	if err != nil {
		contextLogger.Error("SQL 查詢 DTO 轉換失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", email),
			logger.NewField("duration_ms", duration.Milliseconds()),
		)
		return nil, err
	}
	contextLogger.Debug("SQL 查詢成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", email),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return record, nil
}
func (s sqlxMemberSqlite) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*dao.MemberRecord, error) {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.GetAll")
	defer span.End()
	startTime := time.Now()
	query := fmt.Sprintf(querySelectAllBase, pagination.SortBy, pagination.OrderBy)

	members := make([]*sqlx2.MemberSQLXModel, 0)
	err := s.db.SelectContext(repoCtx, &members, query, pagination.Limit, pagination.Offset)
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
	records := make([]*dao.MemberRecord, 0, len(members))
	for _, member := range members {
		record, err := sqlxModelToDTO(member)
		if err != nil {
			contextLogger.Error("SQL 列表查詢 DTO 轉換失敗",
				logger.NewField("error", err),
				logger.NewField("member_id", member.ID),
				logger.NewField("duration_ms", duration.Milliseconds()),
			)
			return nil, err
		}
		records = append(records, record)
	}
	contextLogger.Debug("SQL 列表查詢成功",
		logger.NewField("count", len(members)),
		logger.NewField("duration_ms", duration.Milliseconds()),
	)
	return records, nil
}
func (s sqlxMemberSqlite) CountAll(ctx context.Context) (int, error) {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.CountAll")
	defer span.End()

	startTime := time.Now()

	var count int
	err := s.db.GetContext(repoCtx, &count, queryCountMembers)
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
func (s sqlxMemberSqlite) UpdateProfile(ctx context.Context, m *dao.MemberRecord) (*dao.MemberRecord, error) {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.UpdateProfile")
	defer span.End()

	startTime := time.Now()

	result, err := s.db.ExecContext(repoCtx, queryUpdateMemberProfile, m.Name, m.Email, m.ID)
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
func (s sqlxMemberSqlite) UpdateEmail(ctx context.Context, id int, email string) error {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.UpdateEmail")
	defer span.End()

	startTime := time.Now()

	result, err := s.db.ExecContext(repoCtx, queryUpdateMemberEmail, email, id)
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
func (s sqlxMemberSqlite) UpdatePassword(ctx context.Context, id int, password string) error {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.UpdatePassword")
	defer span.End()

	startTime := time.Now()

	result, err := s.db.ExecContext(repoCtx, queryUpdateMemberPassword, password, id)
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
func (s sqlxMemberSqlite) Delete(ctx context.Context, id int) error {
	// 創建帶有 context 的 logger 用於追蹤
	repoCtx, contextLogger, span := createTracedLogger(ctx, s.tracer, s.logger, "Repository.Delete")
	defer span.End()

	startTime := time.Now()

	result, err := s.db.ExecContext(repoCtx, queryDeleteMember, id)
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

func createTracedLogger(ctx context.Context, tr tracer.Tracer, log logger.Logger, operationName string) (context.Context, logger.Logger, tracer.Span) {
	repoCtx, span := tr.Start(ctx, operationName)
	lg := log.WithContext(repoCtx)
	return repoCtx, lg, span
}
