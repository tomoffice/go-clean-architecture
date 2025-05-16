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
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/usecase/input_port"
	"module-clean/internal/shared/common/enum"
	"module-clean/internal/shared/common/pagination"
)

func CreateMemberDTOToEntity(request dto.CreateMemberRequestDTO) *entity.Member {
	return &entity.Member{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}
func GetMemberByIDDTOToEntity(request dto.GetMemberByIDRequestDTO) *entity.Member {
	return &entity.Member{
		ID: request.ID,
	}
}
func GetMemberByEmailDTOToEntity(request dto.GetMemberByEmailRequestDTO) *entity.Member {
	return &entity.Member{
		Email: request.Email,
	}
}
func ListMemberDTOToPagination(request dto.ListMemberRequestDTO) *pagination.Pagination {
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
func UpdateMemberDTOToInputModel(dto dto.UpdateMemberRequestDTO) *input_port.PatchUpdateMemberInput {
	return &input_port.PatchUpdateMemberInput{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
}
func DeleteMemberDTOToEntity(request dto.DeleteMemberRequestDTO) *entity.Member {
	return &entity.Member{
		ID: request.ID,
	}
}
