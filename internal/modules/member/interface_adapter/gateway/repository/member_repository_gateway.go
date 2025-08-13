package repository

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/entity"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/output"
	"github.com/tomoffice/go-clean-architecture/internal/shared/pagination"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"time"
)

type MemberSQLXGateway struct {
	infraRepo sqlx.MemberSQLXRepository
	logger    logger.Logger
	tracer    tracer.Tracer
}

func NewMemberSQLXGateway(infraRepo sqlx.MemberSQLXRepository, log logger.Logger, tracer tracer.Tracer) output.MemberPersistence {
	baseLogger := log.With(logger.NewField("layer", "gateway"))
	return MemberSQLXGateway{
		infraRepo: infraRepo,
		logger:    baseLogger,
		tracer:    tracer,
	}
}
func (g MemberSQLXGateway) Create(ctx context.Context, m *entity.Member) error {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	repoModel := &sqlx.MemberSQLXModel{
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	err := g.infraRepo.Create(gatewayCtx, repoModel)
	if err != nil {
		contextLogger.Error("會員資料庫創建失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", m.Email),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫創建成功",
		logger.NewField("member_email", m.Email),
	)
	return nil
}

func (g MemberSQLXGateway) GetByID(ctx context.Context, id int) (*entity.Member, error) {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	repoModel, err := g.infraRepo.GetByID(gatewayCtx, id)
	if err != nil {
		contextLogger.Error("會員資料庫查詢(ID)失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		contextLogger.Error("會員資料庫時間解析失敗",
			logger.NewField("error", err),
			logger.NewField("created_at", repoModel.CreatedAt),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}

	member := &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}

	contextLogger.Debug("會員資料庫查詢(ID)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", member.Email),
	)
	return member, nil
}

func (g MemberSQLXGateway) GetByEmail(ctx context.Context, email string) (*entity.Member, error) {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	repoModel, err := g.infraRepo.GetByEmail(gatewayCtx, email)
	if err != nil {
		contextLogger.Error("會員資料庫查詢(Email)失敗",
			logger.NewField("error", err),
			logger.NewField("member_email", email),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
	if err != nil {
		contextLogger.Error("會員資料庫時間解析失敗",
			logger.NewField("error", err),
			logger.NewField("created_at", repoModel.CreatedAt),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}

	member := &entity.Member{
		ID:        repoModel.ID,
		Name:      repoModel.Name,
		Email:     repoModel.Email,
		Password:  "",
		CreatedAt: createdAt,
	}

	contextLogger.Debug("會員資料庫查詢(Email)成功",
		logger.NewField("member_id", member.ID),
		logger.NewField("member_email", email),
	)
	return member, nil
}

func (g MemberSQLXGateway) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Member, error) {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	repoModels, err := g.infraRepo.GetAll(gatewayCtx, pagination)
	if err != nil {
		contextLogger.Error("會員資料庫列表查詢失敗",
			logger.NewField("error", err),
			logger.NewField("limit", pagination.Limit),
			logger.NewField("offset", pagination.Offset),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}
	var members []*entity.Member
	for _, repoModel := range repoModels {
		createdAt, err := time.Parse("2006-01-02 15:04:05", repoModel.CreatedAt)
		if err != nil {
			contextLogger.Error("會員資料庫時間解析失敗",
				logger.NewField("error", err),
				logger.NewField("created_at", repoModel.CreatedAt),
			)
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

	contextLogger.Debug("會員資料庫列表查詢成功",
		logger.NewField("count", len(members)),
	)
	return members, nil
}

func (g MemberSQLXGateway) UpdateProfile(ctx context.Context, m *entity.Member) (*entity.Member, error) {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	repoModel := &sqlx.MemberSQLXModel{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	//讀已寫所以不用回傳
	_, err := g.infraRepo.UpdateProfile(gatewayCtx, repoModel)
	if err != nil {
		contextLogger.Error("會員資料庫資料更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", m.ID),
		)
		return nil, MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫資料更新成功",
		logger.NewField("member_id", m.ID),
		logger.NewField("member_email", m.Email),
	)
	return m, nil
}

func (g MemberSQLXGateway) UpdateEmail(ctx context.Context, id int, newEmail string) error {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	err := g.infraRepo.UpdateEmail(gatewayCtx, id, newEmail)
	if err != nil {
		contextLogger.Error("會員資料庫 Email 更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
			logger.NewField("new_email", newEmail),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫 Email 更新成功",
		logger.NewField("member_id", id),
		logger.NewField("new_email", newEmail),
	)
	return nil
}

func (g MemberSQLXGateway) UpdatePassword(ctx context.Context, id int, newPassword string) error {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	err := g.infraRepo.UpdatePassword(gatewayCtx, id, newPassword)
	if err != nil {
		contextLogger.Error("會員資料庫密碼更新失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫密碼更新成功",
		logger.NewField("member_id", id),
	)
	return nil
}

func (g MemberSQLXGateway) Delete(ctx context.Context, id int) error {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	err := g.infraRepo.Delete(gatewayCtx, id)
	if err != nil {
		contextLogger.Error("會員資料庫刪除失敗",
			logger.NewField("error", err),
			logger.NewField("member_id", id),
		)
		return MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫刪除成功",
		logger.NewField("member_id", id),
	)
	return nil
}

func (g MemberSQLXGateway) CountAll(ctx context.Context) (int, error) {
	// 創建帶有 context 的 logger 用於追蹤
	gatewayCtx, contextLogger, span := createTracedLogger(ctx, g.tracer, g.logger)
	defer span.End()


	count, err := g.infraRepo.CountAll(gatewayCtx)
	if err != nil {
		contextLogger.Error("會員資料庫總數查詢失敗",
			logger.NewField("error", err),
		)
		return 0, MapInfraErrorToUsecaseError(err)
	}

	contextLogger.Debug("會員資料庫總數查詢成功",
		logger.NewField("count", count),
	)
	return count, nil
}

func createTracedLogger(ctx context.Context, tr tracer.Tracer, log logger.Logger) (context.Context, logger.Logger, tracer.Span) {
	gatewayCtx, span := tr.Start(ctx, "")
	lg := log.WithContext(gatewayCtx)
	return gatewayCtx, lg, span
}
