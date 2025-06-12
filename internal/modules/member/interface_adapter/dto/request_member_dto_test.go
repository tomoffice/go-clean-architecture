package dto

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteMemberRequestDTO_Validate(t *testing.T) {
	type fields struct {
		ID int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid ID",
			fields: fields{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid ID",
			fields: fields{
				ID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &DeleteMemberRequestDTO{
				ID: tt.fields.ID,
			}
			err := dto.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestGetMemberByEmailRequestDTO_Validate(t *testing.T) {
	type fields struct {
		Email string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Email",
			fields: fields{
				Email: "test@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			fields: fields{
				Email: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &GetMemberByEmailRequestDTO{
				Email: tt.fields.Email,
			}
			if err := dto.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMemberByIDRequestDTO_Validate(t *testing.T) {
	type fields struct {
		ID int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid ID",
			fields: fields{
				ID: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid ID",
			fields: fields{
				ID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &GetMemberByIDRequestDTO{
				ID: tt.fields.ID,
			}
			if err := dto.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListMemberRequestDTO_Validate(t *testing.T) {
	type fields struct {
		Page    int
		Limit   int
		SortBy  string
		OrderBy string
	}

	tests := []struct {
		name         string
		fields       fields
		wantErrField *string
	}{
		{
			name: "Valid Request",
			fields: fields{
				Page:    1,
				Limit:   1,
				SortBy:  "id",
				OrderBy: "asc",
			},
			wantErrField: nil,
		},
		{
			name: "Invalid Page - zero",
			fields: fields{
				Page:    0,
				Limit:   1,
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "Page"),
		},
		{
			name: "Invalid Page - negative",
			fields: fields{
				Page:    -1,
				Limit:   1,
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "Page"),
		},
		{
			name: "Invalid Limit - zero",
			fields: fields{
				Page:    1,
				Limit:   0,
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "Limit"),
		},
		{
			name: "Invalid Limit - negative",
			fields: fields{
				Page:    1,
				Limit:   -1,
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "Limit"),
		},
		{
			name: "Invalid Limit - too large",
			fields: fields{
				Page:    1,
				Limit:   101, // assuming max limit is 100
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "Limit"),
		},
		{
			name: "Invalid SortBy - wrong format",
			fields: fields{
				Page:    1,
				Limit:   1,
				SortBy:  "gg",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "SortBy"),
		},
		{
			name: "Invalid OrderBy - wrong format",
			fields: fields{
				Page:    1,
				Limit:   1,
				SortBy:  "id",
				OrderBy: "gg",
			},
			wantErrField: stringPtr(t, "OrderBy"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &ListMemberRequestDTO{
				Page:    tt.fields.Page,
				Limit:   tt.fields.Limit,
				SortBy:  tt.fields.SortBy,
				OrderBy: tt.fields.OrderBy,
			}
			err := dto.Validate()
			if tt.wantErrField != nil {
				assert.Error(t, err)
				// 檢查錯誤欄位
				var validationErrField validator.ValidationErrors
				if errors.As(err, &validationErrField) {
					assert.Equal(t, *tt.wantErrField, validationErrField[0].Field())
					t.Logf("Expected error field: %s, got: %s", *tt.wantErrField, validationErrField[0].Field())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterMemberRequestDTO_Validate(t *testing.T) {
	type fields struct {
		Name     string
		Email    string
		Password string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErrField *string
	}{
		{
			name: "Valid Request",
			fields: fields{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "test123",
			},
			wantErrField: nil,
		},
		{
			name: "Invalid Name - empty",
			fields: fields{
				Name:     "",
				Email:    "test@gmail.com",
				Password: "test123",
			},
			wantErrField: stringPtr(t, "Name"),
		},
		{
			name: "Invalid Email - missing '@'",
			fields: fields{
				Name:     "test",
				Email:    "testgmail.com", // missing '@'
				Password: "test123",
			},
			wantErrField: stringPtr(t, "Email"),
		},
		{
			name: "Invalid Email - empty",
			fields: fields{
				Name:     "test",
				Email:    "",
				Password: "test123",
			},
			wantErrField: stringPtr(t, "Email"),
		},
		{
			name: "Invalid Password too short",
			fields: fields{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "123", // too short
			},
			wantErrField: stringPtr(t, "Password"),
		},
		{
			name: "Invalid Password empty",
			fields: fields{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "",
			},
			wantErrField: stringPtr(t, "Password"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &RegisterMemberRequestDTO{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			err := dto.Validate()
			if tt.wantErrField != nil {
				assert.Error(t, err)
				var validationErrField validator.ValidationErrors
				if errors.As(err, &validationErrField) {
					assert.Equal(t, *tt.wantErrField, validationErrField[0].Field())
					t.Logf("Expected error field: %s, got: %s", *tt.wantErrField, validationErrField[0].Field())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMemberEmailRequestDTO_Validate(t *testing.T) {
	type fields struct {
		ID       int
		NewEmail string
		Password string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErrField *string
	}{
		{
			name: "Valid Request",
			fields: fields{
				ID:       1,
				NewEmail: "test@gmail.com",
				Password: "test123",
			},
			wantErrField: nil,
		},
		{
			name: "Invalid ID - zero",
			fields: fields{
				ID:       0,
				NewEmail: "",
				Password: "",
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid ID - negative",
			fields: fields{
				ID:       -1,
				NewEmail: "",
				Password: "",
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid NewEmail - empty",
			fields: fields{
				ID:       1,
				NewEmail: "",
				Password: "",
			},
			wantErrField: stringPtr(t, "NewEmail"),
		},
		{
			name: "Invalid NewEmail - wrong format",
			fields: fields{
				ID:       1,
				NewEmail: "testgmail.com", // missing '@'
				Password: "",
			},
			wantErrField: stringPtr(t, "NewEmail"),
		},
		{
			name: "Invalid Password - empty",
			fields: fields{
				ID:       1,
				NewEmail: "test@gmail.com",
				Password: "",
			},
			wantErrField: stringPtr(t, "Password"),
		},
		{
			name: "Invalid Password - too short",
			fields: fields{
				ID:       1,
				NewEmail: "test@gmail.com",
				Password: "123", // too short
			},
			wantErrField: stringPtr(t, "Password"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &UpdateMemberEmailRequestDTO{
				ID:       tt.fields.ID,
				NewEmail: tt.fields.NewEmail,
				Password: tt.fields.Password,
			}
			err := dto.Validate()
			if tt.wantErrField != nil {
				assert.Error(t, err)
				var validationErrField validator.ValidationErrors
				if errors.As(err, &validationErrField) {
					assert.Equal(t, *tt.wantErrField, validationErrField[0].Field())
					t.Logf("Expected error field: %s, got: %s", *tt.wantErrField, validationErrField[0].Field())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMemberPasswordRequestDTO_Validate(t *testing.T) {
	type fields struct {
		ID          int
		OldPassword string
		NewPassword string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErrField *string
	}{
		{
			name: "Valid Request",
			fields: fields{
				ID:          1,
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			wantErrField: nil,
		},
		{
			name: "Invalid ID - zero",
			fields: fields{
				ID:          0,
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid ID - negative",
			fields: fields{
				ID:          -1,
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid OldPassword - empty",
			fields: fields{
				ID:          1,
				OldPassword: "",
				NewPassword: "newpassword",
			},
			wantErrField: stringPtr(t, "OldPassword"),
		},
		{
			name: "Invalid OldPassword - too short",
			fields: fields{
				ID:          1,
				OldPassword: "123",
				NewPassword: "newpassword",
			},
			wantErrField: stringPtr(t, "OldPassword"),
		},
		{
			name: "Invalid NewPassword - empty",
			fields: fields{
				ID:          1,
				OldPassword: "newpassword",
				NewPassword: "",
			},
			wantErrField: stringPtr(t, "NewPassword"),
		},
		{
			name: "Invalid NewPassword - too short",
			fields: fields{
				ID:          1,
				OldPassword: "oldpassword",
				NewPassword: "123", // too short
			},
			wantErrField: stringPtr(t, "NewPassword"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &UpdateMemberPasswordRequestDTO{
				ID:          tt.fields.ID,
				OldPassword: tt.fields.OldPassword,
				NewPassword: tt.fields.NewPassword,
			}
			err := dto.Validate()
			if tt.wantErrField != nil {
				assert.Error(t, err)
				var validationErrField validator.ValidationErrors
				if errors.As(err, &validationErrField) {
					assert.Equal(t, *tt.wantErrField, validationErrField[0].Field())
					t.Logf("Expected error field: %s, got: %s", *tt.wantErrField, validationErrField[0].Field())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMemberProfileRequestDTO_Validate(t *testing.T) {
	type fields struct {
		ID   int
		Name *string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErrField *string
	}{
		{
			name: "Valid Request",
			fields: fields{
				ID:   1,
				Name: stringPtr(t, "testuser"),
			},
			wantErrField: nil,
		},
		{
			name: "Invalid ID - zero",
			fields: fields{
				ID:   0,
				Name: stringPtr(t, "testuser"),
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid ID - negative",
			fields: fields{
				ID:   -1,
				Name: stringPtr(t, "testuser"),
			},
			wantErrField: stringPtr(t, "ID"),
		},
		{
			name: "Invalid Name - too short",
			fields: fields{
				ID:   1,
				Name: stringPtr(t, "ab"), // less than 3 characters
			},
			wantErrField: stringPtr(t, "Name"),
		},
		{
			name: "Invalid Name - too long",
			fields: fields{
				ID:   1,
				Name: stringPtr(t, "a very long name that exceeds the maximum length of twenty characters"), // more than 20 characters
			},
			wantErrField: stringPtr(t, "Name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &UpdateMemberProfileRequestDTO{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			err := dto.Validate()
			if tt.wantErrField != nil {
				assert.Error(t, err)
				var validationErrField validator.ValidationErrors
				if errors.As(err, &validationErrField) {
					assert.Equal(t, *tt.wantErrField, validationErrField[0].Field())
					t.Logf("Expected error field: %s, got: %s", *tt.wantErrField, validationErrField[0].Field())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func stringPtr(t *testing.T, s string) *string {
	t.Helper()
	return &s
}
