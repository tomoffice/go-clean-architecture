package outputmodel

import (
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dto"
	sharedviewmodel "github.com/tomoffice/go-clean-architecture/internal/shared/viewmodel/http"
)

// mockgen 尚未支援泛型，所以用別名展開所有泛型返回類型
type RegisterMemberResponse = sharedviewmodel.HTTPResponse[dto.RegisterMemberResponseDTO]
type GetMemberByIDResponse = sharedviewmodel.HTTPResponse[dto.GetMemberByIDResponseDTO]
type GetMemberByEmailResponse = sharedviewmodel.HTTPResponse[dto.GetMemberByEmailResponseDTO]
type ListMemberResponse = sharedviewmodel.HTTPResponse[dto.ListMemberResponseDTO]
type UpdateMemberProfileResponse = sharedviewmodel.HTTPResponse[dto.UpdateMemberProfileResponseDTO]
type UpdateMemberEmailResponse = sharedviewmodel.HTTPResponse[dto.UpdateMemberEmailResponseDTO]
type UpdateMemberPasswordResponse = sharedviewmodel.HTTPResponse[dto.UpdateMemberPasswordResponseDTO]
type DeleteMemberResponse = sharedviewmodel.HTTPResponse[dto.DeleteMemberResponseDTO]

// 為 any 的情況也必須別名化
type ErrorResponse = sharedviewmodel.HTTPResponse[any]
