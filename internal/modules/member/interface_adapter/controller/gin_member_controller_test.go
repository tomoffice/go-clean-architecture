package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/controller/mock"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/interface_adapter/outputmodel"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/modules/member/usecase/port/input"
	"module-clean/internal/modules/member/usecase/port/output"
	"module-clean/internal/shared/enum"
	"module-clean/internal/shared/errorcode"
	sharedviewmodel "module-clean/internal/shared/viewmodel/http"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMemberController_Delete(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupPort   func(*mock.MockMemberInputPort, *mock.MockMemberPresenter)
		setupGinCtx func(ginCtx *gin.Context)
		want        *outputmodel.DeleteMemberResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						CreatedAt: testTime,
					}, nil)
				p.EXPECT().PresentDeleteMember(gomock.Any()).Return(outputmodel.DeleteMemberResponse{
					Data: dto.DeleteMemberResponseDTO{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						CreatedAt: testTime.Format(time.RFC3339),
					},
					Error: nil,
					Meta:  nil,
					BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
						UnixTimeStamp:    testTime.Unix(),
						RFC3339TimeStamp: testTime.Format(time.RFC3339),
						Status:           enum.APIStatusSuccess,
					},
				})
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: &outputmodel.DeleteMemberResponse{
				Data: dto.DeleteMemberResponseDTO{
					ID:        1,
					Name:      "test",
					Email:     "test@gmail.com",
					CreatedAt: testTime.Format(time.RFC3339),
				},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test binding error",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					})
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid_id"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test binding error",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test validation error",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"}, // 因為validation會檢查ID是否大於0，所以這裡使用-1來觸發驗證錯誤
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test validation error",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberNotFound)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - delete member failed",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberNoEffect)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNoEffect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)

			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 422,
		},
		{
			name: "usecase error - db error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberDBError)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - unexpected error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberUnexpectedError)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrUnexpectedMemberUseCaseError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrUnexpectedMemberUseCaseError),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrUnexpectedMemberUseCaseError),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := tt.fields.usecase.(*mock.MockMemberInputPort)
			mockPresenter := tt.fields.presenter.(*mock.MockMemberPresenter)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupPort(mockUseCase, mockPresenter)
			tt.setupGinCtx(ginCtx)
			c.Delete(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.DeleteMemberResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_GetByEmail(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(*mock.MockMemberInputPort, *mock.MockMemberPresenter)
		want        *outputmodel.GetMemberByEmailResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request, _ = http.NewRequest("GET", "/api/member?email=test@gmail.com", nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        0,
						Name:      "",
						Email:     "",
						CreatedAt: testTime}, nil)
				p.EXPECT().PresentGetMemberByEmail(gomock.Any()).Return(
					outputmodel.GetMemberByEmailResponse{
						Data: dto.GetMemberByEmailResponseDTO{
							ID:        0,
							Name:      "",
							Email:     "",
							CreatedAt: "",
						},
						Error: nil,
						Meta:  nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusSuccess,
						},
					},
				)
			},
			want: &outputmodel.GetMemberByEmailResponse{
				Data: dto.GetMemberByEmailResponseDTO{
					ID:        0,
					Name:      "",
					Email:     "",
					CreatedAt: "",
				},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					})
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 400,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=invalid-email",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test validation error",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test validation error",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=test@gmail.com",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberNotFound)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - db failure",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=test@gmail.com",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberDBError)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := tt.fields.usecase.(*mock.MockMemberInputPort)
			mockPresenter := tt.fields.presenter.(*mock.MockMemberPresenter)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.GetByEmail(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.GetMemberByEmailResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_GetByID(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	ctrl, testTime := portHelper(t)
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupPort   func(*mock.MockMemberInputPort, *mock.MockMemberPresenter)
		setupGinCtx func(ginCtx *gin.Context)
		want        *outputmodel.GetMemberByIDResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByID(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						CreatedAt: testTime,
					}, nil)
				p.EXPECT().PresentGetMemberByID(gomock.Any()).Return(
					outputmodel.GetMemberByIDResponse{
						Data: dto.GetMemberByIDResponseDTO{
							ID:        1,
							Name:      "test",
							Email:     "test@gmail.com",
							CreatedAt: testTime.Format(time.RFC3339),
						},
						Error: nil,
						Meta:  nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusSuccess,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: &outputmodel.GetMemberByIDResponse{
				Data: dto.GetMemberByIDResponseDTO{
					ID:        1,
					Name:      "test",
					Email:     "test@gmail.com",
					CreatedAt: testTime.Format(time.RFC3339),
				},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid_id"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"}, // 因為validation會檢查ID是否大於0，所以這裡使用-1來觸發驗證錯誤
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				// uc錯誤後傳給presenter輸出錯誤
				uc.EXPECT().GetMemberByID(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberNotFound)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - db failure",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByID(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrMemberDBError)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := tt.fields.usecase.(*mock.MockMemberInputPort)
			mockPresenter := tt.fields.presenter.(*mock.MockMemberPresenter)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.GetByID(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.GetMemberByIDResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_List(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	members := []*entity.Member{
		{
			ID:        1,
			Name:      "test1",
			Email:     "test1@gmail.com",
			Password:  "",
			CreatedAt: testTime,
		},
		{
			ID:        2,
			Name:      "test2",
			Email:     "test2@gmail.com",
			Password:  "",
			CreatedAt: testTime.Add(time.Hour),
		},
	}
	type testPagination struct {
		Page    int
		Limit   int
		SortBy  string
		OrderBy string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		setupPagination testPagination
		setupGinCtx     func(ginCtx *gin.Context, testArgs testPagination)
		setupPort       func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter, testArgs testPagination)
		want            *outputmodel.ListMemberResponse
		wantErr         *outputmodel.ErrorResponse
		wantStatus      int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPagination: testPagination{
				Page:    1,
				Limit:   10,
				SortBy:  "name",
				OrderBy: "asc",
			},
			setupGinCtx: func(ginCtx *gin.Context, testArgs testPagination) {
				//[GIN-debug] GET    /api/v1/members           --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).List-fm (2 handlers)
				ginCtx.Request = httptest.NewRequest("GET", "/api/v1/members?page="+strconv.Itoa(testArgs.Page)+"&limit="+strconv.Itoa(testArgs.Limit), nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter, testArgs testPagination) {
				listItem := make([]dto.ListMemberItemDTO, 0, len(members))
				for _, member := range members {
					listItem = append(listItem, dto.ListMemberItemDTO{
						ID:    member.ID,
						Name:  member.Name,
						Email: member.Email,
					})
				}
				uc.EXPECT().ListMembers(gomock.Any(), gomock.Any()).Return(
					members, len(members), nil)
				p.EXPECT().PresentListMembers(gomock.Any(), gomock.Any()).Return(
					outputmodel.ListMemberResponse{
						Data: dto.ListMemberResponseDTO{
							Members: listItem,
						},
						Error: nil,
						Meta: &sharedviewmodel.MetaPayload{
							Total:  len(members),
							Page:   testArgs.Page,
							Limit:  testArgs.Limit,
							Offset: (testArgs.Page - 1) * testArgs.Limit,
						},
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: &outputmodel.ListMemberResponse{
				Data: dto.ListMemberResponseDTO{
					Members: []dto.ListMemberItemDTO{
						{
							ID:    1,
							Name:  "test1",
							Email: "test1@gmail.com",
						},
						{
							ID:    2,
							Name:  "test2",
							Email: "test2@gmail.com",
						},
					},
				},
				Error: nil,
				Meta: &sharedviewmodel.MetaPayload{
					Total:  len(members),
					Page:   1,
					Limit:  10,
					Offset: 0,
				},
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPagination: testPagination{
				Page:    1,
				Limit:   10,
				SortBy:  "name",
				OrderBy: "asc",
			},
			setupGinCtx: func(ginCtx *gin.Context, testArgs testPagination) {
				ginCtx.Request, _ = http.NewRequest("GET", "/api/v1/members?page=invalid&limit="+strconv.Itoa(testArgs.Limit), nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter, testArgs testPagination) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPagination: testPagination{
				Page:    1,
				Limit:   10,
				SortBy:  "name",
				OrderBy: "asc",
			},
			setupGinCtx: func(ginCtx *gin.Context, testArgs testPagination) {
				ginCtx.Request, _ = http.NewRequest("GET", "/api/v1/members?page=-1&limit="+strconv.Itoa(testArgs.Limit), nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter, testArgs testPagination) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - GetAll db failure", // 沒有not found的情況，因為ListMembers有可能回應空切片
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPagination: testPagination{
				Page:    1,
				Limit:   10,
				SortBy:  "name",
				OrderBy: "asc",
			},
			setupGinCtx: func(ginCtx *gin.Context, testArgs testPagination) {
				ginCtx.Request, _ = http.NewRequest("GET", "/api/v1/members?page="+strconv.Itoa(testArgs.Page)+"&limit="+strconv.Itoa(testArgs.Limit), nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter, testArgs testPagination) {
				uc.EXPECT().ListMembers(gomock.Any(), gomock.Any()).Return(nil, 0, usecase.ErrMemberDBError)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)

			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mock.NewMockMemberInputPort(ctrl)
			mockPresenter := mock.NewMockMemberPresenter(ctrl)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx, tt.setupPagination)
			tt.setupPort(mockUseCase, mockPresenter, tt.setupPagination)
			c.List(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.ListMemberResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})

	}
}

func TestMemberController_Register(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter)
		want        *outputmodel.RegisterMemberResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] POST   /api/v1/members           --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).Register-fm (2 handlers)
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","email":"test@gmail.com","password":"test123"}`),
				)
				ginCtx.Request.Header.Set("Content-Type", "application/json")
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().RegisterMember(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						Password:  "test123",
						CreatedAt: testTime,
					}, nil)
				p.EXPECT().PresentRegisterMember(gomock.Any()).Return(
					outputmodel.RegisterMemberResponse{
						Data: dto.RegisterMemberResponseDTO{
							ID:    0,
							Name:  "test",
							Email: "test@gmail.com",
						},
						Error:            nil,
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					})
			},
			want: &outputmodel.RegisterMemberResponse{
				Data: dto.RegisterMemberResponseDTO{
					ID:    0,
					Name:  "test",
					Email: "test@gmail.com",
				},
				Error:            nil,
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","password":"test123"}`),
				)
				ginCtx.Request.Header.Set("Content-Type", "application/json")
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "invalid email format",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "invalid email format",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","email":"invalid-email","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "invalid email format",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrValidationFailed),
					Message: "invalid email format",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member already exists",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","email":"test@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().RegisterMember(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberAlreadyExists,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberAlreadyExists,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberAlreadyExists),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberAlreadyExists),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 409,
		},
		{
			name: "usecase error - db failure",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","email":"test@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().RegisterMember(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberDBError,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)

			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - GetByEmail not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"POST",
					"/api/v1/members",
					strings.NewReader(`{"name":"test","email":"test@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().RegisterMember(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberNotFound,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mock.NewMockMemberInputPort(ctrl)
			mockPresenter := mock.NewMockMemberPresenter(ctrl)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.Register(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.RegisterMemberResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_UpdateProfile(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter)
		want        *outputmodel.UpdateMemberProfileResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberProfile(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        1,
						Name:      "updated_name",
						Email:     "update@gmail.com",
						CreatedAt: testTime,
					}, nil)
				p.EXPECT().PresentUpdateMemberProfile(gomock.Any()).Return(
					outputmodel.UpdateMemberProfileResponse{
						Data: dto.UpdateMemberProfileResponseDTO{
							ID:   0,
							Name: "updated_name",
						},
						Error:            nil,
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: &outputmodel.UpdateMemberProfileResponse{
				Data: dto.UpdateMemberProfileResponseDTO{
					ID:   0,
					Name: "updated_name",
				},
				Error:            nil,
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error - invalid URI",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid-id"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/invalid-id",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "binding error - invalid body",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":123}`), // invalid name type
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - uri invalid",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/invalid-id",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - empty name",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":""}`), // empty name
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrValidationFailed),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - GetByID member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberProfile(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberNotFound,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - db error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberProfile(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberDBError,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBError),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - no effect",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id       --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateProfile-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1",
					strings.NewReader(`{"name":"new_name"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberProfile(gomock.Any(), gomock.Any()).Return(
					nil,
					usecase.ErrMemberNoEffect,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNoEffect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
					Message: "",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mock.NewMockMemberInputPort(ctrl)
			mockPresenter := mock.NewMockMemberPresenter(ctrl)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.UpdateProfile(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.UpdateMemberProfileResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_UpdateEmail(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter)
		want        *outputmodel.UpdateMemberEmailResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"test@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				p.EXPECT().PresentUpdateMemberEmail().Return(
					outputmodel.UpdateMemberEmailResponse{
						Data:  dto.UpdateMemberEmailResponseDTO{},
						Error: nil,
						Meta:  nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusSuccess,
						},
					},
				)
			},
			want: &outputmodel.UpdateMemberEmailResponse{
				Data:  dto.UpdateMemberEmailResponseDTO{},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error - invalid URI",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid-id"},
				}
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "binding error - invalid body",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - uri invalid",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"test@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - invalid email format",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"invalid-email","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrValidationFailed), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - password too short",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"test@gmail.com","password":"6666"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrValidationFailed), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"testupdate@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberNotFound,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberNotFound), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - db error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"testupdate@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberDBError,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberDBError), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - email already used",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"updatetest@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberEmailAlreadyExists,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberEmailAlreadyExists,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberEmailAlreadyExists),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberEmailAlreadyExists), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 409,
		},
		{
			name: "usecase error - same email",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"updateemail@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberUpdateSameEmail,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberUpdateSameEmail,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberUpdateSameEmail),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberUpdateSameEmail), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 409,
		},
		{
			name: "usecase error - password incorrect",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"updatetest@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberPasswordIncorrect,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberPasswordIncorrect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberPasswordIncorrect),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberPasswordIncorrect), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 401,
		},
		{
			name: "usecase error - no effect",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"updatetest@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberNoEffect,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNoEffect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberNoEffect), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 422,
		},
		{
			name: "usecase error - unknown error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/email --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdateEmail-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/email",
					strings.NewReader(`{"new_email":"updatetest@gmail.com","password":"test123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					errors.New("unknown error"),
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberGatewayError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberGatewayError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberGatewayError), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mock.NewMockMemberInputPort(ctrl)
			mockPresenter := mock.NewMockMemberPresenter(ctrl)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.UpdateEmail(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.UpdateMemberEmailResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_UpdatePassword(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter)
		want        *outputmodel.UpdateMemberPasswordResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				p.EXPECT().PresentUpdateMemberPassword().Return(
					outputmodel.UpdateMemberPasswordResponse{
						Data:  dto.UpdateMemberPasswordResponseDTO{},
						Error: nil,
						Meta:  nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusSuccess,
						},
					},
				)
			},
			want: &outputmodel.UpdateMemberPasswordResponse{
				Data:  dto.UpdateMemberPasswordResponseDTO{},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error - invalid URI",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid-id"},
				}
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - uri invalid",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "binding error - invalid body",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrInvalidParams), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - old password too short",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrValidationFailed), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "validation error - new password too short",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrValidationFailed, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrValidationFailed),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrValidationFailed), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - same password",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"oldpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberUpdateSamePassword,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberUpdateSamePassword,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberUpdateSamePassword),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberUpdateSamePassword), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 409,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberNotFound,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberNotFound), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - password incorrect",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"wrongpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberPasswordIncorrect,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberPasswordIncorrect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberPasswordIncorrect),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberPasswordIncorrect), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 401,
		},
		{
			name: "usecase error - no effect",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberNoEffect,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNoEffect,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNoEffect),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberNoEffect), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 422,
		},
		{
			name: "usecase error - db error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				//[GIN-debug] PATCH  /api/v1/members/:id/password --> module-clean/internal/modules/member/interface_adapter/controller.(*MemberController).UpdatePassword-fm (2 handlers)
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
				ginCtx.Request = httptest.NewRequest(
					"PATCH",
					"/api/v1/members/1/password",
					strings.NewReader(`{"old_password":"oldpass123","new_password":"newpass123"}`),
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().UpdateMemberPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					usecase.ErrMemberDBError,
				)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBError,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBError),
							Message: "",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data:             nil,
				Error:            &sharedviewmodel.ErrorPayload{Code: strconv.Itoa(errorcode.ErrMemberDBError), Message: ""},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := mock.NewMockMemberInputPort(ctrl)
			mockPresenter := mock.NewMockMemberPresenter(ctrl)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.UpdatePassword(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want == nil && tt.wantErr == nil {
				t.Fatalf("wantErr should be nil when want is not nil")
			}
			if tt.want != nil {
				var got outputmodel.UpdateMemberPasswordResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestNewMemberController(t *testing.T) {
	ctrl, _ := portHelper(t)
	usecaseGateway := mock.NewMockMemberInputPort(ctrl)
	presenterGateway := mock.NewMockMemberPresenter(ctrl)
	got := NewMemberController(usecaseGateway, presenterGateway)
	assert.NotNil(t, got)
	assert.Equal(t, usecaseGateway, got.usecase)
	assert.Equal(t, presenterGateway, got.presenter)
}
func portHelper(t *testing.T) (*gomock.Controller, time.Time) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(
		func() {
			ctrl.Finish()
		},
	)
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // Set a fixed time for tests
	return ctrl, testTime
}
func GinCtxHelper(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	return ginCtx, w
}
