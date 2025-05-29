// Package dto 定義資料傳輸物件（Data Transfer Object），用於接收來自外部的請求資料，
// 例如 HTTP JSON 或 Query 參數，並對欄位進行格式驗證。
//
// 職責:
// - 定義輸入資料結構
// - 使用 validator 進行欄位驗證
// - 作為 interface_adapter 與 usecase 間的輸入格式橋樑
// - 不包含資料轉換、不依賴 domain 與 usecase
package dto

import (
	"github.com/go-playground/validator/v10"
)

type CreateMemberRequestDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (dto *CreateMemberRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}

type GetMemberByIDRequestDTO struct {
	ID int `validate:"required,numeric"`
}

func (dto *GetMemberByIDRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}

type GetMemberByEmailRequestDTO struct {
	Email string `validate:"required,email"`
}

func (dto *GetMemberByEmailRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}

type ListMemberRequestDTO struct {
	Page    int    ` validate:"required,min=1"`
	Limit   int    ` validate:"required,min=1,max=100"`
	SortBy  string ` validate:"omitempty,oneof=id name email created_at"`
	OrderBy string ` validate:"omitempty,oneof=asc desc"`
}

func (dto *ListMemberRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}

type UpdateMemberRequestDTO struct {
	ID       int     `json:"id" validate:"required,numeric"`
	Name     *string `json:"name,omitempty" validate:"omitempty"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

func (dto *UpdateMemberRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}

type DeleteMemberRequestDTO struct {
	ID int `validate:"required,numeric,gte=1"`
}

func (dto *DeleteMemberRequestDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}
