package http

import (
	"module-clean/internal/member/domain/entities"
	"module-clean/internal/member/interface_adapters/dto"
)

// ToMemberResponse 轉換單一 member entity 為回應 DTO
func ToMemberResponse(member *entities.Member) dto.MemberResponseDTO {
	return dto.MemberResponseDTO{
		ID:    member.ID,
		Name:  member.Name,
		Email: member.Email,
	}
}

// ToMemberListResponse 轉換多筆 member entity 為回應 DTO 陣列
func ToMemberListResponse(members []*entities.Member) []dto.MemberResponseDTO {
	result := make([]dto.MemberResponseDTO, len(members))
	for i, m := range members {
		result[i] = ToMemberResponse(m)
	}
	return result
}
