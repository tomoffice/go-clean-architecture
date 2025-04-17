package sqlx

import (
	"context"
	"crud-clean/internal/domain/common"
	"crud-clean/internal/domain/entities"
	"crud-clean/internal/domain/repository"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type sqlxMemberRepo struct {
	db *sqlx.DB
}

func NewSQLXMemberRepo(db *sqlx.DB) repository.MemberRepository {
	return &sqlxMemberRepo{
		db: db,
	}
}

func (s sqlxMemberRepo) Create(ctx context.Context, m *entities.Member) error {
	// 執行 INSERT 並取得結果
	result, err := s.db.ExecContext(ctx, `INSERT INTO members (name, email) VALUES (?, ?)`, m.Name, m.Email)
	if err != nil {
		return err
	}

	// 把 auto-increment 的 ID 填入 domain entity
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
	err := s.db.GetContext(ctx, member, `SELECT * FROM members WHERE id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
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
	err := s.db.GetContext(ctx, member, `SELECT * FROM members WHERE email = ?`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}
	return member, nil
}

func (s sqlxMemberRepo) GetAll(ctx context.Context, pagination common.Pagination) ([]*entities.Member, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	members := make([]*entities.Member, 0)
	limit := pagination.Limit
	offset := pagination.Offset
	orderBy := pagination.OrderBy
	query := fmt.Sprintf("SELECT * FROM member ORDER BY %s %s LIMIT ? OFFSET ?", orderBy, pagination.SortBy)
	err := s.db.SelectContext(ctx, &members, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}
	return members, nil
}

func (s sqlxMemberRepo) Update(ctx context.Context, m *entities.Member) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	result, err := s.db.ExecContext(ctx, `UPDATE members SET name = ?, email = ? WHERE id = ?`, m.Name, m.Email, m.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrRecordNotFound
	}
	return nil
}

func (s sqlxMemberRepo) Delete(ctx context.Context, id int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM members WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return repository.ErrRecordNotFound
	}
	return nil
}
