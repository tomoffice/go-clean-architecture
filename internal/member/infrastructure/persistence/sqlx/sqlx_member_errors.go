package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ( // 通用底層db錯誤
	ErrDBRecordNotFound  = errors.New("db: record not found")
	ErrDBDuplicateKey    = errors.New("db: duplicate key")
	ErrDBUpdateNoEffect  = errors.New("db: update affected 0 rows")
	ErrDBDeleteNoEffect  = errors.New("db: delete affected 0 rows")
	ErrDBContextTimeout  = errors.New("db: context deadline exceeded")
	ErrDBContextCanceled = errors.New("db: context canceled")
)

func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrDBContextTimeout
	}
	if errors.Is(err, context.Canceled) {
		return ErrDBContextCanceled
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrDBRecordNotFound
	}
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return ErrDBDuplicateKey
	}
	return err
}
