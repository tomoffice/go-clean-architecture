package repository

import "errors"

// 這裡定義的是 gateway 層會往 usecase 丟的錯誤型別，維護時常用 errors.Is 來判斷。

var (
	// -------infra to gateway error mapping-------
	// ErrGatewayMemberNotFound 查不到這個 member，通常是 id 或 email 不存在時用這個。
	ErrGatewayMemberNotFound = errors.New("gateway: member not found")
	// ErrGatewayMemberAlreadyExists 嘗試建立或更新時遇到唯一值衝突（通常 DB 唯一鍵違反），就會回這個。
	ErrGatewayMemberAlreadyExists = errors.New("gateway: member already exists")
	// ErrGatewayMemberNoEffect 有執行 update 或 delete，但 rows affected = 0，代表沒有東西被改到（不管是不存在還是內容一樣都算）。
	ErrGatewayMemberNoEffect = errors.New("gateway: member operation had no effect")
	// ErrGatewayMemberDBError DB 操作異常，但不是上面那些語意錯誤，像連線掛掉、tx 異常、context timeout 之類都丟這個。
	ErrGatewayMemberDBError = errors.New("gateway: member db operation error")
	// ------- gateway 內部業務語意 -------
	// ErrGatewayMemberUnexpectedError 真的遇到沒辦法預期的錯誤，或沒被前面 case catch 到時 fallback 用這個（ex: 第三方套件怪錯、panic）。
	ErrGatewayMemberUnexpectedError = errors.New("gateway: member gateway unexpected error")
	// ErrGatewayMemberMappingError 轉換 repo model（DB 回來的原始資料）到 entity 失敗就用這個，像是型別不符、時間格式有誤等都會丟這個。
	ErrGatewayMemberMappingError = errors.New("gateway: mapping repo model to entity failed")
)