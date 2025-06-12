package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	sqlx2 "module-clean/internal/modules/member/driver/persistence/sqlx"
	"module-clean/internal/shared/pagination"
)

type sqlxMemberRepo struct {
	db *sqlx.DB
}

func NewSQLXMemberRepo(db *sqlx.DB) sqlx2.MemberSQLXRepository {
	return &sqlxMemberRepo{db: db}
}
func (s sqlxMemberRepo) Create(ctx context.Context, m *sqlx2.MemberSQLXModel) error {
	_, err := s.db.ExecContext(ctx, queryInsertMember, m.Name, m.Email, m.Password)
	if err != nil {
		return mapSQLError(err)
	}

	//id, err := result.LastInsertId()
	//if err != nil {
	//	return err
	//}
	//m.ID = int(id)
	return nil
}
func (s sqlxMemberRepo) GetByID(ctx context.Context, id int) (*sqlx2.MemberSQLXModel, error) {
	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(ctx, member, querySelectByID, id)
	if err != nil {
		return nil, mapSQLError(err)
	}
	return member, nil
}
func (s sqlxMemberRepo) GetByEmail(ctx context.Context, email string) (*sqlx2.MemberSQLXModel, error) {
	member := &sqlx2.MemberSQLXModel{}
	err := s.db.GetContext(ctx, member, querySelectByEmail, email)
	if err != nil {
		return nil, mapSQLError(err)
	}
	return member, nil
}
func (s sqlxMemberRepo) GetAll(ctx context.Context, pagination pagination.Pagination) ([]*sqlx2.MemberSQLXModel, error) {
	members := make([]*sqlx2.MemberSQLXModel, 0)
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
func (s sqlxMemberRepo) UpdateProfile(ctx context.Context, m *sqlx2.MemberSQLXModel) (*sqlx2.MemberSQLXModel, error) {
	result, err := s.db.ExecContext(ctx, queryUpdateMemberProfile, m.Name, m.Email, m.ID)
	if err != nil {
		return nil, mapSQLError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrDBNoEffect
	}
	return m, nil
}
func (s sqlxMemberRepo) UpdateEmail(ctx context.Context, id int, email string) error {
	result, err := s.db.ExecContext(ctx, queryUpdateMemberEmail, email, id)
	if err != nil {
		return mapSQLError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDBNoEffect
	}
	return nil
}
func (s sqlxMemberRepo) UpdatePassword(ctx context.Context, id int, password string) error {
	result, err := s.db.ExecContext(ctx, queryUpdateMemberPassword, password, id)
	if err != nil {
		return mapSQLError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDBNoEffect
	}
	return nil
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
		return ErrDBNoEffect
	}
	return nil
}
