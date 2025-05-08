package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// mapSQLError 將常見的 SQL 錯誤轉換為結構化錯誤
func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return wrap(err, ErrDBRecordNotFound)
	}
	if errors.Is(err, sql.ErrConnDone) {
		return wrap(err, ErrDBConnectionClosed)
	}
	if errors.Is(err, sql.ErrTxDone) {
		return wrap(err, ErrDBTransactionDone)
	}
	//db.ExecContext()、db.QueryContext()、db.QueryRowContext()
	if errors.Is(err, context.DeadlineExceeded) {
		return wrap(err, ErrDBContextTimeout)
	}
	//db.ExecContext()、db.QueryContext()、db.QueryRowContext()
	if errors.Is(err, context.Canceled) {
		return wrap(err, ErrDBContextCanceled)
	}
	// sqlite 特有
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return wrap(err, ErrDBDuplicateKey)
	}
	return wrap(err, ErrDBUnknown)
}
func wrap(rawErr, customErr error) *DBError {
	return &DBError{
		CustomError: customErr,
		RawError:    rawErr,
	}
}
