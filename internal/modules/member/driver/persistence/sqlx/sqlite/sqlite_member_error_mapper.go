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
	// 檢查是否為自定義錯誤
	if errors.Is(err, sql.ErrNoRows) {
		return wrap(err, ErrDBRecordNotFound)
	}
	if errors.Is(err, sql.ErrConnDone) {
		return wrap(err, ErrDBConnectionClosed)
	}
	if errors.Is(err, sql.ErrTxDone) {
		return wrap(err, ErrDBTransactionDone)
	}
	// 任何帶 context 的 DB 操作（ExecContext, QueryContext, QueryRowContext, GetContext, SelectContext 等）都可能收到以下錯誤
	if errors.Is(err, context.DeadlineExceeded) {
		return wrap(err, ErrDBContextTimeout)
	}
	if errors.Is(err, context.Canceled) {
		return wrap(err, ErrDBContextCanceled)
	}
	// sqlite 特有
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return wrap(err, ErrDBDuplicateKey)
	}
	return wrap(err, ErrDBUnexpectedError)
}
func wrap(rawErr, customErr error) *DBError {
	return &DBError{
		CustomError: customErr,
		RawError:    rawErr,
	}
}
