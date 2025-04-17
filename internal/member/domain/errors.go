package domain

import "errors"

var ErrRecordNotFound = errors.New("record not found")     // 資料不存在
var ErrAlreadyExists = errors.New("record already exists") // 資料已存在
var ErrUpdateFailed = errors.New("update failed")          // 如果 rows != 1 代表更新失敗
var ErrDeleteFailed = errors.New("delete failed")          // 如果 rows != 1 代表刪除失敗
