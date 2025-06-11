package usecase

import (
	"errors"
)

// MemberUseCase 錯誤碼
var (
	//------- gateway error mapping -------
	// ErrUseCaseMemberNotFound 查不到這個 member，通常是 repo 查無資料時會拋。
	ErrUseCaseMemberNotFound = errors.New("usecase: member not found")
	// ErrUseCaseMemberAlreadyExists 唯一值衝突，像 email/account 已存在時用這個。
	ErrUseCaseMemberAlreadyExists = errors.New("usecase: member already exists")
	// ErrUseCaseMemberNoEffect 有執行 update/delete 但 rows affected = 0，代表沒東西被改到。
	ErrUseCaseMemberNoEffect = errors.New("usecase: member operation had no effect")
	// ErrUseCaseMemberDBError 遇到 DB 或 infra 技術性問題（像連線失敗、tx 掛掉等）。
	ErrUseCaseMemberDBError = errors.New("usecase: member db operation error")
	// ErrUseCaseMemberUnexpectedError 萬一真的遇到不知道怎麼歸類的 bug，就收斂到這個。
	ErrUseCaseMemberUnexpectedError = errors.New("usecase: member usecase unexpected error")
	// ErrUseCaseMemberGatewayError 外部服務（像第三方 API/gateway）調用失敗，目前都用 DBError 了所以基本上沒在用。暫時沒用到
	ErrUseCaseMemberGatewayError = errors.New("usecase: member gateway error")

	// ------- usecase 內部的業務語意 -------
	// ErrUseCaseMemberUpdateSameEmail 嘗試改 email 結果新舊 email 一樣。
	ErrUseCaseMemberUpdateSameEmail = errors.New("usecase: member use same email")
	// ErrUseCaseMemberEmailAlreadyExists 想換 email，但新 email 已經被別人註冊了。
	ErrUseCaseMemberEmailAlreadyExists = errors.New("usecase: member email already exists")
	// ErrUseCaseMemberUpdateSamePassword 嘗試改密碼結果新舊密碼一樣。
	ErrUseCaseMemberUpdateSamePassword = errors.New("usecase: member use same password")
	// ErrUseCaseMemberPasswordIncorrect 密碼驗證沒過（ex: 修改 email/密碼時比對舊密碼不對）。
	ErrUseCaseMemberPasswordIncorrect = errors.New("usecase: member password incorrect")

)
