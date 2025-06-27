package usecase

import (
	"errors"
)

// MemberUseCase 錯誤碼
var (
	//------- gateway error mapping -------
	// ErrMemberNotFound 查不到這個 member，通常是 repo 查無資料時會拋。
	ErrMemberNotFound = errors.New("usecase: member not found")
	// ErrMemberAlreadyExists 唯一值衝突，像 email/account 已存在時用這個。
	ErrMemberAlreadyExists = errors.New("usecase: member already exists")
	// ErrMemberNoEffect 有執行 update/delete 但 rows affected = 0，代表沒東西被改到。
	ErrMemberNoEffect = errors.New("usecase: member operation had no effect")
	// ErrMemberDBError 遇到 DB 或 infra 技術性問題（像連線失敗、tx 掛掉等）。
	ErrMemberDBError = errors.New("usecase: member db operation error")
	// ErrMemberUnexpectedError 萬一真的遇到不知道怎麼歸類的 bug，就收斂到這個。
	ErrMemberUnexpectedError = errors.New("usecase: member usecase unexpected error")
	// ErrMemberMappingError 從 repo model 轉換到 entity 時發生錯誤，像是型別不符、時間格式有誤等。
	ErrMemberMappingError = errors.New("usecase: member mapping repo model to entity failed")

	// ------- usecase 內部的業務語意 -------
	// ErrMemberUpdateSameEmail 嘗試改 email 結果新舊 email 一樣。
	ErrMemberUpdateSameEmail = errors.New("usecase: member use same email")
	// ErrMemberEmailAlreadyExists 想換 email，但新 email 已經被別人註冊了。
	ErrMemberEmailAlreadyExists = errors.New("usecase: member email already exists")
	// ErrMemberUpdateSamePassword 嘗試改密碼結果新舊密碼一樣。
	ErrMemberUpdateSamePassword = errors.New("usecase: member use same password")
	// ErrMemberPasswordIncorrect 密碼驗證沒過（ex: 修改 email/密碼時比對舊密碼不對）。
	ErrMemberPasswordIncorrect = errors.New("usecase: member password incorrect")
)
