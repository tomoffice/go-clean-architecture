package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"module-clean/internal/common/pagination"
	"module-clean/internal/member/domain/entities"
	"module-clean/internal/member/domain/repository"
	dbError "module-clean/internal/member/infrastructure/errors"
)

type sqlxMemberRepo struct {
	db *sqlx.DB
}

func NewSQLXMemberRepo(db *sqlx.DB) repository.MemberRepository {
	return &sqlxMemberRepo{db: db}
}

func (s sqlxMemberRepo) Create(ctx context.Context, m *entities.Member) error {
	result, err := s.db.ExecContext(ctx, queryInsertMember, m.Name, m.Email, m.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return dbError.ErrDuplicateKey
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}

func (s sqlxMemberRepo) GetByID(ctx context.Context, id int) (*entities.Member, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	member := &entities.Member{}
	err := s.db.GetContext(ctx, member, querySelectByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbError.ErrRecordNotFound
		}
		return nil, err
	}
	return member, nil
}

func (s sqlxMemberRepo) GetByEmail(ctx context.Context, email string) (*entities.Member, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	member := &entities.Member{}
	err := s.db.GetContext(ctx, member, querySelectByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbError.ErrRecordNotFound
		}
		return nil, err
	}
	return member, nil
}

func (s sqlxMemberRepo) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	members := make([]*entities.Member, 0)
	query := fmt.Sprintf(querySelectAllBase, pagination.OrderBy, pagination.SortBy)
	err := s.db.SelectContext(ctx, &members, query, pagination.Limit, pagination.Offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbError.ErrRecordNotFound
		}
		return nil, err
	}
	return members, nil
}

func (s sqlxMemberRepo) Update(ctx context.Context, m *entities.Member) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	result, err := s.db.ExecContext(ctx, queryUpdateMember, m.Name, m.Email, m.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return dbError.ErrUpdateNoEffect
	}
	return nil
}

func (s sqlxMemberRepo) Delete(ctx context.Context, id int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	result, err := s.db.ExecContext(ctx, queryDeleteMember, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return dbError.ErrDeleteNoEffect
	}
	return nil
}
