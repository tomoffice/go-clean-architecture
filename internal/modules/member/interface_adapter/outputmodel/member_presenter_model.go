package outputmodel

import (
	"module-clean/internal/modules/member/interface_adapter/dto"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
)

// mockgen 尚未支援泛型，所以用別名展開所有泛型返回類型
type CreateMemberResponse     = sharedviewmodel.HTTPResponse[dto.CreateMemberResponseDTO]
type GetMemberByIDResponse    = sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO]
type GetMemberByEmailResponse = sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO]
type ListMemberResponse       = sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO]
type UpdateMemberResponse     = sharedviewmodel.HTTPResponse[dto.UpdateMemberResponseDTO]
type DeleteMemberResponse     = sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO]

// 為 any 的情況也必須別名化
type ErrorResponse = sharedviewmodel.HTTPResponse[any]