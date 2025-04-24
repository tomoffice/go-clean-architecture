package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type DBError struct {
	Code    string
	Message string
	Err     error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// 定義錯誤代碼常數
const (
	RecordNotFound   = "DB_RECORD_NOT_FOUND"
	DuplicateKey     = "DB_DUPLICATE_KEY"
	UpdateNoEffect   = "DB_UPDATE_NO_EFFECT"
	DeleteNoEffect   = "DB_DELETE_NO_EFFECT"
	ContextTimeout   = "DB_CONTEXT_TIMEOUT"
	ContextCanceled  = "DB_CONTEXT_CANCELED"
	ConnectionClosed = "DB_CONNECTION_CLOSED"
	TransactionDone  = "DB_TX_DONE"
	Unknown          = "DB_UNKNOWN"
)

// 公開的錯誤實例，可供外部使用 Is/As 判斷
var (
	ErrDBRecordNotFound   = &DBError{Code: RecordNotFound, Message: "db: record not found"}
	ErrDBDuplicateKey     = &DBError{Code: DuplicateKey, Message: "db: duplicate key conflict"}
	ErrDBUpdateNoEffect   = &DBError{Code: UpdateNoEffect, Message: "db: update affected 0 rows"}
	ErrDBDeleteNoEffect   = &DBError{Code: DeleteNoEffect, Message: "db: delete affected 0 rows"}
	ErrDBContextTimeout   = &DBError{Code: ContextTimeout, Message: "db: context deadline exceeded"}
	ErrDBContextCanceled  = &DBError{Code: ContextCanceled, Message: "db: context canceled"}
	ErrDBConnectionClosed = &DBError{Code: ConnectionClosed, Message: "db: connection is already closed"}
	ErrDBTransactionDone  = &DBError{Code: TransactionDone, Message: "db: transaction already committed or rolled back"}
	ErrDBUnknown          = &DBError{Code: Unknown, Message: "db: unknown error"}
)

// mapSQLError 將常見的 SQL 錯誤轉換為結構化錯誤
func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return wrap(err, ErrDBRecordNotFound)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return wrap(err, ErrDBContextTimeout)
	}
	if errors.Is(err, context.Canceled) {
		return wrap(err, ErrDBContextCanceled)
	}
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return wrap(err, ErrDBDuplicateKey)
	}
	if errors.Is(err, sql.ErrConnDone) {
		return wrap(err, ErrDBConnectionClosed)
	}
	if errors.Is(err, sql.ErrTxDone) {
		return wrap(err, ErrDBTransactionDone)
	}
	return wrap(err, ErrDBUnknown)
}
func wrap(err error, dbErr *DBError) *DBError {
	return &DBError{
		Code:    dbErr.Code,
		Message: dbErr.Message,
		Err:     err,
	}
}
