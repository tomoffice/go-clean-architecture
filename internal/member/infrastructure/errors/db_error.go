package errors

import "errors"

var ( // 通用底層錯誤
	ErrRecordNotFound = errors.New("db: record not found")       // 查無資料（SELECT, DELETE）
	ErrInsertFailed   = errors.New("db: insert failed")          // INSERT 發生錯誤（不常用）
	ErrUpdateFailed   = errors.New("db: update affected 0 rows") // UPDATE 沒有影響資料
	ErrDeleteFailed   = errors.New("db: delete affected 0 rows") // DELETE 沒有影響資料
	ErrDuplicateKey   = errors.New("db: duplicate key conflict") // UNIQUE KEY 衝突（email 已存在）
	ErrUpdateNoEffect = errors.New("db: update no effect")       // UPDATE 沒有影響資料（不常用）
	ErrDeleteNoEffect = errors.New("db: delete no effect")       // DELETE 沒有影響資料（不常用）
)
