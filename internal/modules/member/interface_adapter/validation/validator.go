package validation

//go:generate mockgen -source=validator.go -destination=../../interface_adapter/controller/mock/mock_validator.go -package=mock
import "github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dto"

type Validator interface {
	ValidateRegisterMember(dto.RegisterMemberRequestDTO) error
	ValidateGetMemberByID(dto.GetMemberByIDRequestDTO) error
	ValidateGetMemberByEmail(dto.GetMemberByEmailRequestDTO) error
	ValidateListMember(dto.ListMemberRequestDTO) error
	ValidateUpdateProfile(dto.UpdateMemberProfileRequestDTO) error
	ValidateUpdateEmail(dto.UpdateMemberEmailRequestDTO) error
	ValidateUpdatePassword(dto.UpdateMemberPasswordRequestDTO) error
	ValidateDeleteMember(dto.DeleteMemberRequestDTO) error
}
