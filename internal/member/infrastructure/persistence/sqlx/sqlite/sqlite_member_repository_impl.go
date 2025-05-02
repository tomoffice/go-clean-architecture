package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"module-clean/internal/common/pagination"
	"module-clean/internal/member/domain/entities"
	"module-clean/internal/member/domain/repository"
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
		return mapSQLError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}

func (s sqlxMemberRepo) GetByID(ctx context.Context, id int) (*entities.Member, error) {
	member := &entities.Member{}
	err := s.db.GetContext(ctx, member, querySelectByID, id)
	if err != nil {
		return nil, mapSQLError(err)
	}
	return member, nil
}

func (s sqlxMemberRepo) GetByEmail(ctx context.Context, email string) (*entities.Member, error) {
	member := &entities.Member{}
	err := s.db.GetContext(ctx, member, querySelectByEmail, email)
	if err != nil {
		return nil, mapSQLError(err)
	}
	return member, nil
}

func (s sqlxMemberRepo) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entities.Member, error) {
	members := make([]*entities.Member, 0)
	query := fmt.Sprintf(querySelectAllBase, pagination.SortBy, pagination.OrderBy)
	err := s.db.SelectContext(ctx, &members, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, mapSQLError(err)
	}
	return members, nil
}
func (s sqlxMemberRepo) CountAll(ctx context.Context) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, queryCountMembers)
	if err != nil {
		return 0, mapSQLError(err)
	}
	return count, nil
}
func (s sqlxMemberRepo) Update(ctx context.Context, m *entities.Member) (*entities.Member, error) {
	result, err := s.db.ExecContext(ctx, queryUpdateMember, m.Name, m.Email, m.ID)
	if err != nil {
		return nil, mapSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrDBUpdateNoEffect
	}
	return m, nil
}

func (s sqlxMemberRepo) Delete(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, queryDeleteMember, id)
	if err != nil {
		return mapSQLError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrDBDeleteNoEffect
	}
	return nil
}
