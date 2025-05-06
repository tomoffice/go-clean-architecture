// Package mapper 負責處理資料轉換的相關邏輯，
// 是 interface_adapter 層與 usecase/domain 層之間的橋樑，
// 專門將外部資料（如 DTO）轉換成 usecase 所需的資料結構，
// 包含：input_model、domain entity、或共用格式如 common.pagination。
//
// 職責:
// - 不進行資料驗證與商業邏輯處理
// - 專注於欄位對應與結構轉換
// - 確保 usecase 層收到乾淨、一致的輸入格式
package mapper

import (
	"module-clean/internal/modules/member/domain/entities"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/shared/enum"
	"module-clean/internal/shared/pagination"
)

func CreateDTOtoEntity(request dto.CreateMemberRequestDTO) *entities.Member {
	return &entities.Member{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}
func GetMemberByIDDTOToEntity(request dto.GetMemberByIDRequestDTO) *entities.Member {
	return &entities.Member{
		ID: request.ID,
	}
}
func GetMemberByEmailDTOToEntity(request dto.GetMemberByEmailRequestDTO) *entities.Member {
	return &entities.Member{
		Email: request.Email,
	}
}
func ListMemberToPagination(request dto.ListMemberRequestDTO) *pagination.Pagination {
	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	var orderBy enum.OrderBy
	switch request.OrderBy {
	case "asc":
		orderBy = enum.OrderByAsc
	case "desc":
		orderBy = enum.OrderByDesc
	default:
		orderBy = enum.OrderByAsc
	}
	offset := (request.Page - 1) * request.Limit
	return &pagination.Pagination{
		Limit:   request.Limit,
		Offset:  offset,
		SortBy:  sortBy,
		OrderBy: orderBy,
	}
}
func UpdateDTOToInputModel(dto dto.UpdateMemberRequestDTO) *usecase.PatchUpdateMemberInput {
	return &usecase.PatchUpdateMemberInput{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}
func DeleteDTOToEntity(request dto.DeleteMemberRequestDTO) *entities.Member {
	return &entities.Member{
		ID: request.ID,
	}
}
