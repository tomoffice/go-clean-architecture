package repository

import "errors"

// 這裡定義的是 gateway 層會往 usecase 丟的錯誤型別，維護時常用 errors.Is 來判斷。

var (
	// ------- gateway 內部業務語意 -------
	// ErrGatewayMemberUnexpectedError 真的遇到沒辦法預期的錯誤，或沒被前面 case catch 到時 fallback 用這個（ex: 第三方套件怪錯、panic）。
	ErrGatewayMemberUnexpectedError = errors.New("gateway: member gateway unexpected error")
	// ErrGatewayMemberMappingError 轉換 repo model（DB 回來的原始資料）到 entity 失敗就用這個，像是型別不符、時間格式有誤等都會丟這個。
	ErrGatewayMemberMappingError = errors.New("gateway: mapping repo model to entity failed")
)
