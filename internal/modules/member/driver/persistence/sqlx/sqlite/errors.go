package sqlite

import (
	"errors"
)

// 公開的錯誤實例，可供外部使用 Is/As 判斷
var (
	// ErrDBRecordNotFound 查不到資料，通常是 SELECT 沒結果時用這個。
	ErrDBRecordNotFound = errors.New("db: record not found")

	// ErrDBDuplicateKey 違反唯一鍵（像 UNIQUE constraint 之類），大多出現在 INSERT/UPDATE。
	ErrDBDuplicateKey = errors.New("db: duplicate key conflict")

	// ErrDBNoEffect 操作沒有異動任何 row（像 UPDATE/DELETE 都沒改到東西）。
	ErrDBNoEffect = errors.New("db: operation had no effect")

	// ErrDBContextTimeout context 超時，通常是查太久、或 DB 回不來。
	ErrDBContextTimeout = errors.New("db: context deadline exceeded")

	// ErrDBContextCanceled context 被取消，像是 API 請求被中斷或主動取消。
	ErrDBContextCanceled = errors.New("db: context canceled")

	// ErrDBConnectionClosed 連線斷掉了（可能被關閉），這時多半要重連或直接回報錯誤。
	ErrDBConnectionClosed = errors.New("db: connection closed")

	// ErrDBTransactionDone 這個 transaction 已經 commit 或 rollback，不能再用。
	ErrDBTransactionDone = errors.New("db: transaction done")

	// ErrDBUnexpectedError 不知道怎麼歸類的 DB 錯誤（像第三方套件 bug、panic 等）。
	ErrDBUnexpectedError = errors.New("db: unexpected error")

	// ErrDBLastInsertId INSERT 沒拿到 last insert id，理論上很少遇到，但有時候 driver 會這樣。
	ErrDBLastInsertId = errors.New("db: failed to get last insert id")
)
