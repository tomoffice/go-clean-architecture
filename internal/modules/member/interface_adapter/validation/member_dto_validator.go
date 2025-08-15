package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/dto"
)

type MemberValidator struct {
	validator *validator.Validate
}

func NewMemberValidator() *MemberValidator {
	return &MemberValidator{
		validator: validator.New(),
	}
}
func (v *MemberValidator) ValidateRegisterMember(dto dto.RegisterMemberRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateGetMemberByID(dto dto.GetMemberByIDRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateGetMemberByEmail(dto dto.GetMemberByEmailRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateListMember(dto dto.ListMemberRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateUpdateProfile(dto dto.UpdateMemberProfileRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateUpdateEmail(dto dto.UpdateMemberEmailRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateUpdatePassword(dto dto.UpdateMemberPasswordRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
func (v *MemberValidator) ValidateDeleteMember(dto dto.DeleteMemberRequestDTO) error {
	if err := v.validator.Struct(dto); err != nil {
		return err
	}
	return nil
}
