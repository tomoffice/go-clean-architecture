// Package mapper 負責處理資料轉換的相關邏輯，
// 是 interface_adapters 層與 usecase/domain 層之間的橋樑，
// 專門將外部資料（如 DTO）轉換成 usecase 所需的資料結構，
// 包含：input_model、domain entity、或共用格式如 common.pagination。
//
// 職責:
// - 不進行資料驗證與商業邏輯處理
// - 專注於欄位對應與結構轉換
// - 確保 usecase 層收到乾淨、一致的輸入格式
package mapper

import (
	"module-clean/internal/common/pagination"
	"module-clean/internal/member/domain"
	"module-clean/internal/member/interface_adapters/dto"
	"module-clean/internal/member/usecase"
	"module-clean/internal/shared/enum"
)

func CreateDTOtoEntity(request dto.CreateMemberRequest) *domain.Member {
	return &domain.Member{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}
func GetMemberByIDToEntity(request dto.GetMemberByIDRequest) *domain.Member {
	return &domain.Member{
		ID: request.ID,
	}
}
func GetMemberByEmailToEntity(request dto.GetMemberByEmailRequest) *domain.Member {
	return &domain.Member{
		Email: request.Email,
	}
}
func ListMemberToPagination(request dto.ListMemberRequest) *pagination.Pagination {
	var orderBy enum.OrderBy
	offset := (request.Page - 1) * request.Limit
	switch request.OrderBy {
	case "asc":
		orderBy = enum.OrderByAsc
	case "desc":
		orderBy = enum.OrderByDesc
	default:
		orderBy = enum.OrderByAsc
	}
	return &pagination.Pagination{
		Limit:   request.Limit,
		Offset:  offset,
		SortBy:  request.SortBy,
		OrderBy: orderBy,
	}
}
func UpdateDTOToInputModel(dto dto.UpdateMemberRequest) *usecase.PatchUpdateMemberInput {
	return &usecase.PatchUpdateMemberInput{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}
func DeleteDTOToEntity(request dto.DeleteMemberRequest) *domain.Member {
	return &domain.Member{
		ID: request.ID,
	}
}
