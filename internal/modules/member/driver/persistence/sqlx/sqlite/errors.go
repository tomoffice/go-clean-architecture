package sqlite

import (
	"errors"
)

// 公開的錯誤實例，可供外部使用 Is/As 判斷
var (
	ErrDBRecordNotFound   = errors.New("db: record not found")
	ErrDBDuplicateKey     = errors.New("db: duplicate key conflict")
	ErrDBUpdateNoEffect   = errors.New("db: update no effect")
	ErrDBDeleteNoEffect   = errors.New("db: delete no effect")
	ErrDBContextTimeout   = errors.New("db: context deadline exceeded")
	ErrDBContextCanceled  = errors.New("db: context canceled")
	ErrDBConnectionClosed = errors.New("db: connection closed")
	ErrDBTransactionDone  = errors.New("db: transaction done")
	ErrDBUnknown          = errors.New("db: unknown errordefs")
	ErrDBLastInsertId     = errors.New("db: failed to get last insert id")
)
