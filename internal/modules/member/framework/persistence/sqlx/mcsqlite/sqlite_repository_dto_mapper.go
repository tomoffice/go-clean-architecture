package mcsqlite

import (
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dao"
	"time"

	"github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx"
)

func sqlxModelToDTO(model *sqlx.MemberSQLXModel) (*dao.MemberRecord, error) {
	if model == nil {
		return nil, ErrMapperTimeParseFailed
	}
	if model.CreatedAt == "" {
		return nil, ErrMapperTimeParseFailed
	}
	daoCreateAt, err := time.Parse("2006-01-02 15:04:05", model.CreatedAt)
	if err != nil {
		return nil, ErrMapperTimeParseFailed
	}
	return &dao.MemberRecord{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: daoCreateAt,
	}, nil
}
