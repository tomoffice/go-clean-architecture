package repository

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dao"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/output"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
)

type MemberRepoGateway struct {
	dao    dao.MemberDAO
	logger logger.Logger
	tracer tracer.Tracer
}

func NewMemberRepoGateway(dao dao.MemberDAO, log logger.Logger, tracer tracer.Tracer) output.MemberPersistence {
	baseLogger := log.With(logger.NewField("layer", "gateway"))
	return MemberRepoGateway{
		dao:    dao,
		logger: baseLogger,
		tracer: tracer,
	}
}
func (g MemberRepoGateway) Create(ctx context.Context, m *entity.Member) error {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.Create")
	defer span.End()

	record := &dao.MemberRecord{
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	}
	if err := g.dao.Create(gatewayCtx, record); err != nil {
		traceLogger.Error("會員資料庫創建失敗", logger.NewField("error", err), logger.NewField("member_email", m.Email))
		return MapInfraErrorToUsecaseError(err)
	}
	traceLogger.Debug("會員資料庫創建成功", logger.NewField("member_email", m.Email))
	return nil
}

func (g MemberRepoGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.GetByID")
	defer span.End()

	record, err := g.dao.GetByID(gatewayCtx, id)
	if err != nil {
		traceLogger.Error("會員資料庫查詢(ID)失敗", logger.NewField("error", err), logger.NewField("member_id", id))
		return nil, MapInfraErrorToUsecaseError(err)
	}
	member := &entity.Member{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		Password:  "",
		CreatedAt: record.CreatedAt,
	}
	traceLogger.Debug("會員資料庫查詢(ID)成功", logger.NewField("member_id", member.ID), logger.NewField("member_email", member.Email))
	return member, nil
}

func (g MemberRepoGateway) GetByEmail(ctx context.Context, email string) (*entity.Member, error) {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.GetByEmail")
	defer span.End()

	record, err := g.dao.GetByEmail(gatewayCtx, email)
	if err != nil {
		traceLogger.Error("會員資料庫查詢(Email)失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", email),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}

	member := &entity.Member{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		Password:  "",
		CreatedAt: record.CreatedAt,
	}

	traceLogger.Debug("會員資料庫查詢(Email)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", email),
	)
	return member, nil
}

func (g MemberRepoGateway) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, error) {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.GetAll")
	defer span.End()

	records, err := g.dao.GetAll(gatewayCtx, pagination)
	if err != nil {
		traceLogger.Error("會員資料庫列表查詢失敗",
			logger.NewField("error", err),
			logger.NewField("limit", pagination.Limit),
			logger.NewField("offset", pagination.Offset),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}
	members := make([]*entity.Member, 0, len(records))
	for _, record := range records {
		members = append(members, &entity.Member{
			ID:        record.ID,
			Name:      record.Name,
			Email:     record.Email,
			Password:  "",
			CreatedAt: record.CreatedAt,
		})
	}
	traceLogger.Debug("會員資料庫列表查詢成功",
		logger.NewField("count", len(members)),
	)
	return members, nil
}

func (g MemberRepoGateway) UpdateProfile(ctx context.Context, m *entity.Member) (*entity.Member, error) {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.UpdateProfile")
	defer span.End()

	record := &dao.MemberRecord{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	}
	//讀已寫所以不用回傳
	_, err := g.dao.UpdateProfile(gatewayCtx, record)
	if err != nil {
		traceLogger.Error("會員資料庫資料更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", m.ID),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}

	traceLogger.Debug("會員資料庫資料更新成功",
		logger.NewField("member_id", m.ID),
		logger.NewField("member_email", m.Email),
	)
	return m, nil
}

func (g MemberRepoGateway) UpdateEmail(ctx context.Context, id int, newEmail string) error {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.UpdateEmail")
	defer span.End()

	err := g.dao.UpdateEmail(gatewayCtx, id, newEmail)
	if err != nil {
		traceLogger.Error("會員資料庫 Email 更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("new_email", newEmail),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	traceLogger.Debug("會員資料庫 Email 更新成功",
		logger.NewField("member_id", id),
		logger.NewField("new_email", newEmail),
	)
	return nil
}

func (g MemberRepoGateway) UpdatePassword(ctx context.Context, id int, newPassword string) error {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.UpdatePassword")
	defer span.End()

	err := g.dao.UpdatePassword(gatewayCtx, id, newPassword)
	if err != nil {
		traceLogger.Error("會員資料庫密碼更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	traceLogger.Debug("會員資料庫密碼更新成功",
		logger.NewField("member_id", id),
	)
	return nil
}

func (g MemberRepoGateway) Delete(ctx context.Context, id int) error {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.Delete")
	defer span.End()

	err := g.dao.Delete(gatewayCtx, id)
	if err != nil {
		traceLogger.Error("會員資料庫刪除失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	traceLogger.Debug("會員資料庫刪除成功",
		logger.NewField("member_id", id),
	)
	return nil
}

func (g MemberRepoGateway) CountAll(ctx context.Context) (int, error) {
	// 創建帶有 trace 的 logger 用於追蹤
	gatewayCtx, traceLogger, span := createTraceLogger(ctx, g.tracer, g.logger, "Gateway.CountAll")
	defer span.End()

	count, err := g.dao.CountAll(gatewayCtx)
	if err != nil {
		traceLogger.Error("會員資料庫總數查詢失敗",
			logger.NewField("error", err),
		)
		return 0, MapInfraErrorToUsecaseError(err)
	}

	traceLogger.Debug("會員資料庫總數查詢成功",
		logger.NewField("count", count),
	)
	return count, nil
}

// createTraceLogger 在 Gateway 層建立帶 Trace 的 Logger
func createTraceLogger(ctx context.Context, tr tracer.Tracer, log logger.Logger, operationName string) (context.Context, logger.Logger, tracer.Span) {
	gatewayCtx, span := tr.Start(ctx, operationName)
	lg := log.WithContext(gatewayCtx)
	return gatewayCtx, lg, span
}
